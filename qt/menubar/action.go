// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package menubar

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Action represents an action on toolbar or menubar.
type Action struct {
	widgets.QAction

	_ func() `constructor:"init"`

	InBrowser *widgets.QAction
	CopyURL   *widgets.QAction
	Quit      *widgets.QAction
	Toolbar   *widgets.QAction
	GistList  *widgets.QAction
	Settings  *widgets.QAction
	Sync      *widgets.QAction
	NewGist   *widgets.QAction
}

func init() {
	Action_QRegisterMetaType()
}

func (a *Action) init() {
	a.InBrowser = widgets.NewQAction2("In Browser", a)
	a.InBrowser.SetObjectName("ActionInBrowser")

	a.CopyURL = widgets.NewQAction2("Copy URL", a)
	a.CopyURL.SetObjectName("ActionCopyURL")

	a.Quit = widgets.NewQAction2("&Quit", a)
	a.Quit = widgets.NewQAction2("&Quit", a)
	a.Quit.SetObjectName("ActionQuit")
	a.Quit.SetShortcut(gui.QKeySequence_FromString("Ctrl+Q", 0))

	a.NewGist = widgets.NewQAction2("NewGist", a)
	a.NewGist.SetObjectName("ActionNew")
	a.NewGist.SetShortcut(gui.QKeySequence_FromString("Ctrl+N", 0))

	a.Toolbar = widgets.NewQAction2("Toolbar", a)
	a.Toolbar.SetObjectName("ActionToolbar")
	a.Toolbar.SetCheckable(true)
	a.Toolbar.SetChecked(true)

	a.GistList = widgets.NewQAction2("Gist List", a)
	a.GistList.SetObjectName("ActionGistList")
	a.GistList.SetCheckable(true)
	a.GistList.SetChecked(true)

	a.Settings = widgets.NewQAction2("Settings", a)
	a.Settings.SetObjectName("ActionSettings")

	a.Sync = widgets.NewQAction2("Sync", a)
	a.Sync.SetObjectName("ActionSync")
}
