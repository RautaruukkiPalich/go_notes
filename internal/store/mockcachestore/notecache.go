package mockcachestore

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

type NoteCache struct {
	client map[string][]byte
}

func newNoteCache() *NoteCache {
	return &NoteCache{
		client: make(map[string][]byte),
	}
}

func (n *NoteCache) Get() {}

func (n *NoteCache) GetNoteById(id int)(model.Note, error) { 
	var note model.Note

	key := fmt.Sprintf("note:%d", id)
	data:= n.client[key]
	if err := json.Unmarshal([]byte(data), &note); err != nil {
		return note, err
	}
	return note, nil
}

func (n *NoteCache) Set(note *model.Note) error { 
	id := "note:"+strconv.Itoa(note.ID)
	json_data, err := json.Marshal(&note)
	if err != nil {
		return err
	}
	n.client[id] = json_data
	return nil
}

func (n *NoteCache) SetNotes([]model.Note) error		{ return nil }
func (n *NoteCache) Delete(int) error					{ return nil}