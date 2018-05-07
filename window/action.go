// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type appAction struct {
	widgets.QAction

	_ func() `constructor:"init"`

	actionInBrowser *widgets.QAction
	actionCopyURL   *widgets.QAction
	actionQuit      *widgets.QAction
	actionToolbar   *widgets.QAction
	actionGistList  *widgets.QAction
	actionSettings  *widgets.QAction
	actionSync      *widgets.QAction
}

func init() {
	appAction_QRegisterMetaType()
}

func (a *appAction) init() {
	a.actionInBrowser = widgets.NewQAction2("In Browser", a)
	a.actionInBrowser.SetObjectName("actionInBrowser")
	a.actionQuit = widgets.NewQAction2("&Quit", a)
	a.actionCopyURL = widgets.NewQAction2("Copy URL", a)
	a.actionCopyURL.SetObjectName("actionCopyURL")
	a.actionQuit = widgets.NewQAction2("&Quit", a)
	a.actionQuit.SetObjectName("actionQuit")
	a.actionToolbar = widgets.NewQAction2("Toolbar", a)
	a.actionToolbar.SetObjectName("actionToolbar")
	a.actionToolbar.SetCheckable(true)
	a.actionToolbar.SetChecked(true)
	a.actionGistList = widgets.NewQAction2("Gist List", a)
	a.actionGistList.SetObjectName("actionGistList")
	a.actionGistList.SetCheckable(true)
	a.actionGistList.SetChecked(true)
	a.actionSettings = widgets.NewQAction2("Settings", a)
	a.actionSettings.SetObjectName("actionSettings")
	a.actionSync = widgets.NewQAction2("Sync", a)
	a.actionSync.SetObjectName("actionSync")
	a.actionQuit.SetShortcut(gui.QKeySequence_FromString("Ctrl+Q", 0))
}
