#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RAFT_COMPOSE="docker compose -f ${ROOT_DIR}/deploy/docker-compose.raft.yml"
GARNET_COMPOSE="docker compose -f ${ROOT_DIR}/deploy/docker-compose.garnet.yml"

DURATION="${DURATION:-8s}"
WORKERS="${WORKERS:-128}"
START_RPS="${START_RPS:-200}"
STEP_RPS="${STEP_RPS:-200}"
MAX_RPS="${MAX_RPS:-3000}"
KEYSPACE="${KEYSPACE:-10000}"
SEED_KEYS="${SEED_KEYS:-10000}"
ERROR_THRESHOLD="${ERROR_THRESHOLD:-0.01}"
P95_THRESHOLD_MS="${P95_THRESHOLD_MS:-250}"
RYOW_REQUESTS="${RYOW_REQUESTS:-1000}"
RYOW_CONCURRENCY="${RYOW_CONCURRENCY:-64}"
SHARD_LIST="${SHARD_LIST:-3 6 9}"
ROUTES="${ROUTES:-lb}"

vCPUS="${vCPUS:-${APP_CPUS:-1.0}}"
APP_MEM="${APP_MEM:-256m}"
LOADGEN_CPUS="${LOADGEN_CPUS:-1.5}"
LOADGEN_MEM="${LOADGEN_MEM:-1g}"
LB_CPUS="${LB_CPUS:-1.0}"
LB_MEM="${LB_MEM:-256m}"
GARNET_CPUS="${GARNET_CPUS:-1.5}"
GARNET_MEM="${GARNET_MEM:-1g}"

OUT_DIR="${ROOT_DIR}/results/matrix"
mkdir -p "${OUT_DIR}"
rm -f "${OUT_DIR}"/*.json

RAFT_NODE_URLS="http://raft1:8080,http://raft2:8080,http://raft3:8080,http://raft4:8080,http://raft5:8080,http://raft6:8080,http://raft7:8080,http://raft8:8080,http://raft9:8080"
CURRENT_SHARDS=3

raft_compose() {
  SHARD_COUNT="${CURRENT_SHARDS}" RAFT_REPLICATION_FACTOR=3 vCPUS="${vCPUS}" APP_MEM="${APP_MEM}" LOADGEN_CPUS="${LOADGEN_CPUS}" LOADGEN_MEM="${LOADGEN_MEM}" LB_CPUS="${LB_CPUS}" LB_MEM="${LB_MEM}" \
    ${RAFT_COMPOSE} "$@"
}

wait_for_strong_smoke() {
  local base_url="$1"
  local key="startup-smoke"
  local value="ok-$(date +%s%N)"
  local attempts=0
  while (( attempts < 40 )); do
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

run_raft_case() {
  local shards="$1"
  local route="$2"   # lb | shard
  local route_args=()

  if [[ "${route}" == "lb" ]]; then
    route_args=(
      -routing single
      -base-url http://lb
    )
  else
    route_args=(
      -routing shard
      -base-url http://raft1:8080
      -node-urls "${RAFT_NODE_URLS}"
      -route-shard-count "${shards}"
      -route-replication-factor 3
    )
  fi

  raft_compose run --rm --entrypoint bench loadgen \
    "${route_args[@]}" \
    -backend "raft-${route}-s${shards}" \
    -scenario ryw \
    -ryw-requests "${RYOW_REQUESTS}" \
    -ryw-concurrency "${RYOW_CONCURRENCY}" \
    -report "/results/matrix/raft-${route}-s${shards}-ryw.json"

  raft_compose run --rm --entrypoint bench loadgen \
    "${route_args[@]}" \
    -backend "raft-${route}-s${shards}" \
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
    -report "/results/matrix/raft-${route}-s${shards}-get.json"

  raft_compose run --rm --entrypoint bench loadgen \
    "${route_args[@]}" \
    -backend "raft-${route}-s${shards}" \
    -scenario post \
    -duration "${DURATION}" \
    -workers "${WORKERS}" \
    -start-rps "${START_RPS}" \
    -step-rps "${STEP_RPS}" \
    -max-rps "${MAX_RPS}" \
    -keyspace "${KEYSPACE}" \
    -error-threshold "${ERROR_THRESHOLD}" \
    -p95-threshold-ms "${P95_THRESHOLD_MS}" \
    -report "/results/matrix/raft-${route}-s${shards}-post.json"
}

run_garnet_case() {
  ${GARNET_COMPOSE} down -v || true
  vCPUS="${vCPUS}" APP_MEM="${APP_MEM}" LOADGEN_CPUS="${LOADGEN_CPUS}" LOADGEN_MEM="${LOADGEN_MEM}" LB_CPUS="${LB_CPUS}" LB_MEM="${LB_MEM}" GARNET_CPUS="${GARNET_CPUS}" GARNET_MEM="${GARNET_MEM}" \
    ${GARNET_COMPOSE} up -d --build

  "${ROOT_DIR}/scripts/wait_for_http.sh" "http://localhost:28080/healthz" 300
  wait_for_strong_smoke "http://localhost:28080"

  ${GARNET_COMPOSE} run --rm --entrypoint bench loadgen \
    -routing single \
    -base-url http://lb \
    -backend garnet-lb \
    -scenario ryw \
    -ryw-requests "${RYOW_REQUESTS}" \
    -ryw-concurrency "${RYOW_CONCURRENCY}" \
    -report "/results/matrix/garnet-lb-ryw.json"

  ${GARNET_COMPOSE} run --rm --entrypoint bench loadgen \
    -routing single \
    -base-url http://lb \
    -backend garnet-lb \
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
    -report "/results/matrix/garnet-lb-get.json"

  ${GARNET_COMPOSE} run --rm --entrypoint bench loadgen \
    -routing single \
    -base-url http://lb \
    -backend garnet-lb \
    -scenario post \
    -duration "${DURATION}" \
    -workers "${WORKERS}" \
    -start-rps "${START_RPS}" \
    -step-rps "${STEP_RPS}" \
    -max-rps "${MAX_RPS}" \
    -keyspace "${KEYSPACE}" \
    -error-threshold "${ERROR_THRESHOLD}" \
    -p95-threshold-ms "${P95_THRESHOLD_MS}" \
    -report "/results/matrix/garnet-lb-post.json"

  ${GARNET_COMPOSE} down -v
}

generate_report() {
  local report_path="${OUT_DIR}/matrix-report.md"
  {
    echo "# RAFT Matrix (RF=3) vs Garnet"
    echo
    echo "Generated: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo
    echo "Config: DURATION=${DURATION} WORKERS=${WORKERS} START_RPS=${START_RPS} STEP_RPS=${STEP_RPS} MAX_RPS=${MAX_RPS}"
    echo
    echo "## Sustainable Throughput (best point under thresholds)"
    echo
    echo "| System | Route | Shards | Scenario | Target RPS | Achieved RPS | p95 (ms) | p99 (ms) | Error Rate |"
    echo "|---|---|---:|---|---:|---:|---:|---:|---:|"

    for f in "${OUT_DIR}"/*.json; do
      name="$(basename "$f")"
      system="raft"
      route="-"
      shards="-"
      if [[ "$name" == garnet-* ]]; then
        system="garnet"
        route="lb"
      else
        route="$(echo "$name" | cut -d- -f2)"
        shards="$(echo "$name" | sed -E 's/^raft-[^-]+-s([0-9]+)-.*$/\1/')"
      fi
      scenario="$(echo "$name" | sed -E 's/.*-(get|post|ryw)\.json$/\1/')"
      if [[ "$scenario" == "ryw" ]]; then
        continue
      fi
      best_target="$(jq -r '.best.target_rps // "n/a"' "$f")"
      best_ach="$(jq -r '.best.achieved_rps // "n/a"' "$f")"
      best_p95="$(jq -r '.best.p95_ms // "n/a"' "$f")"
      best_p99="$(jq -r '.best.p99_ms // "n/a"' "$f")"
      best_err="$(jq -r '.best.error_rate // "n/a"' "$f")"
      printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s |\n" \
        "$system" "$route" "$shards" "$scenario" "$best_target" "$best_ach" "$best_p95" "$best_p99" "$best_err"
    done | sort

    echo
    echo "## Peak Achieved Throughput (regardless of latency)"
    echo
    echo "| System | Route | Shards | Scenario | Peak Achieved RPS | Target RPS At Peak | p95 At Peak (ms) | Dropped | Failures |"
    echo "|---|---|---:|---|---:|---:|---:|---:|---:|"

    for f in "${OUT_DIR}"/*.json; do
      name="$(basename "$f")"
      system="raft"
      route="-"
      shards="-"
      if [[ "$name" == garnet-* ]]; then
        system="garnet"
        route="lb"
      else
        route="$(echo "$name" | cut -d- -f2)"
        shards="$(echo "$name" | sed -E 's/^raft-[^-]+-s([0-9]+)-.*$/\1/')"
      fi
      scenario="$(echo "$name" | sed -E 's/.*-(get|post|ryw)\.json$/\1/')"
      if [[ "$scenario" == "ryw" ]]; then
        continue
      fi
      peak_json="$(jq -c '.steps | max_by(.achieved_rps)' "$f")"
      peak_ach="$(echo "$peak_json" | jq -r '.achieved_rps')"
      peak_target="$(echo "$peak_json" | jq -r '.target_rps')"
      peak_p95="$(echo "$peak_json" | jq -r '.p95_ms')"
      peak_dropped="$(echo "$peak_json" | jq -r '.dropped')"
      peak_failures="$(echo "$peak_json" | jq -r '.failures')"
      printf "| %s | %s | %s | %s | %s | %s | %s | %s | %s |\n" \
        "$system" "$route" "$shards" "$scenario" "$peak_ach" "$peak_target" "$peak_p95" "$peak_dropped" "$peak_failures"
    done | sort

    echo
    echo "## Read-Your-Own-Write"
    echo
    echo "| System | Route | Shards | Requests | Success | Failures | Violations |"
    echo "|---|---|---:|---:|---:|---:|---:|"

    for f in "${OUT_DIR}"/*-ryw.json; do
      name="$(basename "$f")"
      system="raft"
      route="-"
      shards="-"
      if [[ "$name" == garnet-* ]]; then
        system="garnet"
        route="lb"
      else
        route="$(echo "$name" | cut -d- -f2)"
        shards="$(echo "$name" | sed -E 's/^raft-[^-]+-s([0-9]+)-.*$/\1/')"
      fi
      requests="$(jq -r '.ryw.requests // 0' "$f")"
      success="$(jq -r '.ryw.success // 0' "$f")"
      failures="$(jq -r '.ryw.failures // 0' "$f")"
      violations="$(jq -r '.ryw.violations // 0' "$f")"
      printf "| %s | %s | %s | %s | %s | %s | %s |\n" \
        "$system" "$route" "$shards" "$requests" "$success" "$failures" "$violations"
    done | sort
  } >"${report_path}"

  echo "matrix report: ${report_path}"
}

raft_compose down -v || true
${GARNET_COMPOSE} down -v || true

for shards in ${SHARD_LIST}; do
  CURRENT_SHARDS="${shards}"
  raft_compose up -d --build

  "${ROOT_DIR}/scripts/wait_for_http.sh" "http://localhost:18080/healthz" 300
  wait_for_strong_smoke "http://localhost:18080"

  for route in ${ROUTES}; do
    run_raft_case "${shards}" "${route}"
  done

  raft_compose down -v

done

run_garnet_case

generate_report

echo "done: ${OUT_DIR}"
