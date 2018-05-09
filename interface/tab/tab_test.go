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

func TestTab(t *testing.T) { tRunner.Run(func() { testTab(t) }) }
func testTab(t *testing.T) {
	tab := NewTab(widgets.NewQWidget(nil, 0))
	if tab.saveButton == nil {
		t.Error("tab.saveButton cannot be nil")
	}
	if tab.files == nil {
		t.Error("tab.files cannot be nil")
	}
	if tab.vBoxLayout == nil {
		t.Error("tab.vBoxLayout cannot be nil")
	}
	if tab.description == nil {
		t.Error("tab.description cannot be nil")
	}
}

func TestTabSetGistContents(t *testing.T) { tRunner.Run(func() { testTabSetGistContents(t) }) }
func testTabSetGistContents(t *testing.T) {
	var (
		content     = "q8poAXdCWLD0mO"
		description = "Dp8Szlq"
		fileName    = "4NWNWV"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	g := &gist.Gist{
		Description: description,
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
		},
	}
	tab.ShowGist(tabWidget, g)
	if tab.description.Text() != description {
		t.Errorf("tab.description.Text() = %s, want %s", tab.description.Text(), description)
	}

	file := tab.files[0]
	if file.Content().ToPlainText() != content {
		t.Errorf("file.Content().ToPlainText() = %s, want %s", file.Content().ToPlainText(), content)
	}
	if file.fileName.Text() != fileName {
		t.Errorf("file.fileName.Text() = %s, want %s", file.fileName.Text(), fileName)
	}
}

func TestSaveButton(t *testing.T) {
	tRunner.Run(func() {
		newContent := "0W8D6NEweKwlA3QZ"
		testSaveButton(t, "content", func(tab *Tab, file *File) {
			file.Content().SetText(newContent)
		})
		testSaveButton(t, "filename", func(tab *Tab, file *File) {
			file.SetFileName(newContent)
		})
		testSaveButton(t, "description", func(tab *Tab, file *File) {
			tab.description.SetText(newContent)
		})
	})
}
func testSaveButton(t *testing.T, name string, apply func(*Tab, *File)) {
	var (
		content  = "ox9Plo0zVXVMb7vZlygUoGdcR3g4ZRpo5f7pLBQWxDJY1hzw5v"
		id       = "34WvKeKjLx3Ol7BGNY"
		fileName = "yYWaPM"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	if tab.saveButton == nil {
		t.Errorf("%s: saveButton = nil", name)
		return
	}
	g := &gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
		},
	}
	var called bool
	tab.ConnectUpdateGist(func(*gist.Gist) {
		called = true
	})

	tab.ShowGist(tabWidget, g)
	if tab.saveButton.IsEnabled() {
		t.Errorf("%s: saveButton is already enabled", name)
	}

	file := tab.files[0]
	apply(tab, file)
	if !tab.saveButton.IsEnabled() {
		t.Errorf("%s: saveButton is not enabled", name)
	}

	tab.saveButton.Click()
	if !called {
		t.Errorf("%s: Save button didn't fire the signal", name)
	}

	if tab.saveButton.IsEnabled() {
		t.Errorf("%s: saveButton is still enabled", name)
	}
}

func TestSaveCheckContents(t *testing.T) { tRunner.Run(func() { testSaveCheckContents(t) }) }
func testSaveCheckContents(t *testing.T) {
	var (
		id          = "Botswana"
		fileName1   = "FMJupHSn1"
		content1    = "R4IlF7WB8"
		fileName2   = "L2DP0VD66axKD6pSwh"
		content2    = "Zx4A479U"
		description = "7G67dMf2"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	g := &gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			fileName1: gist.File{Content: content1},
		},
	}
	tab.ConnectUpdateGist(func(g *gist.Gist) {
		if _, ok := g.Files[fileName2]; !ok {
			t.Errorf("%s not found in files: %v", fileName2, g.Files)
		}
		if _, ok := g.Files[fileName1]; !ok {
			t.Errorf("%s is still in files: %v", fileName1, g.Files)
		}
		if g.Files[fileName2].Content != content2 {
			t.Errorf("New Content = %s, want %s", g.Files[fileName2].Content, content2)
		}
		if g.Description != description {
			t.Errorf("g.Description = %s, want %s", g.Description, description)
		}
	})

	tab.ShowGist(tabWidget, g)
	file := tab.files[0]
	file.Content().SetText(content2)
	file.SetFileName(fileName2)
	tab.description.SetText(description)

	tab.saveButton.Click()
}

func TestNewGist(t *testing.T) { tRunner.Run(func() { testNewGist(t) }) }
func testNewGist(t *testing.T) {
	label := "clO6lNMXnzBQN"
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(tabWidget)
	tab.NewGist(tabWidget, label)
	if !tab.saveButton.IsEnabled() {
		t.Error("SaveButton is disabled")
	}
	index := tabWidget.IndexOf(tab)
	if tabWidget.TabText(index) != label {
		t.Errorf("tabWidget.TabText(index) = %s, want %s", tabWidget.TabText(index), label)
	}
	if tabWidget.CurrentWidget().Pointer() != tab.Pointer() {
		t.Errorf("tabWidget.CurrentWidget() = %v, want %v", tabWidget.CurrentWidget().Pointer(), tab.Pointer())
	}
	if tab.gist == nil {
		t.Error("tab.gist cannot be nil")
	}
}

func TestNewGistCheckContents(t *testing.T) { tRunner.Run(func() { testNewGistCheckContents(t) }) }
func testNewGistCheckContents(t *testing.T) {
	var (
		called      bool
		label       = "RvaWxISBZmSk4M5"
		fileName    = "ChPmF2l3tu8SeeNf5z3w"
		content     = "aw3loDKGa6BP"
		description = "lKHF"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	tab.ConnectCreateGist(func(g *gist.Gist) {
		called = true
		if _, ok := g.Files[fileName]; !ok {
			t.Errorf("%s not found in files: %v", fileName, g.Files)
		}
		if g.Files[fileName].Content != content {
			t.Errorf("New Content = %s, want %s", g.Files[fileName].Content, content)
		}
		if g.Description != description {
			t.Errorf("g.Description = %s, want %s", g.Description, description)
		}
	})

	tab.NewGist(tabWidget, label)
	file := tab.files[0]
	file.Content().SetText(content)
	file.SetFileName(fileName)
	tab.description.SetText(description)

	tab.saveButton.Click()
	if !called {
		t.Error("didn't trigger the signal")
	}
}
