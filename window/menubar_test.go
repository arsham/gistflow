// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/testlib"
	"github.com/therecipe/qt/widgets"

	"testing"
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
	if m.optionsMenu == nil {
		t.Error("m.optionsMenu = nil, want *widgets.QMenu")
		return false
	}
	if m.quitAction == nil {
		t.Error("m.quitAction = nil, want *widgets.QAction")
		return false
	}
	actions := m.optionsMenu.Actions()
	if len(actions) == 0 {
		t.Error("len(m.optionsMenu.Actions()) = 0, want at least 1")
		return false
	}
	var foundIt bool
	for _, a := range actions {
		if a.Pointer() == m.quitAction.Pointer() {
			foundIt = true
			break
		}
	}
	if !foundIt {
		t.Error("m.quitAction not found in actions")
	}
	return true
}

func testCtrlQ(t *testing.T) bool {
	var called bool
	window := widgets.NewQMainWindow(nil, 0)
	m := NewMenuBar(window)
	window.Show()
	app.SetActiveWindow(window)
	m.quitAction.ConnectEvent(func(e *core.QEvent) bool {
		called = true
		return true
	})

	event := testlib.NewQTestEventList()
	event.AddKeyClick(core.Qt__Key_Q, core.Qt__ControlModifier, -1)
	event.Simulate(window)

	if !called {
		t.Error("Ctrl+Q didn't trigger the quitAction")
	}
	return true
}

func testToggle(t *testing.T) bool {
	name := "test"
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	mw.setupInteractions()
	app.SetActiveWindow(mw.window)
	mw.window.Show()

	if mw.window.IsHidden() {
		t.Error("window is not shown")
	}
	mw.sysTray.Activated(widgets.QSystemTrayIcon__Trigger)
	if !mw.window.IsHidden() {
		t.Error("window is not hidden")
	}
	mw.sysTray.Activated(widgets.QSystemTrayIcon__Trigger)
	if mw.window.IsHidden() {
		t.Error("window is not shown")
	}
	return true
}
