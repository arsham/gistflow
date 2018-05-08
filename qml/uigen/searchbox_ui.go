// WARNING! All changes made in this file will be lost!
package uigen

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type UISearchboxDialog struct {
	VerticalLayout *widgets.QVBoxLayout
	Input *widgets.QLineEdit
	Results *widgets.QListView
}

func (this *UISearchboxDialog) SetupUI(Dialog *widgets.QDialog) {
	Dialog.SetObjectName("Dialog")
	Dialog.SetWindowModality(core.Qt__ApplicationModal)
	Dialog.SetGeometry(core.NewQRect4(0, 0, 400, 300))
	this.VerticalLayout = widgets.NewQVBoxLayout2(Dialog)
	this.VerticalLayout.SetObjectName("verticalLayout")
	this.VerticalLayout.SetContentsMargins(0, 0, 0, 0)
	this.VerticalLayout.SetSpacing(0)
	this.Input = widgets.NewQLineEdit(Dialog)
	this.Input.SetObjectName("Input")
	this.VerticalLayout.AddWidget(this.Input, 0, 0)
	this.Results = widgets.NewQListView(Dialog)
	this.Results.SetObjectName("Results")
	this.VerticalLayout.AddWidget(this.Results, 0, 0)


    this.RetranslateUi(Dialog)

}

func (this *UISearchboxDialog) RetranslateUi(Dialog *widgets.QDialog) {
    _translate := core.QCoreApplication_Translate
	Dialog.SetWindowTitle(_translate("Dialog", "Dialog", "", -1))
}
