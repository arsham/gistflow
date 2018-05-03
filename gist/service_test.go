// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package gist_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/arsham/gisty/gist"
)

func anyError(err error, list []error) bool {
	for _, e := range list {
		if e == err {
			return true
		}
	}
	return false
}

func TestGistListErrors(t *testing.T) {
	tcs := []struct {
		name     string
		username string
		token    string
		perPage  int
		page     int
		err      []error // any of errors
	}{
		{"no input", "", "", 10, 100, []error{gist.ErrEmptyUsername, gist.ErrEmptyToken}},
		{"no username", "", "XfJu", 10, 100, []error{gist.ErrEmptyUsername}},
		{"no token", "AdthCCaIXhhN", "", 10, 100, []error{gist.ErrEmptyToken}},
		{"spaces in username", "UMgEziO jLGLkhKcjG", "NbkGUkRlQNmIX", 10, 100, []error{gist.ErrBadUsername}},
		{"zero per page", "UMgEziOjLGLadfsdfsfdf", "NbkGUkRlQNmIX", 0, 10, []error{gist.ErrPagination}},
		{"negative per page", "UMgEziOjLGLkhKcjG", "NbkGUkRlQNmIX", 10, -100, []error{gist.ErrPagination}},
		{"negative page", "UMgEziOjLGsdfLkhKcjG", "NbkGUkRlQNmIX", -10, 100, []error{gist.ErrPagination}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g := &gist.Service{
				Username: tc.username,
				Token:    tc.token,
			}
			_, err := g.List(tc.perPage, tc.page)
			if !anyError(err, tc.err) {
				t.Errorf("g.List(): err = %v, want any of %v", err, tc.err)
			}
		})
	}
}

func TestGistList(t *testing.T) {
	d, err := ioutil.ReadFile("testdata/gist1.txt")
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(d)
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "sometoken",
		API:      ts.URL,
	}
	l, err := s.List(10, 1)
	if err != nil {
		t.Errorf("g.List(): err = %v, want nil", err)
	}
	if l == nil {
		t.Error("g.List(): l = nil, want []Response")
	}
}

func TestIter(t *testing.T) {
	d, err := ioutil.ReadFile("testdata/gist1.txt")
	if err != nil {
		t.Fatal(err)
	}
	size := 10
	input := make(chan []byte, size)
	for i := 0; i < size; i++ {
		input <- d
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(<-input)
	}))
	close(input) // because we don't want the server waiting
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "sometoken",
		API:      ts.URL,
	}
	done := make(chan struct{})
	go func() {
		count := 0
		for r := range s.Iter() {
			count++
			if r.ID != "1b212f0843127d2d061f0d53fb581680" {
				t.Errorf("r.ID = %s, want %s", r.ID, "1b212f0843127d2d061f0d53fb581680")
			}
		}
		if count != size {
			t.Errorf("got %d iteration, want %d", count, size)
		}
		close(done)
	}()

	select {
	case <-time.After(time.Second):
		t.Error("Iter did not finish")
	case <-done:
	}
}

func TestGistGetError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{"files": {
    "filename": {
        "content": "something"
    }
}}
            `))
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "sometoken",
		API:      ts.URL,
	}
	if _, err := s.Get(""); err == nil {
		t.Error("g.Get(): err = nil, want error")
	}

	g, err := s.Get("someID")
	if err != nil {
		t.Errorf("g.Get(): err = %v, want nil", err)
	}

	if g.Files["filename"].Content != "something" {
		t.Errorf("content = %s, want `something`", g.Files["filename"].Content)
	}
}

func TestGistGetNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`
{"message":"Not Found","documentation_url":"https://developer.github.com/v3/gists/#get-a-single-gist"}
            `))
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "sometoken",
		API:      ts.URL,
	}

	if _, err := s.Get("someID"); err == nil {
		t.Error("g.Get(): err = nil, want error", err)
	}
}

func TestLoadFromCache(t *testing.T) {
	var (
		calls     = 0
		currentID string
		id1       = "NPrUmNnyLrgFcwIghuu"
		id2       = "DFxIrjJLcneZbqcpR"
	)

	loc, err := ioutil.TempDir("", "gisty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(loc)

	gists := map[string]gist.ResponseGist{
		id1: gist.ResponseGist{
			Files: map[string]gist.ResponseFile{
				"file1": gist.ResponseFile{"WCLwqKzLvzg"},
			},
		},
		id2: gist.ResponseGist{
			Files: map[string]gist.ResponseFile{
				"file1": gist.ResponseFile{"TLsplcHpevo"},
				"file2": gist.ResponseFile{"mbcFO"},
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		resp, err := json.Marshal(gists[currentID])
		if err != nil {
			t.Fatal(err)
		}
		w.Write(resp)
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "sometoken",
		API:      ts.URL,
		CacheDir: loc,
	}
	currentID = id1
	if _, err := s.Get(currentID); err != nil {
		t.Error("g.Get(): err = nil, want error")
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}

	currentID = id2
	r2, err := s.Get(currentID)
	if err != nil {
		t.Error("g.Get(): err = nil, want error")
	}
	if calls != 2 {
		t.Errorf("calls = %d, want 2", calls)
	}

	currentID = id2
	r3, err := s.Get(currentID)
	if err != nil {
		t.Error("g.Get(): err = nil, want error")
	}
	if !reflect.DeepEqual(r2, r3) {
		t.Errorf("r2 = %v, want %v", r2, r3)
	}
	if calls != 2 {
		t.Errorf("calls = %d, want 2", calls)
	}
}
