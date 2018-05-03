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

	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
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
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	mw.setupUI()
	mw.GistService.API = gistTs.URL
	currentIndex := mw.tabWidget.CurrentIndex()
	mw.openGist(id1)
	index1 := currentIndex + 1
	mw.openGist(id2)
	index2 := currentIndex + 2

	if mw.tabIDFromIndex(999) != "" {
		t.Errorf("mw.tabIDFromIndex(%d) = %s, want empty string", 999, mw.tabIDFromIndex(999))
	}

	if mw.tabIDFromIndex(index1) != id1 {
		t.Errorf("mw.tabIDFromIndex(%d) = %s, want %s", index1, mw.tabIDFromIndex(index1), id1)
	}

	if mw.tabIDFromIndex(index2) != id2 {
		t.Errorf("mw.tabIDFromIndex(%d) = %s, want %s", index2, mw.tabIDFromIndex(index2), id2)
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

	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	mw.setupUI()
	mw.setupInteractions()
	app.SetActiveWindow(mw.window)
	mw.show()

	leftTab := NewTab(mw.tabWidget)
	leftTab.showGist(mw.tabWidget, g1)
	rightTab := NewTab(mw.tabWidget)
	rightTab.showGist(mw.tabWidget, g2)

	leftIndex := mw.tabWidget.IndexOf(leftTab)
	rightIndex := mw.tabWidget.IndexOf(rightTab)
	if leftIndex > rightIndex {
		rightIndex = leftIndex
	}
	mw.tabWidget.SetCurrentIndex(rightIndex)
	mw.tabWidget.SetFocus2()

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageUp), core.Qt__ControlModifier, "", false, 1)
	mw.tabWidget.KeyPressEvent(event)

	event = gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageDown), core.Qt__ControlModifier, "", false, 1)
	mw.tabWidget.KeyPressEvent(event)
	if mw.tabWidget.CurrentIndex() != rightIndex {
		t.Errorf("mw.tabWidget.CurrentIndex() = %d, want %d", mw.tabWidget.CurrentIndex(), rightIndex)
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

	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	mw.setupUI()
	mw.setupInteractions()
	app.SetActiveWindow(mw.window)
	mw.show()

	leftTab := NewTab(mw.tabWidget)
	leftTab.showGist(mw.tabWidget, g1)
	rightTab := NewTab(mw.tabWidget)
	rightTab.showGist(mw.tabWidget, g2)

	leftIndex := mw.tabWidget.IndexOf(leftTab)
	rightIndex := mw.tabWidget.IndexOf(rightTab)
	if leftIndex > rightIndex {
		// I just want to position them properly.
		leftIndex, rightIndex = rightIndex, leftIndex
		leftTab, rightTab = rightTab, leftTab
	}
	mw.tabWidget.SetCurrentIndex(rightIndex)
	mw.tabWidget.SetFocus2()

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageUp), core.Qt__ControlModifier+core.Qt__ShiftModifier, "", false, 1)
	mw.tabWidget.KeyPressEvent(event)
	if mw.tabWidget.IndexOf(rightTab) != leftIndex {
		t.Errorf("mw.tabWidget.IndexOf(rightTab) = %d, want %d", mw.tabWidget.IndexOf(rightTab), leftIndex)
	}
	if mw.tabWidget.CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Error("focus is still on leftTab, want rightTab")
	}

	// now we are swichng back
	event = gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageDown), core.Qt__ControlModifier+core.Qt__ShiftModifier, "", false, 1)
	mw.tabWidget.KeyPressEvent(event)
	if mw.tabWidget.IndexOf(rightTab) != rightIndex {
		t.Errorf("mw.tabWidget.IndexOf(rightTab) = %d, want %d", mw.tabWidget.IndexOf(rightTab), rightIndex)
	}
	if mw.tabWidget.CurrentWidget().Pointer() != rightTab.Pointer() {
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
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	mw.setupUI()
	mw.setupInteractions()
	mw.GistService.API = gistTs.URL

	currentLen := len(mw.tabGistList)
	mw.openGist(id1)
	if len(mw.tabGistList) != currentLen+1 {
		t.Errorf("len(mw.tabGistList) = %d, want %d", len(mw.tabGistList), currentLen+1)
	}
	if _, ok := mw.tabGistList[id1]; !ok {
		t.Errorf("%s not found in mw.tabGistList", id1)
	}
	mw.tabWidget.TabCloseRequested(mw.tabWidget.CurrentIndex())
	if len(mw.tabGistList) != currentLen {
		t.Errorf("len(mw.tabGistList) = %d, want %d", len(mw.tabGistList), currentLen)
	}
	if _, ok := mw.tabGistList[id1]; ok {
		t.Errorf("%s is still in mw.tabGistList", id1)
	}
	mw.openGist(id2)
	mw.openGist(id3)
	if err := mw.openGist(id1); err != nil {
		t.Errorf("mw.openGist(%s) = %v, want nil", id1, mw.openGist(id1))
	}
	if len(mw.tabGistList) != currentLen+3 {
		t.Errorf("len(mw.tabGistList) = %d, want %d", len(mw.tabGistList), currentLen+3)
	}
	if _, ok := mw.tabGistList[id1]; !ok {
		t.Errorf("%s not found in mw.tabGistList", id1)
	}

	index := mw.tabWidget.IndexOf(mw.tabGistList[id1])
	if mw.tabWidget.CurrentIndex() != index {
		t.Errorf("mw.tabWidget.CurrentIndex() = %d, want %d", mw.tabWidget.CurrentIndex(), index)
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

	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	mw.setupUI()
	mw.setupInteractions()
	app.SetActiveWindow(mw.window)
	mw.show()

	tab := NewTab(mw.tabWidget)
	if tab == nil {
		t.Error("NewTab(mw.tabWidget) = nil, want *Tab")
		return
	}

	tab.showGist(mw.tabWidget, g)
	currentSize := mw.tabWidget.Count()
	index := mw.tabWidget.IndexOf(tab)
	mw.tabWidget.ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("i = %d, want %d", i, index)
	})

	mw.tabWidget.TabCloseRequested(index)

	if !called {
		t.Error("didn't close the tab")
	}
	if mw.tabWidget.Count() != currentSize-1 {
		t.Errorf("mw.tabWidget.Count() = %d, want %d", mw.tabWidget.Count(), currentSize-1)
	}
	if mw.tabWidget.IndexOf(tab) != -1 {
		t.Errorf("mw.tabWidget.IndexOf(tab) = %d, want %d", mw.tabWidget.IndexOf(tab), -1)
	}

	if _, ok := mw.tabGistList[g.id]; ok {
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

	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	mw.setupUI()
	mw.setupInteractions()
	app.SetActiveWindow(mw.window)
	mw.show()

	tab := NewTab(mw.tabWidget)
	tab.showGist(mw.tabWidget, g)
	index := mw.tabWidget.IndexOf(tab)
	currentSize := mw.tabWidget.Count()

	mw.tabWidget.SetCurrentWidget(tab)
	mw.tabWidget.SetFocus2()

	mw.tabWidget.ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("i = %d, want %d", i, index)
	})

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_W), core.Qt__ControlModifier, "", false, 1)
	mw.tabWidget.KeyPressEvent(event)

	if !called {
		t.Error("didn't close the tab")
	}
	if mw.tabWidget.Count() != currentSize-1 {
		t.Errorf("mw.tabWidget.Count() = %d, want %d", mw.tabWidget.Count(), currentSize-1)
	}
	if mw.tabWidget.IndexOf(tab) != -1 {
		t.Errorf("mw.tabWidget.IndexOf(tab) = %d, want %d", mw.tabWidget.IndexOf(tab), -1)
	}

	if _, ok := mw.tabGistList[g.id]; ok {
		t.Errorf("%s was not removed from the list", g.id)
	}
}
