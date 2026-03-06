# RAFT Matrix (RF=3) vs Garnet

Generated: 2026-03-03T05:45:09Z

Config: DURATION=3s WORKERS=256 START_RPS=1000 STEP_RPS=500 MAX_RPS=8000

## Sustainable Throughput (best point under thresholds)

| System | Route | Shards | Scenario | Target RPS | Achieved RPS | p95 (ms) | p99 (ms) | Error Rate |
|---|---|---:|---|---:|---:|---:|---:|---:|
| garnet | lb | - | get | 6000 | 4323.4834102454615 | 54.191718 | 605.366275 | 0 |
| garnet | lb | - | post | 6500 | 6485.331702215522 | 51.345066 | 63.553123 | 0 |
| raft | lb | 3 | get | 6000 | 5986.634540022294 | 59.850976 | 82.993087 | 0 |
| raft | lb | 3 | post | 6500 | 5256.053165711138 | 83.252301 | 286.952638 | 0 |

## Peak Achieved Throughput (regardless of latency)

| System | Route | Shards | Scenario | Peak Achieved RPS | Target RPS At Peak | p95 At Peak (ms) | Dropped | Failures |
|---|---|---:|---|---:|---:|---:|---:|---:|
| garnet | lb | - | get | 5489.258505988326 | 5500 | 1.404639 | 0 | 0 |
| garnet | lb | - | post | 6485.331702215522 | 6500 | 51.345066 | 0 | 0 |
| raft | lb | 3 | get | 5986.634540022294 | 6000 | 59.850976 | 0 | 0 |
| raft | lb | 3 | post | 5764.557755757914 | 6000 | 86.790989 | 233 | 0 |

## Read-Your-Own-Write

| System | Route | Shards | Requests | Success | Failures | Violations |
|---|---|---:|---:|---:|---:|---:|
| garnet | lb | - | 1000 | 1000 | 0 | 0 |
| raft | lb | 3 | 1000 | 1000 | 0 | 0 |
