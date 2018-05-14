// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package toolbar

import (
	"github.com/arsham/gisty/qt/menubar"
	"github.com/therecipe/qt/widgets"
)

// Toolbar holds all actions that should show on the main toolbar.
type Toolbar struct {
	widgets.QToolBar

	_ func() `constructor:"init"`

	action *menubar.Action
}

func init() {
	Toolbar_QRegisterMetaType()
}

func (t *Toolbar) init() {
	t.SetObjectName("toolBar")

	sizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Preferred, widgets.QSizePolicy__DefaultType)
	sizePolicy.SetHorizontalStretch(0)
	sizePolicy.SetVerticalStretch(45)
	sizePolicy.SetHeightForWidth(t.SizePolicy().HasHeightForWidth())
	t.SetSizePolicy(sizePolicy)
	t.SetMinimumSize2(0, 45)
	t.SetBaseSize2(0, 45)
	t.SetFloatable(false)

}

// SetAction adds some of the menubar actions.
func (t *Toolbar) SetAction(a *menubar.Action) {
	t.AddActions([]*widgets.QAction{
		a.NewGist,
		a.InBrowser,
		a.CopyURL,
	})
}
