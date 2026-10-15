#!/bin/sh
# This script runs inside the Linux toolchain container, not on the host.
set -eu

if [ "$#" -gt 1 ]; then
  echo "usage: lab-go {versions|build|check|format|bootstrap-smoke}" >&2
  exit 2
fi

versions() {
  go version
  golangci-lint --version
  go version -m "$(command -v goimports)"
  govulncheck -version
}

build() {
  mkdir -p bin
  for service in orders inventory payments; do
    CGO_ENABLED=0 go build -trimpath -o "bin/$service" "./cmd/$service"
    go version -m "bin/$service"
  done
}

source_files() {
  file_list=$(mktemp)
  trap 'rm -f "$file_list"' EXIT
  find . -type d \( -name .git -o -name .local -o -name vendor -o -name bin -o -name tmp -o -name artifacts \) -prune \
    -o -type f -name '*.go' -print0 > "$file_list"
}

case "${1:-versions}" in
  versions)
    versions
    ;;
  build)
    build
    ;;
  check)
    versions
    source_files
    formatting=$(xargs -0 -r gofmt -l < "$file_list")
    imports=$(xargs -0 -r goimports -l < "$file_list")
    if [ -n "$formatting$imports" ]; then
      printf 'Formatting/import changes required:\n%s\n%s\n' "$formatting" "$imports" >&2
      exit 1
    fi
    go mod tidy -diff
    go mod verify
    go build ./...
    golangci-lint run --timeout=5m
    go test -count=1 ./...
    go test -race -count=1 ./...
    govulncheck ./...
    ;;
  format)
    source_files
    xargs -0 -r goimports -w < "$file_list"
    ;;
  bootstrap-smoke)
    build
    for service in orders inventory payments; do
      "bin/$service"
    done
    ;;
  *)
    echo "usage: lab-go {versions|build|check|format|bootstrap-smoke}" >&2
    exit 2
    ;;
esac
