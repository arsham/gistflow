// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"github.com/therecipe/qt/widgets"
)

type menuBar struct {
	widgets.QMenuBar

	_ func() `constructor:"init"`
	_ func() `signal:"quit"`

	menuOptions *widgets.QMenu
	menuWindow  *widgets.QMenu
	menuAction  *widgets.QMenu
	action      *appAction
}

func (m *menuBar) init() {
	m.action = NewAppAction(m)
	m.menuOptions = m.AddMenu2("&Options")
	m.menuOptions.SetObjectName("menuOptions")

	m.menuOptions.AddActions([]*widgets.QAction{
		m.action.actionSettings,
		m.action.actionSync,
		m.AddSeparator(),
		m.action.actionQuit,
	})

	m.menuAction = m.AddMenu2("&Actions")
	m.menuAction.SetObjectName("menuAction")
	m.menuAction.AddActions([]*widgets.QAction{
		m.action.actionCopyURL,
		m.action.actionClipboard,
	})

	m.menuWindow = m.AddMenu2("&Window")
	m.menuWindow.SetObjectName("menuWindow")
	m.menuWindow.AddActions([]*widgets.QAction{
		m.action.actionToolbar,
		m.action.actionGistList,
	})

	m.action.actionQuit.ConnectTriggered(func(bool) {
		m.Quit()
	})
}
