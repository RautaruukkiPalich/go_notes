package mockstore

import (
	_ "github.com/lib/pq"
	"github.com/rautaruukkipalich/go_notes/internal/model"
	"github.com/rautaruukkipalich/go_notes/internal/store"
)

type (
	Store struct {
		noteRepo *NoteRepo
	}

	NoteRepo struct {
		store *Store
		notes map[string]*model.Note
	}
)

func New() (*Store, error) {
	store := &Store{}

	noteRepo, err := newNoteRepo(store)
	if err != nil {
		return nil, err
	}

	store.noteRepo = noteRepo
	
	return store, nil
}

func newNoteRepo(s *Store) (*NoteRepo, error) {
	return &NoteRepo{
		store: s,
		notes: make(map[string]*model.Note),
	}, nil
}

func (s *Store) Note() store.NoteRepo {
	return s.noteRepo
}