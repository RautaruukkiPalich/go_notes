package mockstore

import (
	"errors"
	"strings"
	"time"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

func (r *NoteRepo) GetNoteById(id int) (model.Note, error) {
	if r.notes[id] == nil {
		return model.Note{}, errors.New("note not found")
	}
	return *r.notes[id], nil

}

func checkIsValidWithFilters(n model.Note, userId int, filterAuthor int) bool {
	if n.IsPublic{
		return filterAuthor == 0 || filterAuthor == n.AuthorID
	} 

	if filterAuthor == 0 {
		return userId == n.AuthorID
	}

	return (userId == n.AuthorID) && (n.AuthorID == filterAuthor) 

}

func (r *NoteRepo) GetNotes(userId int, filter_body string, filter_author int, limit int, offset int) ([]model.Note, error) {
	output := []model.Note{}

	for _, note := range r.notes {

		if !strings.Contains(note.Body, filter_body) {
			continue
		}

		if !checkIsValidWithFilters(*note, userId, filter_author){
			continue
		}

		output = append(output, *note)
	}
	
	length := len(output)

	if length < offset {
		return []model.Note{}, errors.New("out of limit")
	}

	if length > offset + limit {
		length = offset + limit
	}

	return output[offset:length], nil
}

func (r *NoteRepo) HeatCache() ([]model.Note, error) {
	return []model.Note{}, nil
}

func (r *NoteRepo) Set(n *model.Note) (*model.Note, error) {
	if r.notes == nil {
		return nil, errors.New("Note already exists")
	}

	utcNow := time.Now().UTC()
	r.notes[r.currentId] = n
	n.ID = r.currentId
	n.CreatedAt = utcNow
	n.UpdatedAt = utcNow
	
	r.currentId++

	return n, nil
}

func (r *NoteRepo) Patch(n *model.Note) error {
	for key, note := range r.notes {
		if note.ID == n.ID {
			n.UpdatedAt = time.Now().UTC()
			r.notes[key] = n
			return nil
		}
	}
	return errors.New("note not in map")
}

func (r *NoteRepo) Delete(id int) error {
	for key, note := range r.notes {
		if note.ID == id {
			delete(r.notes, key)
			return nil
		}
	}
	return errors.New("ID not in map")
}
