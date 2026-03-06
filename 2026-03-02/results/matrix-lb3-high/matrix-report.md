# RAFT Matrix (RF=3) vs Garnet

Generated: 2026-03-03T05:53:05Z

Config: DURATION=3s WORKERS=256 START_RPS=1000 STEP_RPS=500 MAX_RPS=8000

## Sustainable Throughput (best point under thresholds)

| System | Route | Shards | Scenario | Target RPS | Achieved RPS | p95 (ms) | p99 (ms) | Error Rate |
|---|---|---:|---|---:|---:|---:|---:|---:|
| garnet | lb | - | get | 6000 | 4910.995048607256 | 17.944653 | 422.430844 | 0 |
| garnet | lb | - | post | 6500 | 5498.370169148074 | 58.485956 | 472.550913 | 0 |
| raft | lb | 3 | get | 6000 | 5987.941785888778 | 69.618881 | 98.311093 | 0 |
| raft | lb | 3 | post | 6500 | 1887.1439895663898 | 181.343684 | 208.304118 | 0 |

## Peak Achieved Throughput (regardless of latency)

| System | Route | Shards | Scenario | Peak Achieved RPS | Target RPS At Peak | p95 At Peak (ms) | Dropped | Failures |
|---|---|---:|---|---:|---:|---:|---:|---:|
| garnet | lb | - | get | 5488.861415813564 | 5500 | 1.465787 | 0 | 0 |
| garnet | lb | - | post | 5988.309992232828 | 6000 | 2.146489 | 0 | 0 |
| raft | lb | 3 | get | 5987.941785888778 | 6000 | 69.618881 | 0 | 0 |
| raft | lb | 3 | post | 2629.07389057781 | 3500 | 146.672716 | 1069 | 0 |

## Read-Your-Own-Write

| System | Route | Shards | Requests | Success | Failures | Violations |
|---|---|---:|---:|---:|---:|---:|
| garnet | lb | - | 1000 | 1000 | 0 | 0 |
| raft | lb | 3 | 1000 | 1000 | 0 | 0 |
