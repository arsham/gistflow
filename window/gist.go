package window

import "github.com/therecipe/qt/core"

var (
	GistID      = int(core.Qt__UserRole)
	Description = int(core.Qt__DisplayRole)
)

type GistModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*Gist                  `property:"gists"`

	_ func(*Gist) `slot:"addGist"`
}

type Gist struct {
	core.QObject

	_ string `property:"gistID"`
	_ string `property:"description"`
}

func init() {
	Gist_QRegisterMetaType()
}

func (m *GistModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		GistID:      core.NewQByteArray2("gistID", len("gistID")),
		Description: core.NewQByteArray2("description", len("description")),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddGist(m.addGist)
}

func (m *GistModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(m.Gists()) {
		return core.NewQVariant()
	}

	var p = m.Gists()[index.Row()]

	switch role {
	case GistID:
		return core.NewQVariant14(p.GistID())

	case Description:
		return core.NewQVariant14(p.Description())

	default:
		return core.NewQVariant()
	}
}

func (m *GistModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Gists())
}

func (m *GistModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *GistModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *GistModel) addGist(p *Gist) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Gists()), len(m.Gists()))
	m.SetGists(append(m.Gists(), p))
	m.EndInsertRows()
}
