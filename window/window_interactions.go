// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"strings"
	"unicode"

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
	tab := NewTabFromPointer(widget.Pointer())
	m.App().Clipboard().SetText(tab.url(), gui.QClipboard__Clipboard)
}

func (m *MainWindow) copyToClipboard(bool) {
	widget := m.TabsWidget().CurrentWidget()
	tab := NewTabFromPointer(widget.Pointer())
	m.App().Clipboard().SetText(tab.content(), gui.QClipboard__Clipboard)
}

func (m *MainWindow) sysTrayClick(widgets.QSystemTrayIcon__ActivationReason) {
	if m.IsVisible() {
		m.Hide()
	} else {
		m.Show()
	}
}

func (m *MainWindow) userInputChange(event *gui.QKeyEvent) {
	if event.Key() == int(core.Qt__Key_Up) || event.Key() == int(core.Qt__Key_Down) {
		m.GistList().SetFocus2()
	}
	m.userInput.KeyPressEventDefault(event)
}

func (m *MainWindow) gistListKeyReleaseEvent(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Enter, core.Qt__Key_Return:
		index := m.GistList().CurrentIndex()
		err := m.openGist(index.Data(gistID).ToString())
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
	index := m.GistList().CurrentIndex()
	err := m.openGist(index.Data(gistID).ToString())
	if err != nil {
		m.logger.Error(err.Error())
	}
}

func (m *MainWindow) tabWidgetKeyPressEvent(event *gui.QKeyEvent) {
	// Closing tab
	if event.Modifiers() == core.Qt__ControlModifier {
		index := m.TabsWidget().CurrentIndex()

		switch core.Qt__Key(event.Key()) {
		case core.Qt__Key_PageDown:
			m.TabsWidget().SetCurrentIndex(index + 1)
		case core.Qt__Key_PageUp:
			m.TabsWidget().SetCurrentIndex(index - 1)
		case core.Qt__Key_W:
			m.TabsWidget().TabCloseRequested(index)
		}
	}

	// Moving left and right
	if event.Modifiers() == core.Qt__ShiftModifier+core.Qt__ControlModifier {
		widget := m.TabsWidget().CurrentWidget()
		index := m.TabsWidget().CurrentIndex()
		text := m.TabsWidget().TabText(index)

		switch core.Qt__Key(event.Key()) {
		case core.Qt__Key_PageDown:
			m.TabsWidget().RemoveTab(index)
			m.TabsWidget().InsertTab(index+1, widget, text)
		case core.Qt__Key_PageUp:
			m.TabsWidget().RemoveTab(index)
			m.TabsWidget().InsertTab(index-1, widget, text)
		}
		m.TabsWidget().SetCurrentWidget(widget)
	}
}
