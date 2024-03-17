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
		GetNoteById(int) (*model.Note, error)
		GetNotes() ([]*model.Note, error)
	}

	NoteRepoCreater interface {
		Create()
	}

	NoteRepoSetter interface {
		Set()
	}

	NoteRepoDeletter interface {
		Delete()
	}
)

type (
	NoteCacheGetter interface {
		Get()
	}

	NoteCacheSetter interface {
		Set()
	}

	NoteCacheDeletter interface {
		Delete()
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
