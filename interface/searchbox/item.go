// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package searchbox

import (
	"github.com/therecipe/qt/core"
)

// ListModel is used in SearchBox's Results.
type ListModel struct {
	core.QAbstractListModel

	_ func()          `constructor:"init"`
	_ func(*ListItem) `slot:"addGist"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*ListItem              `property:"gists"`
}

// ListItem is one row in the QListView. This is a different gist than
// gist.Gist, this one does not have enough information as it was received by
// asking for user's gist list.
type ListItem struct {
	core.QObject

	_ string `property:"GistID"`
	_ string `property:"GistURL"`
	_ string `property:"Description"`
}

func init() {
	ListItem_QRegisterMetaType()
}

func (m *ListModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		gistID:      core.NewQByteArray2("GistID", len("GistID")),
		gistURL:     core.NewQByteArray2("GistURL", len("GistURL")),
		description: core.NewQByteArray2("Description", len("Description")),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddGist(m.addGist)
}

func (m *ListModel) data(index *core.QModelIndex, role int) *core.QVariant {
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

func (m *ListModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Gists())
}

func (m *ListModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *ListModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *ListModel) addGist(p *ListItem) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Gists()), len(m.Gists()))
	m.SetGists(append(m.Gists(), p))
	m.EndInsertRows()
}
