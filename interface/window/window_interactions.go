// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"strings"
	"unicode"

	"github.com/arsham/gisty/interface/tab"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (m *MainWindow) userInputTextChange(text string) {
	newText := strings.Split(text, "")
	m.proxy.SetFilterWildcard(strings.Join(newText, "*"))
}

func (m *MainWindow) copyURLToClipboard(bool) {
	widget := m.TabsWidget().CurrentWidget()
	tab := tab.NewTabFromPointer(widget.Pointer())
	m.clipboard().SetText(tab.URL(), gui.QClipboard__Clipboard)
	m.showNotification("URL has been copied to clipboard")
}

func (m *MainWindow) openInBrowser(bool) {
	widget := m.TabsWidget().CurrentWidget()
	tab := tab.NewTabFromPointer(widget.Pointer())
	gui.QDesktopServices_OpenUrl(core.NewQUrl3(tab.HTMLURL(), 0))
}

func (m *MainWindow) sysTrayClick(widgets.QSystemTrayIcon__ActivationReason) {
	if m.IsVisible() {
		m.Hide()
	} else {
		m.Show()
	}
}

func (m *MainWindow) userInputChange(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Up, core.Qt__Key_Down:
		m.GistList().SetFocus2()
	}
	m.userInput.KeyPressEventDefault(event)
}

func (m *MainWindow) gistListKeyReleaseEvent(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Delete, core.Qt__Key_Space:
		fallthrough
	case core.Qt__Key_Left, core.Qt__Key_Right:
		m.userInput.SetFocus2()
	default:
		char := event.Text()
		for _, c := range char {
			if unicode.IsPrint(c) {
				m.userInput.SetText(m.userInput.Text() + char)
				m.userInput.SetFocus2()
			}
			break
		}
	}
}

func (m *MainWindow) openSelectedGist(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Enter, core.Qt__Key_Return:
		index := m.GistList().CurrentIndex()
		err := m.openGist(index.Data(tab.GistID).ToString())
		if err != nil {
			m.logger.Error(err.Error())
		}
		event.Accept()
	}
}

func (m *MainWindow) gistListDoubleClickEvent(*core.QModelIndex) {
	index := m.GistList().CurrentIndex()
	err := m.openGist(index.Data(tab.GistID).ToString())
	if err != nil {
		m.logger.Error(err.Error())
	}
}

func (m *MainWindow) showNotification(msg string) {
	m.sysTray.ShowMessage("Info", msg, widgets.QSystemTrayIcon__Information, 4000)
}

func (m *MainWindow) tabMovementEventFilter() *core.QObject {
	var filterObject = core.NewQObject(nil)
	filterObject.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {

		if event.Type() == core.QEvent__KeyPress {
			var keyEvent = gui.NewQKeyEventFromPointer(event.Pointer())

			// moving tabs
			if keyEvent.Modifiers() == core.Qt__ShiftModifier+core.Qt__ControlModifier {
				widget := m.TabsWidget().CurrentWidget()
				index := m.TabsWidget().CurrentIndex()
				text := m.TabsWidget().TabText(index)

				switch core.Qt__Key(keyEvent.Key()) {
				case core.Qt__Key_PageDown:
					m.TabsWidget().RemoveTab(index)
					m.TabsWidget().InsertTab(index+1, widget, text)
				case core.Qt__Key_PageUp:
					m.TabsWidget().RemoveTab(index)
					m.TabsWidget().InsertTab(index-1, widget, text)
				}

				switch core.Qt__Key(keyEvent.Key()) {
				case core.Qt__Key_PageDown, core.Qt__Key_PageUp:
					m.TabsWidget().SetCurrentWidget(widget)
					return true
				}
			}

			// switching and closing tabs
			if keyEvent.Modifiers() == core.Qt__ControlModifier {
				index := m.TabsWidget().CurrentIndex()

				switch core.Qt__Key(keyEvent.Key()) {
				case core.Qt__Key_PageDown:
					m.TabsWidget().SetCurrentIndex(index + 1)
				case core.Qt__Key_PageUp:
					m.TabsWidget().SetCurrentIndex(index - 1)
				case core.Qt__Key_W:
					m.TabsWidget().TabCloseRequested(index)
				}

				switch core.Qt__Key(keyEvent.Key()) {
				case core.Qt__Key_PageDown, core.Qt__Key_PageUp, core.Qt__Key_W:
					return true
				}
			}
		}

		return false
	})
	return filterObject
}
