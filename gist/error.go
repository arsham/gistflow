// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

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
