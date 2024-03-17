package mockcachestore

import "time"

type UserCache struct {
	client map[string]any
}

func newUserCache() *UserCache {
	return &UserCache{
		client: make(map[string]any),
	}
}

func (u *UserCache) Get(token string) ([]byte, error) {
	return nil, nil
}
func (u *UserCache) Set(token string, data []byte, expire time.Time) error {
	return nil
}
func (u *UserCache) Delete(token string) error {
	return nil
}