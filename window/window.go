// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

// Package window shows all kinds of windows and dialogs.
package window

import (
	"os"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/qtlib"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// https://github.com/therecipe/advanced-examples/tree/master/test

// Service is in charge of user interaction with the dialog.
type Service struct {
	GistService gist.Service
	app         *widgets.QApplication
	window      *widgets.QMainWindow
	dialog      *widgets.QWidget
	layout      *widgets.QVBoxLayout
	sysTray     *widgets.QSystemTrayIcon
	listView    *widgets.QListView
	userInput   *widgets.QLineEdit
}

// MainWindow shows the main window. prefix is the path prefix.
func (s *Service) MainWindow() error {
	err := s.displayMainWindow()
	if err != nil {
		return err
	}
	widgets.QApplication_Exec()
	return nil
}

func (s *Service) displayMainWindow() (err error) {
	core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	s.app = widgets.NewQApplication(len(os.Args), os.Args)
	s.window = widgets.NewQMainWindow(nil, 0)
	s.dialog, err = s.mainDialog()
	if err != nil {
		return err
	}

	icon := gui.NewQIcon5("./qml/app.ico")
	s.sysTray = widgets.NewQSystemTrayIcon(s.dialog)
	s.sysTray.SetIcon(icon)
	s.sysTray.SetVisible(true)
	s.sysTray.SetToolTip("Gisty")

	s.window.SetupUi(s.dialog)
	s.window.SetWindowIcon(icon)
	s.dialog.Show()

	model := NewGistModel(nil)
	s.setupUI(model)
	s.loadSettings()
	go s.populate(model)
	return nil
}

func (s *Service) setupUI(model *GistModel) {
	proxy := core.NewQSortFilterProxyModel(nil)
	proxy.SetFilterCaseSensitivity(core.Qt__CaseInsensitive)
	proxy.SetSourceModel(model)

	s.listView = s.listViewWidget()
	s.listView.SetModel(proxy)
	s.userInput = s.userInputWidget(proxy)

	s.layout.AddWidget(s.userInput, 0, 0)
	s.layout.AddWidget(s.listView, 0, 0)

	mainMenu := s.mainMenu()
	s.sysTray.SetContextMenu(mainMenu)
	s.sysTray.ConnectActivated(func(widgets.QSystemTrayIcon__ActivationReason) {
		if s.dialog.IsVisible() {
			s.dialog.Hide()
		} else {
			s.dialog.Show()
		}
	})
	s.window.SetTabOrder(s.userInput, s.listView)
}

func (s *Service) loadSettings() {
	settings := core.NewQSettings3(
		core.QSettings__NativeFormat,
		core.QSettings__UserScope,
		"gisty",
		"app_settings",
		nil,
	)
	tmp := widgets.NewQWidget(nil, 0)
	tmp.SetGeometry2(100, 100, 600, 600)
	defSize := tmp.SaveGeometry()
	sizeVar := settings.Value("mainWindowGeometry", core.NewQVariant15(defSize))
	s.dialog.RestoreGeometry(sizeVar.ToByteArray())

	s.app.ConnectAboutToQuit(func() {
		current := s.dialog.SaveGeometry()
		currentVar := core.NewQVariant15(current.QByteArray_PTR())
		settings.SetValue("mainWindowGeometry", currentVar)
		settings.Sync()
	})
}

func (s *Service) populate(model *GistModel) {
	var foundOne bool
	for item := range s.GistService.Iter() {
		foundOne = true
		var gg = NewGist(nil)
		gg.SetGistID(item.ID)
		gg.SetGistURL(item.URL)
		gg.SetDescription(item.Description)
		model.AddGist(gg)
	}
	if !foundOne {
		messagebox(s.dialog).error("didn't find any gists")
	}
}

func (s *Service) gistDialog(index *core.QModelIndex) error {
	var content string
	id := index.Data(GistID).ToString()
	url := index.Data(GistURL).ToString()
	dialog := widgets.NewQMainWindow(s.dialog, 0)
	ui, err := qtlib.LoadResource(dialog, "./qml/gist.ui")
	if err != nil {
		return err
	}
	g, err := s.GistService.Get(id)
	if err != nil {
		return err
	}
	view := widgets.NewQPlainTextEditFromPointer(
		ui.FindChild("gist", core.Qt__FindChildrenRecursively).Pointer(),
	)
	for _, f := range g.Files {
		content = f.Content
		break
	}
	view.SetPlainText(content)
	widgets.NewQDialogFromPointer(ui.Pointer()).SetModal(true)
	ok := widgets.NewQPushButtonFromPointer(
		ui.FindChild("ok", core.Qt__FindChildrenRecursively).Pointer(),
	)
	ok.ConnectClicked(func(bool) {
		dialog.Close()
	})
	clipboard := widgets.NewQPushButtonFromPointer(
		ui.FindChild("clipboard", core.Qt__FindChildrenRecursively).Pointer(),
	)
	clipboard.ConnectClicked(func(bool) {
		s.app.Clipboard().SetText(content, gui.QClipboard__Clipboard)
		s.sysTray.ShowMessage("Info", "Gist has been copied to clipboard", widgets.QSystemTrayIcon__Information, 4000)
	})
	browser := widgets.NewQPushButtonFromPointer(
		ui.FindChild("browser", core.Qt__FindChildrenRecursively).Pointer(),
	)
	browser.ConnectClicked(func(bool) {
		gui.QDesktopServices_OpenUrl(core.NewQUrl3(url, 0))
	})

	dialog.Show()
	dialog.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			dialog.Close()
		}
	})
	return nil
}
