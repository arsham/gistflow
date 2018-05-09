// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package gist

// Response is the response coming back from gist API.
type Response struct {
	ID          string          `json:"id"`
	URL         string          `json:"html_url"`
	Description string          `json:"description"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Files       map[string]File `json:"files"`
}

// Gist represents one gist.
type Gist struct {
	ID          string          `json:"id"`
	URL         string          `json:"url"`
	HTMLURL     string          `json:"html_url"`
	Description string          `json:"description"`
	Files       map[string]File `json:"files"`
}

// File is one file in a Gist.
type File struct {
	Content string `json:"content"`
}
