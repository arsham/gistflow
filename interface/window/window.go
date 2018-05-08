// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

// Package window shows the main application.
package window

import (
	"os"
	"path"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/gisty/interface/gistlist"
	"github.com/arsham/gisty/interface/menubar"
	"github.com/arsham/gisty/interface/searchbox"
	"github.com/arsham/gisty/interface/tab"
	"github.com/arsham/gisty/interface/toolbar"
	"github.com/pkg/errors"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	mainWindowGeometry = "mainWindowGeometry"
)

type clipboard interface {
	SetText(string, gui.QClipboard__Mode)
}

// MainWindow is the main window of the application.
type MainWindow struct {
	widgets.QMainWindow

	_ func() `constructor:"setupUI"`

	name        string // namespace in setting file
	app         *widgets.QApplication
	settings    *core.QSettings
	logger      boxLogger
	gistService gist.Service

	menubar    *menubar.MenuBar
	toolBar    *toolbar.Toolbar
	sysTray    *widgets.QSystemTrayIcon
	icon       *gui.QIcon
	statusArea *widgets.QStatusBar // named this way to avoid collision

	searchbox  *searchbox.Dialog
	gistList   *gistlist.Container
	dockWidget *widgets.QDockWidget
	tabsWidget *widgets.QTabWidget

	tabGistList map[string]*tab.Tab // gist id to the tab
	clipboard   func() clipboard
}

func init() {
	MainWindow_QRegisterMetaType()
}

// Display shows the main window.
func (m *MainWindow) Display() error {
	if m.name == "" {
		m.name = "gisty"
	}

	m.settings = getSettings(m.name)
	m.loadSettings()
	go m.populate()
	m.Show()
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
		m.tabGistList = make(map[string]*tab.Tab, 0)
	}

	centralWidget := widgets.NewQWidget(m, core.Qt__Widget)
	centralWidget.SetObjectName("centralWidget")
	m.SetCentralWidget(centralWidget)
	verticalLayout := widgets.NewQVBoxLayout2(centralWidget)
	verticalLayout.SetObjectName("verticalLayout")

	m.tabsWidget = widgets.NewQTabWidget(centralWidget)
	m.tabsWidget.SetObjectName("tabWidget")
	m.tabsWidget.SetTabsClosable(true)
	m.tabsWidget.SetMovable(true)

	tab1 := widgets.NewQWidget(m.tabsWidget, core.Qt__Widget)
	tab1.SetObjectName("Untitled")
	m.tabsWidget.AddTab(tab1, "Untitled")
	m.tabGistList["untitled"] = nil // there is no gist associated to this tab

	verticalLayout.AddWidget(m.tabsWidget, 0, 0)

	m.menubar = menubar.NewMenuBar(m)
	m.menubar.SetObjectName("menubar")
	m.menubar.SetGeometry(core.NewQRect4(0, 0, 1043, 30))
	m.SetMenuBar(m.menubar)

	m.statusArea = widgets.NewQStatusBar(m)
	m.statusArea.SetObjectName("statusarea")
	m.SetStatusBar(m.statusArea)

	m.dockWidget = widgets.NewQDockWidget("Gists", m, 0)
	m.dockWidget.SetObjectName("dockWidget")
	m.dockWidget.SetMinimumSize(core.NewQSize2(100, 130))
	m.dockWidget.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetClosable)
	m.dockWidget.SetAllowedAreas(core.Qt__LeftDockWidgetArea | core.Qt__RightDockWidgetArea)

	dockWidgetContents := widgets.NewQWidget(m.dockWidget, core.Qt__Widget)
	dockWidgetContents.SetObjectName("dockWidgetContents")
	verticalLayout2 := widgets.NewQVBoxLayout2(dockWidgetContents)
	verticalLayout2.SetObjectName("verticalLayout2")

	m.gistList = gistlist.NewContainer(dockWidgetContents)
	m.gistList.SetObjectName("gistList")

	verticalLayout2.AddWidget(m.gistList, 0, 0)
	m.dockWidget.SetWidget(dockWidgetContents)

	m.dockWidget.SetWidget(dockWidgetContents)
	m.AddDockWidget(core.Qt__LeftDockWidgetArea, m.dockWidget)

	m.toolBar = toolbar.NewToolbar("Toolbar", m)
	m.AddToolBar(core.Qt__TopToolBarArea, m.toolBar)
	m.toolBar.SetAction(m.menubar.Actions())

	m.icon = gui.NewQIcon5("./qml/app.ico")
	m.sysTray = widgets.NewQSystemTrayIcon(m)
	m.sysTray.SetIcon(m.icon)
	m.sysTray.SetVisible(true)
	m.sysTray.SetToolTip("Gisty")
	m.sysTray.SetContextMenu(m.menubar.Options())

	m.SetWindowIcon(m.icon)
	filter := m.tabMovementEventFilter()

	m.gistList.ConnectDoubleClicked(m.gistListDoubleClickEvent)
	m.gistList.ConnectKeyReleaseEvent(m.openSelectedGist)
	m.gistList.InstallEventFilter(filter)

	m.tabsWidget.InstallEventFilter(filter)
	m.tabsWidget.ConnectTabCloseRequested(m.closeTab)

	m.menubar.ConnectCopyURLToClipboard(m.copyURLToClipboard)
	m.menubar.ConnectOpenInBrowser(m.openInBrowser)
	m.menubar.ConnectQuit(func() {
		m.app.Quit()
	})

	m.sysTray.ConnectActivated(m.sysTrayClick)
	m.clipboard = func() clipboard {
		return m.app.Clipboard()
	}

	m.searchbox = searchbox.NewDialog(m, 0)
	m.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_P) && event.Modifiers() == core.Qt__ControlModifier {
			m.searchbox.View(m.Geometry())
		}
	})
	m.searchbox.ConnectOpenGist(m.openGistByID)
}

// GistList returns the associated gistList.
func (m *MainWindow) GistList() *gistlist.Container { return m.gistList }

// TabsWidget returns the associated tabsWidget.
func (m *MainWindow) TabsWidget() *widgets.QTabWidget { return m.tabsWidget }

// SetGistService sets the service required for public api interactions.
func (m *MainWindow) SetGistService(g gist.Service) { m.gistService = g }

// SetApp sets the app instance.
func (m *MainWindow) SetApp(app *widgets.QApplication) { m.app = app }

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
	for item := range m.gistService.Iter() {
		foundOne = true
		m.searchbox.Add(item)
		m.gistList.Add(item)
	}
	if !foundOne {
		m.logger.Error("didn't find any gists")
	}
}

func (m *MainWindow) openGist(id string) error {
	if g, ok := m.tabGistList[id]; ok {
		m.tabsWidget.SetCurrentWidget(g)
		return nil
	}
	rg, err := m.gistService.Get(id)
	if err != nil {
		return errors.Wrapf(err, "id: %s", id)
	}

	tab := tab.NewTab(m.tabsWidget)
	tab.ShowGist(m.tabsWidget, &rg)
	m.tabGistList[id] = tab
	tab.ConnectCopyToClipboard(func(text string) {
		m.clipboard().SetText(text, gui.QClipboard__Clipboard)
	})
	return nil
}

func (m *MainWindow) tabIDFromIndex(index int) string {
	tab := m.tabsWidget.Widget(index)
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
	m.tabsWidget.RemoveTab(index)
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
