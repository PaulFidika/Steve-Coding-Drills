#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RAFT_COMPOSE="docker compose -f ${ROOT_DIR}/deploy/docker-compose.raft.yml"
GARNET_COMPOSE="docker compose -f ${ROOT_DIR}/deploy/docker-compose.garnet.yml"

DURATION="${DURATION:-20s}"
WORKERS="${WORKERS:-128}"
START_RPS="${START_RPS:-200}"
STEP_RPS="${STEP_RPS:-200}"
MAX_RPS="${MAX_RPS:-5000}"
KEYSPACE="${KEYSPACE:-10000}"
SEED_KEYS="${SEED_KEYS:-10000}"
ERROR_THRESHOLD="${ERROR_THRESHOLD:-0.01}"
P95_THRESHOLD_MS="${P95_THRESHOLD_MS:-250}"
RYOW_REQUESTS="${RYOW_REQUESTS:-2000}"
RYOW_CONCURRENCY="${RYOW_CONCURRENCY:-64}"

mkdir -p "${ROOT_DIR}/results"

wait_for_strong_smoke() {
  local base_url="$1"
  local key="startup-smoke"
  local value="ok-$(date +%s%N)"
  local attempts=0
  while (( attempts < 30 )); do
    attempts=$((attempts + 1))
    curl -fsS -X POST "${base_url}/state/${key}" \
      -H "content-type: application/json" \
      -d "{\"value\":\"${value}\"}" >/dev/null 2>&1 || true
    body="$(curl -fsS "${base_url}/state/${key}" 2>/dev/null || true)"
    if [[ "${body}" == *"\"value\":\"${value}\""* ]]; then
      echo "startup smoke passed for ${base_url}"
      return 0
    fi
    sleep 1
  done
  echo "startup smoke failed for ${base_url}" >&2
  return 1
}

run_for_backend() {
  local compose_cmd="$1"
  local backend="$2"
  local host_port="$3"

  ${compose_cmd} down -v || true
  ${compose_cmd} up -d --build

  "${ROOT_DIR}/scripts/wait_for_http.sh" "http://localhost:${host_port}/healthz" 300
  wait_for_strong_smoke "http://localhost:${host_port}"

  ${compose_cmd} run --rm --entrypoint bench loadgen \
    -base-url "http://lb" \
    -backend "${backend}" \
    -scenario ryw \
    -ryw-requests "${RYOW_REQUESTS}" \
    -ryw-concurrency "${RYOW_CONCURRENCY}" \
    -report "/results/${backend}-ryw.json"

  ${compose_cmd} run --rm --entrypoint bench loadgen \
    -base-url "http://lb" \
    -backend "${backend}" \
    -scenario get \
    -duration "${DURATION}" \
    -workers "${WORKERS}" \
    -start-rps "${START_RPS}" \
    -step-rps "${STEP_RPS}" \
    -max-rps "${MAX_RPS}" \
    -keyspace "${KEYSPACE}" \
    -seed-keys "${SEED_KEYS}" \
    -error-threshold "${ERROR_THRESHOLD}" \
    -p95-threshold-ms "${P95_THRESHOLD_MS}" \
    -report "/results/${backend}-get.json"

  ${compose_cmd} run --rm --entrypoint bench loadgen \
    -base-url "http://lb" \
    -backend "${backend}" \
    -scenario post \
    -duration "${DURATION}" \
    -workers "${WORKERS}" \
    -start-rps "${START_RPS}" \
    -step-rps "${STEP_RPS}" \
    -max-rps "${MAX_RPS}" \
    -keyspace "${KEYSPACE}" \
    -error-threshold "${ERROR_THRESHOLD}" \
    -p95-threshold-ms "${P95_THRESHOLD_MS}" \
    -report "/results/${backend}-post.json"

  ${compose_cmd} down -v
}

run_for_backend "${RAFT_COMPOSE}" raft 18080
run_for_backend "${GARNET_COMPOSE}" garnet 28080

cd "${ROOT_DIR}"
go run ./cmd/compare \
  -raft-get results/raft-get.json \
  -raft-post results/raft-post.json \
  -raft-ryw results/raft-ryw.json \
  -garnet-get results/garnet-get.json \
  -garnet-post results/garnet-post.json \
  -garnet-ryw results/garnet-ryw.json \
  -out results/comparison.md

echo "done: results/comparison.md"
