package model

import (
	"fmt"
	"time"
)

type User struct {
	ID       int       `json:"id"`
	TokenTTL time.Time `json:"token_ttl"`
}

func (u *User) Str() string {
	return fmt.Sprintf("User: id %v; ttl: %v", u.ID, u.TokenTTL)
}