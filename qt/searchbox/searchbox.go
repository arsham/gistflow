// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package searchbox

import (
	"strings"

	"github.com/arsham/gistflow/gist"
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

// Dialog is shown when user hits Ctrl+P.
type Dialog struct {
	widgets.QDialog

	_ func()            `constructor:"init"`
	_ func(*core.QRect) `slot:"view"`
	_ func(gist.Gist)   `slot:"add"`
	_ func(string)      `signal:"openGist"`

	input   *widgets.QLineEdit
	results *widgets.QListView
	model   *listModel
	proxy   *core.QSortFilterProxyModel
}

func init() {
	Dialog_QRegisterMetaType()
}

func (d *Dialog) init() {
	d.SetWindowFlags(core.Qt__FramelessWindowHint)
	d.SetModal(true)
	d.input = widgets.NewQLineEdit(d)
	d.input.SetObjectName("Input")
	d.results = widgets.NewQListView(d)
	d.results.SetObjectName("Results")
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
	vLayout.AddWidget(d.input, 0, 0)
	vLayout.AddWidget(d.results, 0, 0)
	d.Hide()
	d.ConnectView(d.view)
	d.model = NewListModel(d)
	d.proxy = core.NewQSortFilterProxyModel(d)
	d.proxy.SetSourceModel(d.model)
	d.results.SetModel(d.proxy)

	d.input.ConnectTextChanged(func(text string) {
		text = strings.Join(strings.Split(text, ""), "*")
		d.proxy.SetFilterWildcard(text)
		d.proxy.SetFilterCaseSensitivity(core.Qt__CaseInsensitive)
		d.selectFirstRow()
	})

	d.ConnectKeyPressEvent(d.handleArrowKeys)
	d.results.ConnectKeyPressEvent(d.handleResultsKeyPress)
}

func (d *Dialog) selectFirstRow() {
	index := d.Model().Index(0, 0, core.NewQModelIndex())
	d.results.SelectionModel().Select(index, core.QItemSelectionModel__ClearAndSelect)
	d.results.SetCurrentIndex(index)
}

func (d *Dialog) view(r *core.QRect) {
	c := core.NewQRect4(r.Width()/2-dialogWidth/2, 0, dialogWidth, 500)
	d.SetGeometry(c)
	d.Show()
	d.input.SetFocus2()
	d.selectFirstRow()
}

// Model returns the results' model.
func (d *Dialog) Model() *core.QAbstractItemModel {
	return d.results.Model()
}

func (d *Dialog) add(r gist.Gist) {
	item := NewListItem(d)
	item.GistID = r.ID
	description := r.Description
	if description == "" {
		for name := range r.Files {
			description = name
			break
		}
	}
	item.Description = description
	d.model.AddGist(item)
}

// ID returns the ID of gist at row.
func (d *Dialog) ID(row int) string {
	index := d.results.Model().Index(row, 0, core.NewQModelIndex())
	return index.Data(gistID).ToString()
}

// Description returns the Description of gist at row.
func (d *Dialog) Description(row int) string {
	index := d.results.Model().Index(row, 0, core.NewQModelIndex())
	return index.Data(description).ToString()
}

func (d *Dialog) handleArrowKeys(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Up, core.Qt__Key_Down:
		d.results.SetFocus2()
	case core.Qt__Key_Enter, core.Qt__Key_Return:
		index := d.results.CurrentIndex()
		d.OpenGist(index.Data(gistID).ToString())
	}
}

func (d *Dialog) handleResultsKeyPress(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Up, core.Qt__Key_Down:
		d.results.KeyPressEventDefault(event)
	case core.Qt__Key_Enter, core.Qt__Key_Return:
		index := d.results.CurrentIndex()
		d.OpenGist(index.Data(gistID).ToString())
	default:
		d.input.SetText(d.input.Text() + event.Text())
		d.input.SetFocus2()
	}
}

// Results returns the result list view.
func (d *Dialog) Results() *widgets.QListView { return d.results }

// HasID returns true if the gistID was found in the model
func (d *Dialog) HasID(gistID string) bool {
	for _, p := range d.model.gists {
		if p.GistID == gistID {
			return true
		}
	}
	return false
}

// Remove removes the gist identified by gistID from the model.
func (d *Dialog) Remove(gistID string) {
	d.model.remove(gistID)
}

// Clear removes all data from model.
func (d *Dialog) Clear() {
	d.model.clear()
}
