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

	_ func()           `constructor:"init"`
	_ func(string)     `signal:"copyToClipboard"`
	_ func(*gist.Gist) `signal:"updateGist"`

	_ *widgets.QVBoxLayout `property:"vBoxLayout"`
	_ *widgets.QPushButton `property:"save"`
	_ []File               `property:"files"`
	_ *gist.Gist           `property:"gist"`
}

func init() {
	Tab_QRegisterMetaType()
}

func (t *Tab) init() {
	layout := widgets.NewQVBoxLayout2(t)
	layout.SetObjectName("Inner Layout")
	t.SetVBoxLayout(layout)
	t.SetLayout(layout)

	butttons := widgets.NewQVBoxLayout2(nil)
	hLayout := widgets.NewQHBoxLayout()
	hSpacer := widgets.NewQSpacerItem(40, 20, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Minimum)
	hLayout.AddLayout(butttons, 0)
	hLayout.AddItem(hSpacer)

	line := widgets.NewQFrame(t, core.Qt__Widget)
	line.SetFrameShadow(widgets.QFrame__Sunken)
	line.SetFrameShape(widgets.QFrame__HLine)

	layout.AddItem(hLayout)
	layout.AddWidget(line, 0, 0)
	t.SetSave(widgets.NewQPushButton2("Save Gist", t))
	t.Save().SetDisabled(true)
	butttons.AddWidget(t.Save(), 0, 0)
}

// ShowGist shows each file in a separate container.
func (t *Tab) ShowGist(tabWidget *widgets.QTabWidget, g *gist.Gist) {
	for label, g := range g.Files {
		f := NewFile(t, 0)
		f.Content().SetText(g.Content)
		f.Information().SetText(label)
		t.VBoxLayout().AddWidget(f, 0, 0)
		t.SetFiles(append(t.Files(), f))
		f.ConnectCopyToClipboard(t.CopyToClipboard)
		f.ConnectUpdateGist(func() {
			t.Save().SetEnabled(true)
		})
		f.fileName = label
	}
	for label := range g.Files {
		tabWidget.AddTab(t, label)
		break
	}
	tabWidget.SetCurrentWidget(t)
	t.SetGist(g)
	t.Save().ConnectClicked(func(bool) {
		g := t.Gist()
		for _, f := range t.Files() {
			content := g.Files[f.fileName]
			content.Content = f.Content().ToPlainText()
			g.Files[f.fileName] = content
		}
		t.UpdateGist(g)
		t.Save().SetDisabled(true)
	})
}

// URL returns the main gist's URL.
func (t Tab) URL() string {
	return t.Gist().URL
}

// HTMLURL returns the URL to the html page of the gist.
func (t Tab) HTMLURL() string {
	return t.Gist().HTMLURL
}
