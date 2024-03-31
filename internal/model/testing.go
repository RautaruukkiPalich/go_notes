package model

import (
	"testing"
	"time"
)

func TestNote(t *testing.T) *Note {
	t.Helper()
	return &Note{
		AuthorID: 1,
		Body: "testBody text",
		IsPublic: true,
	}
}

func TestUser(t *testing.T) *User {
	t.Helper()
	return &User{
		ID: 1,
		TokenTTL: time.Now().UTC().Add(time.Minute * 2),
	}
}