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
	if m.Options() == nil {
		t.Error("m.Options() = nil, want *widgets.QMenu")
		return
	}
	if m.Actions().Quit == nil {
		t.Error("m.Actions().Quit = nil, want *widgets.QAction")
		return
	}
	actions := m.Options().Actions()
	if len(actions) == 0 {
		t.Error("len(m.Options().Actions()) = 0, want at least 1")
		return
	}
}

func TestCtrlQ(t *testing.T) { tRunner.Run(func() { testCtrlQ(t) }) }
func testCtrlQ(t *testing.T) {
	var called bool
	window := widgets.NewQMainWindow(nil, 0)
	m := NewMenuBar(window)
	window.Show()
	app.SetActiveWindow(window)
	m.Actions().Quit.ConnectEvent(func(e *core.QEvent) bool {
		called = true
		return true
	})

	event := testlib.NewQTestEventList()
	event.AddKeyClick(core.Qt__Key_Q, core.Qt__ControlModifier, -1)
	event.Simulate(window)

	if !called {
		t.Error("Ctrl+Q didn't trigger the actions.Quit")
	}
	return
}
