// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/widgets"
)

// Tab is a widget shown on the QTabWidget.
type Tab struct {
	widgets.QTabWidget

	_ func()       `constructor:"init"`
	_ func(string) `signal:"copyToClipboard"`

	_ *widgets.QVBoxLayout `property:"vBoxLayout"`
	_ []file               `property:"files"`
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
}

func (t *Tab) showGist(tabWidget *widgets.QTabWidget, g *gist.Gist) {
	for label, g := range g.Files {
		f := NewFile(t, 0)
		f.Content().SetText(g.Content)
		f.Information().SetText(label)
		t.VBoxLayout().AddWidget(f, 0, 0)
		t.SetFiles(append(t.Files(), f))
		f.ConnectCopyToClipboard(t.CopyToClipboard)
	}
	for label := range g.Files {
		tabWidget.AddTab(t, label)
		break
	}
	tabWidget.SetCurrentWidget(t)
	t.SetGist(g)
}

func (t Tab) url() string {
	return t.Gist().URL
}

func (t Tab) htmlURL() string {
	return t.Gist().HTMLURL
}
