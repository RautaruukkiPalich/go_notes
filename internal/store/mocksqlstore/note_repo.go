package mockstore

import "github.com/rautaruukkipalich/go_notes/internal/model"

func (r *NoteRepo) GetNoteById(id int) (*model.Note, error) {
	return &model.Note{}, nil
}
func (r *NoteRepo) GetNotes() ([]*model.Note, error) {
	return []*model.Note{}, nil
}

func (r *NoteRepo) Set() {}

func (r *NoteRepo) Delete() {}

func (r *NoteRepo) Create() {}