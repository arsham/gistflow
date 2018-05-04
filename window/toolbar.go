// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import "github.com/therecipe/qt/widgets"

type appToolbar struct {
	widgets.QToolBar
	_ func() `constructor:"init"`

	action *appAction
}

func (a *appToolbar) init() {
	a.SetObjectName("toolBar")

	sizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Preferred, widgets.QSizePolicy__DefaultType)
	sizePolicy.SetHorizontalStretch(0)
	sizePolicy.SetVerticalStretch(45)
	sizePolicy.SetHeightForWidth(a.SizePolicy().HasHeightForWidth())
	a.SetSizePolicy(sizePolicy)
	a.SetMinimumSize2(0, 45)
	a.SetBaseSize2(0, 45)
	a.SetFloatable(false)
}

func (a *appToolbar) SetAction(action *appAction) {
	a.action = action
	a.AddActions([]*widgets.QAction{
		action.actionClipboard,
		action.actionCopyURL,
		action.actionSync,
	})
}
