package server

import "fmt"

var (
	ErrParseURL            = fmt.Errorf("url parse error")
	ErrNoPermissons        = fmt.Errorf("no permissions")
	ErrNotFound            = fmt.Errorf("not found")
	ErrInternalServerError = fmt.Errorf("internal error")
)