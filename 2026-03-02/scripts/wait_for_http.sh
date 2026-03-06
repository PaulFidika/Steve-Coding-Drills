#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <url> [timeout_seconds]" >&2
  exit 1
fi

URL="$1"
TIMEOUT="${2:-180}"
START="$(date +%s)"

while true; do
  if curl -fsS "$URL" >/dev/null 2>&1; then
    echo "ready: $URL"
    exit 0
  fi
  NOW="$(date +%s)"
  if (( NOW - START > TIMEOUT )); then
    echo "timeout waiting for $URL" >&2
    exit 1
  fi
  sleep 2
done
