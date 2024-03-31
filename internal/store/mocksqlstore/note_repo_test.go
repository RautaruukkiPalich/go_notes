package mockstore_test

import (
	"testing"

	"github.com/rautaruukkipalich/go_notes/internal/model"
	mockstore "github.com/rautaruukkipalich/go_notes/internal/store/mocksqlstore"
	"github.com/stretchr/testify/assert"
)


func TestNoteRepository_Set(t *testing.T) {
	s, _ := mockstore.New()
	n, err := s.Note().Set(model.TestNote(t))
	assert.NoError(t, err)
	assert.NotNil(t, n)
}

func TestNoteRepository_GetNoteById(t *testing.T) {
	s, _ := mockstore.New()

	id := 213
	_, err := s.Note().GetNoteById(id)
	assert.Error(t, err)

	n, _ := s.Note().Set(model.TestNote(t))
	note, err := s.Note().GetNoteById(n.ID)
	assert.NoError(t, err)
	assert.NotNil(t, note)

}

func TestNoteRepository_GetNotes(t *testing.T) {
	s, _ := mockstore.New()


	userId, filter_body, filter_author, limit, offset := 1, "", 2, 2, 3
	_, err := s.Note().GetNotes(userId, filter_body, filter_author, limit, offset)
	assert.Error(t, err)

	for _, tc := range mockstore.Test_notes {
		note := tc
		s.Note().Set(&note)
	}

	userId, filter_body, filter_author, limit, offset = 1, "", 1, 2, 0
	notes, err := s.Note().GetNotes(userId, filter_body, filter_author, limit, offset)
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
	assert.Error(t, err)
	assert.Len(t, notes, 0)
}

func TestNoteRepository_Patch(t *testing.T) {
	s, _ := mockstore.New()

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
	s, _ := mockstore.New()

	err := s.Note().Delete(2)
	assert.Error(t, err)

	for _, tc := range mockstore.Test_notes {
		note := tc
		s.Note().Set(&note)
	}

	err = s.Note().Delete(2)
	assert.NoError(t, err)

	err = s.Note().Delete(2)
	assert.Error(t, err)

	err = s.Note().Delete(4)
	assert.NoError(t, err)

	notesAfter, _ := s.Note().GetNotes(0, "", 0, 10, 0)
	assert.Len(t, notesAfter, 3)
}