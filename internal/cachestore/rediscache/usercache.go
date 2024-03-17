package rediscache

import "github.com/redis/go-redis/v9"

type UserCache struct{
	client *redis.Client
}

func newUserCache(client *redis.Client) *UserCache {
	return &UserCache{
		client: client,
	}
}

func (u *UserCache) Get()    {}
func (u *UserCache) Set()    {}
func (u *UserCache) Delete() {}