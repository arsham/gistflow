// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

// Package gist communicates with api.github.com in order to retrieve and update
// user's gists.
package gist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Service holds the information about the user.
type Service struct {
	Username string
	Token    string
	API      string
}

func (s *Service) api() string {
	if s.API == "" {
		s.API = "https://api.github.com"
	}
	return s.API
}

// List fetches all gists for the user.
func (s *Service) List(perPage, page int) ([]Response, error) {
	if s.Token == "" {
		return nil, ErrEmptyToken
	}
	if s.Username == "" {
		return nil, ErrEmptyUsername
	}
	if strings.Contains(s.Username, " ") {
		return nil, ErrBadUsername
	}
	if perPage <= 0 || page < 0 {
		return nil, ErrPagination
	}
	urlPath := fmt.Sprintf("%s/users/%s/gists", s.api(), s.Username)
	url, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	v := url.Query()
	v.Add("access_token", s.Token)
	v.Add("page", strconv.Itoa(page))
	v.Add("per_page", strconv.Itoa(perPage))
	url.RawQuery = v.Encode()
	r, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var res []Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Iter returns a channel which emits new Response objects. It will follow the
// paginations until it's exhausted.
func (s *Service) Iter() chan Response {
	ch := make(chan Response)
	go func() {
		perPage := 40
		for page := 0; ; page += perPage {
			gs, err := s.List(perPage, page)
			if err != nil || len(gs) == 0 {
				break
			}
			for _, g := range gs {
				ch <- g
			}
		}
		close(ch)
	}()
	return ch
}

// Get gets a gist item by its id.
func (s *Service) Get(id string) (ResponseGist, error) {
	if id == "" {
		return ResponseGist{}, ErrEmptyID
	}

	url := fmt.Sprintf("%s/gists/%s?access_token=%s", s.api(), id, s.Token)
	r, err := http.Get(url)
	if err != nil {
		return ResponseGist{}, err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotFound {
		return ResponseGist{}, ErrGistNotFound
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ResponseGist{}, err
	}

	var g ResponseGist
	err = json.Unmarshal(body, &g)
	if err != nil {
		return ResponseGist{}, err
	}
	return g, nil
}
