// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type menuBar struct {
	widgets.QMenuBar

	_ func() `constructor:"init"`

	optionsMenu *widgets.QMenu
	quitAction  *widgets.QAction
}

func (m *menuBar) init() {
	m.optionsMenu = m.AddMenu2("&Options")
	m.quitAction = widgets.NewQAction2("&Quit", m)
	m.optionsMenu.AddActions([]*widgets.QAction{
		m.quitAction,
	})

	m.quitAction.SetShortcut(gui.QKeySequence_FromString("Ctrl+Q", 0))
}
