#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
build linux amd64 "${APP_NAME}-linux-amd64"
