package cachestore

import (
	"github.com/rautaruukkipalich/go_notes/internal/store"
	"github.com/redis/go-redis/v9"
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

func (r *Redis) Note() store.NoteCache {
	return r.note
}

func (r *Redis) User() store.UserCache {
	return r.user
}