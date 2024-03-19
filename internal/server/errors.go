package server

import "fmt"

var (
	ErrParseURL            = fmt.Errorf("url parse error")
	ErrUnprocessableEntity = fmt.Errorf("data parse error")
	ErrNoPermissons        = fmt.Errorf("no permissions")
	ErrNotFound            = fmt.Errorf("not found")
	ErrInternalServerError = fmt.Errorf("internal error")
	ErrInvalidToken 	   = fmt.Errorf("invalid token")
)