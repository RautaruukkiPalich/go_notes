package model

import (
	"time"
)

type Note struct {
	ID int `json:"id"`
	AuthorID  int  `json:"author_id"`
	Body string  `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	IsPublic bool  `json:"is_public"`
}
