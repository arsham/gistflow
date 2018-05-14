// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"github.com/arsham/gisty/qt/gistlist"
	"github.com/arsham/gisty/qt/tab"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (m *MainWindow) copyURLToClipboard(bool) {
	widget := m.tabsWidget.CurrentWidget()
	tab := tab.NewTabFromPointer(widget.Pointer())
	m.clipboard().SetText(tab.HTMLURL(), gui.QClipboard__Clipboard)
	m.showNotification("URL has been copied to clipboard")
}

func (m *MainWindow) openInBrowser(bool) {
	widget := m.tabsWidget.CurrentWidget()
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

// openSelectedGist opens the gist on GistList widget.
func (m *MainWindow) openSelectedGist(event *gui.QKeyEvent) {
	switch core.Qt__Key(event.Key()) {
	case core.Qt__Key_Enter, core.Qt__Key_Return:
		index := m.gistList.CurrentIndex()
		id := gistlist.NewContainerFromPointer(index.Pointer()).IndexID(index)
		m.openGistByID(id)
		event.Accept()
	}
}

func (m *MainWindow) openGistByID(id string) {
	err := m.openGist(id)
	if err != nil {
		m.logger.Error(err.Error())
	}
	m.searchbox.Hide()
}

func (m *MainWindow) gistListDoubleClickEvent(*core.QModelIndex) {
	index := m.gistList.CurrentIndex()
	id := gistlist.NewContainerFromPointer(index.Pointer()).IndexID(index)
	err := m.openGist(id)
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
				widget := m.tabsWidget.CurrentWidget()
				index := m.tabsWidget.CurrentIndex()
				text := m.tabsWidget.TabText(index)

				switch core.Qt__Key(keyEvent.Key()) {
				case core.Qt__Key_PageDown:
					m.tabsWidget.RemoveTab(index)
					m.tabsWidget.InsertTab(index+1, widget, text)
				case core.Qt__Key_PageUp:
					m.tabsWidget.RemoveTab(index)
					m.tabsWidget.InsertTab(index-1, widget, text)
				}

				switch core.Qt__Key(keyEvent.Key()) {
				case core.Qt__Key_PageDown, core.Qt__Key_PageUp:
					m.tabsWidget.SetCurrentWidget(widget)
					return true
				}
			}

			// switching and closing tabs
			if keyEvent.Modifiers() == core.Qt__ControlModifier {
				index := m.tabsWidget.CurrentIndex()

				switch core.Qt__Key(keyEvent.Key()) {
				case core.Qt__Key_PageDown:
					m.tabsWidget.SetCurrentIndex(index + 1)
				case core.Qt__Key_PageUp:
					m.tabsWidget.SetCurrentIndex(index - 1)
				case core.Qt__Key_W:
					m.tabsWidget.TabCloseRequested(index)
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
