// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"os"
	"testing"

	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/widgets"
)

func TestMain(m *testing.M) {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func TestSaveButton(t *testing.T) { tRunner.Run(func() { testSaveButton(t) }) }
func testSaveButton(t *testing.T) {
	var (
		content    = "ox9Plo0zVXVMb7vZlygUoGdcR3g4ZRpo5f7pLBQWxDJY1hzw5v"
		newContent = "0W8D6NEweKwlA3QZ"
		id         = "34WvKeKjLx3Ol7BGNY"
		fileName   = "yYWaPM"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	if tab.Save() == nil {
		t.Error("Save() = nil")
		return
	}
	g := &gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
		},
	}
	tab.ShowGist(tabWidget, g)
	if tab.Save().IsEnabled() {
		t.Error("Save() is already enabled")
	}
	file := tab.Files()[0]
	file.Content().SetText(newContent)
	if !tab.Save().IsEnabled() {
		t.Error("Save() is not enabled")
	}
	var called bool
	tab.ConnectUpdateGist(func(g2 *gist.Gist) {
		called = true
		if g2 != g {
			t.Errorf("g2 = %v, want %v", g2, g)
		}
		if g2.Files[fileName].Content != newContent {
			t.Errorf("new content = %s, want %s", g2.Files[fileName].Content, newContent)
		}
	})
	tab.Save().Click()
	if !called {
		t.Error("Save button didn't fire the signal")
	}

	if tab.Save().IsEnabled() {
		t.Error("Save() is still enabled")
	}
}
