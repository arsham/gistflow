// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"testing"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/testlib"
	"github.com/therecipe/qt/widgets"
)

func TestMenubar(t *testing.T) {
	tRunner.Run(func() {
		tcs := map[string]func(*testing.T) bool{
			"testMenuBar": testMenuBar,
			"testCtrlQ":   testCtrlQ,
			"testToggle":  testToggle,
		}
		for name, tc := range tcs {
			if !tc(t) {
				t.Errorf("stopped at %s", name)
				return
			}
		}
	})
}

func testMenuBar(t *testing.T) bool {
	parent := widgets.NewQWidget(nil, 0)
	m := NewMenuBar(parent)
	if m.menuOptions == nil {
		t.Error("m.menuOptions = nil, want *widgets.QMenu")
		return false
	}
	if m.action.actionQuit == nil {
		t.Error("m.action.actionQuit = nil, want *widgets.QAction")
		return false
	}
	actions := m.menuOptions.Actions()
	if len(actions) == 0 {
		t.Error("len(m.menuOptions.Actions()) = 0, want at least 1")
		return false
	}
	return true
}

func testCtrlQ(t *testing.T) bool {
	var called bool
	window := widgets.NewQMainWindow(nil, 0)
	m := NewMenuBar(window)
	window.Show()
	app.SetActiveWindow(window)
	m.action.actionQuit.ConnectEvent(func(e *core.QEvent) bool {
		called = true
		return true
	})

	event := testlib.NewQTestEventList()
	event.AddKeyClick(core.Qt__Key_Q, core.Qt__ControlModifier, -1)
	event.Simulate(window)

	if !called {
		t.Error("Ctrl+Q didn't trigger the actions.actionQuit")
	}
	return true
}

func testToggle(t *testing.T) bool {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	window.setupUI()
	window.setupInteractions()
	app.SetActiveWindow(window)
	window.Show()

	if window.IsHidden() {
		t.Error("window is not shown")
	}
	window.sysTray.Activated(widgets.QSystemTrayIcon__Trigger)
	if !window.IsHidden() {
		t.Error("window is not hidden")
	}
	window.sysTray.Activated(widgets.QSystemTrayIcon__Trigger)
	if window.IsHidden() {
		t.Error("window is not shown")
	}
	return true
}
