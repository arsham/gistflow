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

// Service is in charge of user interaction with the dialog.
type Service struct {
	GistService gist.Service
	dialog      *widgets.QWidget
	window      *widgets.QMainWindow
	app         *widgets.QApplication
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
	s.window.SetCentralWidget(s.dialog)
	s.dialog.Show()
	s.setupUI()
	widgets.QApplication_Exec()

	return nil
}

func (s *Service) setupUI() {
	model := NewGistModel(nil)
	lv := widgets.NewQListViewFromPointer(
		s.dialog.FindChild("listView", core.Qt__FindChildrenRecursively).Pointer(),
	)
	lv.SetModel(model)
	go s.populate(model)
	lv.ConnectDoubleClicked(func(i *core.QModelIndex) {
		id := i.Data(GistID).ToString()
		err := s.showGist(id)
		if err != nil {
			messagebox(s.dialog).warningf("%s: %s", err, id)
		}
	})
	quit := widgets.NewQActionFromPointer(
		s.dialog.FindChild("actionQuit", core.Qt__FindChildrenRecursively).Pointer(),
	)
	quit.ConnectTriggered(func(bool) {
		s.app.Quit()
	})

}
func (s *Service) populate(model *GistModel) {
	list, err := s.GistService.List()
	if err != nil {
		messagebox(s.dialog).error(err.Error())
	}
	for _, item := range list {
		var gg = NewGist(nil)
		gg.SetGistID(item.ID)
		gg.SetDescription(item.Description)
		model.AddGist(gg)
	}
}

func (s *Service) showGist(id string) error {
	var content string
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
	})
	dialog.Show()
	dialog.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			dialog.Close()
		}
	})
	return nil
}
