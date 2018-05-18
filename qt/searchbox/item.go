// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package searchbox

import (
	"github.com/therecipe/qt/core"
)

// listModel is used in SearchBox's Results.
type listModel struct {
	core.QAbstractListModel

	_ func()          `constructor:"init"`
	_ func(*listItem) `slot:"addGist"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*listItem              `property:"gists"`
}

// listItem is one row in the QListView. This is a different gist than
// gist.Gist, this one does not have enough information as it was received by
// asking for user's gist list.
type listItem struct {
	core.QObject

	GistID      string
	GistURL     string
	Description string
}

func (l *listModel) init() {
	l.SetRoles(map[int]*core.QByteArray{
		gistID:      core.NewQByteArray2("GistID", len("GistID")),
		gistURL:     core.NewQByteArray2("GistURL", len("GistURL")),
		description: core.NewQByteArray2("Description", len("Description")),
	})

	l.ConnectData(l.data)
	l.ConnectRowCount(l.rowCount)
	l.ConnectColumnCount(l.columnCount)
	l.ConnectRoleNames(l.roleNames)

	l.ConnectAddGist(l.addGist)
}

func (l *listModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(l.Gists()) {
		return core.NewQVariant()
	}

	var p = l.Gists()[index.Row()]
	switch role {
	case gistID:
		return core.NewQVariant14(p.GistID)

	case gistURL:
		return core.NewQVariant14(p.GistURL)

	case description:
		return core.NewQVariant14(p.Description)

	default:
		return core.NewQVariant()
	}
}

func (l *listModel) rowCount(parent *core.QModelIndex) int {
	return len(l.Gists())
}

func (l *listModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (l *listModel) roleNames() map[int]*core.QByteArray {
	return l.Roles()
}

func (l *listModel) addGist(p *listItem) {
	l.BeginInsertRows(core.NewQModelIndex(), len(l.Gists()), len(l.Gists()))
	l.SetGists(append(l.Gists(), p))
	l.EndInsertRows()
}

// remove removes the gist identified by gistID from the list.
func (l *listModel) remove(gistID string) {
	for row, p := range l.Gists() {
		if p.GistID == gistID {
			l.BeginRemoveRows(core.NewQModelIndex(), row, row)
			l.SetGists(append(l.Gists()[:row], l.Gists()[row+1:]...))
			l.EndRemoveRows()
			return
		}
	}
}

func (l *listModel) clear() {
	l.BeginResetModel()
	l.ResetInternalData()
	g := make([]*listItem, 0)
	l.SetGists(g)
	l.EndResetModel()
}
