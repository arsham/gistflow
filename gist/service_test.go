// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package gist_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
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

type logger struct {
	errorFunc   func(string)
	warningFunc func(string)
}

func (l logger) Error(msg string)                         { l.errorFunc(msg) }
func (l logger) Warning(msg string)                       { l.warningFunc(msg) }
func (l logger) Warningf(format string, a ...interface{}) { l.Warning(fmt.Sprintf(format, a...)) }

func getLogger() *logger {
	return &logger{
		errorFunc:   func(string) {},
		warningFunc: func(string) {},
	}
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
		Token:    "s9jO",
		API:      ts.URL,
	}
	l, err := s.List(10, 1)
	if err != nil {
		t.Errorf("g.List(): err = %v, want nil", err)
	}
	if l == nil {
		t.Error("g.List(): l = nil, want []Gist")
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
		Token:    "QnqPp208",
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
		Token:    "dGMt7",
		API:      ts.URL,
		Logger:   getLogger(),
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
		Token:    "9z3PUQ93y3ww",
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

	gists := map[string]gist.Gist{
		id1: gist.Gist{
			Files: map[string]gist.File{
				"file1": gist.File{Content: "WCLwqKzLvzg"},
			},
		},
		id2: gist.Gist{
			Files: map[string]gist.File{
				"file1": gist.File{Content: "TLsplcHpevo"},
				"file2": gist.File{Content: "mbcFO"},
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
		Token:    "xqrIq6LUY2oIIjWGnKG",
		API:      ts.URL,
		CacheDir: loc,
	}
	currentID = id1
	if _, err = s.Get(currentID); err != nil {
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

func TestGistUpdateError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("r.Method = %s, want PATCH", r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "NwXYPlSH8Aga",
		Logger:   getLogger(),
	}
	g := gist.Gist{
		ID:          "koqeTgkrRPkorZUz1hIG",
		URL:         ts.URL,
		Description: "WyQS",
		Files: map[string]gist.File{
			"file1": gist.File{Content: "ydPlYOIZuk"},
		},
	}

	if err := s.Update(g); err == nil {
		t.Error("g.Update(): err = nil, want error")
	}
}

func TestGistUpdate(t *testing.T) {
	var g gist.Gist
	loc, err := ioutil.TempDir("", "gisty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("r.Method = %s, want PATCH", r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b, err := json.Marshal(g)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "4V07gOec7oKO",
		Logger:   getLogger(),
		CacheDir: loc,
	}
	g = gist.Gist{
		ID:          "8piWdiCyQbDRLU4TGc",
		URL:         ts.URL,
		Description: "M81ZYa5m",
		Files: map[string]gist.File{
			"file1": gist.File{Content: "KnL4DiGPsA"},
		},
	}

	if err := s.Update(g); err != nil {
		t.Errorf("g.Update(): err = %v, want nil", err)
	}
}

func TestNewGistBadURLError(t *testing.T) {
	var (
		g     gist.Gist
		url   string
		state bool
	)
	loc, err := ioutil.TempDir("", "gisty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state = strings.Contains(r.URL.Path, "/gists")
		url = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "thQrny",
		Logger:   getLogger(),
		API:      ts.URL,
		CacheDir: loc,
	}
	g = gist.Gist{
		ID:          "zwFaGyRUaPFSLr",
		Description: "CPOQOZhAaTXH",
		Files: map[string]gist.File{
			"file1": gist.File{Content: "KcYhTNKTowyl"},
		},
	}

	if err := s.Create(g); err != nil {
		t.Errorf("g.Create() = %s, want nil", err)
	}
	if !state {
		t.Errorf("%s should contain `gists`", url)
	}
}

func TestNewGistError(t *testing.T) {
	var (
		g      gist.Gist
		called bool
	)
	loc, err := ioutil.TempDir("", "gisty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "93N1Eb",
		Logger:   getLogger(),
		CacheDir: loc,
		API:      ts.URL,
	}
	g = gist.Gist{
		ID:          "amCp6p",
		Description: "JufnaH",
		Files: map[string]gist.File{
			"file1": gist.File{Content: "8jx3KQcr1"},
		},
	}

	if err := s.Create(g); err == nil {
		t.Error("g.Create() = nil, want error")
	}
	if !called {
		t.Error("server wasn't called")
	}
}

func TestNewGist(t *testing.T) {
	testNewGist(t, http.StatusOK)
	testNewGist(t, http.StatusCreated)
}
func testNewGist(t *testing.T, code int) {
	var g gist.Gist
	loc, err := ioutil.TempDir("", "gisty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("r.Method = %s, want POST", r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b, err := json.Marshal(g)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(code)
		w.Write(b)
	}))
	defer ts.Close()
	s := &gist.Service{
		Username: "arsham",
		Token:    "93N1Eb",
		Logger:   getLogger(),
		API:      ts.URL,
		CacheDir: loc,
	}
	g = gist.Gist{
		ID:          "amCp6p",
		Description: "JufnaH",
		Files: map[string]gist.File{
			"file1": gist.File{Content: "8jx3KQcr1"},
		},
	}

	if err := s.Create(g); err != nil {
		t.Errorf("code(%d): g.Create() = %v, want nil", code, err)
	}
}

func TestRemoveFile(t *testing.T) {
	var (
		called bool
		file1  = "file1"
		file2  = "file2"
	)
	var g gist.Gist
	loc, err := ioutil.TempDir("", "gisty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(loc)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != "PATCH" {
			t.Errorf("r.Method = %s, want PATCH", r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		g = gist.Gist{}
		if err := json.Unmarshal(body, &g); err != nil {
			t.Fatal(err)
		}
		if _, ok := g.Files[file1]; ok {
			t.Errorf("%s is in the request: %v", file1, g.Files)
		}
		if _, ok := g.Files[file2]; !ok {
			t.Errorf("%s is not in the request: %v", file2, g.Files)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(g.Files[file2], gist.File{}) {
			t.Errorf("g.Files[%s] = %v, want %v", file2, g.Files[file2], gist.Gist{})
		}
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	s := &gist.Service{
		Username: "arsham",
		Token:    "GonCvyFGU",
		Logger:   getLogger(),
		CacheDir: loc,
	}
	g = gist.Gist{
		ID:          "WjhdgmvZDWdjhd",
		URL:         ts.URL,
		Description: "IOgUKr5",
		Files: map[string]gist.File{
			file1: gist.File{Content: "6yZ9K645eb"},
			file2: gist.File{Content: "t8s3FySsVmYfZH72Qt"},
		},
	}

	if err := s.DeleteFile(g, file2); err != nil {
		t.Errorf("g.DeleteFile(): err = %v, want nil", err)
	}

	if !called {
		t.Error("endpoint wasn't called")
	}
}
