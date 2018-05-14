// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package gist

// Gist represents a gist coming back from a list response or when requesting a
// single gist.
type Gist struct {
	ID          string          `json:"id"`
	URL         string          `json:"url"`
	HTMLURL     string          `json:"html_url"`
	Description string          `json:"description"`
	Public      bool            `json:"public"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Files       map[string]File `json:"files"`
}

// File is one file in a Gist.
type File struct {
	Content string `json:"content"`
}
