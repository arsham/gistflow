// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

// Package window shows all kinds of windows and dialogs.
package window

import (
	"os"
	"strings"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/qtlib"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Service is in charge of user interaction with the dialog.
type Service struct {
	GistService gist.Service
	dialog      *widgets.QWidget
	window      *widgets.QMainWindow
	app         *widgets.QApplication
	sysTray     *widgets.QSystemTrayIcon
}

// MainWindow shows the main window.
func (s *Service) MainWindow() error {
	var err error
	core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	s.app = widgets.NewQApplication(len(os.Args), os.Args)
	s.window = widgets.NewQMainWindow(nil, 0)
	s.dialog, err = qtlib.LoadResource(s.window, "./qml/mainwindow.ui")
	if err != nil {
		return err
	}
	icon := gui.NewQIcon5("./qml/app.ico")

	s.sysTray = widgets.NewQSystemTrayIcon(s.dialog)
	s.sysTray.SetIcon(icon)
	s.sysTray.SetVisible(true)
	s.sysTray.SetToolTip("Gisty")

	s.window.SetCentralWidget(s.dialog)
	s.window.SetWindowIcon(icon)

	s.dialog.Show()
	s.setupUI()
	widgets.QApplication_Exec()

	return nil
}

func (s *Service) setupUI() {
	model := NewGistModel(nil)
	proxy := core.NewQSortFilterProxyModel(nil)
	proxy.SetFilterCaseSensitivity(core.Qt__CaseInsensitive)
	proxy.SetSourceModel(model)

	listView := widgets.NewQListViewFromPointer(
		s.dialog.FindChild("listView", core.Qt__FindChildrenRecursively).Pointer(),
	)
	listView.SetModel(proxy)
	go s.populate(model)
	listView.ConnectDoubleClicked(func(index *core.QModelIndex) {
		err := s.gistDialog(index)
		if err != nil {
			id := index.Data(GistID).ToString()
			messagebox(s.dialog).warningf("%s: %s", err, id)
		}
	})
	quit := widgets.NewQActionFromPointer(
		s.dialog.FindChild("actionQuit", core.Qt__FindChildrenRecursively).Pointer(),
	)
	userInput := widgets.NewQLineEditFromPointer(
		s.dialog.FindChild("userInput", core.Qt__FindChildrenRecursively).Pointer(),
	)
	userInput.SetClearButtonEnabled(true)
	userInput.ConnectTextChanged(func(text string) {
		newText := strings.Split(text, "")
		proxy.SetFilterWildcard(strings.Join(newText, "*"))
	})
	quit.ConnectTriggered(func(bool) {
		s.app.Quit()
	})

	mainMenu := widgets.NewQMenuFromPointer(
		s.dialog.FindChild("mainMenu", core.Qt__FindChildrenRecursively).Pointer(),
	)

	s.sysTray.SetContextMenu(mainMenu)
	s.sysTray.ConnectActivated(func(widgets.QSystemTrayIcon__ActivationReason) {
		if s.dialog.IsVisible() {
			s.dialog.Hide()
		} else {
			s.dialog.Show()
		}
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
