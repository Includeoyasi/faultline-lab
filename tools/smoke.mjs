// Cross-platform host-side Docker orchestration; no Docker socket in containers.
import { execFile } from 'node:child_process';
import { randomUUID } from 'node:crypto';
import { mkdir, writeFile } from 'node:fs/promises';
import path from 'node:path';
import { promisify } from 'node:util';
import { fileURLToPath } from 'node:url';

const exec = promisify(execFile);
const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const runID = randomUUID();
const project = `faultline-smoke-${runID.slice(0, 8)}`;
const directory = path.join(root, 'artifacts', 'smoke', runID);
const services = ['orders', 'inventory', 'payments'];
const report = {
  run_id: runID, started_at: new Date().toISOString(), result: 'FAIL',
  host_platform: process.platform, node: process.version, project,
  services: {}, commands: [], limitations: ['No business operations or capacity measurement'],
};
let phase = 'prerequisites';
let containersMayExist = false;

async function command(program, args, timeout = 60_000) {
  const observation = { program, args, started_at: new Date().toISOString() };
  report.commands.push(observation);
  try {
    const result = await exec(program, args, {
      cwd: root, timeout, maxBuffer: 16 * 1024 * 1024, windowsHide: true,
      env: { ...process.env, FL_RUNTIME_ENV_FILE: 'configs/env/container.env.example' },
    });
    observation.exit_code = 0;
    return result;
  } catch (error) {
    observation.exit_code = error.code ?? 'UNKNOWN';
    observation.error = String(error.code ?? error.name);
    // Keep build/runtime diagnostics without printing configuration or inspect Env.
    await writeFile(path.join(directory, `command-${report.commands.length}.log`), `${error.stdout ?? ''}${error.stderr ?? ''}`);
    throw new Error(`${program} ${args.join(' ')} failed (${observation.exit_code})`, { cause: error });
  } finally {
    observation.finished_at = new Date().toISOString();
  }
}

const docker = args => command('docker', args);
const composeArgs = args => ['compose', '-f', 'compose.yaml', '-p', project, ...args];
const compose = args => docker(composeArgs(args));
const requireValue = (condition, message) => { if (!condition) throw new Error(message); };

async function containerID(service) {
  const { stdout } = await compose(['ps', '--all', '-q', service]);
  const id = stdout.trim();
  requireValue(/^[a-f0-9]{12,64}$/.test(id), `Missing unique container for ${service}`);
  return id;
}

async function inspect(service) {
  const id = await containerID(service);
  const { stdout } = await docker(['inspect', id]);
  return { id, data: JSON.parse(stdout)[0] };
}

async function checkRunning(service) {
  await compose(['exec', '-T', service, '/app/service', 'healthcheck']);
  const { id, data } = await inspect(service);
  const host = data.HostConfig;
  requireValue(data.State.Status === 'running' && data.State.Health?.Status === 'healthy', `${service} is not healthy`);
  requireValue(data.Config.User === '65532:65532', `${service} must run as an unprivileged numeric user`);
  requireValue(host.ReadonlyRootfs && !host.Privileged, `${service} must have a read-only, unprivileged runtime`);
  requireValue(host.CapDrop?.includes('ALL'), `${service} must drop capabilities`);
  requireValue(host.SecurityOpt?.some(value => value.startsWith('no-new-privileges')), `${service} must prevent privilege escalation`);
  requireValue(Object.keys(host.PortBindings ?? {}).length === 0, `${service} unexpectedly publishes a host port`);
  requireValue(data.Mounts.length === 0, `${service} unexpectedly mounts host or volume data`);
  const networkIDs = Object.values(data.NetworkSettings.Networks).map(value => value.NetworkID);
  requireValue(networkIDs.length > 0, `${service} has no runtime network`);
  for (const networkID of networkIDs) {
    const { stdout } = await docker(['network', 'inspect', networkID]);
    requireValue(JSON.parse(stdout)[0].Internal === true, `${service} must use an internal network`);
  }
  report.services[service] = {
    ...report.services[service], id, image_id: data.Image,
    architecture: report.docker_architecture, healthy: true,
    security: { user: data.Config.User, read_only: true, capabilities_dropped: true, no_new_privileges: true, host_ports: [], mounts: [], internal_network: true },
  };
}

async function collectLogs(service) {
  const id = await containerID(service);
  const { stdout, stderr } = await docker(['logs', id]);
  const text = stdout + stderr;
  await writeFile(path.join(directory, `${service}.jsonl`), text);
  return text.split(/\r?\n/).filter(Boolean).map(line => JSON.parse(line));
}

async function checkStopped(service) {
  const { data } = await inspect(service);
  requireValue(data.State.Status === 'exited' && data.State.ExitCode === 0 && !data.State.OOMKilled, `${service} failed normal SIGTERM shutdown`);
  const logs = await collectLogs(service);
  for (const record of logs) {
    requireValue(record.service === service && typeof record.time === 'string' && typeof record.level === 'string' && typeof record.msg === 'string', `${service} emitted an invalid JSON record`);
  }
  const started = logs.findIndex(record => record.msg === 'service started');
  const stopping = logs.findIndex(record => record.msg === 'service stopping' && record.health === 'NOT_SERVING');
  const stopped = logs.findIndex(record => record.msg === 'service stopped' && record.stop_mode === 'graceful');
  requireValue(started >= 0 && stopping > started && stopped > stopping, `${service} lacks the ordered start/readiness/stop records`);
  report.services[service] = { ...report.services[service], stopped: true, exit_code: 0, oom_killed: false, stop_mode: 'graceful' };
}

await mkdir(directory, { recursive: true });
try {
  console.log(`Runtime smoke ${project}: prerequisites`);
  const version = await docker(['version', '--format', '{{.Server.Version}}']);
  const platform = await docker(['info', '--format', '{{.OSType}}/{{.Architecture}}']);
  requireValue(platform.stdout.trim().startsWith('linux/'), 'Docker must run Linux containers');
  report.docker_version = version.stdout.trim();
  report.docker_architecture = platform.stdout.trim();
  report.compose_version = (await docker(['compose', 'version', '--short'])).stdout.trim();
  report.source_revision = (await command('git', ['rev-parse', 'HEAD'])).stdout.trim();
  report.source_dirty = (await command('git', ['status', '--porcelain'])).stdout.trim() !== '';
  phase = 'build';
  console.log('Building three runtime images from the pinned toolchain');
  await compose(['config', '--quiet']);
  await command('docker', composeArgs(['build', ...services]), 600_000);
  phase = 'start';
  containersMayExist = true;
  await compose(['up', '-d', '--wait', '--wait-timeout', '60', ...services]);
  for (const service of services) await checkRunning(service);
  phase = 'isolation';
  await compose(['stop', 'orders']);
  await checkStopped('orders');
  for (const service of ['inventory', 'payments']) await checkRunning(service);
  report.isolation_after_orders_stop = 'PASS';
  phase = 'shutdown';
  await compose(['stop', 'inventory', 'payments']);
  for (const service of ['inventory', 'payments']) await checkStopped(service);
  report.result = 'PASS';
} catch (error) {
  report.result = phase === 'prerequisites' ? 'NOT_RUN' : 'FAIL';
  report.failed_phase = phase;
  report.error = error.message;
  console.error(`${report.result}: ${phase}: ${error.message}`);
} finally {
  if (containersMayExist) {
    for (const service of services) {
      try { await collectLogs(service); } catch (error) {
        report.services[service] = { ...report.services[service], diagnostic_error: error.message };
      }
    }
    try {
      await compose(['down', '--remove-orphans']);
      report.cleanup = 'PASS';
    } catch (error) {
      report.cleanup = 'FAIL';
      report.cleanup_error = error.message;
      report.result = 'FAIL';
    }
  } else report.cleanup = 'NOT_NEEDED';
  report.finished_at = new Date().toISOString();
  await writeFile(path.join(directory, 'report.json'), JSON.stringify(report, null, 2) + '\n');
  console.log(`${report.result}: ${path.relative(root, path.join(directory, 'report.json'))}`);
  process.exitCode = report.result === 'PASS' ? 0 : report.result === 'NOT_RUN' ? 2 : 1;
}
