// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

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

	name string // namespace in setting file
	app  *widgets.QApplication

	gistService gist.Service
	logger      boxLogger
	settings    *core.QSettings

	tabWidget  *widgets.QTabWidget
	menubar    *menuBar
	statusbar  *widgets.QStatusBar
	dockWidget *widgets.QDockWidget
	userInput  *widgets.QLineEdit
	gistList   *widgets.QListView
	toolBar    *appToolbar
	sysTray    *widgets.QSystemTrayIcon
	icon       *gui.QIcon

	tabGistList map[string]*Tab // gist id to the tab

	model *GistModel
	proxy *core.QSortFilterProxyModel
}

// Display shows the main window.
func (m *MainWindow) Display() error {
	if m.name == "" {
		m.name = "gisty"
	}

	m.show()
	// TODO: use singleShot
	m.setModel()
	m.settings = getSettings(m.name)
	m.loadSettings()
	m.populate()
	m.setupInteractions()
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

	m.tabWidget = widgets.NewQTabWidget(centralWidget)
	m.tabWidget.SetObjectName("tabWidget")
	m.tabWidget.SetTabsClosable(true)
	m.tabWidget.SetMovable(true)

	tab1 := widgets.NewQWidget(m.tabWidget, core.Qt__Widget)
	tab1.SetObjectName("Untitled")
	m.tabWidget.AddTab(tab1, "Untitled")
	m.tabGistList["untitled"] = nil // there is no gist associated to this tab

	verticalLayout.AddWidget(m.tabWidget, 0, 0)

	m.menubar = NewMenuBar(m)
	m.menubar.SetObjectName("menubar")
	m.menubar.SetGeometry(core.NewQRect4(0, 0, 1043, 30))
	m.SetMenuBar(m.menubar)

	m.statusbar = widgets.NewQStatusBar(m)
	m.statusbar.SetObjectName("statusbar")
	m.SetStatusBar(m.statusbar)

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

	m.gistList = widgets.NewQListView(dockWidgetContents)
	m.gistList.SetObjectName("gistList")

	verticalLayout2.AddWidget(m.userInput, 0, 0)
	verticalLayout2.AddWidget(m.gistList, 0, 0)
	m.dockWidget.SetWidget(dockWidgetContents)

	m.dockWidget.SetWidget(dockWidgetContents)
	m.AddDockWidget(core.Qt__LeftDockWidgetArea, m.dockWidget)

	m.toolBar = NewAppToolbar("Toolbar", m)
	m.AddToolBar(core.Qt__TopToolBarArea, m.toolBar)
	m.toolBar.SetAction(m.menubar.action)

	m.icon = gui.NewQIcon5("./qml/app.ico")
	m.sysTray = widgets.NewQSystemTrayIcon(m)
	m.sysTray.SetIcon(m.icon)
	m.sysTray.SetVisible(true)
	m.sysTray.SetToolTip("Gisty")
	m.sysTray.SetContextMenu(m.menubar.menuOptions)

	m.SetWindowIcon(m.icon)
	m.menubar.ConnectQuit(func() {
		m.app.Quit()
	})
}

// SetGistService sets the service required for public api interactions.
func (m *MainWindow) SetGistService(g gist.Service) {
	m.gistService = g
}

// SetApp is required to be called in order to be able to control the
// application's quit signals.
func (m *MainWindow) SetApp(app *widgets.QApplication) {
	m.app = app
}

func (m *MainWindow) show() {
	m.Show()
	m.userInput.SetFocus2()
}

func (m *MainWindow) setModel() {
	m.model = NewGistModel(nil)

	m.proxy = core.NewQSortFilterProxyModel(nil)
	m.proxy.SetSourceModel(m.model)
	m.proxy.SetFilterCaseSensitivity(core.Qt__CaseInsensitive)

	m.gistList.SetModel(m.proxy)
}

func (m *MainWindow) loadSettings() {
	tmp := widgets.NewQWidget(nil, 0)
	tmp.SetGeometry2(100, 100, 600, 600)
	defSize := tmp.SaveGeometry()
	sizeVar := m.settings.Value(mainWindowGeometry, core.NewQVariant15(defSize))
	m.RestoreGeometry(sizeVar.ToByteArray())
	m.app.ConnectAboutToQuit(m.saveSettings)
}

func (m *MainWindow) saveSettings() {
	current := m.SaveGeometry()
	currentVar := core.NewQVariant15(current.QByteArray_PTR())
	m.settings.SetValue(mainWindowGeometry, currentVar)
	m.settings.Sync()
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
	var (
		content string
		name    string
	)
	if g, ok := m.tabGistList[id]; ok {
		m.tabWidget.SetCurrentWidget(g)
		return nil
	}
	rg, err := m.gistService.Get(id)
	if err != nil {
		return errors.Wrapf(err, "id: %s", id)
	}

	for n, f := range rg.Files {
		content = f.Content
		name = n
		break
	}

	g := &tabGist{
		id:      id,
		content: content,
		label:   name,
		url:     rg.URL,
	}
	tab := NewTab(m.tabWidget)
	tab.showGist(m.tabWidget, g)
	m.tabGistList[id] = tab
	return nil
}

func (m *MainWindow) tabIDFromIndex(index int) string {
	tab := m.tabWidget.Widget(index)
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
	m.tabWidget.RemoveTab(index)
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
