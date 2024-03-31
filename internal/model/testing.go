package model

import "testing"

func TestNote(t *testing.T) *Note {
	t.Helper()
	return &Note{
		AuthorID: 1,
		Body: "testBody text",
		IsPublic: true,
	}
}