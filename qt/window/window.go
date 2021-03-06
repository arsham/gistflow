// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

// Package window shows the main application.
package window

import (
	"fmt"
	"os"
	"path"

	"github.com/arsham/gistflow/gist"
	"github.com/arsham/gistflow/qt/conf"
	"github.com/arsham/gistflow/qt/gistlist"
	"github.com/arsham/gistflow/qt/menubar"
	"github.com/arsham/gistflow/qt/messagebox"
	"github.com/arsham/gistflow/qt/searchbox"
	"github.com/arsham/gistflow/qt/tab"
	"github.com/arsham/gistflow/qt/toolbar"
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
	settings    *conf.Settings
	logger      messagebox.Message
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

func (m *MainWindow) setupUI() {
	if m.ObjectName() == "" {
		m.SetObjectName("gistflow")
	}
	if m.logger == nil {
		m.logger = messagebox.New(m)
	}
	if m.tabGistList == nil {
		m.tabGistList = make(map[string]*tab.Tab, 0)
	}

	centralWidget := widgets.NewQWidget(m, core.Qt__Widget)
	centralWidget.SetObjectName("centralWidget")
	m.SetCentralWidget(centralWidget)

	m.tabsWidget = widgets.NewQTabWidget(centralWidget)
	m.tabsWidget.SetObjectName("tabWidget")
	m.tabsWidget.SetTabsClosable(true)
	m.tabsWidget.SetMovable(true)

	verticalLayout := widgets.NewQVBoxLayout2(centralWidget)
	verticalLayout.SetObjectName("verticalLayout")
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

	m.icon = gui.NewQIcon5("./qml/v1/app.ico")
	m.sysTray = widgets.NewQSystemTrayIcon(m)
	m.sysTray.SetIcon(m.icon)
	m.sysTray.SetVisible(true)
	m.sysTray.SetToolTip("GistFlow")
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
	m.menubar.ConnectOpenSettings(func(bool) {
		m.showSettings(func() {
			m.gistService.Username = m.settings.Username
			m.gistService.Token = m.settings.Token
			m.gistList.Clear()
			m.searchbox.Clear()
			go m.populate()
		})
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
	m.menubar.ConnectNewGist(m.newGist)

	if m.gistService.Logger == nil {
		m.gistService.Logger = m.logger
	}
	if m.gistService.CacheDir == "" {
		m.gistService.CacheDir = m.cacheDir()
	}

	m.menubar.ConnectToggleToolbar(func(active bool) {
		switch active {
		case true:
			m.toolBar.Show()
		case false:
			m.toolBar.Hide()
		}
	})
	m.menubar.ConnectToggleGistList(func(active bool) {
		switch active {
		case true:
			m.dockWidget.Show()
		case false:
			m.dockWidget.Close()
		}
	})
}

// Display shows the main window.
func (m *MainWindow) Display(app *widgets.QApplication) error {
	var err error
	if m.name == "" {
		m.name = "gistflow"
	}
	m.app = app
	m.setStyleSheet(":/qml/stylesheet.qss")
	m.Show()

	populate := func() {
		m.lastGeometry()
		m.app.ConnectAboutToQuit(m.recordGeometry)
		m.gistService.Username = m.settings.Username
		m.gistService.Token = m.settings.Token
		go m.populate()
	}
	m.settings, err = conf.New(m.name)
	if err != nil {
		m.showSettings(populate)
		return nil
	}

	populate()
	return nil
}

func (m *MainWindow) setStyleSheet(name string) {
	file := core.NewQFile2(name)
	if ok := file.Open(core.QIODevice__ReadOnly); ok {
		defer file.Close()
		contents := file.ReadAll()
		sheet := core.NewQLatin1String5(contents)
		m.SetStyleSheet(sheet.Latin1())
	}
}

// showSettings calls the `callback` after the settings tab is closed.
func (m *MainWindow) showSettings(callback func()) {
	t := conf.NewTab(m.tabsWidget)
	t.SetSettings(m.settings)
	m.tabsWidget.AddTab(t, "Settings")
	m.tabsWidget.SetCurrentWidget(t)
	tabIndex := m.tabsWidget.IndexOf(t)
	m.tabsWidget.ConnectTabCloseRequested(func(index int) {
		if index == tabIndex {
			callback()
		}
	})
}

func (m *MainWindow) lastGeometry() {
	tmp := widgets.NewQWidget(nil, 0)
	tmp.SetGeometry2(100, 100, 600, 600)
	defSize := tmp.SaveGeometry()
	sizeVar := m.settings.Value(mainWindowGeometry, core.NewQVariant15(defSize))
	m.RestoreGeometry(sizeVar.ToByteArray())
}

func (m *MainWindow) recordGeometry() {
	current := m.SaveGeometry()
	currentVar := core.NewQVariant15(current.QByteArray_PTR())
	m.settings.SetValue(mainWindowGeometry, currentVar)
	m.settings.Sync()
}

func (m *MainWindow) populate() {
	var foundOne bool
	for item := range m.gistService.Iter() {
		foundOne = true
		m.searchbox.Add(item)
		m.gistList.Add(item)
	}
	if !foundOne {
		m.logger.Error("didn't find any gists")
	}
}

func (m *MainWindow) newGist(bool) {
	id := nextUntitled(m.tabGistList)
	t := tab.NewTab(m.tabsWidget)
	t.NewGist(m.tabsWidget, id)
	m.tabGistList[id] = t

	t.ConnectCopyToClipboard(func(text string) {
		m.clipboard().SetText(text, gui.QClipboard__Clipboard)
	})
	t.ConnectCreateGist(func(g *gist.Gist) {
		newGist, err := m.gistService.Create(*g)
		if err != nil {
			msg := fmt.Sprintf("Could not create new gist: %s", err)
			m.logger.Error(msg)
			return
		}
		m.showNotification("New gist has been created")
		t.GistCreated(&newGist)
		m.searchbox.Add(newGist)
		m.gistList.Add(newGist)
	})
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

	t := tab.NewTab(m.tabsWidget)
	t.ShowGist(m.tabsWidget, &rg)
	m.tabGistList[id] = t

	t.ConnectCopyToClipboard(func(text string) {
		m.clipboard().SetText(text, gui.QClipboard__Clipboard)
	})

	t.ConnectUpdateGist(func(g *gist.Gist) {
		_, err := m.gistService.Update(*g)
		if err != nil {
			msg := fmt.Sprintf("Could not update the gist: %s", err)
			m.logger.Error(msg)
			return
		}
		m.showNotification("Gist has been updated")
	})

	t.ConnectDeleteFile(func(g *gist.Gist, name string) {
		_, err := m.gistService.DeleteFile(*g, name)
		if err != nil {
			msg := fmt.Sprintf("Could not delete file: %s", err)
			m.logger.Error(msg)
			return
		}
		m.showNotification("File was removed from your gist")
		t.FileDeleted(name)
	})

	t.ConnectDeleteGist(func(g *gist.Gist) {
		err := m.gistService.DeleteGist(g.ID)
		if err != nil {
			msg := fmt.Sprintf("Could not delete gist: %s", err)
			m.logger.Error(msg)
			return
		}
		m.searchbox.Remove(g.ID)
		m.gistList.Remove(g.ID)
		tab := m.tabGistList[g.ID]
		delete(m.tabGistList, g.ID)
		m.showNotification("Gist has been removed")
		tab.DestroyQWidget()
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

func nextUntitled(tabGistList map[string]*tab.Tab) string {
	id := "untitled"
	if _, ok := tabGistList[id]; !ok {
		return id
	}
	for i := 1; i < 100000; i++ {
		id := fmt.Sprintf("untitled-%d", i)
		if _, ok := tabGistList[id]; !ok {
			return id
		}
	}
	return id
}
