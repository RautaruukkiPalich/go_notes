package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rautaruukkipalich/go_notes/internal/store"
)

type (
	Store struct {
		db *sql.DB
		noteRepo *NoteRepo
	}

	NoteRepo struct {
		sqlstore *Store
		stmts []sql.Stmt
	}
)

func New(db *sql.DB) (*Store, error) {
	store := &Store{db: db}

	noteRepo, err := newNoteRepo(store)
	if err != nil {
		return nil, err
	}

	store.noteRepo = noteRepo
	
	return store, nil
}

func newNoteRepo(s *Store) (*NoteRepo, error) {
	return &NoteRepo{
		sqlstore: s,
	}, nil
}

func (s *Store) Note() store.NoteRepo {
	return s.noteRepo
}