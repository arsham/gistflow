package gist_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
		err      []error // any of errors
	}{
		{"no input", "", "", []error{gist.ErrEmptyUsername, gist.ErrEmptyToken}},
		{"no username", "", "XfJu", []error{gist.ErrEmptyUsername}},
		{"no token", "AdthCCaIXhhN", "", []error{gist.ErrEmptyToken}},
		{"spaces in username", "UMgEziO jLGLkhKcjG", "NbkGUkRlQNmIX", []error{gist.ErrBadUsername}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g := &gist.Service{
				Username: tc.username,
				Token:    tc.token,
			}
			_, err := g.List()
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
	l, err := s.List()
	if err != nil {
		t.Errorf("g.List(): err = %v, want nil", err)
	}
	if l == nil {
		t.Error("g.List(): l = nil, want []Response")
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
	_, err := s.Get("")
	if err == nil {
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

	_, err := s.Get("someID")
	if err == nil {
		t.Error("g.Get(): err = nil, want error", err)
	}
}
