package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisProcessedStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisProcessedStore(client *redis.Client, ttl time.Duration) *RedisProcessedStore {
	return &RedisProcessedStore{client: client, ttl: ttl}
}

func (s *RedisProcessedStore) key(orderID string) string {
	return fmt.Sprintf("processed:%s", orderID)
}

func (s *RedisProcessedStore) Seen(ctx context.Context, orderID string) (bool, error) {
	res, err := s.client.Exists(ctx, s.key(orderID)).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (s *RedisProcessedStore) Mark(ctx context.Context, orderID string) error {
	return s.client.Set(ctx, s.key(orderID), "1", s.ttl).Err()
}
