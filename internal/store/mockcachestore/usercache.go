package mockcachestore

import (
	"errors"
	"time"
)

type UserCache struct {
	client map[string][]byte
}

func newUserCache() *UserCache {
	return &UserCache{
		client: make(map[string][]byte),
	}
}

func (u *UserCache) Get(token string) ([]byte, error) {
	return u.client[token], nil
}

func (u *UserCache) Set(token string, data []byte, expire time.Time) error {
	ttl := expire.Sub(time.Now().UTC())
	if ttl < 0 {
		return errors.New("token expired")
	}
	u.client[token] = data
	return nil
}
func (u *UserCache) Delete(token string) error {
	delete(u.client, token)
	return nil
}