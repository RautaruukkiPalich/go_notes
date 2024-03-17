package cachestore

import "github.com/redis/go-redis/v9"

type NoteCache struct{
	client *redis.Client
}

func newNoteCache(client *redis.Client) *NoteCache {
	return &NoteCache{
		client: client,
	}
}

func (n *NoteCache) Get()    {}
func (n *NoteCache) Set()    {}
func (n *NoteCache) Delete() {}