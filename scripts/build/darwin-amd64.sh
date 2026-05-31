#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
build darwin amd64 "${APP_NAME}-darwin-amd64"
