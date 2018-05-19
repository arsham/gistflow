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
	_ func(bool) `signal:"newGist"`
	_ func(bool) `signal:"copyURLToClipboard"`
	_ func(bool) `signal:"openInBrowser"`
	_ func(bool) `signal:"openSettings"`
	_ func(bool) `signal:"toggleToolbar"`
	_ func(bool) `signal:"toggleGistList"`

	options *widgets.QMenu
	window  *widgets.QMenu
	edit    *widgets.QMenu
	actions *Action
}

func init() {
	MenuBar_QRegisterMetaType()
}

func (m *MenuBar) init() {
	m.actions = NewAction(m)

	m.options = m.AddMenu2("&Options")
	m.options.SetObjectName("menuOptions")
	m.options.AddActions([]*widgets.QAction{
		m.actions.NewGist,
		m.AddSeparator(),
		m.actions.Settings,
		m.AddSeparator(),
		m.actions.Quit,
	})

	m.edit = m.AddMenu2("&Edit")
	m.edit.SetObjectName("edit")
	m.edit.AddActions([]*widgets.QAction{
		m.actions.InBrowser,
		m.actions.CopyURL,
	})

	m.window = m.AddMenu2("&Window")
	m.window.SetObjectName("menuWindow")
	m.window.AddActions([]*widgets.QAction{
		m.actions.Toolbar,
		m.actions.GistList,
	})

	m.actions.CopyURL.ConnectTriggered(m.CopyURLToClipboard)
	m.actions.Quit.ConnectTriggered(func(bool) {
		m.Quit()
	})
	m.actions.InBrowser.ConnectTriggered(m.OpenInBrowser)
	m.actions.NewGist.ConnectTriggered(m.NewGist)

	m.actions.Settings.ConnectTriggered(m.OpenSettings)
	m.actions.Toolbar.ConnectToggled(m.ToggleToolbar)
	m.actions.GistList.ConnectToggled(m.ToggleGistList)
}

// Actions returns all actions assigned to this MenuBar.
func (m *MenuBar) Actions() *Action { return m.actions }

// Options returns the options.
func (m *MenuBar) Options() *widgets.QMenu { return m.options }

// Window returns the window.
func (m *MenuBar) Window() *widgets.QMenu { return m.window }

// Edit returns the edit.
func (m *MenuBar) Edit() *widgets.QMenu { return m.edit }
