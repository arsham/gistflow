// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"strings"
	"unicode"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (m *MainWindow) setupInteractions() {
	m.userInput.ConnectTextChanged(func(text string) {
		newText := strings.Split(text, "")
		m.proxy.SetFilterWildcard(strings.Join(newText, "*"))
	})
	m.userInput.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Up) || event.Key() == int(core.Qt__Key_Down) {
			m.gistList.SetFocus2()
		}
		m.userInput.KeyPressEventDefault(event)
	})

	m.gistList.ConnectDoubleClicked(m.gistListDoubleClickEvent)
	m.gistList.ConnectKeyReleaseEvent(m.gistListKeyReleaseEvent)

	m.sysTray.ConnectActivated(func(widgets.QSystemTrayIcon__ActivationReason) {
		if m.IsVisible() {
			m.Hide()
		} else {
			m.Show()
		}
	})

	m.tabWidget.ConnectKeyPressEvent(m.tabWidgetKeyPressEvent)
	m.tabWidget.ConnectTabCloseRequested(m.closeTab)
	m.menubar.action.actionClipboard.ConnectTriggered(func(bool) {
		widget := m.tabWidget.CurrentWidget()
		tab := NewTabFromPointer(widget.Pointer())
		m.app.Clipboard().SetText(tab.content(), gui.QClipboard__Clipboard)
	})
	m.menubar.action.actionCopyURL.ConnectTriggered(func(bool) {
		widget := m.tabWidget.CurrentWidget()
		tab := NewTabFromPointer(widget.Pointer())
		m.app.Clipboard().SetText(tab.url(), gui.QClipboard__Clipboard)
	})
}

func (m *MainWindow) gistListKeyReleaseEvent(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Enter, core.Qt__Key_Return:
		index := m.gistList.CurrentIndex()
		err := m.openGist(index.Data(GistID).ToString())
		if err != nil {
			m.logger.Error(err.Error())
		}
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

func (m *MainWindow) gistListDoubleClickEvent(*core.QModelIndex) {
	index := m.gistList.CurrentIndex()
	err := m.openGist(index.Data(GistID).ToString())
	if err != nil {
		m.logger.Error(err.Error())
	}
}

func (m *MainWindow) tabWidgetKeyPressEvent(event *gui.QKeyEvent) {
	// Closing tab
	if event.Modifiers() == core.Qt__ControlModifier {
		index := m.tabWidget.CurrentIndex()

		switch core.Qt__Key(event.Key()) {
		case core.Qt__Key_PageDown:
			m.tabWidget.SetCurrentIndex(index + 1)
		case core.Qt__Key_PageUp:
			m.tabWidget.SetCurrentIndex(index - 1)
		case core.Qt__Key_W:
			m.tabWidget.TabCloseRequested(index)
		}
	}

	// Moving left and right
	if event.Modifiers() == core.Qt__ShiftModifier+core.Qt__ControlModifier {
		widget := m.tabWidget.CurrentWidget()
		index := m.tabWidget.CurrentIndex()
		text := m.tabWidget.TabText(index)

		switch core.Qt__Key(event.Key()) {
		case core.Qt__Key_PageDown:
			m.tabWidget.RemoveTab(index)
			m.tabWidget.InsertTab(index+1, widget, text)
		case core.Qt__Key_PageUp:
			m.tabWidget.RemoveTab(index)
			m.tabWidget.InsertTab(index-1, widget, text)
		}
		m.tabWidget.SetCurrentWidget(widget)
	}
}
