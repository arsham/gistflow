package gist

import "errors"

// Errors for Gist.
var (
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptyToken    = errors.New("token cannot be empty")
	ErrBadUsername   = errors.New("bad username")
	ErrEmptyID       = errors.New("id cannot be empty")
	ErrGistNotFound  = errors.New("gist not found")
)
