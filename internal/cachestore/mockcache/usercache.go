package mockcache

type UserCache struct{
	client map[string]any
}

func newUserCache() *UserCache {
	return &UserCache{
		client: make(map[string]any),
	}
}

func (u *UserCache) Get()    {}
func (u *UserCache) Set()    {}
func (u *UserCache) Delete() {}