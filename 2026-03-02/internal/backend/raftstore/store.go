package raftstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lni/dragonboat/v4"
	"github.com/lni/dragonboat/v4/client"
	"github.com/lni/dragonboat/v4/config"
	sm "github.com/lni/dragonboat/v4/statemachine"
	"github.com/lni/goutils/random"

	"raft-garnet-bench/internal/backend"
)

type Config struct {
	ReplicaID         uint64
	ShardCount        uint64
	ReplicationFactor int
	RaftAddress       string
	DataDir           string
	RaftPeers         map[uint64]string
	HTTPPeers         map[uint64]string
}

type Store struct {
	cfg           Config
	nh            *dragonboat.NodeHost
	sessions      map[uint64]*client.Session
	sessionRand   random.Source
	peerIDs       []uint64
	shardReplicas map[uint64][]uint64
	ownedShards   map[uint64]struct{}
	mu            sync.Mutex
}

func New(cfg Config) (*Store, error) {
	if cfg.ReplicaID == 0 {
		return nil, errors.New("replica id is required")
	}
	if cfg.ShardCount == 0 {
		cfg.ShardCount = 3
	}
	if cfg.ReplicationFactor == 0 {
		cfg.ReplicationFactor = 3
	}
	if strings.TrimSpace(cfg.RaftAddress) == "" {
		return nil, errors.New("raft address is required")
	}
	if len(cfg.RaftPeers) == 0 {
		return nil, errors.New("raft peers are required")
	}
	if _, ok := cfg.RaftPeers[cfg.ReplicaID]; !ok {
		return nil, fmt.Errorf("replica id %d missing from RAFT_PEERS", cfg.ReplicaID)
	}
	if cfg.ShardCount < 3 || cfg.ShardCount > uint64(len(cfg.RaftPeers)) {
		return nil, fmt.Errorf("SHARD_COUNT must be between 3 and %d", len(cfg.RaftPeers))
	}
	if cfg.ReplicationFactor < 3 || cfg.ReplicationFactor > len(cfg.RaftPeers) {
		return nil, fmt.Errorf("RAFT_REPLICATION_FACTOR must be between 3 and %d", len(cfg.RaftPeers))
	}
	if cfg.DataDir == "" {
		cfg.DataDir = fmt.Sprintf(".runtime/raft-%d", cfg.ReplicaID)
	}
	if err := os.MkdirAll(cfg.DataDir, 0o700); err != nil {
		return nil, fmt.Errorf("create raft dir: %w", err)
	}

	nhCfg := config.NodeHostConfig{
		NodeHostDir:    cfg.DataDir,
		WALDir:         cfg.DataDir,
		RaftAddress:    cfg.RaftAddress,
		RTTMillisecond: 200,
	}
	nh, err := dragonboat.NewNodeHost(nhCfg)
	if err != nil {
		return nil, err
	}

	s := &Store{
		cfg:           cfg,
		nh:            nh,
		sessions:      make(map[uint64]*client.Session, cfg.ShardCount),
		sessionRand:   random.NewLockedRand(),
		shardReplicas: make(map[uint64][]uint64, cfg.ShardCount),
		ownedShards:   make(map[uint64]struct{}),
	}
	s.peerIDs = s.sortedPeerIDs()
	s.buildShardPlan()

	if err := s.startReplicas(); err != nil {
		nh.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Mode() string {
	return "raft"
}

func (s *Store) startReplicas() error {
	owned := s.ownedShardIDs()
	if len(owned) == 0 {
		return errors.New("this replica owns no shards")
	}
	for _, shardID := range owned {
		rc := config.Config{
			ShardID:            shardID,
			ReplicaID:          s.cfg.ReplicaID,
			HeartbeatRTT:       1,
			ElectionRTT:        10,
			CheckQuorum:        true,
			SnapshotEntries:    10000,
			CompactionOverhead: 1000,
		}
		members := s.membersForShard(shardID)
		join := s.cfg.ReplicaID != s.lowestReplicaIDForShard(shardID)
		if join {
			members = map[uint64]string{}
		}
		if err := s.nh.StartReplica(members, join, s.createStateMachine(), rc); err != nil {
			return fmt.Errorf("start replica shard %d: %w", shardID, err)
		}
		s.sessions[shardID] = client.NewNoOPSession(shardID, s.sessionRand)
	}
	return nil
}

func (s *Store) createStateMachine() sm.CreateStateMachineFunc {
	return func(uint64, uint64) sm.IStateMachine {
		return newStateMachine()
	}
}

func (s *Store) membersForShard(shardID uint64) map[uint64]string {
	replicas := s.replicasForShard(shardID)
	members := make(map[uint64]string, len(replicas))
	for _, id := range replicas {
		addr, ok := s.cfg.RaftPeers[id]
		if !ok {
			continue
		}
		if strings.TrimSpace(addr) != "" {
			members[id] = addr
		}
	}
	return members
}

func (s *Store) lowestReplicaIDForShard(shardID uint64) uint64 {
	replicas := s.replicasForShard(shardID)
	if len(replicas) == 0 {
		return s.cfg.ReplicaID
	}
	lowest := replicas[0]
	for _, id := range replicas[1:] {
		if id < lowest {
			lowest = id
		}
	}
	return lowest
}

func (s *Store) Set(ctx context.Context, key, value string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("key is required")
	}
	shardID := computeShardID(key, s.cfg.ShardCount)
	if !s.ownsShard(shardID) {
		return s.routeToReplica(s.preferredReplicaForShard(shardID))
	}
	leaderID, leaderURL, err := s.waitForLeader(ctx, shardID)
	if err != nil {
		return err
	}
	if leaderID != s.cfg.ReplicaID {
		return &backend.NotLeaderError{LeaderID: leaderID, LeaderURL: leaderURL}
	}

	payload, _ := json.Marshal(setCommand{Key: key, Value: value})
	cmd, _ := json.Marshal(command{Type: "set", Payload: payload})
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	_, err = s.nh.SyncPropose(ctx, s.sessionForShard(shardID), cmd)
	return err
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", errors.New("key is required")
	}
	shardID := computeShardID(key, s.cfg.ShardCount)
	if !s.ownsShard(shardID) {
		return "", s.routeToReplica(s.preferredReplicaForShard(shardID))
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	res, err := s.nh.SyncRead(ctx, shardID, query{Type: "get", Key: key})
	if err != nil {
		return "", err
	}
	out, ok := res.(getResult)
	if !ok {
		return "", fmt.Errorf("unexpected query result type %T", res)
	}
	if !out.Found {
		return "", backend.ErrNotFound
	}
	return out.Value, nil
}

func (s *Store) waitForLeader(ctx context.Context, shardID uint64) (uint64, string, error) {
	if !s.ownsShard(shardID) {
		return 0, "", fmt.Errorf("shard %d is not hosted on replica %d", shardID, s.cfg.ReplicaID)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	deadline := time.Now().Add(5 * time.Second)
	if dl, ok := ctx.Deadline(); ok {
		deadline = dl
	}
	for {
		leaderID, _, ok, err := s.nh.GetLeaderID(shardID)
		if err == nil && ok && leaderID != 0 {
			return leaderID, s.cfg.HTTPPeers[leaderID], nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return 0, "", err
			}
			return 0, "", fmt.Errorf("leader unavailable for shard %d", shardID)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *Store) sessionForShard(shardID uint64) *client.Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	session := s.sessions[shardID]
	if session != nil {
		return session
	}
	session = client.NewNoOPSession(shardID, s.sessionRand)
	s.sessions[shardID] = session
	return session
}

func (s *Store) sortedPeerIDs() []uint64 {
	ids := make([]uint64, 0, len(s.cfg.RaftPeers))
	for id := range s.cfg.RaftPeers {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func (s *Store) buildShardPlan() {
	n := len(s.peerIDs)
	shards := int(s.cfg.ShardCount)
	rf := s.cfg.ReplicationFactor
	for shard := 1; shard <= shards; shard++ {
		leaderPos := ((shard - 1) * n) / shards
		members := make([]uint64, 0, rf)
		seen := make(map[uint64]struct{}, rf)
		for offset := 0; offset < n && len(members) < rf; offset++ {
			id := s.peerIDs[(leaderPos+offset)%n]
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			members = append(members, id)
		}
		sid := uint64(shard)
		s.shardReplicas[sid] = members
		for _, rid := range members {
			if rid == s.cfg.ReplicaID {
				s.ownedShards[sid] = struct{}{}
			}
		}
	}
}

func (s *Store) ownedShardIDs() []uint64 {
	out := make([]uint64, 0, len(s.ownedShards))
	for shardID := range s.ownedShards {
		out = append(out, shardID)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (s *Store) replicasForShard(shardID uint64) []uint64 {
	replicas := s.shardReplicas[shardID]
	out := make([]uint64, len(replicas))
	copy(out, replicas)
	return out
}

func (s *Store) ownsShard(shardID uint64) bool {
	_, ok := s.ownedShards[shardID]
	return ok
}

func (s *Store) preferredReplicaForShard(shardID uint64) uint64 {
	replicas := s.replicasForShard(shardID)
	if len(replicas) == 0 {
		return 0
	}
	return replicas[0]
}

func (s *Store) routeToReplica(replicaID uint64) error {
	if replicaID == 0 {
		return errors.New("shard route unavailable")
	}
	leaderURL := strings.TrimSpace(s.cfg.HTTPPeers[replicaID])
	if leaderURL == "" {
		return fmt.Errorf("http peer missing for replica %d", replicaID)
	}
	return &backend.NotLeaderError{LeaderID: replicaID, LeaderURL: leaderURL}
}

func (s *Store) Close() error {
	if s.nh != nil {
		s.nh.Close()
	}
	return nil
}

func computeShardID(key string, shardCount uint64) uint64 {
	if shardCount == 0 {
		return 1
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return uint64(h.Sum32()%uint32(shardCount)) + 1
}
