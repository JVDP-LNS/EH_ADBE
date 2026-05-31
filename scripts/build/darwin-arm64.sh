#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
build darwin arm64 "${APP_NAME}-darwin-arm64"
