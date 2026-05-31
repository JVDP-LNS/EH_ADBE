#!/usr/bin/env bash
set -euo pipefail
DIR="$(dirname "$0")"
for script in darwin-amd64 darwin-arm64 windows-amd64 linux-amd64; do
  "${DIR}/${script}.sh"
done
