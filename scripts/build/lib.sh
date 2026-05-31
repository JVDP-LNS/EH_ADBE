#!/usr/bin/env bash
# Shared build helpers. Source from platform scripts; do not run directly.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND="${ROOT}/backend"
BIN="${ROOT}/bin"
APP_NAME="eh_adbe_backend"

build() {
  local goos="$1"
  local goarch="$2"
  local out_name="$3"

  mkdir -p "${BIN}"
  echo "building ${goos}/${goarch} -> ${BIN}/${out_name}"
  (
    cd "${BACKEND}"
    GOOS="${goos}" GOARCH="${goarch}" CGO_ENABLED=0 go build -ldflags="-s -w" -o "${BIN}/${out_name}" .
  )
}
