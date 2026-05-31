#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
build windows amd64 "${APP_NAME}-windows-amd64.exe"
