package gist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
// http://api.github.com/users/arsham/gists?access_token=4fe4218d35fa707d9a964142bd120ea5d37428e3
func (s *Service) List() ([]Response, error) {
	if s.Token == "" {
		return nil, ErrEmptyToken
	}
	if s.Username == "" {
		return nil, ErrEmptyUsername
	}
	if strings.Contains(s.Username, " ") {
		return nil, ErrBadUsername
	}

	url := fmt.Sprintf("%s/users/%s/gists?access_token=%s", s.api(), s.Username, s.Token)
	r, err := http.Get(url)
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

// Get gets a gist item by its id.
func (s *Service) Get(id string) (Gist, error) {
	if id == "" {
		return Gist{}, ErrEmptyID
	}

	url := fmt.Sprintf("%s/gists/%s?access_token=%s", s.api(), id, s.Token)
	r, err := http.Get(url)
	if err != nil {
		return Gist{}, err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotFound {
		return Gist{}, ErrGistNotFound
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Gist{}, err
	}

	var g Gist
	err = json.Unmarshal(body, &g)
	if err != nil {
		return Gist{}, err
	}
	return g, nil
}
