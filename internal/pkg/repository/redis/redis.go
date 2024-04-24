package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	return &Redis{client: redis.NewClient(&redis.Options{})}
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(context.Background(), key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}
	return res.Val(), nil
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
