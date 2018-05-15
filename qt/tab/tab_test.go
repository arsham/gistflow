// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/therecipe/qt/core"

	"github.com/arsham/gistflow/gist"
	"github.com/therecipe/qt/widgets"
)

type logger struct {
	errorFunc    func(string)
	criticalFunc func(string) widgets.QMessageBox__StandardButton
	warningFunc  func(string)
}

func (l logger) Error(msg string)                                        { l.errorFunc(msg) }
func (l logger) Critical(msg string) widgets.QMessageBox__StandardButton { return l.criticalFunc(msg) }
func (l logger) Warning(msg string)                                      { l.warningFunc(msg) }
func (l logger) Warningf(format string, a ...interface{})                { l.Warning(fmt.Sprintf(format, a...)) }

var app *widgets.QApplication

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func TestTab(t *testing.T) { tRunner.Run(func() { testTab(t) }) }
func testTab(t *testing.T) {
	tab := NewTab(widgets.NewQWidget(nil, 0))
	if tab.saveButton == nil {
		t.Error("tab.saveButton cannot be nil")
	}
	if tab.deleteButton == nil {
		t.Error("tab.deleteButton cannot be nil")
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
	if tab.messageBox == nil {
		t.Error("tab.messageBox cannot be nil")
	}
	if tab.publicCheckBox == nil {
		t.Error("tab.publicCheckBox cannot be nil")
	}
	if tab.addFileButton == nil {
		t.Error("tab.addFileButton cannot be nil")
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
			t.Errorf("%s is removed from files: %v", fileName1, g.Files)
			return
		}
		empty := gist.File{}
		if !reflect.DeepEqual(g.Files[fileName1], empty) {
			t.Errorf("g.Files[%s] = %s, want gist.File{}", fileName1, g.Files[fileName1])
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
	var (
		called bool
		label  = "clO6lNMXnzBQN"
	)
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
	if !tab.publicCheckBox.IsEnabled() {
		t.Error("tab.publicCheckBox is not enabled")
	}

	tab.ConnectCreateGist(func(g *gist.Gist) {
		called = true
	})
	tab.SaveButton().Click()
	if !called {
		t.Error("didn't trigger the signal")
	}

	// see TestOpenGistPublicCheckBox
	if tab.publicCheckBox.IsEnabled() {
		t.Error("tab.publicCheckBox is enabled")
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

func TestGistCreated(t *testing.T) { tRunner.Run(func() { testGistCreated(t) }) }
func testGistCreated(t *testing.T) {
	g := &gist.Gist{
		ID: "rH0xmdXVDsMl0D7a3",
	}
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	tab.NewGist(tabWidget, "Puf3rJc")
	tab.GistCreated(g)
	if !reflect.DeepEqual(tab.gist, g) {
		t.Errorf("tab.gist = %v, want %v", tab.gist, g)
	}
}

// test when we add a new file, it should just remove it, otherwise it should
// ask for permission.
func TestDeleteFileSignals(t *testing.T) { tRunner.Run(func() { testDeleteFileSignals(t) }) }
func testDeleteFileSignals(t *testing.T) {
	var (
		called   bool
		id       = "bciIJQRWq"
		fileName = "1qW6y7O"
		content  = "E4cMSsK75KN5G"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	g := &gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
		},
	}
	tab.ShowGist(tabWidget, g)

	file1 := tab.files[0]
	file1.messageBox = logger{
		criticalFunc: func(string) widgets.QMessageBox__StandardButton {
			called = true
			return widgets.QMessageBox__Ok
		},
	}

	tab.addFileButton.Click()
	file2 := tab.files[1]
	file2.messageBox = logger{
		criticalFunc: func(string) widgets.QMessageBox__StandardButton {
			t.Error("messagebox was shown")
			return widgets.QMessageBox__Ok
		},
	}

	file1.deleteButton.Click()
	if !called {
		t.Error("messagebox was not shown")
	}
	file2.deleteButton.Click()
}

func TestDeleteFileRemoveWidget(t *testing.T) { tRunner.Run(func() { testDeleteFileRemoveWidget(t) }) }
func testDeleteFileRemoveWidget(t *testing.T) {
	var (
		id        = "5XWadk"
		fileName1 = "PQYuZxEl64B0"
		fileName2 = "mBC"
		fileName3 = "3By6gLPs0kdsYd2gC7J"
		content   = "E7ewVidfDIXv"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	g := &gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			fileName1: gist.File{Content: content},
			fileName2: gist.File{Content: content},
			fileName3: gist.File{Content: content},
		},
	}
	tab.ShowGist(tabWidget, g)
	if len(tab.files) != len(g.Files) {
		t.Errorf("len(tab.files) = %d, want %d", len(tab.files), len(g.Files))
	}
	initialLen := len(g.Files)
	initialVBoxLen := tab.vBoxLayout.Count()
	tab.FileDeleted(fileName1)
	if len(tab.files) != initialLen-1 {
		t.Errorf("len(tab.files) = %d, want %d", len(tab.files), initialLen-1)
	}
	if len(tab.Gist().Files) != initialLen-1 {
		t.Errorf("len(tab.Gist().Files) = %d, want %d", len(tab.Gist().Files), initialLen-1)
	}
	if tab.vBoxLayout.Count() != initialVBoxLen-1 {
		t.Errorf("len(tab.vBoxLayout) = %d, want %d", tab.vBoxLayout.Count(), initialVBoxLen-1)
	}

	for _, file := range tab.files {
		if file.fileName.Text() == fileName1 {
			t.Errorf("didn't remove %s", fileName1)
		}
	}
	if c := tab.FindChild(fileName1, core.Qt__FindChildrenRecursively); c.Pointer() != nil {
		t.Errorf("didn't remove %s from layout", fileName1)
	}

	tcs := []struct {
		name   string
		exists bool
	}{
		{fileName1, false},
		{fileName2, true},
		{fileName3, true},
	}
	for _, tc := range tcs {
		g = tab.Gist()
		if _, ok := g.Files[tc.name]; ok != tc.exists {
			t.Errorf("g.Files[%s]: ok = %t, want %t", tc.name, ok, tc.exists)
		}
	}
}

func TestDeleteGistConfirmation(t *testing.T) { tRunner.Run(func() { testDeleteGistConfirmation(t) }) }
func testDeleteGistConfirmation(t *testing.T) {
	var (
		called       bool
		deleteCalled bool
		fileName     = "diO962wnDCMheXu"
		content      = "FxJKKgUP6T7b"
		button       widgets.QMessageBox__StandardButton
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	tab.messageBox = logger{
		criticalFunc: func(string) widgets.QMessageBox__StandardButton {
			called = true
			return button
		},
	}

	g := &gist.Gist{
		ID: "6XAQyCfcefvA",
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
		},
	}
	tab.ConnectDeleteGist(func(gs *gist.Gist) {
		deleteCalled = true
		if gs != g {
			t.Errorf("gs = %v, want %v", gs, g)
		}
	})
	tab.ShowGist(tabWidget, g)

	button = widgets.QMessageBox__Cancel
	tab.deleteButton.Click()
	if !called {
		t.Error("didn't trigger the signal")
	}
	if deleteCalled {
		t.Error("didn't expect to send deletion signal")
	}

	button = widgets.QMessageBox__Ok
	called = false
	tab.deleteButton.Click()
	if !called {
		t.Error("didn't trigger the signal")
	}
	if !deleteCalled {
		t.Error("didn't send deletion signal")
	}
}

func TestOpenGistPublicCheckBox(t *testing.T) { tRunner.Run(func() { testOpenGistPublicCheckBox(t) }) }
func testOpenGistPublicCheckBox(t *testing.T) {
	description := "qxwiWD1Kmlmqd"
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	g := &gist.Gist{
		Description: description,
	}
	tab.ShowGist(tabWidget, g)
	if tab.publicCheckBox.IsChecked() {
		t.Error("Gist is not public, but the checkbox is checked")
	}

	tab.DestroyQWidget()
	tab = NewTab(widgets.NewQWidget(nil, 0))
	g = &gist.Gist{
		Description: description,
		Public:      true,
	}
	tab.ShowGist(tabWidget, g)
	if !tab.publicCheckBox.IsChecked() {
		t.Error("Gist is public, but the checkbox is not checked")
	}
	// because the gist API doesn't have this functionality, it should be
	// disabled for opening existing gists.
	if tab.publicCheckBox.IsEnabled() {
		t.Error("tab.publicCheckBox is not disabled")
	}
}

func TestAddFile(t *testing.T) { tRunner.Run(func() { testAddFile(t) }) }
func testAddFile(t *testing.T) {
	var (
		description = "CyX221C3RpptC"
	)
	tabWidget := widgets.NewQTabWidget(nil)
	tab := NewTab(widgets.NewQWidget(nil, 0))
	g := &gist.Gist{
		Description: description,
	}
	tab.ShowGist(tabWidget, g)
	if tab.publicCheckBox.IsChecked() {
		t.Error("Gist is not public, but the checkbox is checked")
	}

	filesCount := len(tab.files)
	widgetsCount := tab.vBoxLayout.Count()
	tab.addFile()
	if len(tab.files) != filesCount+1 {
		t.Errorf("len(tab.files) = %d, want %d", len(tab.files), filesCount+1)
	}
	if tab.vBoxLayout.Count() != widgetsCount+1 {
		t.Errorf("tab.vBoxLayout.Count() = %d, want %d", tab.vBoxLayout.Count(), widgetsCount+1)
	}

	tab.addFileButton.Click()
	if len(tab.files) != filesCount+2 {
		t.Errorf("len(tab.files) = %d, want %d", len(tab.files), filesCount+2)
	}
	if tab.vBoxLayout.Count() != widgetsCount+2 {
		t.Errorf("tab.vBoxLayout.Count() = %d, want %d", tab.vBoxLayout.Count(), widgetsCount+2)
	}
}
