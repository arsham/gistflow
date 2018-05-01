package window

import (
	"fmt"

	"github.com/therecipe/qt/widgets"
)

type boxLogger interface {
	error(msg string)
	warning(msg string)
	warningf(format string, a ...interface{})
}
type mb struct{ dialog *widgets.QMainWindow }

func messagebox(parent *widgets.QMainWindow) *mb { return &mb{parent} }
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
