// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

func TestTabCreation(t *testing.T) {
	tRunner.Run(func() {
		testTabCreation(t)
		testTabIDFromIndex(t)
		testSwitchTabs(t)
		testMovingTabs(t)
		testWindowClickCloseTab(t)
		testShortcutTabClose(t)
		testRemoveOpenTab(t)
	})
}

func testTabCreation(t *testing.T) {
	name := "test"
	g := &tabGist{
		id:      "uWIkJYdkFuVwYcyy",
		label:   "LpqrRCgBBYY",
		content: "fLGLysiOuxReut\nASUonvyd",
	}

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()

	if window.tabWidget.Count() != 1 {
		t.Errorf("window.tabWidget.Count() = %d, want 1", window.tabWidget.Count())
		return
	}

	tab := NewTab(window.tabWidget)
	if tab == nil {
		t.Error("NewTab(window.tabWidget) = nil, want *Tab")
		return
	}

	tab.showGist(window.tabWidget, g)

	index := window.tabWidget.CurrentIndex()
	shownTab := window.tabWidget.Widget(index)
	if shownTab.Pointer() != tab.Pointer() {
		t.Errorf("shownTab.Pointer() = %v, want %v", shownTab.Pointer(), tab.Pointer())
	}
	if tab.textEdit.ToPlainText() != g.content {
		t.Errorf("content = %s, want %s", tab.textEdit.ToPlainText(), g.content)
	}
	if window.tabWidget.TabText(index) != g.label {
		t.Errorf("TabText(%d) = %s, want %s", index, window.tabWidget.TabText(index), g.label)
	}
}

func testTabIDFromIndex(t *testing.T) {
	var (
		name = "test"
		id1  = "mbzsNwJS"
		id2  = "eulYvWSUHubADRV"
	)

	files := map[string]gist.ResponseFile{
		"GqrqZkTNpw": gist.ResponseFile{Content: "uMWDmwSvLlqtFXZUX"},
	}
	gres := gist.ResponseGist{
		Files: files,
	}

	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(gres)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.setupUI()
	window.gistService.API = gistTs.URL
	currentIndex := window.tabWidget.CurrentIndex()
	window.openGist(id1)
	index1 := currentIndex + 1
	window.openGist(id2)
	index2 := currentIndex + 2

	if window.tabIDFromIndex(999) != "" {
		t.Errorf("window.tabIDFromIndex(%d) = %s, want empty string", 999, window.tabIDFromIndex(999))
	}

	if window.tabIDFromIndex(index1) != id1 {
		t.Errorf("window.tabIDFromIndex(%d) = %s, want %s", index1, window.tabIDFromIndex(index1), id1)
	}

	if window.tabIDFromIndex(index2) != id2 {
		t.Errorf("window.tabIDFromIndex(%d) = %s, want %s", index2, window.tabIDFromIndex(index2), id2)
	}
}

func testSwitchTabs(t *testing.T) {
	var name = "test"
	g1 := &tabGist{
		id:      "uWIkJYdkFuVwYcyy",
		label:   "LpqrRCgBBYY",
		content: "fLGLysiOuxReut\nASUonvyd",
	}
	g2 := &tabGist{
		id:      "FJsPzPqhI",
		label:   "bsDmGRE",
		content: "KuiIIVYnCKycPPkXLibh",
	}

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	window.setupInteractions()
	app.SetActiveWindow(window)
	window.show()

	leftTab := NewTab(window.tabWidget)
	leftTab.showGist(window.tabWidget, g1)
	rightTab := NewTab(window.tabWidget)
	rightTab.showGist(window.tabWidget, g2)

	leftIndex := window.tabWidget.IndexOf(leftTab)
	rightIndex := window.tabWidget.IndexOf(rightTab)
	if leftIndex > rightIndex {
		rightIndex = leftIndex
	}
	window.tabWidget.SetCurrentIndex(rightIndex)
	window.tabWidget.SetFocus2()

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageUp), core.Qt__ControlModifier, "", false, 1)
	window.tabWidget.KeyPressEvent(event)

	event = gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageDown), core.Qt__ControlModifier, "", false, 1)
	window.tabWidget.KeyPressEvent(event)
	if window.tabWidget.CurrentIndex() != rightIndex {
		t.Errorf("window.tabWidget.CurrentIndex() = %d, want %d", window.tabWidget.CurrentIndex(), rightIndex)
	}
}

func testMovingTabs(t *testing.T) {
	var name = "test"
	g1 := &tabGist{
		id:      "uWIkJYdkFuVwYcyy",
		label:   "LpqrRCgBBYY",
		content: "fLGLysiOuxReut\nASUonvyd",
	}
	g2 := &tabGist{
		id:      "FJsPzPqhI",
		label:   "bsDmGRE",
		content: "KuiIIVYnCKycPPkXLibh",
	}

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	window.setupInteractions()
	app.SetActiveWindow(window)
	window.show()

	leftTab := NewTab(window.tabWidget)
	leftTab.showGist(window.tabWidget, g1)
	rightTab := NewTab(window.tabWidget)
	rightTab.showGist(window.tabWidget, g2)

	leftIndex := window.tabWidget.IndexOf(leftTab)
	rightIndex := window.tabWidget.IndexOf(rightTab)
	if leftIndex > rightIndex {
		// I just want to position them properly.
		leftIndex, rightIndex = rightIndex, leftIndex
		leftTab, rightTab = rightTab, leftTab
	}
	window.tabWidget.SetCurrentIndex(rightIndex)
	window.tabWidget.SetFocus2()

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageUp), core.Qt__ControlModifier+core.Qt__ShiftModifier, "", false, 1)
	window.tabWidget.KeyPressEvent(event)
	if window.tabWidget.IndexOf(rightTab) != leftIndex {
		t.Errorf("window.tabWidget.IndexOf(rightTab) = %d, want %d", window.tabWidget.IndexOf(rightTab), leftIndex)
	}
	if window.tabWidget.CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Error("focus is still on leftTab, want rightTab")
	}

	// now we are swichng back
	event = gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageDown), core.Qt__ControlModifier+core.Qt__ShiftModifier, "", false, 1)
	window.tabWidget.KeyPressEvent(event)
	if window.tabWidget.IndexOf(rightTab) != rightIndex {
		t.Errorf("window.tabWidget.IndexOf(rightTab) = %d, want %d", window.tabWidget.IndexOf(rightTab), rightIndex)
	}
	if window.tabWidget.CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Error("focus is still on leftTab, want rightTab")
	}
}

func testRemoveOpenTab(t *testing.T) {
	var (
		name = "test"
		id1  = "mbzsNwJS"
		id2  = "eulYvWSUHubADRV"
		id3  = "zdXyiCAdDkG"
	)

	files := map[string]gist.ResponseFile{
		"GqrqZkTNpw": gist.ResponseFile{Content: "uMWDmwSvLlqtFXZUX"},
	}
	gres := gist.ResponseGist{
		Files: files,
	}

	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(gres)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.setupUI()
	window.setupInteractions()
	window.gistService.API = gistTs.URL

	currentLen := len(window.tabGistList)
	window.openGist(id1)
	if len(window.tabGistList) != currentLen+1 {
		t.Errorf("len(window.tabGistList) = %d, want %d", len(window.tabGistList), currentLen+1)
	}
	if _, ok := window.tabGistList[id1]; !ok {
		t.Errorf("%s not found in window.tabGistList", id1)
	}
	window.tabWidget.TabCloseRequested(window.tabWidget.CurrentIndex())
	if len(window.tabGistList) != currentLen {
		t.Errorf("len(window.tabGistList) = %d, want %d", len(window.tabGistList), currentLen)
	}
	if _, ok := window.tabGistList[id1]; ok {
		t.Errorf("%s is still in window.tabGistList", id1)
	}
	window.openGist(id2)
	window.openGist(id3)
	if err := window.openGist(id1); err != nil {
		t.Errorf("window.openGist(%s) = %v, want nil", id1, window.openGist(id1))
	}
	if len(window.tabGistList) != currentLen+3 {
		t.Errorf("len(window.tabGistList) = %d, want %d", len(window.tabGistList), currentLen+3)
	}
	if _, ok := window.tabGistList[id1]; !ok {
		t.Errorf("%s not found in window.tabGistList", id1)
	}

	index := window.tabWidget.IndexOf(window.tabGistList[id1])
	if window.tabWidget.CurrentIndex() != index {
		t.Errorf("window.tabWidget.CurrentIndex() = %d, want %d", window.tabWidget.CurrentIndex(), index)
	}
}

func testWindowClickCloseTab(t *testing.T) {
	var (
		name   = "test"
		called bool
	)

	g := &tabGist{
		id:      "uWIkJYdkFuVwYcyy",
		label:   "LpqrRCgBBYY",
		content: "fLGLysiOuxReut\nASUonvyd",
	}

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	window.setupInteractions()
	app.SetActiveWindow(window)
	window.show()

	tab := NewTab(window.tabWidget)
	if tab == nil {
		t.Error("NewTab(window.tabWidget) = nil, want *Tab")
		return
	}

	tab.showGist(window.tabWidget, g)
	currentSize := window.tabWidget.Count()
	index := window.tabWidget.IndexOf(tab)
	window.tabWidget.ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("i = %d, want %d", i, index)
	})

	window.tabWidget.TabCloseRequested(index)

	if !called {
		t.Error("didn't close the tab")
	}
	if window.tabWidget.Count() != currentSize-1 {
		t.Errorf("window.tabWidget.Count() = %d, want %d", window.tabWidget.Count(), currentSize-1)
	}
	if window.tabWidget.IndexOf(tab) != -1 {
		t.Errorf("window.tabWidget.IndexOf(tab) = %d, want %d", window.tabWidget.IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.id]; ok {
		t.Errorf("%s was not removed from the list", g.id)
	}
}

func testShortcutTabClose(t *testing.T) {
	var (
		name   = "test"
		called bool
	)
	g := &tabGist{
		id:      "TGtyIHIK",
		label:   "hPtRE",
		content: "quIwMlPsoVaNr",
	}

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	window.setupInteractions()
	app.SetActiveWindow(window)
	window.show()

	tab := NewTab(window.tabWidget)
	tab.showGist(window.tabWidget, g)
	index := window.tabWidget.IndexOf(tab)
	currentSize := window.tabWidget.Count()

	window.tabWidget.SetCurrentWidget(tab)
	window.tabWidget.SetFocus2()

	window.tabWidget.ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("i = %d, want %d", i, index)
	})

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_W), core.Qt__ControlModifier, "", false, 1)
	window.tabWidget.KeyPressEvent(event)

	if !called {
		t.Error("didn't close the tab")
	}
	if window.tabWidget.Count() != currentSize-1 {
		t.Errorf("window.tabWidget.Count() = %d, want %d", window.tabWidget.Count(), currentSize-1)
	}
	if window.tabWidget.IndexOf(tab) != -1 {
		t.Errorf("window.tabWidget.IndexOf(tab) = %d, want %d", window.tabWidget.IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.id]; ok {
		t.Errorf("%s was not removed from the list", g.id)
	}
}
