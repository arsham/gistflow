// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// Tab is a widget shown on the QTabWidget.
type Tab struct {
	widgets.QTabWidget

	_ func()                   `constructor:"init"`
	_ func(string)             `signal:"copyToClipboard"`
	_ func(*gist.Gist, string) `signal:"deleteFile"`
	_ func(string)             `slot:"fileDeleted"`
	_ func(*gist.Gist)         `signal:"updateGist"`
	_ func(*gist.Gist)         `signal:"createGist"`

	// TODO: add dirty property
	vBoxLayout  *widgets.QVBoxLayout
	saveButton  *widgets.QPushButton
	description *widgets.QLineEdit
	files       []*File
	gist        *gist.Gist
}

func init() {
	Tab_QRegisterMetaType()
}

func (t *Tab) init() {
	layout := widgets.NewQVBoxLayout2(t)
	layout.SetObjectName("Inner Layout")
	t.vBoxLayout = layout
	t.SetLayout(layout)

	t.description = widgets.NewQLineEdit(t)
	t.description.SetToolTip("Set the gist's description")
	t.description.SetPlaceholderText("Description")
	butttons := widgets.NewQVBoxLayout2(nil)
	hLayout := widgets.NewQHBoxLayout()
	hLayout.AddWidget(t.description, 0, 0)
	hLayout.AddLayout(butttons, 0)

	line := widgets.NewQFrame(t, core.Qt__Widget)
	line.SetFrameShadow(widgets.QFrame__Sunken)
	line.SetFrameShape(widgets.QFrame__HLine)

	layout.AddItem(hLayout)
	layout.AddWidget(line, 0, 0)
	t.saveButton = widgets.NewQPushButton2("Save Gist", t)
	butttons.AddWidget(t.saveButton, 0, 0)
	t.files = make([]*File, 0)
	t.description.ConnectTextChanged(func(string) {
		t.saveButton.SetEnabled(true)
	})
	t.ConnectFileDeleted(t.removeFile)
}

// ShowGist shows each file in a separate container.
func (t *Tab) ShowGist(tabWidget *widgets.QTabWidget, g *gist.Gist) {
	for label, gf := range g.Files {
		f := NewFile(t, 0)
		f.SetObjectName(label)
		f.Content().SetText(gf.Content)
		t.vBoxLayout.AddWidget(f, 0, 0)
		t.files = append(t.files, f)
		f.ConnectCopyToClipboard(t.CopyToClipboard)
		f.ConnectDeleteFile(func(name string) {
			t.DeleteFile(g, name)
		})
		f.ConnectUpdateGist(func() {
			t.saveButton.SetEnabled(true)
		})
		f.SetFileName(label)
	}
	for label := range g.Files {
		tabWidget.AddTab(t, label)
		break
	}
	t.description.SetText(g.Description)
	tabWidget.SetCurrentWidget(t)
	t.gist = g
	t.saveButton.ConnectClicked(func(bool) {
		g := t.gist
		g.Description = t.description.Text()
		names := make(map[string]struct{}, len(t.files))
		for _, f := range t.files {
			content := g.Files[f.FileName()]
			content.Content = f.Content().ToPlainText()
			g.Files[f.FileName()] = content
			names[f.FileName()] = struct{}{}
		}
		for name := range g.Files {
			if _, ok := names[name]; !ok {
				g.Files[name] = gist.File{}
			}
		}
		t.UpdateGist(g)
		t.saveButton.SetDisabled(true)
	})
	t.saveButton.SetDisabled(true)
}

// NewGist opens a new tab for creating a new gist.
func (t *Tab) NewGist(tabWidget *widgets.QTabWidget, label string) {
	f := NewFile(t, 0)
	t.vBoxLayout.AddWidget(f, 0, 0)
	t.files = append(t.files, f)
	f.ConnectCopyToClipboard(t.CopyToClipboard)
	tabWidget.AddTab(t, label)
	tabWidget.SetCurrentWidget(t)
	t.saveButton.SetEnabled(true)
	t.gist = new(gist.Gist)

	t.saveButton.ConnectClicked(func(bool) {
		g := t.gist
		g.Description = t.description.Text()
		g.Files = make(map[string]gist.File, len(t.files))
		for _, f := range t.files {
			content := g.Files[f.FileName()]
			content.Content = f.Content().ToPlainText()
			g.Files[f.FileName()] = content
		}
		t.CreateGist(g)
		t.saveButton.SetDisabled(true)
	})
}

// URL returns the main gist's URL.
func (t Tab) URL() string {
	return t.gist.URL
}

// HTMLURL returns the URL to the html page of the gist.
func (t Tab) HTMLURL() string {
	return t.gist.HTMLURL
}

// Files returns the *File slice.
func (t *Tab) Files() []*File { return t.files }

// Gist returns Gist.
func (t *Tab) Gist() *gist.Gist { return t.gist }

// SaveButton returns SaveButton.
func (t *Tab) SaveButton() *widgets.QPushButton { return t.saveButton }

// SetDescription sets the description
func (t *Tab) SetDescription(text string) { t.description.SetText(text) }

// removeFile removes the file section that corresponds to name from the layout.
func (t *Tab) removeFile(name string) {
	for i, f := range t.files {
		if f.fileName.Text() == name {
			t.files = append(t.files[:i], t.files[i+1:]...)
			break
		}
	}
	c := t.FindChild(name, core.Qt__FindChildrenRecursively).Pointer()
	f := NewFileFromPointer(c)
	f.DestroyQWidget()

	delete(t.gist.Files, name)
}
