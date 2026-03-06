package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Backend               string
	ListenAddr            string
	NodeID                uint64
	ShardCount            uint64
	RaftReplicationFactor int
	RaftAddress           string
	RaftDataDir           string
	RaftPeers             map[uint64]string
	HTTPPeers             map[uint64]string
	GarnetAddr            string
	GarnetPassword        string
	GarnetDB              int
}

func Load() (Config, error) {
	cfg := Config{
		Backend:               envOrDefault("BACKEND", "raft"),
		ListenAddr:            envOrDefault("HTTP_ADDR", ":8080"),
		NodeID:                mustUint64(envOrDefault("NODE_ID", "1")),
		ShardCount:            mustUint64(envOrDefault("SHARD_COUNT", "3")),
		RaftReplicationFactor: mustInt(envOrDefault("RAFT_REPLICATION_FACTOR", "3")),
		RaftAddress:           os.Getenv("RAFT_ADDR"),
		RaftDataDir:           envOrDefault("RAFT_DATA_DIR", ".runtime/raft"),
		GarnetAddr:            envOrDefault("GARNET_ADDR", "localhost:6379"),
		GarnetPassword:        os.Getenv("GARNET_PASSWORD"),
		GarnetDB:              mustInt(envOrDefault("GARNET_DB", "0")),
	}

	if cfg.ShardCount == 0 {
		cfg.ShardCount = 1
	}
	if cfg.NodeID == 0 {
		return Config{}, fmt.Errorf("NODE_ID must be >= 1")
	}

	var err error
	cfg.RaftPeers, err = parsePeerMap(os.Getenv("RAFT_PEERS"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid RAFT_PEERS: %w", err)
	}
	cfg.HTTPPeers, err = parsePeerMap(os.Getenv("HTTP_PEERS"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_PEERS: %w", err)
	}

	if cfg.Backend == "raft" {
		if strings.TrimSpace(cfg.RaftAddress) == "" {
			return Config{}, fmt.Errorf("RAFT_ADDR is required for raft backend")
		}
		if len(cfg.RaftPeers) == 0 {
			return Config{}, fmt.Errorf("RAFT_PEERS is required for raft backend")
		}
		if len(cfg.HTTPPeers) == 0 {
			return Config{}, fmt.Errorf("HTTP_PEERS is required for raft backend")
		}
		if cfg.ShardCount < 3 || cfg.ShardCount > uint64(len(cfg.RaftPeers)) {
			return Config{}, fmt.Errorf("SHARD_COUNT must be between 3 and %d", len(cfg.RaftPeers))
		}
		if cfg.RaftReplicationFactor < 3 || cfg.RaftReplicationFactor > len(cfg.RaftPeers) {
			return Config{}, fmt.Errorf("RAFT_REPLICATION_FACTOR must be between 3 and %d", len(cfg.RaftPeers))
		}
	}

	if cfg.Backend != "raft" && cfg.Backend != "garnet" {
		return Config{}, fmt.Errorf("BACKEND must be raft or garnet")
	}

	return cfg, nil
}

func parsePeerMap(input string) (map[uint64]string, error) {
	out := map[uint64]string{}
	input = strings.TrimSpace(input)
	if input == "" {
		return out, nil
	}
	pairs := strings.Split(input, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected id=value pair, got %q", pair)
		}
		id, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil || id == 0 {
			return nil, fmt.Errorf("invalid id %q", parts[0])
		}
		value := strings.TrimSpace(parts[1])
		if value == "" {
			return nil, fmt.Errorf("missing address for id %d", id)
		}
		out[id] = value
	}
	return out, nil
}

func envOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func mustUint64(v string) uint64 {
	n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func mustInt(v string) int {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0
	}
	return n
}
