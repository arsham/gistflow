// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the LGPL-v3 license
// License that can be found in the LICENSE file.

package menubar

import (
	"os"
	"testing"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/testlib"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func TestMenuBar(t *testing.T) { tRunner.Run(func() { testMenuBar(t) }) }
func testMenuBar(t *testing.T) {
	parent := widgets.NewQWidget(nil, 0)
	m := NewMenuBar(parent)
	if m.options == nil {
		t.Error("m.options = nil, want *widgets.QMenu")
		return
	}
	if m.actions.Quit == nil {
		t.Error("m.actions.Quit = nil, want *widgets.QAction")
		return
	}
	if m.actions.NewGist == nil {
		t.Error("m.actions.NewGist = nil, want *widgets.QAction")
		return
	}
	actions := m.options.Actions()
	if len(actions) == 0 {
		t.Error("len(m.options.actions) = 0, want at least 1")
	}
}

func TestCtrlQ(t *testing.T) { tRunner.Run(func() { testCtrlQ(t) }) }
func testCtrlQ(t *testing.T) {
	var called bool
	window := widgets.NewQMainWindow(nil, 0)
	m := NewMenuBar(window)
	window.Show()
	defer window.Hide()
	app.SetActiveWindow(window)
	m.actions.Quit.ConnectEvent(func(e *core.QEvent) bool {
		called = true
		return true
	})

	event := testlib.NewQTestEventList()
	event.AddKeyClick(core.Qt__Key_Q, core.Qt__ControlModifier, -1)
	event.Simulate(window)

	if !called {
		t.Error("Ctrl+Q didn't trigger the actions.Quit")
	}
}

func TestCtrlN(t *testing.T) { tRunner.Run(func() { testCtrlN(t) }) }
func testCtrlN(t *testing.T) {
	var called bool
	window := widgets.NewQMainWindow(nil, 0)
	m := NewMenuBar(window)
	window.Show()
	defer window.Hide()
	app.SetActiveWindow(window)
	m.actions.NewGist.ConnectEvent(func(e *core.QEvent) bool {
		called = true
		return true
	})

	event := testlib.NewQTestEventList()
	event.AddKeyClick(core.Qt__Key_N, core.Qt__ControlModifier, -1)
	event.Simulate(window)

	if !called {
		t.Error("Ctrl+N didn't trigger the actions.NewGist")
	}
}

func TestToggleToolbar(t *testing.T) { tRunner.Run(func() { testToggleToolbar(t) }) }
func testToggleToolbar(t *testing.T) {
	var called bool
	window := widgets.NewQMainWindow(nil, 0)
	m := NewMenuBar(window)
	window.Show()
	defer window.Hide()
	app.SetActiveWindow(window)

	m.ConnectToggleToolbar(func(bool) {
		called = true
	})

	m.actions.Toolbar.Trigger()
	if !called {
		t.Error("didn't send the ToggleToolbar signal")
	}
}

func TestToggleGistList(t *testing.T) { tRunner.Run(func() { testToggleGistList(t) }) }
func testToggleGistList(t *testing.T) {
	var called bool
	window := widgets.NewQMainWindow(nil, 0)
	m := NewMenuBar(window)
	window.Show()
	defer window.Hide()
	app.SetActiveWindow(window)

	m.ConnectToggleGistList(func(bool) {
		called = true
	})

	m.actions.GistList.Trigger()
	if !called {
		t.Error("didn't send the ToggleGistList signal")
	}
}
