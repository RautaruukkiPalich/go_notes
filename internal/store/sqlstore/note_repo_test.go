package sqlstore_test

import (
	"testing"

	"github.com/rautaruukkipalich/go_notes/internal/model"
	sqlstore "github.com/rautaruukkipalich/go_notes/internal/store/sqlstore"
	mockstore "github.com/rautaruukkipalich/go_notes/internal/store/mocksqlstore"
	"github.com/stretchr/testify/assert"
)

func TestNoteRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown("notes")
	s, _ := sqlstore.New(db)

	n, err := s.Note().Set(model.TestNote(t))
	assert.NoError(t, err)
	assert.NotNil(t, n)
}

func TestNoteRepository_GetNoteById(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown("notes")
	s, _ := sqlstore.New(db)

	id := 213
	_, err := s.Note().GetNoteById(id)
	assert.Error(t, err)

	n, _ := s.Note().Set(model.TestNote(t))
	note, err := s.Note().GetNoteById(n.ID)
	assert.NoError(t, err)
	assert.NotNil(t, note)
}

func TestNoteRepository_GetNotes(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown("notes")
	s, _ := sqlstore.New(db)

	userId, filter_body, filter_author, limit, offset := 1, "", 2, 2, 3
	notes, err := s.Note().GetNotes(userId, filter_body, filter_author, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, notes, 0)

	for _, tc := range mockstore.Test_notes {
		note := tc
		_, err := s.Note().Set(&note)
		if err != nil {
			t.Fatal(err)
		}
	}

	userId, filter_body, filter_author, limit, offset = 1, "", 1, 2, 0
	notes, err = s.Note().GetNotes(userId, filter_body, filter_author, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, notes, 2)

	userId, filter_body, filter_author, limit, offset = 0, "", 2, 2, 0
	notes, err = s.Note().GetNotes(userId, filter_body, filter_author, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, notes, 2)

	userId, filter_body, filter_author, limit, offset = 1, "2", 0, 2, 0
	notes, err = s.Note().GetNotes(userId, filter_body, filter_author, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, notes, 2)

	userId, filter_body, filter_author, limit, offset = 1, "", 0, 2, 7
	notes, err = s.Note().GetNotes(userId, filter_body, filter_author, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, notes, 0)
}

func TestNoteRepository_Patch(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown("notes")
	s, _ := sqlstore.New(db)

	err := s.Note().Patch(model.TestNote(t))
	assert.Error(t, err)

	n, _ := s.Note().Set(model.TestNote(t))

	note1 := *n
	note1.Body = "newText"
	err = s.Note().Patch(&note1)
	assert.NoError(t, err)

	note2 := *n
	note2.ID = 123
	err = s.Note().Patch(&note2)
	assert.Error(t, err)
}

func TestNoteRepository_Delete(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown("notes")
	s, _ := sqlstore.New(db)

	err := s.Note().Delete(2)
	assert.Error(t, err)

	for _, tc := range mockstore.Test_notes {
		note := tc
		s.Note().Set(&note)
	}

	notes, _ := s.Note().GetNotes(0, "", 0, 5, 0)
	ids := []int{}
	for _, note := range notes {
		ids = append(ids, note.ID)
	}

	err = s.Note().Delete(ids[len(ids) - 1] + 1)
	assert.Error(t, err)

	err = s.Note().Delete(ids[len(ids) - 1] + 5)
	assert.Error(t, err)

	err = s.Note().Delete(ids[len(ids) - 1])
	assert.NoError(t, err)

	notesAfter, _ := s.Note().GetNotes(0, "", 0, 5, 0)
	assert.Len(t, notesAfter, 4)
}