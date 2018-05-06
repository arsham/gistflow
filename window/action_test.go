// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/gui"
)

func TestActions(t *testing.T) {
	tRunner.Run(func() {
		tcs := []func(*testing.T){
			testCopyContents,
			testCopyURL,
		}
		for _, tc := range tcs {
			tc(t)
		}
	})
}

func testCopyContents(t *testing.T) {
	var (
		name     = "test"
		id1      = "hxwjPAUr"
		content1 = "nIHxqYurtVgPxhJnoGvxOXPBqde"
		id2      = "JgXLvhmbEnSoIBAO"
		content2 = "RqnWSPAzagjzccGqpggWi"
		content  string
	)

	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gres := gist.Gist{
			Files: map[string]gist.File{
				"vtsmQN": gist.File{Content: content},
			},
		}
		b, err := json.Marshal(gres)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.setupUI()
	window.gistService.API = gistTs.URL

	clipboard := app.Clipboard()
	c := window.menubar.Actions().actionClipboard
	content = content1
	if err := window.openGist(id1); err != nil {
		t.Errorf("window.openGist(%s) = %v, want nil", id1, err)
	}
	content = content2
	if err := window.openGist(id2); err != nil {
		t.Errorf("window.openGist(%s) = %v, want nil", id2, err)
	}

	tab1, tab2 := window.tabGistList[id1], window.tabGistList[id2]
	window.TabsWidget().SetCurrentWidget(tab1)

	c.Trigger()
	if clipboard.Text(gui.QClipboard__Clipboard) != content1 {
		t.Errorf("clipboard.Text(gui.QClipboard__Clipboard) = `%s`, want `%s`", clipboard.Text(gui.QClipboard__Clipboard), content1)
	}

	window.TabsWidget().SetCurrentWidget(tab2)
	c.Trigger()
	if clipboard.Text(gui.QClipboard__Clipboard) != content2 {
		t.Errorf("clipboard.Text(gui.QClipboard__Clipboard) = `%s`, want `%s`", clipboard.Text(gui.QClipboard__Clipboard), content2)
	}
}

func testCopyURL(t *testing.T) {
	var (
		name     = "test"
		id1      = "wqWKsfoQevEbGjhmz"
		content  = "AFMQydAKTiJLa"
		id2      = "yuaosJCTsGUqEldvigi"
		api, url string
	)

	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = fmt.Sprintf("%s%s", api, r.URL.Path)
		gres := gist.Gist{
			Files: map[string]gist.File{
				"vtsmQN": gist.File{Content: content},
			},
		}
		b, err := json.Marshal(gres)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.setupUI()
	window.gistService.API = gistTs.URL
	api = gistTs.URL

	clipboard := app.Clipboard()
	c := window.menubar.Actions().actionCopyURL
	if err := window.openGist(id1); err != nil {
		t.Errorf("window.openGist(%s) = %v, want nil", id1, err)
	}
	url1 := url
	if err := window.openGist(id2); err != nil {
		t.Errorf("window.openGist(%s) = %v, want nil", id2, err)
	}
	url2 := url

	tab1, tab2 := window.tabGistList[id1], window.tabGistList[id2]
	window.TabsWidget().SetCurrentWidget(tab1)

	c.Trigger()
	if clipboard.Text(gui.QClipboard__Clipboard) != url1 {
		t.Errorf("clipboard.Text(gui.QClipboard__Clipboard) = `%s`, want `%s`", clipboard.Text(gui.QClipboard__Clipboard), url1)
	}

	window.TabsWidget().SetCurrentWidget(tab2)
	c.Trigger()
	if clipboard.Text(gui.QClipboard__Clipboard) != url2 {
		t.Errorf("clipboard.Text(gui.QClipboard__Clipboard) = `%s`, want `%s`", clipboard.Text(gui.QClipboard__Clipboard), url2)
	}
}
