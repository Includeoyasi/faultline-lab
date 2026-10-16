import { readFile, readdir } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { containsMachinePath } from './doc_safety.mjs';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const excluded = new Set(['.git', 'node_modules', '.local', 'artifacts', 'tmp']);
const decoder = new TextDecoder('utf-8', { fatal: true });
const errors = [];

async function collect(directory) {
  const entries = await readdir(directory, { withFileTypes: true });
  const result = [];
  for (const entry of entries) {
    if (excluded.has(entry.name)) continue;
    const absolute = path.join(directory, entry.name);
    if (entry.isDirectory()) result.push(...await collect(absolute));
    else if (entry.isFile()) result.push(absolute);
  }
  return result;
}

const files = await collect(root);
const knownPaths = new Set(files.map(file => path.relative(root, file).split(path.sep).join('/')));
const markdown = new Map();
for (const absolute of files) {
  if (!/\.(md|json|svg|ps1|mjs|ya?ml)$/.test(absolute)) continue;
  const relative = path.relative(root, absolute).split(path.sep).join('/');
  try {
    const text = decoder.decode(await readFile(absolute));
    if (text.includes('\r')) errors.push(`${relative}: use LF line endings for reproducible fingerprints`);
    if (containsMachinePath(text)) errors.push(`${relative}: contains a machine-specific user path`);
    if (/-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/.test(text)) errors.push(`${relative}: contains a private key`);
    if (relative.endsWith('.md')) markdown.set(relative, text);
    if (relative.endsWith('.json')) JSON.parse(text);
  } catch (error) {
    errors.push(`${relative}: ${error.message}`);
  }
}

function anchors(text) {
  const result = new Set([...text.matchAll(/<a\s+id="([^"]+)"/g)].map(match => match[1]));
  const seen = new Map();
  let fence = null;
  for (const line of text.split('\n')) {
    const marker = line.match(/^\s{0,3}(`{3,}|~{3,})/);
    if (marker) {
      if (!fence) fence = marker[1][0];
      else if (marker[1][0] === fence) fence = null;
      continue;
    }
    if (fence) continue;
    const heading = line.match(/^#{1,6}\s+(.+?)(?:\s+#+)?$/);
    if (!heading) continue;
    const slug = heading[1].toLowerCase().replace(/<[^>]*>/g, '').replace(/[^\p{L}\p{N}_\-\s]/gu, '').replace(/ /g, '-');
    const duplicate = seen.get(slug) ?? 0;
    seen.set(slug, duplicate + 1);
    result.add(duplicate ? `${slug}-${duplicate}` : slug);
  }
  return result;
}

let linkCount = 0;
let diagrams = 0;
for (const [relative, text] of markdown) {
  let fence = null;
  const visible = [];
  for (const line of text.split('\n')) {
    const marker = line.match(/^\s{0,3}(`{3,}|~{3,})(.*)$/);
    if (marker) {
      if (!fence) {
        fence = { char: marker[1][0], length: marker[1].length };
        if (marker[2].trim() === 'mermaid') diagrams++;
      } else if (marker[1][0] === fence.char && marker[1].length >= fence.length && !marker[2].trim()) {
        fence = null;
      }
      continue;
    }
    if (!fence) visible.push(line);
  }
  if (fence) errors.push(`${relative}: unclosed code fence`);
  const content = visible.join('\n');
  const links = [
    ...[...content.matchAll(/!?\[[^\]\n]*\]\(([^)\n]+)\)/g)].map(match => match[1]),
    ...[...content.matchAll(/(?:href|src)="([^"]+)"/g)].map(match => match[1]),
  ];
  for (const raw of links) {
    const target = raw.replace(/^<|>$/g, '').split(/\s+"/)[0];
    if (/^(?:[a-z][a-z0-9+.-]*:|\/\/)/i.test(target)) continue;
    linkCount++;
    const [targetPath, fragment] = target.split('#', 2);
    const normalized = path.posix.normalize(path.posix.join(path.posix.dirname(relative), decodeURIComponent(targetPath || path.posix.basename(relative))));
    if (normalized.startsWith('../') || !knownPaths.has(normalized)) {
      errors.push(`${relative}: missing local target ${target}`);
      continue;
    }
    if (fragment && markdown.has(normalized) && !anchors(markdown.get(normalized)).has(decodeURIComponent(fragment))) {
      errors.push(`${relative}: missing anchor ${target}`);
    }
  }
}

if (errors.length) {
  console.error(errors.join('\n'));
  process.exitCode = 1;
} else {
  console.log(`Documentation OK: ${markdown.size} Markdown files; ${linkCount} local links; ${diagrams} Mermaid blocks (structure only).`);
}
