package rediscache

import (
	"github.com/redis/go-redis/v9"
	"github.com/rautaruukkipalich/go_notes/internal/cachestore"
)

type (
	Redis struct {
		user *UserCache
		note *NoteCache
	}
)

func New(client *redis.Client) (*Redis, error) {
	return &Redis{
		user: newUserCache(client),
		note: newNoteCache(client),
	}, nil
}

func (r *Redis) Note() cachestore.NoteCache {
	return r.note
}

func (r *Redis) User() cachestore.UserCache {
	return r.user
}