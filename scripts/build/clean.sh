#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"

if [[ -d "${BIN}" ]]; then
  rm -rf "${BIN}"
  echo "removed ${BIN}"
else
  echo "nothing to clean (${BIN} does not exist)"
fi
