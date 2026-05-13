package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/cureeeeee/order-service/internal/domain"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{client: client, ttl: ttl}
}

func (c *RedisCache) key(orderID string) string {
	return fmt.Sprintf("order:%s", orderID)
}

func (c *RedisCache) Get(ctx context.Context, orderID string) (domain.Order, error) {
	data, err := c.client.Get(ctx, c.key(orderID)).Result()
	if err != nil {
		return domain.Order{}, err
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (c *RedisCache) Set(ctx context.Context, order domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key(order.ID), data, c.ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, orderID string) error {
	return c.client.Del(ctx, c.key(orderID)).Err()
}
