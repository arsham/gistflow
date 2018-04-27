package window

import (
	"strings"
	"unicode"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (s *Service) mainDialog() (*widgets.QWidget, error) {
	dialog := widgets.NewQWidget(nil, 0)
	s.layout = widgets.NewQVBoxLayout()
	dialog.SetLayout(s.layout)
	return dialog, nil
}

func (s *Service) listViewWidget() *widgets.QListView {
	openGist := func(index *core.QModelIndex) {
		err := s.gistDialog(index)
		if err != nil {
			id := index.Data(GistID).ToString()
			messagebox(s.dialog).warningf("%s: %s", err, id)
		}
	}

	listView := widgets.NewQListView(s.dialog)
	listView.ConnectDoubleClicked(openGist)
	listView.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		switch event.Key() {
		case int(core.Qt__Key_Enter), int(core.Qt__Key_Return):
			index := listView.CurrentIndex()
			openGist(index)
		case int(core.Qt__Key_Up), int(core.Qt__Key_Down):
		default:
			char := event.Text()
			for _, c := range char {
				if unicode.IsPrint(c) {
					s.userInput.SetText(s.userInput.Text() + char)
					s.userInput.SetFocus2()
				}
				break
			}
		}
	})

	return listView
}

func (s *Service) userInputWidget(proxy *core.QSortFilterProxyModel) *widgets.QLineEdit {
	userInput := widgets.NewQLineEdit(s.dialog)
	userInput.SetClearButtonEnabled(true)
	userInput.ConnectTextChanged(func(text string) {
		newText := strings.Split(text, "")
		proxy.SetFilterWildcard(strings.Join(newText, "*"))
	})
	userInput.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Up) || event.Key() == int(core.Qt__Key_Down) {
			s.listView.SetFocus2()
		}
		userInput.KeyPressEventDefault(event)
	})

	return userInput
}

func (s *Service) mainMenu() *widgets.QMenu {
	quit := widgets.NewQAction(s.window)
	quit.SetObjectName("actionQuit")
	quit.ConnectTriggered(func(bool) {
		s.app.Quit()
	})
	quit.SetShortcut(gui.QKeySequence_FromString("Ctrl+Q", 0))
	quit.SetText("Quit")

	mainMenu := widgets.NewQMenu2("Option", s.window)
	mainMenu.AddActions([]*widgets.QAction{quit})

	menuBar := widgets.NewQMenuBar(s.window)
	menuBar.SetObjectName("menuBar")
	menuBar.AddMenu(mainMenu)

	s.layout.SetMenuBar(menuBar)
	return mainMenu
}
