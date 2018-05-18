// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arsham/gistflow/gist"
	"github.com/arsham/gistflow/qt/tab"
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
		content = "fLGLysiOuxReut\nASUonvyd"
		label   = "LpqrRCgBBYY"
		g       = getGist("uWIkJYdkFuVwYcyy", label, content)
	)

	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	initialCount := window.tabsWidget.Count()
	tab := tab.NewTab(window.tabsWidget)
	if tab == nil {
		t.Error("NewTab(window.tabsWidget) = nil, want *Tab")
		return
	}

	tab.ShowGist(window.tabsWidget, g)
	if window.tabsWidget.Count() != initialCount+1 {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), initialCount+1)
		return
	}

	index := window.tabsWidget.CurrentIndex()
	shownTab := window.tabsWidget.Widget(index)
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
	if window.tabsWidget.TabText(index) != label {
		t.Errorf("TabText(%d) = %s, want %s", index, window.tabsWidget.TabText(index), label)
	}
}

func TestTabIDFromIndex(t *testing.T) { tRunner.Run(func() { testTabIDFromIndex(t) }) }
func testTabIDFromIndex(t *testing.T) {
	var (
		id1 = "mbzsNwJS"
		id2 = "eulYvWSUHubADRV"
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
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.gistService.API = gistTs.URL
	currentIndex := window.tabsWidget.CurrentIndex()
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
		return window.tabsWidget
	})
	testSwitchTabsOnWidget(t, "GistList", func(window *MainWindow) keyReleaseEventWidget {
		return window.gistList
	})
}

func testSwitchTabsOnWidget(t *testing.T, name string, f func(*MainWindow) keyReleaseEventWidget) {
	_, window, cleanup, err := setup(t, appName, nil, 0)
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

	leftTab := tab.NewTab(window.tabsWidget)
	leftTab.ShowGist(window.tabsWidget, g1)
	rightTab := tab.NewTab(window.tabsWidget)
	rightTab.ShowGist(window.tabsWidget, g2)

	leftIndex := window.tabsWidget.IndexOf(leftTab)
	rightIndex := window.tabsWidget.IndexOf(rightTab)
	if leftIndex > rightIndex {
		rightIndex = leftIndex
	}
	window.tabsWidget.SetCurrentIndex(rightIndex)

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageUp, core.Qt__ControlModifier, -1)
	event.Simulate(w)
	if window.tabsWidget.CurrentIndex() != leftIndex {
		t.Errorf("%s: window.tabsWidget.CurrentIndex() = %d, want %d", name, window.tabsWidget.CurrentIndex(), leftIndex)
	}

	event = testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageDown, core.Qt__ControlModifier, -1)
	event.Simulate(w)
	if window.tabsWidget.CurrentIndex() != rightIndex {
		t.Errorf("%s: window.tabsWidget.CurrentIndex() = %d, want %d", name, window.tabsWidget.CurrentIndex(), rightIndex)
	}
}

func TestMovingTabs(t *testing.T) { tRunner.Run(func() { testMovingTabs(t) }) }
func testMovingTabs(t *testing.T) {
	testMovingTabsOnTabWidget(t, "TabsWidget", func(window *MainWindow) keyReleaseEventWidget {
		return window.tabsWidget
	})
	testMovingTabsOnTabWidget(t, "GistList", func(window *MainWindow) keyReleaseEventWidget {
		return window.gistList
	})
}

func testMovingTabsOnTabWidget(t *testing.T, name string, f func(*MainWindow) keyReleaseEventWidget) {
	_, window, cleanup, err := setup(t, appName, nil, 0)
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

	leftTab := tab.NewTab(window.tabsWidget)
	leftTab.ShowGist(window.tabsWidget, g1)
	rightTab := tab.NewTab(window.tabsWidget)
	rightTab.ShowGist(window.tabsWidget, g2)

	leftIndex := window.tabsWidget.IndexOf(leftTab)
	rightIndex := window.tabsWidget.IndexOf(rightTab)
	if leftIndex > rightIndex {
		// I just want to position them properly.
		leftIndex, rightIndex = rightIndex, leftIndex
		leftTab, rightTab = rightTab, leftTab
	}
	window.tabsWidget.SetCurrentIndex(rightIndex)

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageUp, core.Qt__ControlModifier+core.Qt__ShiftModifier, -1)
	event.Simulate(w)
	if window.tabsWidget.IndexOf(rightTab) != leftIndex {
		t.Errorf("%s: window.tabsWidget.IndexOf(rightTab) = %d, want %d", name, window.tabsWidget.IndexOf(rightTab), leftIndex)
	}
	if window.tabsWidget.CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Errorf("%s: focus is still on leftTab, want rightTab", name)
	}

	// now we are swichng back
	event = testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_PageDown, core.Qt__ControlModifier+core.Qt__ShiftModifier, -1)
	event.Simulate(w)
	if window.tabsWidget.IndexOf(rightTab) != rightIndex {
		t.Errorf("%s: window.tabsWidget.IndexOf(rightTab) = %d, want %d", name, window.tabsWidget.IndexOf(rightTab), rightIndex)
	}
	if window.tabsWidget.CurrentWidget().Pointer() != rightTab.Pointer() {
		t.Errorf("%s: focus is still on leftTab, want rightTab", name)
	}
}

func TestRemoveOpenTab(t *testing.T) { tRunner.Run(func() { testRemoveOpenTab(t) }) }
func testRemoveOpenTab(t *testing.T) {
	var (
		id1 = "mbzsNwJS"
		id2 = "eulYvWSUHubADRV"
		id3 = "zdXyiCAdDkG"
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
	_, window, cleanup, err := setup(t, appName, nil, 0)
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
	window.tabsWidget.TabCloseRequested(window.tabsWidget.CurrentIndex())
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

	index := window.tabsWidget.IndexOf(window.tabGistList[id1])
	if window.tabsWidget.CurrentIndex() != index {
		t.Errorf("window.tabsWidget.CurrentIndex() = %d, want %d", window.tabsWidget.CurrentIndex(), index)
	}
}

func TestWindowClickCloseTab(t *testing.T) { tRunner.Run(func() { testWindowClickCloseTab(t) }) }
func testWindowClickCloseTab(t *testing.T) {
	var called bool
	g := getGist("uWIkJYdkFuVwYcyy", "LpqrRCgBBYY", "fLGLysiOuxReut\nASUonvyd")

	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	app.SetActiveWindow(window)
	window.Show()

	tab := tab.NewTab(window.tabsWidget)
	if tab == nil {
		t.Error("NewTab(window.tabsWidget) = nil, want *Tab")
		return
	}

	tab.ShowGist(window.tabsWidget, g)
	currentSize := window.tabsWidget.Count()
	index := window.tabsWidget.IndexOf(tab)
	window.tabsWidget.ConnectTabCloseRequested(func(i int) {
		if i == index {
			called = true
			return
		}
		t.Errorf("i = %d, want %d", i, index)
	})

	window.tabsWidget.TabCloseRequested(index)

	if !called {
		t.Errorf("%s: didn't close the tab", appName)
	}
	if window.tabsWidget.Count() != currentSize-1 {
		t.Errorf("%s: window.tabsWidget.Count() = %d, want %d", appName, window.tabsWidget.Count(), currentSize-1)
	}
	if window.tabsWidget.IndexOf(tab) != -1 {
		t.Errorf("%s: window.tabsWidget.IndexOf(tab) = %d, want %d", appName, window.tabsWidget.IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.ID]; ok {
		t.Errorf("%s was not removed from the list", g.ID)
	}
}

func TestShortcutTabClose(t *testing.T) { tRunner.Run(func() { testShortcutTabClose(t) }) }
func testShortcutTabClose(t *testing.T) {
	testShortcutTabCloseWidget(t, "TabsWidget", func(window *MainWindow) keyReleaseEventWidget {
		return window.tabsWidget
	})
	testShortcutTabCloseWidget(t, "GistList", func(window *MainWindow) keyReleaseEventWidget {
		return window.gistList
	})
}

func testShortcutTabCloseWidget(t *testing.T, name string, f func(*MainWindow) keyReleaseEventWidget) {
	var called bool
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	app.SetActiveWindow(window)
	window.Show()
	w := f(window)
	g := getGist("TGtyIHIK", "hPtRE", "quIwMlPsoVaNr")

	tab := tab.NewTab(window.tabsWidget)
	tab.ShowGist(window.tabsWidget, g)
	index := window.tabsWidget.IndexOf(tab)
	currentSize := window.tabsWidget.Count()

	window.tabsWidget.SetCurrentWidget(tab)
	w.SetFocus2()

	window.tabsWidget.ConnectTabCloseRequested(func(i int) {
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
	if window.tabsWidget.Count() != currentSize-1 {
		t.Errorf("%s: window.tabsWidget.Count() = %d, want %d", name, window.tabsWidget.Count(), currentSize-1)
	}
	if window.tabsWidget.IndexOf(tab) != -1 {
		t.Errorf("%s: window.tabsWidget.IndexOf(tab) = %d, want %d", name, window.tabsWidget.IndexOf(tab), -1)
	}

	if _, ok := window.tabGistList[g.ID]; ok {
		t.Errorf("%s: %s was not removed from the list", name, g.ID)
	}
}

func TestMultipleFileGist(t *testing.T) { tRunner.Run(func() { testMultipleFileGist(t) }) }
func testMultipleFileGist(t *testing.T) {
	var (
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

	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	tab := tab.NewTab(window.tabsWidget)
	tab.ShowGist(window.tabsWidget, g)
	index := window.tabsWidget.CurrentIndex()

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
	if window.tabsWidget.TabText(index) != label1 {
		t.Errorf("TabText(%d) = %s, want %s", index, window.tabsWidget.TabText(index), label1)
	}
	file2 := tab.Files()[1]
	if file2.Content().ToPlainText() != content2 {
		t.Errorf("content2 = %s, want %s", file2.Content().ToPlainText(), content2)
	}
}

func TestCopyContents(t *testing.T) { tRunner.Run(func() { testCopyContents(t) }) }
func testCopyContents(t *testing.T) {
	var (
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
	_, window, cleanup, err := setup(t, appName, nil, 0)
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
	f1.CopyButton().Click()
	if clpText != content1 {
		t.Errorf("clpText = %s, want %s", clpText, content1)
	}

	f2 := tab.Files()[1]
	f2.CopyButton().Click()
	if clpText != content2 {
		t.Errorf("clpText = %s, want %s", clpText, content2)
	}
}

func TestDeleteFile(t *testing.T) { tRunner.Run(func() { testDeleteFile(t) }) }
func testDeleteFile(t *testing.T) {
	var (
		content1 = "WiVxf9eFeqQtdAm12wl"
		content2 = "FRPdPlkV"
		file1    = "Vea1UGy61WK4rEL"
		file2    = "lr0sO9Ep"
		called   bool
		signaled bool
		id       = "vLr0dPOjaqRsTR"
		g        = &gist.Gist{
			ID: id,
			Files: map[string]gist.File{
				file1: gist.File{Content: content1},
				file2: gist.File{Content: content2},
			},
		}
	)
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		b, err := json.Marshal(g)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	updateGistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		g := gist.Gist{}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		defer r.Body.Close()
		if err := json.Unmarshal(data, &g); err != nil {
			t.Error(err)
			return
		}
		if _, ok := g.Files[file1]; !ok {
			t.Errorf("%s was not in the request: %v", file1, g.Files)
		}
		if _, ok := g.Files[file2]; ok {
			t.Errorf("%s was in the request: %v", file2, g.Files)
		}
		w.Write([]byte("{}"))
	}))
	defer updateGistTs.Close()

	window.gistService.API = gistTs.URL
	window.openGist(id)
	if !called {
		t.Error("didn't send the open request")
	}
	tab := window.tabGistList[id]
	tab.ConnectFileDeleted(func(name string) {
		signaled = true
		if name != file1 {
			t.Errorf("name = %s, want %s", name, file1)
		}
	})
	g.URL = updateGistTs.URL

	called = false
	tab.DeleteFile(g, file1)
	if !called {
		t.Error("didn't send the remove request")
	}
	if !signaled {
		t.Error("didn't send the deleted signal")
	}
}

func TestDeleteGist(t *testing.T) { tRunner.Run(func() { testDeleteGist(t) }) }
func testDeleteGist(t *testing.T) {
	var (
		id        = "sRtW06Rs8u"
		content   = "cF9KLDx46o"
		file      = "QEcsU5y7"
		called    bool
		destroyed bool
		g         = &gist.Gist{
			ID: id,
			Files: map[string]gist.File{
				file: gist.File{Content: content},
			},
		}
		r = gist.Gist{
			ID: id,
			Files: map[string]gist.File{
				file: gist.File{Content: content},
			},
		}
	)
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		b, err := json.Marshal(g)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	deleteGistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		req := strings.Split(r.URL.Path, "/")
		reqID := req[len(req)-1]
		if reqID != id {
			t.Errorf("reqID = %s, want %s", reqID, id)
		}
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("{}"))
	}))
	defer deleteGistTs.Close()

	window.gistService.API = gistTs.URL
	window.openGist(id)
	if !called {
		t.Error("didn't send the open request")
	}
	tab := window.tabGistList[id]
	window.gistService.API = deleteGistTs.URL
	window.searchbox.Add(r)
	window.gistList.Add(r)

	if !window.gistList.HasID(id) {
		t.Errorf("%s is not in gistList", file)
	}
	if !window.searchbox.HasID(g.ID) {
		t.Errorf("%s is not in searchbox", g.ID)
	}

	called = false
	tab.ConnectDestroyed(func(*core.QObject) {
		destroyed = true
	})

	tab.DeleteGist(g)
	if !called {
		t.Error("didn't send the remove request")
	}
	if window.gistList.HasID(id) {
		t.Errorf("%s is still in gistList", file)
	}
	if window.gistList.Count() > 0 {
		t.Errorf("%s is still in gistList", file)
	}
	if window.searchbox.HasID(g.ID) {
		t.Errorf("%s is still in searchbox", g.ID)
	}

	if _, ok := window.tabGistList[id]; ok {
		t.Errorf("window.tabGistList still contains %s", id)
	}
	if !destroyed {
		t.Error("Tab is not removed")
	}
}
