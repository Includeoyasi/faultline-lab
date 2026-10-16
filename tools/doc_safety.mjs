// Deliberately narrow hygiene check, not a general secret or filesystem scanner.
export function containsMachinePath(text) {
  const withoutWebURLs = text.replace(/\bhttps?:\/\/[^\s"'<>`]+/gi, '');
  return /\b[A-Za-z]:[\\/]+Users[\\/]+/i.test(withoutWebURLs)
    || /(?:^|[^A-Za-z0-9_.-])\/(?:Users|home)\/[^/\s"'<>()[\]{}]+/m.test(withoutWebURLs);
}
