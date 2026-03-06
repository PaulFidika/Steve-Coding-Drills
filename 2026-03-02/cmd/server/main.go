package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"raft-garnet-bench/internal/api"
	"raft-garnet-bench/internal/backend"
	"raft-garnet-bench/internal/backend/garnetstore"
	"raft-garnet-bench/internal/backend/raftstore"
	appcfg "raft-garnet-bench/internal/config"
)

func main() {
	cfg, err := appcfg.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	var store backend.Store
	switch cfg.Backend {
	case "raft":
		store, err = raftstore.New(raftstore.Config{
			ReplicaID:         cfg.NodeID,
			ShardCount:        cfg.ShardCount,
			ReplicationFactor: cfg.RaftReplicationFactor,
			RaftAddress:       cfg.RaftAddress,
			DataDir:           cfg.RaftDataDir,
			RaftPeers:         cfg.RaftPeers,
			HTTPPeers:         cfg.HTTPPeers,
		})
		if err != nil {
			log.Fatalf("raft init error: %v", err)
		}
	case "garnet":
		store = garnetstore.New(cfg.GarnetAddr, cfg.GarnetPassword, cfg.GarnetDB)
	default:
		log.Fatalf("unsupported backend %q", cfg.Backend)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("store close error: %v", err)
		}
	}()

	srv := api.New(cfg.ListenAddr, cfg.NodeID, store)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received signal %s", sig)
	case err := <-errCh:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
