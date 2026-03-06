# RAFT (Dragonboat) vs Microsoft Garnet Benchmark Harness

This project provides a 9-node webserver cluster in two modes:

- `raft`: each webserver embeds Dragonboat RAFT with shard placement and replication factor 3 (`1` leader + `2` followers per shard).
- `garnet`: 9 stateless webservers use a single Microsoft Garnet instance as the source of truth.

Both modes expose the same API:

- `POST /state/{key}` with JSON body `{"value":"..."}`
- `GET /state/{key}`
- `GET /healthz`

For RAFT mode, followers forward reads/writes to the current leader for the key shard to maintain strong consistency semantics. For Garnet mode, all reads/writes go directly to Garnet.

## Benchmark Snapshot (2026-03-03)

This repository includes a completed benchmark run comparing:

- 9 RAFT-backed webservers (`dragonboat`)
- 9 Garnet-backed webservers (`microsoft/garnet` as Redis-compatible source of truth)

### Test Setup

- API tested: `GET /state/{key}`, `POST /state/{key}`, and read-your-own-write (`ryw`)
- Load path: `loadgen` service -> nginx load balancer -> random webserver
- Consistency target: strong read-your-own-write across random node selection
- Selection rule for "best sustainable throughput":
  - `error_rate <= 0.01`
  - `p95 <= 250ms`

Note:

- This snapshot was produced before switching RAFT from full-mesh shard replication to the current `replication_factor=3` shard placement model. Re-run benchmarks after topology changes.

Resource constraints used:

- Webserver (`vCPUS=1.0`, `APP_MEM=256m`)
- Load generator (`LOADGEN_CPUS=1.0`, `LOADGEN_MEM=512m`)
- Load balancer (`LB_CPUS=0.25`, `LB_MEM=128m`)
- Garnet (`GARNET_CPUS=1.0`, `GARNET_MEM=512m`)

Benchmark sweep settings:

- `DURATION=10s`
- `WORKERS=128`
- `START_RPS=200`
- `STEP_RPS=200`
- `MAX_RPS=3000`
- `KEYSPACE=10000`
- `SEED_KEYS=10000`
- `RYOW_REQUESTS=2000`
- `RYOW_CONCURRENCY=64`

### Results

Generated at: `2026-03-03T01:53:58Z` (`results/comparison.md`)

| Scenario | Backend | Target RPS | Achieved RPS | p50 (ms) | p95 (ms) | p99 (ms) | Error Rate |
|---|---:|---:|---:|---:|---:|---:|---:|
| GET | raft | 400 | 399.7 | 2.71 | 53.21 | 74.76 | 0.0000 |
| GET | garnet | 400 | 399.8 | 0.93 | 1.41 | 1.91 | 0.0000 |
| POST | raft | 400 | 399.3 | 16.45 | 33.85 | 76.54 | 0.0000 |
| POST | garnet | 400 | 399.8 | 1.04 | 1.52 | 2.06 | 0.0000 |

Read-your-own-write consistency:

| Backend | Requests | Success | Failures | Violations |
|---|---:|---:|---:|---:|
| raft | 2000 | 2000 | 0 | 0 |
| garnet | 2000 | 2000 | 0 | 0 |

### Interpretation

- Both modes satisfied read-your-own-write in this run (`0` violations each).
- At equivalent sustainable throughput in this environment, Garnet showed substantially lower latency for both reads and writes.
- RAFT write latency is expectedly higher due to consensus replication before acknowledgment.

## Compose Topology

- `deploy/docker-compose.raft.yml`
  - `raft1..raft9`
  - `lb` (nginx round-robin)
  - `loadgen` (benchmark runner service)
- `deploy/docker-compose.garnet.yml`
  - `garnet`
  - `gweb1..gweb9`
  - `lb` (nginx round-robin)
  - `loadgen`

## Resource Constraints

Use environment variables when running compose:

- `vCPUS`, `APP_MEM`
- `LOADGEN_CPUS`, `LOADGEN_MEM`
- `LB_CPUS`, `LB_MEM`
- `GARNET_CPUS`, `GARNET_MEM` (garnet stack)
- `SHARD_COUNT` (RAFT only, supported: `3..9`)
- `RAFT_REPLICATION_FACTOR` (RAFT only, default: `3`)

Example:

```bash
vCPUS=1.0 APP_MEM=256m GARNET_CPUS=1.0 GARNET_MEM=512m ./scripts/run_compare.sh
```

## Benchmark Workloads

The benchmark CLI (`cmd/bench`) supports:

- `scenario=get`: latency and throughput sweep for reads
- `scenario=post`: latency and throughput sweep for writes
- `scenario=ryw`: read-your-own-write validation under concurrent clients

Throughput sweep returns the best sustainable point under thresholds:

- `error_rate <= error-threshold`
- `p95 <= p95-threshold-ms`

## End-to-End Run

```bash
./scripts/run_compare.sh
```

Outputs:

- `results/raft-ryw.json`
- `results/raft-get.json`
- `results/raft-post.json`
- `results/garnet-ryw.json`
- `results/garnet-get.json`
- `results/garnet-post.json`
- `results/comparison.md`

## Manual Runs

Start RAFT stack:

```bash
docker compose -f deploy/docker-compose.raft.yml up -d --build
curl http://localhost:18080/healthz
```

Run load test inside compose network:

```bash
docker compose -f deploy/docker-compose.raft.yml run --rm --entrypoint bench loadgen \
  -base-url http://lb -backend raft -scenario get -report /results/raft-get.json
```

Stop stack:

```bash
docker compose -f deploy/docker-compose.raft.yml down -v
```

Start Garnet stack:

```bash
docker compose -f deploy/docker-compose.garnet.yml up -d --build
curl http://localhost:28080/healthz
```

## Notes

- RAFT uses static peer membership from `RAFT_PEERS`, modeled after the Dragonboat setup pattern in `~/cozy/gen-orchestrator`.
- RAFT shard placement uses deterministic node groups per shard.
- RAFT defaults to `SHARD_COUNT=3`, `RAFT_REPLICATION_FACTOR=3` (recommended shard range: `3..9` on this 9-node setup).
