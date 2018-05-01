// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"testing"
)

func TestTabCreation(t *testing.T) {
	tRunner.Run(func() {
		testViewGistTab(t)
	})
}

func testViewGistTab(t *testing.T) {
	name := "test"
	g := &tabGist{
		id:      "uWIkJYdkFuVwYcyy",
		label:   "LpqrRCgBBYY",
		content: "fLGLysiOuxReut\nASUonvyd",
	}

	_, mw, cleanup := setup(t, name, nil, 0)
	defer cleanup()
	mw.setupUI()

	if mw.tabWidget.Count() != 1 {
		t.Errorf("mw.tabWidget.Count() = %d, want 1", mw.tabWidget.Count())
		return
	}

	tab := NewTab(mw.tabWidget)
	if tab == nil {
		t.Error("NewTab(mw.tabWidget) = nil, want *Tab")
		return
	}

	tab.showGist(mw.tabWidget, g)

	index := mw.tabWidget.CurrentIndex()
	shownTab := mw.tabWidget.Widget(index)
	if shownTab.Pointer() != tab.Pointer() {
		t.Errorf("shownTab.Pointer() = %v, want %v", shownTab.Pointer(), tab.Pointer())
	}
	if tab.textEdit.ToPlainText() != g.content {
		t.Errorf("content = %s, want %s", tab.textEdit.ToPlainText(), g.content)
	}
	if mw.tabWidget.TabText(index) != g.label {
		t.Errorf("TabText(%d) = %s, want %s", index, mw.tabWidget.TabText(index), g.label)
	}
	return
}
