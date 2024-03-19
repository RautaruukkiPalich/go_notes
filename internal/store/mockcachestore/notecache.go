package mockcachestore

import "github.com/rautaruukkipalich/go_notes/internal/model"

type NoteCache struct {
	client map[string]any
}

func newNoteCache() *NoteCache {
	return &NoteCache{
		client: make(map[string]any),
	}
}

func (n *NoteCache) Get()                               {}
func (n *NoteCache) GetNoteById(int)(model.Note, error) { return model.Note{}, nil}
func (n *NoteCache) Set(*model.Note) error              { return nil }
func (n *NoteCache) SetNotes([]model.Note) error		{ return nil }
func (n *NoteCache) Delete(int) error					{ return nil}