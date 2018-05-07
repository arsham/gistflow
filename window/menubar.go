// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"github.com/therecipe/qt/widgets"
)

type menuBar struct {
	widgets.QMenuBar

	_ func()     `constructor:"init"`
	_ func()     `signal:"quit"`
	_ func(bool) `signal:"copyURLToClipboard"`
	_ func(bool) `signal:"openInBrowser"`

	_ *widgets.QMenu `property:"options"`
	_ *widgets.QMenu `property:"window"`
	_ *widgets.QMenu `property:"edit"`
	_ *appAction     `property:"actions"`
}

func init() {
	menuBar_QRegisterMetaType()
}

func (m *menuBar) init() {
	m.SetActions(NewAppAction(m))

	m.SetOptions(m.AddMenu2("&Options"))
	m.Options().SetObjectName("menuOptions")
	m.Options().AddActions([]*widgets.QAction{
		m.Actions().actionSettings,
		m.Actions().actionSync,
		m.AddSeparator(),
		m.Actions().actionQuit,
	})

	m.SetEdit(m.AddMenu2("&Edit"))
	m.Edit().SetObjectName("edit")
	m.Edit().AddActions([]*widgets.QAction{
		m.Actions().actionInBrowser,
		m.Actions().actionCopyURL,
	})

	m.SetWindow(m.AddMenu2("&Window"))
	m.Window().SetObjectName("menuWindow")
	m.Window().AddActions([]*widgets.QAction{
		m.Actions().actionToolbar,
		m.Actions().actionGistList,
	})

	m.Actions().actionCopyURL.ConnectTriggered(m.CopyURLToClipboard)
	m.Actions().actionQuit.ConnectTriggered(func(bool) {
		m.Quit()
	})
	m.Actions().actionInBrowser.ConnectTriggered(m.OpenInBrowser)
}
