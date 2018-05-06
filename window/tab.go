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

	_ func()                  `constructor:"init"`
	_ *widgets.QPlainTextEdit `property:"editor"`
	_ *gist.Gist              `property:"gist"`
}

func init() {
	Tab_QRegisterMetaType()
}

func (t *Tab) init() {
	layout := widgets.NewQVBoxLayout()
	t.SetLayout(layout)
	t.SetEditor(widgets.NewQPlainTextEdit(t))
	t.Editor().SetObjectName("content")

	layout.AddWidget(t.Editor(), 0, 0)
}

func (t *Tab) showGist(tabWidget *widgets.QTabWidget, g *gist.Gist) {
	for label, g := range g.Files {
		t.Editor().SetPlainText(g.Content)
		tabWidget.AddTab(t, label)
		break
	}
	tabWidget.SetCurrentWidget(t)
	t.SetGist(g)
}

func (t Tab) content() string {
	return t.Editor().ToPlainText()
}

func (t Tab) url() string {
	return t.Gist().URL
}
