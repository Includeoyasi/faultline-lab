import assert from 'node:assert/strict';
import test from 'node:test';
import { containsMachinePath } from './doc_safety.mjs';

test('detects synthetic Windows, macOS and Linux homes', () => {
  const windows = ['C:', 'Users', 'example', 'private.txt'].join('\\');
  const mac = ['', 'Users', 'example', 'private.txt'].join('/');
  const linux = ['', 'home', 'example', 'private.txt'].join('/');
  for (const value of [windows, JSON.stringify({ path: windows }), windows.replaceAll('\\', '/'), mac, linux, `path="${mac}"`, `file://${linux}`, linux.slice(0, linux.lastIndexOf('/'))]) {
    assert.equal(containsMachinePath(value), true, value);
  }
});

test('allows container paths, relative paths, regular URLs and ordinary prose', () => {
  for (const value of ['/app/service', '/workspace', '/root/.cache/go-build', 'docs/Users/guide.md', 'src/home/example', 'https://example.test/home/example/report', 'https://example.test/Users/example', 'Users and home directories']) {
    assert.equal(containsMachinePath(value), false, value);
  }
});
