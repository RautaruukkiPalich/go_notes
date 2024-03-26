package model

import "time"

type User struct {
	ID       int       `json:"id"`
	TokenTTL time.Time `json:"token_ttl"`
}