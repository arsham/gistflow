// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package searchbox

import (
	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	description = int(core.Qt__DisplayRole)
	gistID      = int(core.Qt__UserRole) + 1<<iota
	gistURL
	dialogWidth = 500
)

// Dialog is shown when user hists Ctrl+P.
type Dialog struct {
	widgets.QDialog

	_ func()              `constructor:"init"`
	_ func(*core.QRect)   `slot:"view"`
	_ func(gist.Response) `slot:"add"`
	_ func(string)        `signal:"openGist"`

	_ *widgets.QLineEdit `property:"input"`
	_ *widgets.QListView `property:"results"`

	model *ListModel
	proxy *core.QSortFilterProxyModel
}

func init() {
	Dialog_QRegisterMetaType()
}

func (d *Dialog) init() {
	d.SetWindowFlags(core.Qt__FramelessWindowHint)
	d.SetModal(true)
	d.SetInput(widgets.NewQLineEdit(d))
	d.Input().SetObjectName("Input")
	d.SetResults(widgets.NewQListView(d))
	d.Results().SetObjectName("Results")
	d.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			d.Hide()
		}
	})

	d.ConnectAdd(d.add)

	vLayout := widgets.NewQVBoxLayout2(d)
	vLayout.SetObjectName("verticalLayout")
	vLayout.SetContentsMargins(0, 0, 0, 0)
	vLayout.SetSpacing(0)
	vLayout.AddWidget(d.Input(), 0, 0)
	vLayout.AddWidget(d.Results(), 0, 0)
	d.Hide()
	d.ConnectView(d.view)
	d.model = NewListModel(d)
	d.proxy = core.NewQSortFilterProxyModel(d)
	d.proxy.SetSourceModel(d.model)
	d.Results().SetModel(d.proxy)

	d.Input().ConnectTextChanged(func(text string) {
		d.proxy.SetFilterWildcard(text)
		d.selectFirstRow()
	})

	d.ConnectKeyPressEvent(d.handleArrowKeys)
	d.Results().ConnectActivated(func(index *core.QModelIndex) {
		d.OpenGist(index.Data(gistID).ToString())
	})
}

func (d *Dialog) selectFirstRow() {
	index := d.Model().Index(0, 0, core.NewQModelIndex())
	d.Results().SelectionModel().Select(index, core.QItemSelectionModel__ClearAndSelect)
	d.Results().SetCurrentIndex(index)
}

func (d *Dialog) view(r *core.QRect) {
	c := core.NewQRect4(r.Width()/2-dialogWidth/2, 0, dialogWidth, 300)
	d.SetGeometry(c)
	d.Show()
	d.Input().SetFocus2()
	d.selectFirstRow()
}

// Model returns the results' model.
func (d *Dialog) Model() *core.QAbstractItemModel {
	return d.Results().Model()
}

func (d *Dialog) add(r gist.Response) {
	item := NewListItem(d)
	item.SetGistID(r.ID)
	item.SetGistURL(r.URL)
	description := r.Description
	item.SetDescription(description)
	d.model.AddGist(item)
}

// ID returns the ID of gist at row.
func (d *Dialog) ID(row int) string {
	index := d.Results().Model().Index(row, 0, core.NewQModelIndex())
	return index.Data(gistID).ToString()
}

// Description returns the Description of gist at row.
func (d *Dialog) Description(row int) string {
	index := d.Results().Model().Index(row, 0, core.NewQModelIndex())
	return index.Data(description).ToString()
}

func (d *Dialog) handleArrowKeys(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Up, core.Qt__Key_Down:
		d.Results().SetFocus2()
	}
}
