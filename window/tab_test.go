// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

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

func getGist(id, label, content string) *gist.Gist {
	return &gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			label: gist.File{Content: content},
		},
	}
}
func testTabCreation(t *testing.T) {
	var (
		name    = "test"
		content = "fLGLysiOuxReut\nASUonvyd"
		label   = "LpqrRCgBBYY"
		g       = getGist("uWIkJYdkFuVwYcyy", label, content)
	)

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()

	if window.TabsWidget().Count() != 1 {
		t.Errorf("window.TabsWidget().Count() = %d, want 1", window.TabsWidget().Count())
		return
	}

	tab := NewTab(window.TabsWidget())
	if tab == nil {
		t.Error("NewTab(window.TabsWidget()) = nil, want *Tab")
		return
	}

	tab.showGist(window.TabsWidget(), g)

	index := window.TabsWidget().CurrentIndex()
	shownTab := window.TabsWidget().Widget(index)
	if shownTab.Pointer() != tab.Pointer() {
		t.Errorf("shownTab.Pointer() = %v, want %v", shownTab.Pointer(), tab.Pointer())
	}
	if tab.Editor().ToPlainText() != content {
		t.Errorf("content = %s, want %s", tab.Editor().ToPlainText(), content)
	}
	if window.TabsWidget().TabText(index) != label {
		t.Errorf("TabText(%d) = %s, want %s", index, window.TabsWidget().TabText(index), label)
	}
}

func testTabIDFromIndex(t *testing.T) {
	var (
		name = "test"
		id1  = "mbzsNwJS"
		id2  = "eulYvWSUHubADRV"
	)

	files := map[string]gist.File{
		"GqrqZkTNpw": gist.File{Content: "uMWDmwSvLlqtFXZUX"},
	}
	gres := gist.Gist{
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
	currentIndex := window.TabsWidget().CurrentIndex()
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
	g1 := getGist("uWIkJYdkFuVwYcyy", "LpqrRCgBBYY", "fLGLysiOuxReut\nASUonvyd")
	g2 := getGist("FJsPzPqhI", "bsDmGRE", "KuiIIVYnCKycPPkXLibh")

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	app.SetActiveWindow(window)
	window.show()

	leftTab := NewTab(window.TabsWidget())
	leftTab.showGist(window.TabsWidget(), g1)
	rightTab := NewTab(window.TabsWidget())
	rightTab.showGist(window.TabsWidget(), g2)

	leftIndex := window.TabsWidget().IndexOf(leftTab)
	rightIndex := window.TabsWidget().IndexOf(rightTab)
	if leftIndex > rightIndex {
		rightIndex = leftIndex
	}
	window.TabsWidget().SetCurrentIndex(rightIndex)
	window.TabsWidget().SetFocus2()

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageUp), core.Qt__ControlModifier, "", false, 1)
	window.TabsWidget().KeyPressEvent(event)

	event = gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageDown), core.Qt__ControlModifier, "", false, 1)
	window.TabsWidget().KeyPressEvent(event)
	if window.TabsWidget().CurrentIndex() != rightIndex {
		t.Errorf("window.TabsWidget().CurrentIndex() = %d, want %d", window.TabsWidget().CurrentIndex(), rightIndex)
	}
}

func testMovingTabs(t *testing.T) {
	var name = "test"
	g1 := getGist("uWIkJYdkFuVwYcyy", "LpqrRCgBBYY", "fLGLysiOuxReut\nASUonvyd")
	g2 := getGist("FJsPzPqhI", "bsDmGRE", "KuiIIVYnCKycPPkXLibh")

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	app.SetActiveWindow(window)
	window.show()

	leftTab := NewTab(window.TabsWidget())
	leftTab.showGist(window.TabsWidget(), g1)
	rightTab := NewTab(window.TabsWidget())
	rightTab.showGist(window.TabsWidget(), g2)

	leftIndex := window.TabsWidget().IndexOf(leftTab)
	rightIndex := window.TabsWidget().IndexOf(rightTab)
	if leftIndex > rightIndex {
		// I just want to position them properly.
		leftIndex, rightIndex = rightIndex, leftIndex
		leftTab, rightTab = rightTab, leftTab
	}
	window.TabsWidget().SetCurrentIndex(rightIndex)
	window.TabsWidget().SetFocus2()

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageUp), core.Qt__ControlModifier+core.Qt__ShiftModifier, "", false, 1)
	window.TabsWidget().KeyPressEvent(event)
	if window.TabsWidget().IndexOf(rightTab) != leftIndex {
		t.Errorf("window.TabsWidget().IndexOf(rightTab) = %d, want %d", window.TabsWidget().IndexOf(rightTab), leftIndex)
	}
	if window.TabsWidget().CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Error("focus is still on leftTab, want rightTab")
	}

	// now we are swichng back
	event = gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_PageDown), core.Qt__ControlModifier+core.Qt__ShiftModifier, "", false, 1)
	window.TabsWidget().KeyPressEvent(event)
	if window.TabsWidget().IndexOf(rightTab) != rightIndex {
		t.Errorf("window.TabsWidget().IndexOf(rightTab) = %d, want %d", window.TabsWidget().IndexOf(rightTab), rightIndex)
	}
	if window.TabsWidget().CurrentWidget().Pointer() != rightTab.Pointer() {
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

	files := map[string]gist.File{
		"GqrqZkTNpw": gist.File{Content: "uMWDmwSvLlqtFXZUX"},
	}
	gres := gist.Gist{
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

	currentLen := len(window.tabGistList)
	window.openGist(id1)
	if len(window.tabGistList) != currentLen+1 {
		t.Errorf("len(window.tabGistList) = %d, want %d", len(window.tabGistList), currentLen+1)
	}
	if _, ok := window.tabGistList[id1]; !ok {
		t.Errorf("%s not found in window.tabGistList", id1)
	}
	window.TabsWidget().TabCloseRequested(window.TabsWidget().CurrentIndex())
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

	index := window.TabsWidget().IndexOf(window.tabGistList[id1])
	if window.TabsWidget().CurrentIndex() != index {
		t.Errorf("window.TabsWidget().CurrentIndex() = %d, want %d", window.TabsWidget().CurrentIndex(), index)
	}
}

func testWindowClickCloseTab(t *testing.T) {
	var (
		name   = "test"
		called bool
	)

	g := getGist("uWIkJYdkFuVwYcyy", "LpqrRCgBBYY", "fLGLysiOuxReut\nASUonvyd")

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	app.SetActiveWindow(window)
	window.show()

	tab := NewTab(window.TabsWidget())
	if tab == nil {
		t.Error("NewTab(window.TabsWidget()) = nil, want *Tab")
		return
	}

	tab.showGist(window.TabsWidget(), g)
	currentSize := window.TabsWidget().Count()
	index := window.TabsWidget().IndexOf(tab)
	window.TabsWidget().ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("i = %d, want %d", i, index)
	})

	window.TabsWidget().TabCloseRequested(index)

	if !called {
		t.Error("didn't close the tab")
	}
	if window.TabsWidget().Count() != currentSize-1 {
		t.Errorf("window.TabsWidget().Count() = %d, want %d", window.TabsWidget().Count(), currentSize-1)
	}
	if window.TabsWidget().IndexOf(tab) != -1 {
		t.Errorf("window.TabsWidget().IndexOf(tab) = %d, want %d", window.TabsWidget().IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.ID]; ok {
		t.Errorf("%s was not removed from the list", g.ID)
	}
}

func testShortcutTabClose(t *testing.T) {
	var (
		name   = "test"
		called bool
	)
	g := getGist("TGtyIHIK", "hPtRE", "quIwMlPsoVaNr")

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setupUI()
	app.SetActiveWindow(window)
	window.show()

	tab := NewTab(window.TabsWidget())
	tab.showGist(window.TabsWidget(), g)
	index := window.TabsWidget().IndexOf(tab)
	currentSize := window.TabsWidget().Count()

	window.TabsWidget().SetCurrentWidget(tab)
	window.TabsWidget().SetFocus2()

	window.TabsWidget().ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("i = %d, want %d", i, index)
	})

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(core.Qt__Key_W), core.Qt__ControlModifier, "", false, 1)
	window.TabsWidget().KeyPressEvent(event)

	if !called {
		t.Error("didn't close the tab")
	}
	if window.TabsWidget().Count() != currentSize-1 {
		t.Errorf("window.TabsWidget().Count() = %d, want %d", window.TabsWidget().Count(), currentSize-1)
	}
	if window.TabsWidget().IndexOf(tab) != -1 {
		t.Errorf("window.TabsWidget().IndexOf(tab) = %d, want %d", window.TabsWidget().IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.ID]; ok {
		t.Errorf("%s was not removed from the list", g.ID)
	}
}
