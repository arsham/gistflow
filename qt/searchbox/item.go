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

	roles map[int]*core.QByteArray
	gists []*listItem
}

// listItem is one row in the QListView. This is a different gist than
// gist.Gist, this one does not have enough information as it was received by
// asking for user's gist list.
type listItem struct {
	core.QObject

	GistID      string
	Description string
}

func (l *listModel) init() {
	l.roles = map[int]*core.QByteArray{
		gistID:      core.NewQByteArray2("GistID", len("GistID")),
		description: core.NewQByteArray2("Description", len("Description")),
	}

	l.ConnectData(l.data)
	l.ConnectRowCount(func(parent *core.QModelIndex) int {
		return len(l.gists)
	})
	l.ConnectColumnCount(func(parent *core.QModelIndex) int {
		return 1
	})
	l.ConnectRoleNames(func() map[int]*core.QByteArray {
		return l.roles
	})

	l.ConnectAddGist(l.addGist)
}

func (l *listModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(l.gists) {
		return core.NewQVariant()
	}

	var p = l.gists[index.Row()]
	switch role {
	case gistID:
		return core.NewQVariant14(p.GistID)

	case description:
		return core.NewQVariant14(p.Description)

	default:
		return core.NewQVariant()
	}
}

func (l *listModel) addGist(p *listItem) {
	l.BeginInsertRows(core.NewQModelIndex(), len(l.gists), len(l.gists))
	l.gists = append(l.gists, p)
	l.EndInsertRows()
}

// remove removes the gist identified by gistID from the list.
func (l *listModel) remove(gistID string) {
	for row, p := range l.gists {
		if p.GistID == gistID {
			l.BeginRemoveRows(core.NewQModelIndex(), row, row)
			l.gists = append(l.gists[:row], l.gists[row+1:]...)
			l.EndRemoveRows()
			return
		}
	}
}

func (l *listModel) clear() {
	l.BeginResetModel()
	l.ResetInternalData()
	g := make([]*listItem, 0)
	l.gists = g
	l.EndResetModel()
}
