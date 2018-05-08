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
	"github.com/arsham/gisty/interface/tab"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/testlib"
	"github.com/therecipe/qt/widgets"
)

func getGist(id, label, content string) *gist.Gist {
	return &gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			label: gist.File{Content: content},
		},
	}
}

func TestTabCreation(t *testing.T) { tRunner.Run(func() { testTabCreation(t) }) }
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

	if window.TabsWidget().Count() != 1 {
		t.Errorf("window.TabsWidget().Count() = %d, want 1", window.TabsWidget().Count())
		return
	}

	tab := tab.NewTab(window.TabsWidget())
	if tab == nil {
		t.Error("NewTab(window.TabsWidget()) = nil, want *Tab")
		return
	}

	tab.ShowGist(window.TabsWidget(), g)

	index := window.TabsWidget().CurrentIndex()
	shownTab := window.TabsWidget().Widget(index)
	if shownTab.Pointer() != tab.Pointer() {
		t.Errorf("shownTab.Pointer() = %v, want %v", shownTab.Pointer(), tab.Pointer())
	}

	if tab.Files() == nil {
		t.Error("tab.Files() = nil")
		return
	}
	if len(tab.Files()) != 1 {
		t.Errorf("len(tab.Files()) = %d, want 1", len(tab.Files()))
		return
	}

	file := tab.Files()[0]
	if file.Content().ToPlainText() != content {
		t.Errorf("content = %s, want %s", file.Content().ToPlainText(), content)
	}
	if window.TabsWidget().TabText(index) != label {
		t.Errorf("TabText(%d) = %s, want %s", index, window.TabsWidget().TabText(index), label)
	}
}

func TestTabIDFromIndex(t *testing.T) { tRunner.Run(func() { testTabIDFromIndex(t) }) }
func testTabIDFromIndex(t *testing.T) {
	var (
		name = "test"
		id1  = "mbzsNwJS"
		id2  = "eulYvWSUHubADRV"
	)

	gres := gist.Gist{
		Files: map[string]gist.File{
			"GqrqZkTNpw": gist.File{Content: "uMWDmwSvLlqtFXZUX"},
		},
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

type keyReleaseEventWidget interface {
	widgets.QWidget_ITF
	KeyReleaseEvent(gui.QKeyEvent_ITF)
	SetFocus2()
}

func TestSwitchTabs(t *testing.T) { tRunner.Run(func() { testSwitchTabs(t) }) }
func testSwitchTabs(t *testing.T) {
	testSwitchTabsOnWidget(t, "TabsWidget", func(window *MainWindow) keyReleaseEventWidget {
		return window.TabsWidget()
	})
	testSwitchTabsOnWidget(t, "GistList", func(window *MainWindow) keyReleaseEventWidget {
		return window.GistList()
	})
}

func testSwitchTabsOnWidget(t *testing.T, name string, f func(*MainWindow) keyReleaseEventWidget) {
	_, window, cleanup, err := setup(t, "test", nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	app.SetActiveWindow(window)
	window.Show()
	w := f(window)
	g1 := getGist("uWIkJYdkFuVwYcyy", "LpqrRCgBBYY", "fLGLysiOuxReut\nASUonvyd")
	g2 := getGist("FJsPzPqhI", "bsDmGRE", "KuiIIVYnCKycPPkXLibh")

	leftTab := tab.NewTab(window.TabsWidget())
	leftTab.ShowGist(window.TabsWidget(), g1)
	rightTab := tab.NewTab(window.TabsWidget())
	rightTab.ShowGist(window.TabsWidget(), g2)

	leftIndex := window.TabsWidget().IndexOf(leftTab)
	rightIndex := window.TabsWidget().IndexOf(rightTab)
	if leftIndex > rightIndex {
		rightIndex = leftIndex
	}
	window.TabsWidget().SetCurrentIndex(rightIndex)

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageUp, core.Qt__ControlModifier, -1)
	event.Simulate(w)
	if window.TabsWidget().CurrentIndex() != leftIndex {
		t.Errorf("%s: window.TabsWidget().CurrentIndex() = %d, want %d", name, window.TabsWidget().CurrentIndex(), leftIndex)
	}

	event = testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageDown, core.Qt__ControlModifier, -1)
	event.Simulate(w)
	if window.TabsWidget().CurrentIndex() != rightIndex {
		t.Errorf("%s: window.TabsWidget().CurrentIndex() = %d, want %d", name, window.TabsWidget().CurrentIndex(), rightIndex)
	}
}

func TestMovingTabs(t *testing.T) { tRunner.Run(func() { testMovingTabs(t) }) }
func testMovingTabs(t *testing.T) {
	testMovingTabsOnTabWidget(t, "TabsWidget", func(window *MainWindow) keyReleaseEventWidget {
		return window.TabsWidget()
	})
	testMovingTabsOnTabWidget(t, "GistList", func(window *MainWindow) keyReleaseEventWidget {
		return window.GistList()
	})
}

func testMovingTabsOnTabWidget(t *testing.T, name string, f func(*MainWindow) keyReleaseEventWidget) {
	_, window, cleanup, err := setup(t, "test", nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	app.SetActiveWindow(window)
	window.Show()
	w := f(window)
	g1 := getGist("uWIkJYdkFuVwYcyy", "LpqrRCgBBYY", "fLGLysiOuxReut\nASUonvyd")
	g2 := getGist("FJsPzPqhI", "bsDmGRE", "KuiIIVYnCKycPPkXLibh")

	leftTab := tab.NewTab(window.TabsWidget())
	leftTab.ShowGist(window.TabsWidget(), g1)
	rightTab := tab.NewTab(window.TabsWidget())
	rightTab.ShowGist(window.TabsWidget(), g2)

	leftIndex := window.TabsWidget().IndexOf(leftTab)
	rightIndex := window.TabsWidget().IndexOf(rightTab)
	if leftIndex > rightIndex {
		// I just want to position them properly.
		leftIndex, rightIndex = rightIndex, leftIndex
		leftTab, rightTab = rightTab, leftTab
	}
	window.TabsWidget().SetCurrentIndex(rightIndex)

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageUp, core.Qt__ControlModifier+core.Qt__ShiftModifier, -1)
	event.Simulate(w)
	if window.TabsWidget().IndexOf(rightTab) != leftIndex {
		t.Errorf("%s: window.TabsWidget().IndexOf(rightTab) = %d, want %d", name, window.TabsWidget().IndexOf(rightTab), leftIndex)
	}
	if window.TabsWidget().CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Errorf("%s: focus is still on leftTab, want rightTab", name)
	}

	// now we are swichng back
	event = testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageDown, core.Qt__ControlModifier+core.Qt__ShiftModifier, -1)
	event.Simulate(w)
	if window.TabsWidget().IndexOf(rightTab) != rightIndex {
		t.Errorf("%s: window.TabsWidget().IndexOf(rightTab) = %d, want %d", name, window.TabsWidget().IndexOf(rightTab), rightIndex)
	}
	if window.TabsWidget().CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Errorf("%s: focus is still on leftTab, want rightTab", name)
	}
}

func TestRemoveOpenTab(t *testing.T) { tRunner.Run(func() { testRemoveOpenTab(t) }) }
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

func TestWindowClickCloseTab(t *testing.T) { tRunner.Run(func() { testWindowClickCloseTab(t) }) }
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
	app.SetActiveWindow(window)
	window.Show()

	tab := tab.NewTab(window.TabsWidget())
	if tab == nil {
		t.Error("NewTab(window.TabsWidget()) = nil, want *Tab")
		return
	}

	tab.ShowGist(window.TabsWidget(), g)
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
		t.Errorf("%s: didn't close the tab", name)
	}
	if window.TabsWidget().Count() != currentSize-1 {
		t.Errorf("%s: window.TabsWidget().Count() = %d, want %d", name, window.TabsWidget().Count(), currentSize-1)
	}
	if window.TabsWidget().IndexOf(tab) != -1 {
		t.Errorf("%s: window.TabsWidget().IndexOf(tab) = %d, want %d", name, window.TabsWidget().IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.ID]; ok {
		t.Errorf("%s was not removed from the list", g.ID)
	}
}

func TestShortcutTabClose(t *testing.T) { tRunner.Run(func() { testShortcutTabClose(t) }) }
func testShortcutTabClose(t *testing.T) {
	testShortcutTabCloseWidget(t, "TabsWidget", func(window *MainWindow) keyReleaseEventWidget {
		return window.TabsWidget()
	})
	testShortcutTabCloseWidget(t, "GistList", func(window *MainWindow) keyReleaseEventWidget {
		return window.GistList()
	})
}

func testShortcutTabCloseWidget(t *testing.T, name string, f func(*MainWindow) keyReleaseEventWidget) {
	var called bool
	_, window, cleanup, err := setup(t, "test", nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	app.SetActiveWindow(window)
	window.Show()
	w := f(window)
	g := getGist("TGtyIHIK", "hPtRE", "quIwMlPsoVaNr")

	tab := tab.NewTab(window.TabsWidget())
	tab.ShowGist(window.TabsWidget(), g)
	index := window.TabsWidget().IndexOf(tab)
	currentSize := window.TabsWidget().Count()

	window.TabsWidget().SetCurrentWidget(tab)
	w.SetFocus2()

	window.TabsWidget().ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("%s: i = %d, want %d", name, i, index)
	})

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_W, core.Qt__ControlModifier, -1)
	event.Simulate(w)

	if !called {
		t.Errorf("%s: didn't close the tab", name)
	}
	if window.TabsWidget().Count() != currentSize-1 {
		t.Errorf("%s: window.TabsWidget().Count() = %d, want %d", name, window.TabsWidget().Count(), currentSize-1)
	}
	if window.TabsWidget().IndexOf(tab) != -1 {
		t.Errorf("%s: window.TabsWidget().IndexOf(tab) = %d, want %d", name, window.TabsWidget().IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.ID]; ok {
		t.Errorf("%s: %s was not removed from the list", name, g.ID)
	}
}

func TestMultipleFileGist(t *testing.T) { tRunner.Run(func() { testMultipleFileGist(t) }) }
func testMultipleFileGist(t *testing.T) {
	var (
		name     = "test"
		content1 = "fLGLysiOuxReutASUonvyd"
		content2 = "zXLpDTgdCZtmxiZDqDQJAcEZ"
		label1   = "LpqrRCgBBYY"
		label2   = "ORfVfQPH"
		g        = &gist.Gist{
			ID: "vgCWaGVbqWtHaH",
			Files: map[string]gist.File{
				label1: gist.File{Content: content1},
				label2: gist.File{Content: content2},
			},
		}
	)

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	tab := tab.NewTab(window.TabsWidget())
	tab.ShowGist(window.TabsWidget(), g)
	index := window.TabsWidget().CurrentIndex()

	if tab.Files() == nil {
		t.Error("tab.Files() = nil")
		return
	}
	if len(tab.Files()) != 2 {
		t.Errorf("len(tab.Files()) = %d, want 2", len(tab.Files()))
		return
	}

	file1 := tab.Files()[0]
	if file1.Content().ToPlainText() != content1 {
		t.Errorf("content1 = %s, want %s", file1.Content().ToPlainText(), content1)
	}
	if window.TabsWidget().TabText(index) != label1 {
		t.Errorf("TabText(%d) = %s, want %s", index, window.TabsWidget().TabText(index), label1)
	}
	file2 := tab.Files()[1]
	if file2.Content().ToPlainText() != content2 {
		t.Errorf("content2 = %s, want %s", file2.Content().ToPlainText(), content2)
	}
}

func TestCopyContents(t *testing.T) { tRunner.Run(func() { testCopyContents(t) }) }
func testCopyContents(t *testing.T) {
	var (
		name     = "test"
		content1 = "fLGLysiOuxReutASUonvyd"
		content2 = "zXLpDTgdCZtmxiZDqDQJAcEZ"
		label1   = "LpqrRCgBBYY"
		label2   = "ORfVfQPH"
		g        = &gist.Gist{
			ID: "vgCWaGVbqWtHaH",
			Files: map[string]gist.File{
				label1: gist.File{Content: content1},
				label2: gist.File{Content: content2},
			},
		}
		clpText string
	)
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(g)
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

	window.gistService.API = gistTs.URL
	window.clipboard = func() clipboard {
		return &fakeClipboard{
			textFunc: func(text string, mode gui.QClipboard__Mode) {
				clpText = text
			},
		}
	}
	id := "XxAV5V0GAbN9j1cha8"
	window.openGist(id)
	tab := window.tabGistList[id]

	f1 := tab.Files()[0]
	f1.Copy().Click()
	if clpText != content1 {
		t.Errorf("clpText = %s, want %s", clpText, content1)
	}

	f2 := tab.Files()[1]
	f2.Copy().Click()
	if clpText != content2 {
		t.Errorf("clpText = %s, want %s", clpText, content2)
	}
}
