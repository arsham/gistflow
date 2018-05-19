// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

// Package gist communicates with api.github.com in order to retrieve and update
// user's gists.
package gist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

type boxLogger interface {
	Warning(msg string)
	Warningf(format string, a ...interface{})
}

// Service holds the information about the user.
type Service struct {
	Username string
	Token    string
	API      string
	CacheDir string
	Logger   boxLogger
}

func (s *Service) api() string {
	if s.API == "" {
		s.API = "https://api.github.com"
	}
	return s.API
}

// List fetches all gists for the user.
func (s *Service) List(perPage, page int) ([]Gist, error) {
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

	var res []Gist
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Iter returns a channel which emits new Gist objects. It will follow the
// paginations until it's exhausted.
func (s *Service) Iter() chan Gist {
	ch := make(chan Gist)
	go func() {
		perPage := 40
		for page := 1; ; page++ {
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
func (s *Service) Get(id string) (Gist, error) {
	var g Gist
	if id == "" {
		return Gist{}, ErrEmptyID
	}

	body, err := fromCache(s.CacheDir, id)
	switch err {
	case nil, ErrCacheNotExists, ErrEmptyCacheLoc:
	default:
		s.Logger.Warningf("reading from cache: %s", err.Error())
	}

	gistURL := fmt.Sprintf("%s/gists/%s", s.api(), id)
	if body != nil {
		if err = json.Unmarshal(body, &g); err == nil {
			g.URL = gistURL
			return g, nil
		}
		s.Logger.Warning(err.Error())
	}

	url := s.withToken(gistURL)
	r, err := http.Get(url)
	if err != nil {
		return Gist{}, err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotFound {
		return Gist{}, ErrGistNotFound
	}

	g, err = s.readAndCache(r.Body, id)
	if err != nil {
		return g, err
	}
	g.URL = gistURL
	return g, nil
}

// Update returns an error if the remote API responds other than 200.
func (s *Service) Update(g Gist) (Gist, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return Gist{}, err
	}
	url := s.withToken(g.URL)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(b))
	if err != nil {
		return Gist{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return Gist{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		reason, _ := ioutil.ReadAll(res.Body)
		return Gist{}, fmt.Errorf("error updating gist: %s", reason)
	}

	return s.readAndCache(res.Body, g.ID)
}

// Create returns an error if the remote API responds other than 200.
func (s *Service) Create(g Gist) (Gist, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return Gist{}, err
	}
	url := s.withToken(s.API + "/gists")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return Gist{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return Gist{}, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		reason, _ := ioutil.ReadAll(res.Body)
		return Gist{}, fmt.Errorf("error updating gist: %s", reason)
	}
	return s.readAndCache(res.Body, g.ID)
}

// DeleteFile sends a request to the server to remove a file.
func (s *Service) DeleteFile(g Gist, name string) (Gist, error) {
	url := s.withToken(g.URL)
	g = Gist{
		Files: map[string]File{
			name: File{},
		},
	}
	b, err := json.Marshal(g)
	if err != nil {
		return Gist{}, err
	}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(b))
	if err != nil {
		return Gist{}, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Gist{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		reason, _ := ioutil.ReadAll(res.Body)
		return Gist{}, fmt.Errorf("error updating gist: %s", reason)
	}
	return s.readAndCache(res.Body, g.ID)
}

// DeleteGist sends a request to the server to remove a gist.
func (s *Service) DeleteGist(id string) error {
	gistURL := fmt.Sprintf("%s/gists/%s", s.api(), id)
	url := s.withToken(gistURL)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error removing gist: %d", res.StatusCode)
	}
	return deleteCache(s.CacheDir, id)
}

func (s *Service) withToken(url string) string {
	return fmt.Sprintf("%s?access_token=%s", url, s.Token)
}

// readAndCache reads from r and fills out the returning Gist, while updating
// the cache entry for the Gist.
func (s *Service) readAndCache(r io.ReadCloser, id string) (Gist, error) {
	g := Gist{}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return Gist{}, err
	}
	err = json.Unmarshal(body, &g)
	if err != nil {
		return Gist{}, err
	}
	if id == "" {
		id = g.ID
	}
	if err = saveCache(s.CacheDir, id, body); err != nil {
		s.Logger.Warning(err.Error())
	}
	return g, nil
}

func fromCache(location, id string) ([]byte, error) {
	if location == "" {
		return nil, ErrEmptyCacheLoc
	}
	name := path.Join(location, id)
	file, err := os.Open(name)
	if err != nil {
		return nil, ErrCacheNotExists
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func deleteCache(location, id string) error {
	if location == "" {
		return ErrEmptyCacheLoc
	}
	name := path.Join(location, id)
	os.Remove(name)
	return nil
}

func saveCache(location, id string, contents []byte) error {
	if location == "" {
		return ErrEmptyCacheLoc
	}
	name := path.Join(location, id)
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(contents)
	return err
}
