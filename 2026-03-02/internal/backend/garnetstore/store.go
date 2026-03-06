package garnetstore

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"raft-garnet-bench/internal/backend"
)

type Store struct {
	client *redis.Client
}

func New(addr string, password string, db int) *Store {
	return &Store{client: redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		DialTimeout:     3 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolSize:        1024,
		MinIdleConns:    64,
		PoolTimeout:     4 * time.Second,
	})}
}

func (s *Store) Mode() string {
	return "garnet"
}

func (s *Store) Set(ctx context.Context, key, value string) error {
	return s.client.Set(ctx, key, value, 0).Err()
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	value, err := s.client.Get(ctx, key).Result()
	if err == nil {
		return value, nil
	}
	if errors.Is(err, redis.Nil) {
		return "", backend.ErrNotFound
	}
	return "", err
}

func (s *Store) Close() error {
	return s.client.Close()
}
