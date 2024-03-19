package mockstore

import "github.com/rautaruukkipalich/go_notes/internal/model"

func (r *NoteRepo) GetNoteById(id int) (model.Note, error) {
	return model.Note{}, nil
}
func (r *NoteRepo) GetNotes(userId int, filter_body string, filter_author int, limit int, offset int) ([]model.Note, error) {
	return []model.Note{}, nil
}

func (r *NoteRepo) HeatCache() ([]model.Note, error) {
	return []model.Note{}, nil
}

func (r *NoteRepo) SetNotes() error {
	return nil
}

func (r *NoteRepo) Set(note *model.Note) error {
	return nil
}

func (r *NoteRepo) Patch(n *model.Note) error {
	return nil
}

func (r *NoteRepo) Delete(id int) error { return nil }
