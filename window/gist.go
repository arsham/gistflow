// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const (
	description = int(core.Qt__DisplayRole)
	gistID      = int(core.Qt__UserRole) + 1<<iota
	gistURL
)

type listGistModel struct {
	core.QAbstractListModel

	_ func()              `constructor:"init"`
	_ func(*listGistItem) `slot:"addGist"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*listGistItem          `property:"gists"`
}

// listGistItem is one row in the QListView. This is a different gist than
// gist.Gist, this one does not have enough information as it was received by
// asking for user's gist list.
type listGistItem struct {
	core.QObject

	_ string `property:"gistID"`
	_ string `property:"gistURL"`
	_ string `property:"description"`
}

func init() {
	listGistItem_QRegisterMetaType()
}

func (m *listGistModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		gistID:      core.NewQByteArray2("gistID", len("gistID")),
		gistURL:     core.NewQByteArray2("gistURL", len("gistURL")),
		description: core.NewQByteArray2("description", len("description")),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddGist(m.addGist)
}

func (m *listGistModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(m.Gists()) {
		return core.NewQVariant()
	}

	var p = m.Gists()[index.Row()]
	switch role {
	case gistID:
		return core.NewQVariant14(p.GistID())

	case gistURL:
		return core.NewQVariant14(p.GistURL())

	case description:
		return core.NewQVariant14(p.Description())

	default:
		return core.NewQVariant()
	}
}

func (m *listGistModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Gists())
}

func (m *listGistModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *listGistModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *listGistModel) addGist(p *listGistItem) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Gists()), len(m.Gists()))
	m.SetGists(append(m.Gists(), p))
	m.EndInsertRows()
}

// file represents one file in a gist.
type file struct {
	widgets.QWidget

	_ func()       `constructor:"init"`
	_ func(string) `signal:"copyToClipboard"`

	_ *widgets.QLabel      `property:"information"`
	_ *widgets.QTextEdit   `property:"content"`
	_ *widgets.QPushButton `property:"copy"`
}

func (f *file) init() {
	f.SetObjectName("File")
	vLayout := widgets.NewQVBoxLayout2(f)
	hLayout := widgets.NewQHBoxLayout()
	f.SetInformation(widgets.NewQLabel(f, core.Qt__Widget))
	hLayout.AddWidget(f.Information(), 0, 0)
	hSpacer := widgets.NewQSpacerItem(40, 20, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Minimum)
	hLayout.AddItem(hSpacer)
	f.SetCopy(widgets.NewQPushButton(f))
	f.Copy().SetText("Copy")
	hLayout.AddWidget(f.Copy(), 0, 0)
	vLayout.AddLayout(hLayout, 0)
	f.SetContent(widgets.NewQTextEdit(f))
	f.Content().SetObjectName("content")
	vLayout.AddWidget(f.Content(), 0, 0)

	f.Copy().ConnectClicked(func(bool) {
		f.CopyToClipboard(f.Content().ToPlainText())
	})
}
