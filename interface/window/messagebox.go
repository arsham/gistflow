// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"fmt"

	"github.com/therecipe/qt/widgets"
)

type boxLogger interface {
	Error(msg string)
	Warning(msg string)
	Warningf(format string, a ...interface{})
}
type mb struct{ dialog *MainWindow }

func messagebox(parent *MainWindow) *mb { return &mb{parent} }
func (m mb) Error(msg string) {
	qmb := widgets.NewQMessageBox(nil)
	qmb.Critical(m.dialog, "Warning", msg, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func (m mb) Warning(msg string) {
	qmb := widgets.NewQMessageBox(nil)
	qmb.Warning(m.dialog, "Warning", msg, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func (m mb) Warningf(format string, a ...interface{}) {
	m.Warning(fmt.Sprintf(format, a...))
}
