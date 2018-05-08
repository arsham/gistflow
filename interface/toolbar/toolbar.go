// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package toolbar

import (
	"github.com/arsham/gisty/interface/menubar"
	"github.com/therecipe/qt/widgets"
)

// Toolbar holds all actions that should show on the main toolbar.
type Toolbar struct {
	widgets.QToolBar

	_ func()          `constructor:"init"`
	_ *menubar.Action `property:"action"`
}

func init() {
	Toolbar_QRegisterMetaType()
}

func (a *Toolbar) init() {
	a.SetObjectName("toolBar")

	sizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Preferred, widgets.QSizePolicy__DefaultType)
	sizePolicy.SetHorizontalStretch(0)
	sizePolicy.SetVerticalStretch(45)
	sizePolicy.SetHeightForWidth(a.SizePolicy().HasHeightForWidth())
	a.SetSizePolicy(sizePolicy)
	a.SetMinimumSize2(0, 45)
	a.SetBaseSize2(0, 45)
	a.SetFloatable(false)

	a.ConnectSetAction(func(action *menubar.Action) {
		a.AddActions([]*widgets.QAction{
			action.InBrowser,
			action.CopyURL,
			action.Sync,
		})
	})
}