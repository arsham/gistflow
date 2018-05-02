// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

// Package window shows all kinds of windows and dialogs.
package window

import (
	"os"
	"strings"
	"unicode"

	"github.com/arsham/gisty/gist"
	"github.com/pkg/errors"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// https://github.com/therecipe/advanced-examples/tree/master/test

const (
	mainWindowGeometry = "mainWindowGeometry"
)

// MainWindow is in charge of user interaction with the dialog.
type MainWindow struct {
	GistService gist.Service
	logger      boxLogger
	ConfName    string // namespace in setting file
	settings    *core.QSettings

	app       *widgets.QApplication
	window    *widgets.QMainWindow
	sysTray   *widgets.QSystemTrayIcon
	menubar   *menuBar
	statusbar *widgets.QStatusBar
	icon      *gui.QIcon

	userInput   *widgets.QLineEdit
	tabWidget   *widgets.QTabWidget
	tabGistList map[string]*Tab // gist id to the tab
	dockWidget  *widgets.QDockWidget
	listView    *widgets.QListView

	model *GistModel
	proxy *core.QSortFilterProxyModel
}

// Display shows the main window.
func (m *MainWindow) Display() error {
	core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	m.app = widgets.NewQApplication(len(os.Args), os.Args)
	if m.ConfName == "" {
		m.ConfName = "gisty"
	}
	err := m.setupUI()
	if err != nil {
		return err
	}
	m.window.Show()
	// use singleShot
	m.setModel()
	m.settings = getSettings(m.ConfName)
	m.loadSettings()
	m.populate()
	m.setupInteractions()
	widgets.QApplication_Exec()
	return nil
}

func (m *MainWindow) setupUI() (err error) {
	if m.logger == nil {
		m.logger = messagebox(m.window)
	}
	m.window = widgets.NewQMainWindow(nil, 0)
	m.window.SetGeometry(core.NewQRect4(0, 0, 1043, 600))
	centralWidget := widgets.NewQWidget(m.window, core.Qt__Widget)
	vLayout := widgets.NewQVBoxLayout2(centralWidget)
	vLayout.SetObjectName("verticalLayout")

	if m.tabGistList == nil {
		m.tabGistList = make(map[string]*Tab, 0)
	}
	m.tabWidget = widgets.NewQTabWidget(centralWidget)
	m.tabWidget.SetObjectName("tabWidget")
	m.tabWidget.SetTabsClosable(true)
	tab1 := widgets.NewQWidget(m.tabWidget, core.Qt__Widget)
	tab1.SetObjectName("Untitled")
	m.tabWidget.AddTab(tab1, "Untitled")
	m.tabGistList["untitled"] = nil // there is no gist associated to this tab
	m.userInput = widgets.NewQLineEdit(m.window)
	vLayout.AddWidget(m.userInput, 0, 0)
	vLayout.AddWidget(m.tabWidget, 0, 0)
	m.window.SetCentralWidget(centralWidget)

	m.tabWidget.ConnectTabCloseRequested(m.closeTab)

	m.menubar = NewMenuBar(m.window)
	m.menubar.SetObjectName("Menubar")
	m.menubar.SetGeometry(core.NewQRect4(0, 0, 1043, 30))
	m.window.SetMenuBar(m.menubar)

	m.statusbar = widgets.NewQStatusBar(m.window)
	m.statusbar.SetObjectName("Statusbar")
	m.window.SetStatusBar(m.statusbar)

	m.dockWidget = widgets.NewQDockWidget("Gists", m.window, 0)
	m.dockWidget.SetObjectName("DockWidget")
	m.dockWidget.SetMinimumSize(core.NewQSize2(100, 130))
	m.dockWidget.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetClosable)
	m.dockWidget.SetAllowedAreas(core.Qt__LeftDockWidgetArea | core.Qt__RightDockWidgetArea)

	widgetContent := widgets.NewQWidget(m.dockWidget, core.Qt__Widget)
	widgetContent.SetObjectName("DockWidgetContents")
	vLayout2 := widgets.NewQVBoxLayout2(widgetContent)
	vLayout2.SetObjectName("verticalLayout_2")
	vLayout2.SetContentsMargins(0, 0, 0, 0)
	vLayout2.SetSpacing(0)
	m.listView = widgets.NewQListView(widgetContent)
	m.listView.SetObjectName("ListView")
	vLayout2.AddWidget(m.listView, 0, 0)
	m.dockWidget.SetWidget(m.listView)

	m.tabWidget.SetCurrentIndex(0)
	m.window.AddDockWidget(core.Qt__LeftDockWidgetArea, m.dockWidget)

	m.icon = gui.NewQIcon5("./qml/app.ico")
	m.sysTray = widgets.NewQSystemTrayIcon(m.window)
	m.sysTray.SetIcon(m.icon)
	m.sysTray.SetVisible(true)
	m.sysTray.SetToolTip("Gisty")
	m.sysTray.SetContextMenu(m.menubar.optionsMenu)

	m.window.SetWindowIcon(m.icon)
	// m.window.SetTabOrder(m.userInput, m.listView)

	m.menubar.quitAction.ConnectTriggered(func(bool) {
		m.app.Quit()
	})
	return nil
}

func (m *MainWindow) setModel() {
	m.model = NewGistModel(nil)

	m.proxy = core.NewQSortFilterProxyModel(nil)
	m.proxy.SetSourceModel(m.model)
	m.proxy.SetFilterCaseSensitivity(core.Qt__CaseInsensitive)

	m.listView.SetModel(m.proxy)
}

func (m *MainWindow) loadSettings() {
	tmp := widgets.NewQWidget(nil, 0)
	tmp.SetGeometry2(100, 100, 600, 600)
	defSize := tmp.SaveGeometry()
	sizeVar := m.settings.Value(mainWindowGeometry, core.NewQVariant15(defSize))
	m.window.RestoreGeometry(sizeVar.ToByteArray())
	m.app.ConnectAboutToQuit(m.saveSettings)
}

func (m *MainWindow) saveSettings() {
	current := m.window.SaveGeometry()
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
	// go goroutine
	for item := range m.GistService.Iter() {
		foundOne = true
		var g = NewGistItem(nil)
		g.SetGistID(item.ID)
		g.SetGistURL(item.URL)
		g.SetDescription(item.Description)
		m.model.AddGist(g)
	}
	if !foundOne {
		m.logger.error("didn't find any gists")
	}
}

func (m *MainWindow) setupInteractions() {
	m.userInput.ConnectTextChanged(func(text string) {
		newText := strings.Split(text, "")
		m.proxy.SetFilterWildcard(strings.Join(newText, "*"))
	})
	m.userInput.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Up) || event.Key() == int(core.Qt__Key_Down) {
			m.listView.SetFocus2()
		}
		m.userInput.KeyPressEventDefault(event)
	})
	// listView.ConnectDoubleClicked(openGist)

	m.listView.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		switch core.Qt__Key(event.Key()) {
		case core.Qt__Key_Enter, core.Qt__Key_Return:
			index := m.listView.CurrentIndex()
			err := m.openGist(index.Data(GistID).ToString())
			if err != nil {
				m.logger.error(err.Error())
			}
		case core.Qt__Key_Delete, core.Qt__Key_Space:
			fallthrough
		case core.Qt__Key_Left, core.Qt__Key_Right:
			m.userInput.SetFocus2()
		default:
			char := event.Text()
			for _, c := range char {
				if unicode.IsPrint(c) {
					// m.userInput.SetText(m.userInput.Text() + char)
					m.userInput.SetFocus2()
				}
				break
			}
		}
	})
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
	rg, err := m.GistService.Get(id)
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

func (m *MainWindow) gistDialog(index *core.QModelIndex) error {
	// widgets.NewQDialogFromPointer(dialog.Pointer()).SetModal(true)
	// ok := ui.ok
	// ok.ConnectClicked(func(bool) {
	// 	dialog.Close()
	// })
	// clipboard := ui.clipboard
	// clipboard.ConnectClicked(func(bool) {
	// 	m.app.Clipboard().SetText(content, gui.QClipboard__Clipboard)
	// 	m.sysTray.ShowMessage("Info", "Gist has been copied to clipboard", widgets.QSystemTrayIcon__Information, 4000)
	// })
	// browser := ui.browser
	// browser.ConnectClicked(func(bool) {
	// 	gui.QDesktopServices_OpenUrl(core.NewQUrl3(url, 0))
	// })

	// dialog.Show()
	// dialog.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
	// 	if event.Key() == int(core.Qt__Key_Escape) {
	// 		dialog.Close()
	// 	}
	// })
	return nil
}
