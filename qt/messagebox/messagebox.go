// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package messagebox

import (
	"fmt"

	"github.com/therecipe/qt/widgets"
)

// Message is a contract for showing messages in a modal QMessageBox.
type Message interface {
	Error(msg string)
	Critical(msg string) widgets.QMessageBox__StandardButton
	Warning(msg string)
	Warningf(format string, a ...interface{})
}

// MessageBox implements Message.
type MessageBox struct {
	parent widgets.QWidget_ITF
}

// New returns a new instance of MessageBox.
func New(parent widgets.QWidget_ITF) *MessageBox { return &MessageBox{parent} }

func (m MessageBox) Error(msg string) {
	qmb := widgets.NewQMessageBox(nil)
	qmb.Critical(m.parent, "Error", msg, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

// Critical offers an ok and a cancel button. The default is cancel.
func (m MessageBox) Critical(msg string) widgets.QMessageBox__StandardButton {
	qmb := widgets.NewQMessageBox(nil)
	return qmb.Critical(m.parent, "Critical", msg, widgets.QMessageBox__Ok|widgets.QMessageBox__Cancel, widgets.QMessageBox__Cancel)
}

// Warning shows a warning message.
func (m MessageBox) Warning(msg string) {
	qmb := widgets.NewQMessageBox(nil)
	qmb.Warning(m.parent, "Warning", msg, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

// Warningf shows a warning message with formatting.
func (m MessageBox) Warningf(format string, a ...interface{}) {
	m.Warning(fmt.Sprintf(format, a...))
}
