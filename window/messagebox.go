package window

import (
	"fmt"

	"github.com/therecipe/qt/widgets"
)

type mb struct{ dialog *widgets.QWidget }

func messagebox(parent *widgets.QWidget) *mb {
	return &mb{parent}
}

func (m mb) error(msg string) {
	qmb := widgets.NewQMessageBox(nil)
	qmb.Critical(m.dialog, "Warning", msg, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func (m mb) warning(msg string) {
	qmb := widgets.NewQMessageBox(nil)
	qmb.Warning(m.dialog, "Warning", msg, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func (m mb) warningf(format string, a ...interface{}) {
	m.warning(fmt.Sprintf(format, a...))
}
