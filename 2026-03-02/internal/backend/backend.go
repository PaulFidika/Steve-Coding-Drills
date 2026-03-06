package backend

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("key not found")

type NotLeaderError struct {
	LeaderID  uint64
	LeaderURL string
}

func (e *NotLeaderError) Error() string {
	if e == nil {
		return "not leader"
	}
	if e.LeaderURL != "" {
		return "not leader: forward to " + e.LeaderURL
	}
	return "not leader"
}

type Store interface {
	Mode() string
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Close() error
}
