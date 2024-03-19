package cachestore

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct{
	client *redis.Client
}

func newUserCache(client *redis.Client) *UserCache {
	return &UserCache{
		client: client,
	}
}

func (u *UserCache) Get(token string) ([]byte, error) {
	ctx := context.Background()
	data, err := u.client.Get(ctx, token).Result()
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func (u *UserCache) Set(token string, data []byte, expire time.Time) error {
	ctx := context.Background()
	ttl := expire.Sub(time.Now().UTC())
	if ttl < 0 {
		return fmt.Errorf("token expired: %v seconds ago", math.Abs(math.Round(ttl.Seconds())))
	}
	return u.client.Set(ctx, token, data, ttl).Err()
}

func (u *UserCache) Delete(token string) error {
	ctx := context.Background()
	return u.client.Del(ctx, token).Err()
}
