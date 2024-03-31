package store

import (
	"time"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

type (
	NoteRepo interface {
		NoteRepoGetter
		NoteRepoSetter
		NoteRepoDeletter
	}

	NoteCache interface {
		NoteCacheGetter
		NoteCacheSetter
		NoteCacheDeletter
	}

	UserCache interface {
		UserCacheGetter
		UserCacheSetter
		UserCacheDeletter
	}
)

type (
	NoteRepoGetter interface {
		GetNoteById(int) (model.Note, error)
		GetNotes(int,string, int, int, int) ([]model.Note, error)
		HeatCache() ([]model.Note, error)
	}

	NoteRepoSetter interface {
		Set(*model.Note) (*model.Note, error)
		Patch(*model.Note) error
	}

	NoteRepoDeletter interface {
		Delete(int) error
	}
)

type (
	NoteCacheGetter interface {
		GetNoteById(int) (model.Note, error)
	}

	NoteCacheSetter interface {
		Set(*model.Note) error
		SetNotes([]model.Note) error
	}

	NoteCacheDeletter interface {
		Delete(int) error
	}
)

type (
	UserCacheGetter interface {
		Get(string) ([]byte, error)
	}

	UserCacheSetter interface {
		Set(string, []byte, time.Time) error
	}

	UserCacheDeletter interface {
		Delete(string) error
	}
)
