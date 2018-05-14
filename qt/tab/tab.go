// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"github.com/arsham/gisty/gist"
	"github.com/arsham/gisty/qt/messagebox"
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
	_ func(*gist.Gist)         `signal:"GistCreated"`
	_ func(*gist.Gist)         `signal:"deleteGist"`

	// TODO: add dirty property
	messageBox messagebox.Message
	files      []*File
	gist       *gist.Gist

	description    *widgets.QLineEdit
	vBoxLayout     *widgets.QVBoxLayout // layout on gist level operations.
	saveButton     *widgets.QPushButton
	publicCheckBox *widgets.QCheckBox
	deleteButton   *widgets.QPushButton
	addFileButton  *widgets.QPushButton
}

func init() {
	Tab_QRegisterMetaType()
}

func (t *Tab) init() {
	layout := widgets.NewQVBoxLayout2(t)
	layout.SetObjectName("Inner Layout")
	t.vBoxLayout = layout
	t.SetLayout(layout)
	t.messageBox = messagebox.New(t)

	t.saveButton = widgets.NewQPushButton2("Save Gist", t)
	t.saveButton.SetToolTip("Saves the gist on github")
	t.deleteButton = widgets.NewQPushButton2("Delete", t)
	t.deleteButton.SetToolTip("Deletes the gist on github. This action is irreversible.")
	t.publicCheckBox = widgets.NewQCheckBox2("Public", t)
	t.addFileButton = widgets.NewQPushButton2("Add File", t)

	t.description = widgets.NewQLineEdit(t)
	t.description.SetToolTip("Set the gist's description")
	t.description.SetPlaceholderText("Description")
	butttons := widgets.NewQHBoxLayout()
	hLayout := widgets.NewQHBoxLayout()
	hLayout.AddWidget(t.publicCheckBox, 0, 0)
	hLayout.AddWidget(t.description, 0, 0)
	hLayout.AddLayout(butttons, 0)

	line := widgets.NewQFrame(t, core.Qt__Widget)
	line.SetFrameShadow(widgets.QFrame__Sunken)
	line.SetFrameShape(widgets.QFrame__HLine)

	layout.AddItem(hLayout)
	layout.AddWidget(line, 0, 0)
	butttons.AddWidget(t.deleteButton, 0, 0)
	butttons.AddWidget(t.addFileButton, 0, 0)
	butttons.AddWidget(t.saveButton, 0, 0)
	t.files = make([]*File, 0)
	t.description.ConnectTextChanged(func(string) {
		t.saveButton.SetEnabled(true)
	})
	t.ConnectFileDeleted(t.removeFile)
	t.deleteButton.ConnectClicked(func(bool) {
		b := t.messageBox.Critical("Are you sure you want to delete this gist?")
		if b == widgets.QMessageBox__Ok {
			t.DeleteGist(t.gist)
		}
	})
	t.addFileButton.ConnectClicked(func(bool) {
		f := t.addFile()
		// disabling the dialog.
		f.deleteButton.DisconnectClicked()
		f.deleteButton.ConnectClicked(func(bool) {
			f.DestroyQWidget()
		})
	})
}

// ShowGist shows each file in a separate container.
func (t *Tab) ShowGist(tabWidget *widgets.QTabWidget, g *gist.Gist) {
	for label, gf := range g.Files {
		f := t.addFile()
		f.SetObjectName(label)
		f.Content().SetText(gf.Content)
		f.ConnectDeleteFile(func(name string) {
			t.DeleteFile(g, name)
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
	if g.Public {
		t.publicCheckBox.SetChecked(true)
	}
	t.publicCheckBox.SetDisabled(true)
	t.publicCheckBox.SetToolTip("Unfortunately the API doesn't allow us to change this. You need to update it from your browser.")
}

// NewGist opens a new tab for creating a new gist.
func (t *Tab) NewGist(tabWidget *widgets.QTabWidget, label string) {
	t.addFile()
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
		g.Public = t.publicCheckBox.IsChecked()
		t.publicCheckBox.SetDisabled(true)
		t.CreateGist(g)
		t.saveButton.SetDisabled(true)
	})
	t.ConnectGistCreated(func(g *gist.Gist) {
		t.gist = g
		index := tabWidget.IndexOf(t)
		for label := range g.Files {
			tabWidget.SetTabText(index, label)
			break
		}
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
	// TODO: Protect this logic. If the name is not enlisted, it should show a
	// message.
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

func (t *Tab) addFile() *File {
	f := NewFile(t, 0)
	t.vBoxLayout.AddWidget(f, 0, 0)
	t.files = append(t.files, f)
	f.ConnectCopyToClipboard(t.CopyToClipboard)
	f.ConnectUpdateGist(func() {
		t.saveButton.SetEnabled(true)
	})
	return f
}
