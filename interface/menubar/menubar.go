// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package menubar

import (
	"github.com/therecipe/qt/widgets"
)

// MenuBar holds actions for toolbars and menubar.
type MenuBar struct {
	widgets.QMenuBar

	_ func()     `constructor:"init"`
	_ func()     `signal:"quit"`
	_ func(bool) `signal:"copyURLToClipboard"`
	_ func(bool) `signal:"openInBrowser"`

	_ *widgets.QMenu `property:"options"`
	_ *widgets.QMenu `property:"window"`
	_ *widgets.QMenu `property:"edit"`
	_ *Action        `property:"actions"`
}

func init() {
	MenuBar_QRegisterMetaType()
}

func (m *MenuBar) init() {
	m.SetActions(NewAction(m))

	m.SetOptions(m.AddMenu2("&Options"))
	m.Options().SetObjectName("menuOptions")
	m.Options().AddActions([]*widgets.QAction{
		m.Actions().Settings,
		m.Actions().Sync,
		m.AddSeparator(),
		m.Actions().Quit,
	})

	m.SetEdit(m.AddMenu2("&Edit"))
	m.Edit().SetObjectName("edit")
	m.Edit().AddActions([]*widgets.QAction{
		m.Actions().InBrowser,
		m.Actions().CopyURL,
	})

	m.SetWindow(m.AddMenu2("&Window"))
	m.Window().SetObjectName("menuWindow")
	m.Window().AddActions([]*widgets.QAction{
		m.Actions().Toolbar,
		m.Actions().GistList,
	})

	m.Actions().CopyURL.ConnectTriggered(m.CopyURLToClipboard)
	m.Actions().Quit.ConnectTriggered(func(bool) {
		m.Quit()
	})
	m.Actions().InBrowser.ConnectTriggered(m.OpenInBrowser)
}
