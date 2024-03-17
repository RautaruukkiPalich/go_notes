package mockcache

import (
	"github.com/rautaruukkipalich/go_notes/internal/cachestore"
)

type (
	Cache struct {
		user *UserCache
		note *NoteCache
	}
)

func New() (*Cache, error) {
	return &Cache{
		user: newUserCache(),
		note: newNoteCache(),
	}, nil
}

func (c *Cache) Note() cachestore.NoteCache {
	return c.note
}

func (c *Cache) User() cachestore.UserCache {
	return c.user
}