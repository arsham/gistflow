package window

import (
	"fmt"
	"log"
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
	s.dialog, err = qtlib.LoadResource(s.window, "./window/qml/mainwindow.ui")
	if err != nil {
		log.Fatal(err)
	}
	s.window.SetCentralWidget(s.dialog)

	list, err := s.GistService.List()
	if err != nil {
		return err
	}
	model := core.NewQStringListModel(s.dialog)
	m := make([]string, len(list))

	lv := widgets.NewQListViewFromPointer(
		s.dialog.FindChild("listView", core.Qt__FindChildrenRecursively).Pointer(),
	)

	for i, item := range list {
		m[i] = item.Description
	}
	model.SetStringList(m)

	lv.ConnectDoubleClicked(func(i *core.QModelIndex) {
		err := s.showGist(list[i.Row()].ID)
		if err != nil {
			fmt.Println(err) // TODO: show a dialog
		}
	})
	lv.SetModel(model)
	s.dialog.Show()
	widgets.QApplication_Exec()

	return nil
}

func (s *Service) showGist(id string) error {
	var content string
	dialog := widgets.NewQMainWindow(s.dialog, 0)
	ui, err := qtlib.LoadResource(dialog, "./window/qml/gist.ui")
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
	return nil
}
