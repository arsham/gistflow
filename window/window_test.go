// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func TestMain(m *testing.M) {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func TestWindowFetch(t *testing.T) {
	tRunner.Run(func() {
		testWindowFetch(t)
	})
}

func testWindowFetch(t *testing.T) {
	var (
		gres gist.Response
		done bool
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if done {
			w.Write([]byte("[\n]"))
			return
		}
		b, err := json.Marshal([]gist.Response{gres})
		if err != nil {
			t.Fatal(err)
		}
		done = true
		w.Write(b)
	}))
	defer ts.Close()
	gres = gist.Response{
		ID:          "QXhJNchXAK",
		URL:         fmt.Sprintf("%s/gists/%s", ts.URL, gres.ID),
		Description: "kfxLTwoCOkqEuPlp",
	}

	ws := Service{
		GistService: gist.Service{
			Username: "arsham",
			Token:    "token",
			API:      ts.URL,
		},
	}
	err := ws.displayMainWindow("../qml")
	if err != nil {
		t.Error(err)
		return
	}
	defer ws.window.Hide()
	lv := widgets.NewQListViewFromPointer(
		ws.window.FindChild("listView", core.Qt__FindChildrenRecursively).Pointer(),
	)

	model := lv.Model()
	item := model.Index(0, 0, core.NewQModelIndex())
	desc := item.Data(Description).ToString()
	id := item.Data(GistID).ToString()
	if desc != gres.Description {
		t.Errorf("Display = %s, want %s", desc, gres.Description)
	}
	if id != gres.ID {
		t.Errorf("Display = %s, want %s", id, gres.ID)
	}
}
