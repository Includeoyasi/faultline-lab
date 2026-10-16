#!/bin/sh
# This script runs inside the Linux toolchain container, not on the host.
set -eu

if [ "$#" -gt 1 ]; then
  echo "usage: lab-go {versions|build|check|format|service-smoke|smoke-build|smoke-cross-check|smoke-cli-tests|data-build|data-cross-check}" >&2
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

host_build() {
  command_name=$1
  target_os=${GOOS:-$(go env GOOS)}
  target_arch=${GOARCH:-$(go env GOARCH)}
  case "$target_os/$target_arch" in
    windows/amd64|linux/amd64|linux/arm64|darwin/amd64|darwin/arm64) ;;
    *) echo "unsupported smoke target" >&2; exit 2 ;;
  esac
  extension=
  if [ "$target_os" = windows ]; then extension=.exe; fi
  mkdir -p bin
  CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" go build -trimpath -o "bin/$command_name-$target_os-$target_arch$extension" "./cmd/$command_name"
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
    go test -count=1 -timeout=3m ./...
    go test -race -count=1 -timeout=5m ./...
    govulncheck ./...
    ;;
  format)
    source_files
    xargs -0 -r goimports -w < "$file_list"
    ;;
  service-smoke)
    build
    LAB_BIN_DIR=/workspace/bin go test -count=1 -v -timeout=2m ./tests/smoke
    ;;
  smoke-build)
    host_build smoke
    ;;
  smoke-cross-check)
    for target in windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
      GOOS=${target%/*} GOARCH=${target#*/} host_build smoke
    done
    ;;
  data-build)
    host_build labdata
    ;;
  data-cross-check)
    for target in windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
      GOOS=${target%/*} GOARCH=${target#*/} host_build labdata
    done
    ;;
  smoke-cli-tests)
    mkdir -p bin
    CGO_ENABLED=0 go test -c -o bin/smoke-cli-tests ./tests/smokecli
    ;;
  *)
    echo "usage: lab-go {versions|build|check|format|service-smoke|smoke-build|smoke-cross-check|smoke-cli-tests|data-build|data-cross-check}" >&2
    exit 2
    ;;
esac
