// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"github.com/therecipe/qt/widgets"
)

// Tab is a widget shown on the QTabWidget.
type Tab struct {
	widgets.QTabWidget

	_ func() `constructor:"init"`

	textEdit *widgets.QPlainTextEdit
	gist     *tabGist
}

func (t *Tab) init() {
	layout := widgets.NewQVBoxLayout()
	t.SetLayout(layout)
	t.textEdit = widgets.NewQPlainTextEdit(t)
	t.textEdit.SetObjectName("content")

	layout.AddWidget(t.textEdit, 0, 0)
}

func (t *Tab) showGist(tabWidget *widgets.QTabWidget, g *tabGist) {
	t.textEdit.SetPlainText(g.content)
	tabWidget.AddTab(t, g.label)
	tabWidget.SetCurrentWidget(t)
	t.gist = g
}

func (t Tab) content() string {
	return t.gist.content
}

func (t Tab) url() string {
	return t.gist.url
}
