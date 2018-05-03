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
type mb struct{ dialog *widgets.QMainWindow }

func messagebox(parent *widgets.QMainWindow) *mb { return &mb{parent} }
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
