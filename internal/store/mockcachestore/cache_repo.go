package mockcachestore

import (
	"github.com/rautaruukkipalich/go_notes/internal/store"
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

func (c *Cache) Note() store.NoteCache {
	return c.note
}

func (c *Cache) User() store.UserCache {
	return c.user
}