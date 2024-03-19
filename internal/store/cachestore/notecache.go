package cachestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/rautaruukkipalich/go_notes/internal/model"
	"github.com/redis/go-redis/v9"
)

type NoteCache struct{
	client *redis.Client
}

func newNoteCache(client *redis.Client) *NoteCache {
	return &NoteCache{
		client: client,
	}
}

func (n *NoteCache) GetNoteById(id int) (model.Note, error) {
	var note model.Note

	ctx := context.Background()

	key := fmt.Sprintf("note:%d", id)
	data, err := n.client.Get(ctx, key).Result()
	if err != nil {
		return note, err
	}
	if err := json.Unmarshal([]byte(data), &note); err != nil {
		return note, err
	}
	return note, nil
}

func (n *NoteCache) Set(note *model.Note) error {

	ctx := context.Background()
	dur := time.Minute * 30

	id := "note:"+strconv.Itoa(note.ID)
	json_data, err := json.Marshal(&note)
	if err != nil {
		return err
	}
	return n.client.Set(ctx, id, json_data, dur).Err()
}

func (n *NoteCache) SetNotes(notes []model.Note) error {
	for _, note := range notes {
		if err := n.Set(&note); err != nil { return err}
	}
	return nil
}
func (n *NoteCache) Delete(id int) error {
	ctx := context.Background()
	key := "note:"+strconv.Itoa(id)
	return n.client.Del(ctx, key).Err()
}