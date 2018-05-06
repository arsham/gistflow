// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

// Package window shows the main application.
package window

import (
	"os"
	"path"

	"github.com/arsham/gisty/gist"
	"github.com/pkg/errors"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	mainWindowGeometry = "mainWindowGeometry"
)

// MainWindow is the main window of the application.
type MainWindow struct {
	widgets.QMainWindow

	_ func() `constructor:"setupUI"`
	_ func() `slot:"userInputChange"`
	_ func() `slot:"userInputTextChange"`
	_ func() `slot:"gistListDoubleClickEvent"`
	_ func() `slot:"gistListKeyReleaseEvent"`
	_ func() `slot:"openSelectedGist"`
	_ func() `slot:"sysTrayClick"`
	_ func() `slot:"tabWidgetKeyPressEvent"`
	_ func() `slot:"closeTab"`
	_ func() `slot:"copyToClipboard"`
	_ func() `slot:"copyURLToClipboard"`
	_ func() `slot:"openInBrowser"`

	_ *widgets.QApplication `property:"app"`
	_ *core.QSettings       `property:"settings"`
	_ *widgets.QTabWidget   `property:"tabsWidget"`
	_ *widgets.QStatusBar   `property:"statusArea"` // named this way to avoid collision
	_ *widgets.QListView    `property:"gistList"`

	name string // namespace in setting file

	gistService gist.Service
	logger      boxLogger

	menubar    *menuBar
	dockWidget *widgets.QDockWidget
	userInput  *widgets.QLineEdit
	toolBar    *appToolbar
	sysTray    *widgets.QSystemTrayIcon
	icon       *gui.QIcon

	tabGistList map[string]*Tab // gist id to the tab

	model *listGistModel
	proxy *core.QSortFilterProxyModel
}

func init() {
	MainWindow_QRegisterMetaType()
}

// Display shows the main window.
func (m *MainWindow) Display() error {
	if m.name == "" {
		m.name = "gisty"
	}

	m.show()
	// TODO: use singleShot
	m.setModel()
	m.SetSettings(getSettings(m.name))
	m.loadSettings()
	m.populate()
	widgets.QApplication_Exec()
	return nil
}

func (m *MainWindow) setupUI() {
	if m.ObjectName() == "" {
		m.SetObjectName("gisty")
	}
	if m.logger == nil {
		m.logger = messagebox(m)
	}
	if m.tabGistList == nil {
		m.tabGistList = make(map[string]*Tab, 0)
	}

	centralWidget := widgets.NewQWidget(m, core.Qt__Widget)
	centralWidget.SetObjectName("centralWidget")
	m.SetCentralWidget(centralWidget)
	verticalLayout := widgets.NewQVBoxLayout2(centralWidget)
	verticalLayout.SetObjectName("verticalLayout")

	m.SetTabsWidget(widgets.NewQTabWidget(centralWidget))
	m.TabsWidget().SetObjectName("tabWidget")
	m.TabsWidget().SetTabsClosable(true)
	m.TabsWidget().SetMovable(true)

	tab1 := widgets.NewQWidget(m.TabsWidget(), core.Qt__Widget)
	tab1.SetObjectName("Untitled")
	m.TabsWidget().AddTab(tab1, "Untitled")
	m.tabGistList["untitled"] = nil // there is no gist associated to this tab

	verticalLayout.AddWidget(m.TabsWidget(), 0, 0)

	m.menubar = NewMenuBar(m)
	m.menubar.SetObjectName("menubar")
	m.menubar.SetGeometry(core.NewQRect4(0, 0, 1043, 30))
	m.SetMenuBar(m.menubar)

	m.SetStatusArea(widgets.NewQStatusBar(m))
	m.StatusArea().SetObjectName("statusarea")
	m.SetStatusBar(m.StatusArea())

	m.dockWidget = widgets.NewQDockWidget("Gists", m, 0)
	m.dockWidget.SetObjectName("dockWidget")
	m.dockWidget.SetMinimumSize(core.NewQSize2(100, 130))
	m.dockWidget.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetClosable)
	m.dockWidget.SetAllowedAreas(core.Qt__LeftDockWidgetArea | core.Qt__RightDockWidgetArea)

	dockWidgetContents := widgets.NewQWidget(m.dockWidget, core.Qt__Widget)
	dockWidgetContents.SetObjectName("dockWidgetContents")
	verticalLayout2 := widgets.NewQVBoxLayout2(dockWidgetContents)
	verticalLayout2.SetObjectName("verticalLayout2")

	m.userInput = widgets.NewQLineEdit(dockWidgetContents)
	m.userInput.SetObjectName("userInput")
	m.userInput.SetClearButtonEnabled(true)

	m.SetGistList(widgets.NewQListView(dockWidgetContents))
	m.GistList().SetObjectName("gistList")

	verticalLayout2.AddWidget(m.userInput, 0, 0)
	verticalLayout2.AddWidget(m.GistList(), 0, 0)
	m.dockWidget.SetWidget(dockWidgetContents)

	m.dockWidget.SetWidget(dockWidgetContents)
	m.AddDockWidget(core.Qt__LeftDockWidgetArea, m.dockWidget)

	m.toolBar = NewAppToolbar("Toolbar", m)
	m.AddToolBar(core.Qt__TopToolBarArea, m.toolBar)
	m.toolBar.SetAction(m.menubar.Actions())

	m.icon = gui.NewQIcon5("./qml/app.ico")
	m.sysTray = widgets.NewQSystemTrayIcon(m)
	m.sysTray.SetIcon(m.icon)
	m.sysTray.SetVisible(true)
	m.sysTray.SetToolTip("Gisty")
	m.sysTray.SetContextMenu(m.menubar.Options())

	m.SetWindowIcon(m.icon)

	m.GistList().ConnectKeyReleaseEvent(m.gistListKeyReleaseEvent)
	m.GistList().ConnectDoubleClicked(m.gistListDoubleClickEvent)
	m.userInput.ConnectKeyPressEvent(m.userInputChange)
	m.sysTray.ConnectActivated(m.sysTrayClick)
	m.TabsWidget().ConnectKeyPressEvent(m.tabWidgetKeyPressEvent)
	m.TabsWidget().ConnectTabCloseRequested(m.closeTab)
	m.userInput.ConnectTextChanged(m.userInputTextChange)
	m.menubar.ConnectCopyToClipboard(m.copyToClipboard)
	m.menubar.ConnectCopyURLToClipboard(m.copyURLToClipboard)
	m.menubar.ConnectOpenInBrowser(m.openInBrowser)
	m.GistList().ConnectKeyReleaseEvent(m.openSelectedGist)
	m.userInput.ConnectKeyReleaseEvent(m.openSelectedGist)

	m.menubar.ConnectQuit(func() {
		m.App().Quit()
	})
}

// SetGistService sets the service required for public api interactions.
func (m *MainWindow) SetGistService(g gist.Service) {
	m.gistService = g
}

func (m *MainWindow) show() {
	m.Show()
	m.userInput.SetFocus2()
}

func (m *MainWindow) setModel() {
	m.model = NewListGistModel(nil)

	m.proxy = core.NewQSortFilterProxyModel(nil)
	m.proxy.SetSourceModel(m.model)
	m.proxy.SetFilterCaseSensitivity(core.Qt__CaseInsensitive)

	m.GistList().SetModel(m.proxy)
}

func (m *MainWindow) loadSettings() {
	tmp := widgets.NewQWidget(nil, 0)
	tmp.SetGeometry2(100, 100, 600, 600)
	defSize := tmp.SaveGeometry()
	sizeVar := m.Settings().Value(mainWindowGeometry, core.NewQVariant15(defSize))
	m.RestoreGeometry(sizeVar.ToByteArray())
	m.App().ConnectAboutToQuit(m.saveSettings)
}

func (m *MainWindow) saveSettings() {
	current := m.SaveGeometry()
	currentVar := core.NewQVariant15(current.QByteArray_PTR())
	m.Settings().SetValue(mainWindowGeometry, currentVar)
	m.Settings().Sync()
}

func getSettings(name string) *core.QSettings {
	return core.NewQSettings3(
		core.QSettings__NativeFormat,
		core.QSettings__UserScope,
		"gisty",
		name,
		nil,
	)
}

func (m *MainWindow) populate() {
	var foundOne bool
	if m.gistService.Logger == nil {
		m.gistService.Logger = m.logger
	}
	if m.gistService.CacheDir == "" {
		m.gistService.CacheDir = m.cacheDir()
	}
	// TODO: populate in background.
	for item := range m.gistService.Iter() {
		foundOne = true
		var g = NewGistItem(nil)
		g.SetGistID(item.ID)
		g.SetGistURL(item.URL)
		g.SetDescription(item.Description)
		m.model.AddGist(g)
	}
	if !foundOne {
		m.logger.Error("didn't find any gists")
	}
}

func (m *MainWindow) openGist(id string) error {
	if g, ok := m.tabGistList[id]; ok {
		m.TabsWidget().SetCurrentWidget(g)
		return nil
	}
	rg, err := m.gistService.Get(id)
	if err != nil {
		return errors.Wrapf(err, "id: %s", id)
	}

	tab := NewTab(m.TabsWidget())
	tab.showGist(m.TabsWidget(), &rg)
	m.tabGistList[id] = tab
	return nil
}

func (m *MainWindow) tabIDFromIndex(index int) string {
	tab := m.TabsWidget().Widget(index)
	if tab.Pointer() == nil {
		return ""
	}

	for name, o := range m.tabGistList {
		if o.Pointer() == tab.Pointer() {
			return name
		}
	}
	return ""
}

func (m *MainWindow) closeTab(index int) {
	id := m.tabIDFromIndex(index)
	m.TabsWidget().RemoveTab(index)
	delete(m.tabGistList, id)
}

func (m *MainWindow) cacheDir() string {
	loc := core.QStandardPaths_StandardLocations(core.QStandardPaths__GenericCacheLocation)[0]
	cacheDir := path.Join(loc, m.name)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.Mkdir(cacheDir, 0740); err != nil {
			m.logger.Warningf("Creating cache dir: %s", err)
		}
	}
	return cacheDir
}
