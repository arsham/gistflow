// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/arsham/gistflow/gist"
	"github.com/arsham/gistflow/qt/conf"
	"github.com/arsham/gistflow/qt/gistlist"
	"github.com/arsham/gistflow/qt/searchbox"
	"github.com/arsham/gistflow/qt/tab"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/testlib"
	"github.com/therecipe/qt/widgets"
)

// testing vanilla setup.
func TestWindowStartupWidgets(t *testing.T) { tRunner.Run(func() { testWindowStartupWidgets(t) }) }
func testWindowStartupWidgets(t *testing.T) {
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.clipboard = func() clipboard {
		return &fakeClipboard{textFunc: func(text string, mode gui.QClipboard__Mode) {}}
	}

	oldLogger := window.logger
	window.logger = nil
	window.setupUI()
	if window.logger == nil {
		t.Error("window.logger = nil, want boxLogger")
	}
	window.logger = oldLogger

	if window.tabsWidget == nil {
		t.Error("window.tabsWidget = nil, want *widgets.QTabWidget")
		return
	}
	if !window.tabsWidget.IsMovable() {
		t.Error("window.tabsWidget is not movable")
	}
	if window.menubar == nil {
		t.Error("window.menubar = nil, want *widgets.QMenuBar")
		return
	}
	if window.statusArea == nil {
		t.Error("window.statusArea = nil, want *widgets.QStatusBar")
	}
	if window.StatusBar().Pointer() != window.statusArea.Pointer() {
		t.Errorf("window.StatusBar().Pointer() = %v, want %v", window.StatusBar().Pointer(), window.statusArea.Pointer())
	}
	if window.dockWidget == nil {
		t.Error("window.dockWidget = nil, want *widgets.QDockWidget")
	}
	if window.tabGistList == nil {
		t.Error("window.tabGistList = nil, want []*tabGist")
	}
	if window.gistList == nil {
		t.Error("window.gistList = nil, want *gistlist.Container")
	}
	if window.toolBar == nil {
		t.Error("window.toolBar = nil, want *toolbar.Toolbar")
	}
	if window.icon == nil {
		t.Error("window.icon = nil, want *gui.QIcon")
	}
	if window.sysTray == nil {
		t.Error("window.sysTray = nil, want *widgets.QSystemTrayIcon")
		return
	}
	if window.sysTray.Icon() == nil {
		t.Error("window.sysTray.Icon() = nil, want *gui.QIcon")
	}
	if window.sysTray.Icon().Name() != window.icon.Name() {
		t.Errorf("window.sysTray.Icon().Name() = %v, want %v", window.sysTray.Icon().Name(), window.icon.Name())
	}
	if window.sysTray.ContextMenu().Pointer() != window.menubar.Options().Pointer() {
		t.Errorf("window.sysTray.ContextMenu().Pointer() = %v, want %v",
			window.sysTray.ContextMenu().Pointer(),
			window.menubar.Options().Pointer(),
		)
	}
	if window.clipboard() == nil {
		t.Error("window.clipboard() = nil, want *gui.QClipboard")
	}
	if window.searchbox == nil {
		t.Error("window.searchbox = nil, want *searchbox.Dialog")
	}
	if window.gistService.Logger == nil {
		t.Error("window.gistService.Logger is not assigned")
		return
	}
	if window.gistService.CacheDir == "" {
		t.Error("gistService.CacheDir wasn't initialised")
	}
}

func TestPopulateError(t *testing.T) { tRunner.Run(func() { testPopulateError(t) }) }
func testPopulateError(t *testing.T) {
	called := make(chan struct{})
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	window.gistService.Logger = nil
	defer cleanup()

	window.logger = &logger{
		errorFunc: func(str string) {
			close(called)
		},
	}
	window.populate()
	model := window.searchbox.Model()
	index := core.NewQModelIndex()
	if c := model.RowCount(index); c != 0 {
		t.Errorf("model.RowCount() = %d, want 0", c)
	}
	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Error("expected an error, didn't register the error")
	}
}

func TestPopulate(t *testing.T) { tRunner.Run(func() { testPopulate(t) }) }
func testPopulate(t *testing.T) {
	size := 5
	gres := gist.Gist{
		ID:          "QXhJNchXAK",
		Description: "kfxLTwoCOkqEuPlp",
	}
	ts, window, cleanup, err := setup(t, appName, []gist.Gist{gres}, size)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	gres.URL = fmt.Sprintf("%s/gists/%s", ts.URL, gres.ID)
	window.populate()

	c := gistlist.NewContainerFromPointer(window.gistList.Pointer())
	if c.Description(0) != gres.Description {
		t.Errorf("c.Description(0) = %s, want %s", c.Description(0), gres.Description)
	}
	if c.ID(0) != gres.ID {
		t.Errorf("c.ID(0) = %s, want %s", c.ID(0), gres.ID)
	}

	s := searchbox.NewDialogFromPointer(window.searchbox.Pointer())
	if s.Description(0) != gres.Description {
		t.Errorf("s.Description(0) = %s, want %s", s.Description(0), gres.Description)
	}
	if s.ID(0) != gres.ID {
		t.Errorf("s.ID(0) = %s, want %s", s.ID(0), gres.ID)
	}
}

func TestLoadingGeometry(t *testing.T) { tRunner.Run(func() { testLoadingGeometry(t) }) }
func testLoadingGeometry(t *testing.T) {
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	x, y, w, h := 400, 500, 600, 700
	tmpObj := widgets.NewQWidget(nil, 0)
	tmpObj.SetGeometry2(x, y, w, h)
	settings, cleanup2 := testSettings(appName)
	defer cleanup2()
	size := tmpObj.SaveGeometry()
	settings.SetValue(mainWindowGeometry, core.NewQVariant15(size))
	settings.SetValue(conf.AccessToken, core.NewQVariant17("zYhdGyNiTRuzxcW2j"))
	settings.SetValue(conf.Username, core.NewQVariant17("fB3RZg"))
	settings.Sync()

	window.settings, err = conf.New(appName)
	if err != nil {
		t.Errorf("getting settings: %v", err)
		return
	}
	window.lastGeometry()
	geometry := window.Geometry()
	check := func(name string, size, with int) {
		if size != with {
			t.Errorf("%s = %d, want %d", name, size, with)
		}
	}
	check("loading: geometry.X()", geometry.X(), x)
	check("loading: geometry.Y()", geometry.Y(), y)
	check("loading: geometry.Width()", geometry.Width(), w)
	check("loading: geometry.Height()", geometry.Height(), h)

	x, y, w, h = 500, 600, 700, 800
	tmpObj = widgets.NewQWidget(nil, 0)
	tmpObj.SetGeometry2(x, y, w, h)
	newGeometry := tmpObj.SaveGeometry()
	window.RestoreGeometry(newGeometry)
	check("to make sure: geometry.X()", window.Geometry().X(), x)
	check("to make sure: geometry.Y()", window.Geometry().Y(), y)
	check("to make sure: geometry.Width()", window.Geometry().Width(), w)
	check("to make sure: geometry.Height()", window.Geometry().Height(), h)
}

func TestViewGist(t *testing.T) { tRunner.Run(func() { testViewGist(t) }) }
func testViewGist(t *testing.T) {
	var (
		id       = "uWIkJYdkFuVwYcyy"
		badID    = "kJuZxkDCBp"
		fileName = "LpqrRCgBBYY"
		content  = "fLGLysiOuxReut\nASUonvyd"
	)

	files := map[string]gist.File{
		fileName: gist.File{Content: content},
	}
	gres := gist.Gist{
		Files: files,
	}
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, id) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// check the URL
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

	if window.tabsWidget.Count() != 0 {
		t.Errorf("window.tabsWidget.Count() = %d, want 0", window.tabsWidget.Count())
		return
	}
	if err := window.openGist(badID); err == nil {
		t.Errorf("window.openGist(%s) = nil, want error", badID)
	}
	if err := window.openGist(id); err != nil {
		t.Errorf("window.openGist(%s) = %s, want nil", id, err)
	}

	newIndex := 1
	if window.tabsWidget.Count() != newIndex {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), newIndex)
		return
	}

	index := window.tabsWidget.CurrentIndex()
	tab := tab.NewTabFromPointer(window.tabsWidget.Widget(index).Pointer())
	if tab.Files() == nil {
		t.Error("tab.Files() = nil")
		return
	}
	if len(tab.Files()) == 0 {
		t.Error("len(tab.Files()) = 0")
		return
	}

	guts := tab.Files()[0].Content()
	if guts.ToPlainText() != content {
		t.Errorf("content = %s, want %s", guts.ToPlainText(), content)
	}
	if window.tabsWidget.TabText(index) != fileName {
		t.Errorf("TabText(%d) = %s, want %s", index, window.tabsWidget.TabText(index), fileName)
	}
}

func TestClickViewGist(t *testing.T) { tRunner.Run(func() { testClickViewGist(t) }) }
func testClickViewGist(t *testing.T) {
	var (
		called    bool
		errCalled bool
	)
	gres := gist.Gist{
		ID:          "QXhJNchXAK",
		Description: "kfxLTwoCOkqEuPlp",
	}
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, gres.ID) {
			t.Errorf("r.URL.Path = %s, want %s in it", r.URL.Path, gres.ID)
			return
		}
		// check the URL
		called = true
		b, err := json.Marshal(gres)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()

	gres.URL = fmt.Sprintf("%s/gists/%s", gistTs.URL, gres.ID)
	_, window, cleanup, err := setup(t, appName, []gist.Gist{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.populate()
	window.gistService.API = gistTs.URL
	app.SetActiveWindow(window)

	window.logger = &logger{
		errorFunc:   func(str string) { errCalled = true },
		warningFunc: func(str string) { errCalled = true },
	}

	// with no selection, it should error because there is no item selectedisd,
	// hence no id.
	event := testlib.NewQTestEventList()
	event.AddKeyRelease(core.Qt__Key_Down, core.Qt__NoModifier, -1)
	event.AddKeyRelease(core.Qt__Key_Enter, core.Qt__NoModifier, -1)
	event.Simulate(window.gistList)
	if !errCalled {
		t.Error("didn't show error")
	}

	for _, key := range []core.Qt__Key{core.Qt__Key_Enter, core.Qt__Key_Return} {
		called = false
		window.gistList.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(core.Qt__Key_Down, core.Qt__NoModifier, -1)
		event.AddKeyRelease(key, core.Qt__NoModifier, -1)
		event.Simulate(window.gistList)

		if !called {
			t.Error("didn't call for gist")
		}
		delete(window.tabGistList, gres.ID)
		os.RemoveAll(window.gistService.CacheDir)
		os.Mkdir(window.gistService.CacheDir, 0777)
	}
}

func TestOpeningGistTwice(t *testing.T) { tRunner.Run(func() { testOpeningGistTwice(t) }) }
func testOpeningGistTwice(t *testing.T) {
	var (
		id1 = "mbzsNwJS"
		id2 = "eulYvWSUHubADRV"
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

	startingSize := window.tabsWidget.Count()
	startingIndex := window.tabsWidget.CurrentIndex()
	window.openGist(id1)
	if window.tabsWidget.Count() != startingSize+1 {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), startingSize+1)
	}
	if window.tabsWidget.CurrentIndex() != startingIndex+1 {
		t.Errorf("window.tabsWidget.CurrentIndex() = %d, want %d", window.tabsWidget.CurrentIndex(), startingIndex+1)
	}

	id1Index := window.tabsWidget.CurrentIndex()
	window.openGist(id2)
	if window.tabsWidget.Count() != startingSize+2 {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), startingSize+2)
	}
	if window.tabsWidget.CurrentIndex() != startingIndex+2 {
		t.Errorf("window.tabsWidget.CurrentIndex() = %d, want %d", window.tabsWidget.CurrentIndex(), startingIndex+2)
	}

	window.openGist(id1)
	if window.tabsWidget.Count() != startingSize+2 {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), startingSize+2)
	}
	if window.tabsWidget.CurrentIndex() != id1Index {
		t.Errorf("window.tabsWidget.CurrentIndex() = %d, want %d", window.tabsWidget.CurrentIndex(), id1Index)
	}
}

func TestToggle(t *testing.T) { tRunner.Run(func() { testToggle(t) }) }
func testToggle(t *testing.T) {
	window := NewMainWindow(nil, 0)
	window.name = appName
	defer window.Hide()

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
}

func TestCopyURL(t *testing.T) { tRunner.Run(func() { testCopyURL(t) }) }
func testCopyURL(t *testing.T) {
	var (
		id1      = "wqWKsfoQevEbGjhmz"
		content  = "AFMQydAKTiJLa"
		id2      = "yuaosJCTsGUqEldvigi"
		api, url string
		clpText  string
	)

	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = fmt.Sprintf("%s%s", api, r.URL.Path)
		gres := gist.Gist{
			HTMLURL: url,
			Files: map[string]gist.File{
				"vtsmQN": gist.File{Content: content},
			},
		}
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
	window.clipboard = func() clipboard {
		return &fakeClipboard{
			textFunc: func(text string, mode gui.QClipboard__Mode) {
				clpText = text
			},
		}
	}
	api = gistTs.URL
	c := window.menubar.Actions().CopyURL
	if err := window.openGist(id1); err != nil {
		t.Errorf("window.openGist(%s) = %v, want nil", id1, err)
	}
	url1 := url
	if err := window.openGist(id2); err != nil {
		t.Errorf("window.openGist(%s) = %v, want nil", id2, err)
	}

	url2 := url
	tab1, tab2 := window.tabGistList[id1], window.tabGistList[id2]
	window.tabsWidget.SetCurrentWidget(tab1)

	c.Trigger()
	if clpText != url1 {
		t.Errorf("clpText = `%s`, want `%s`", clpText, url1)
	}

	window.tabsWidget.SetCurrentWidget(tab2)
	c.Trigger()
	if clpText != url2 {
		t.Errorf("clpText = `%s`, want `%s`", clpText, url2)
	}
}

func TestEmptyDescription(t *testing.T) { tRunner.Run(func() { testEmptyDescription(t) }) }
func testEmptyDescription(t *testing.T) {
	var (
		content  = "CNF5EmQJxiGvzwedbmTME3p0Y"
		fileName = "84nkJJG0"
	)
	gres := gist.Gist{
		ID: "QXhJNchXAK",
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
		},
	}
	_, window, cleanup, err := setup(t, appName, []gist.Gist{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.populate()
	model := window.gistList.Model()
	item := model.Index(0, 0, core.NewQModelIndex())
	desc := item.Data(int(core.Qt__DisplayRole)).ToString()
	if desc != fileName {
		t.Errorf("Display = %s, want %s", desc, fileName)
	}
}

func TestOpenSearchBox(t *testing.T) { tRunner.Run(func() { testOpenSearchBox(t) }) }
func testOpenSearchBox(t *testing.T) {
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	app.SetActiveWindow(window)
	window.Show()
	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_P, core.Qt__ControlModifier, -1)
	tcs := []struct {
		name   string
		widget widgets.QWidget_ITF
	}{
		{"MainWindow", window},
		{"Tabs", window.tabsWidget},
		{"GistList", window.gistList},
	}
	for _, tc := range tcs {
		window.searchbox.Hide()
		event.Simulate(tc.widget)
		if !window.searchbox.IsVisible() {
			t.Errorf("%s: SearchBox is not shown", tc.name)
		}
	}
}

func TestClickOpenGist(t *testing.T) { tRunner.Run(func() { testClickOpenGist(t) }) }
func testClickOpenGist(t *testing.T) {
	var called bool
	gres := gist.Gist{
		ID:          "QXhJNchXAK",
		Description: "kfxLTwoCOkqEuPlp",
	}
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, gres.ID) {
			t.Errorf("r.URL.Path = %s, want %s in it", r.URL.Path, gres.ID)
			return
		}
		// check the URL
		called = true
		b, err := json.Marshal(gres)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()

	gres.URL = fmt.Sprintf("%s/gists/%s", gistTs.URL, gres.ID)
	_, window, cleanup, err := setup(t, appName, []gist.Gist{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.populate()
	window.gistService.API = gistTs.URL

	app.SetActiveWindow(window)
	window.Show()

	window.logger = &logger{
		errorFunc:   func(str string) {},
		warningFunc: func(str string) {},
	}

	for _, key := range []core.Qt__Key{core.Qt__Key_Enter, core.Qt__Key_Return} {
		called = false
		window.searchbox.Show()
		event := testlib.NewQTestEventList()
		event.AddKeyClick(core.Qt__Key_Down, core.Qt__NoModifier, -1)
		event.Simulate(window.searchbox)

		event = testlib.NewQTestEventList()
		event.AddKeyClick(key, core.Qt__NoModifier, -1)
		event.Simulate(window.searchbox.Results())

		if !called {
			t.Error("didn't call for gist")
		}
		// checking the searchbox is closed after
		if window.searchbox.IsVisible() {
			t.Error("searchbox wasn't closed")
		}

		delete(window.tabGistList, gres.ID)
		os.RemoveAll(window.gistService.CacheDir)
		os.Mkdir(window.gistService.CacheDir, 0777)
	}
}

func TestUpdateGistError(t *testing.T) { tRunner.Run(func() { testUpdateGistError(t) }) }
func testUpdateGistError(t *testing.T) {
	var (
		id         = "nb4X55PupEo0bmwM"
		content    = "XzdlfdVudcyYfpm"
		fileName   = "PichzTJDNn"
		newContent = "kRPtHlDFJH9dqzX"
		called     bool
		errored    bool
	)
	updateTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer updateTS.Close()

	gres := gist.Gist{
		ID:  id,
		URL: updateTS.URL,
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
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
	gres.URL = fmt.Sprintf("%s/gists/%s", gistTs.URL, gres.ID)
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.gistService.API = gistTs.URL
	window.logger = &logger{
		errorFunc: func(str string) {
			errored = true
		},
		warningFunc: func(str string) {},
	}

	err = window.openGist(id)
	if err != nil {
		t.Error(err)
		return
	}

	// current tab is newTab
	tabWidget := window.tabGistList[id]
	file := tabWidget.Files()[0]
	// hijacking the url to the new place
	tabWidget.Gist().URL = updateTS.URL
	file.Content().SetText(newContent)
	tabWidget.SaveButton().Click()
	if !called {
		t.Error("didn't call the server")
	}
	if !errored {
		t.Error("didn't show the error")
	}
}

func TestUpdateGist(t *testing.T) { tRunner.Run(func() { testUpdateGist(t) }) }
func testUpdateGist(t *testing.T) {
	var (
		id         = "D4cFvlqRnVg"
		content    = "uKkKeExm8yyJJZEvNFcj"
		fileName   = "OaBdnMbHtq1Y6"
		newContent = "77fVwXPT0lM"
		called     bool
	)
	updateTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer updateTS.Close()

	gres := gist.Gist{
		ID:  id,
		URL: updateTS.URL,
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
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
	gres.URL = fmt.Sprintf("%s/gists/%s", gistTs.URL, gres.ID)
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.gistService.API = gistTs.URL

	err = window.openGist(id)
	if err != nil {
		t.Error(err)
		return
	}

	// current tab is newTab
	tabWidget := window.tabGistList[id]
	file := tabWidget.Files()[0]
	// hijacking the url to the new place
	tabWidget.Gist().URL = updateTS.URL
	file.Content().SetText(newContent)
	tabWidget.SaveButton().Click()
	if !called {
		t.Error("didn't call the server")
	}
}

func TestNewGist(t *testing.T) { tRunner.Run(func() { testNewGist(t) }) }
func testNewGist(t *testing.T) {
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	app.SetActiveWindow(window)
	window.Show()
	currentLen := len(window.tabGistList)
	window.newGist(true)
	if len(window.tabGistList) != currentLen+1 {
		t.Errorf("didn't open the tab: window.tabGistList = %d, want %d", len(window.tabGistList), currentLen+1)
		return
	}

	currentLen = len(window.tabGistList)
	event := testlib.NewQTestEventList()
	event.AddKeyClick(core.Qt__Key_N, core.Qt__ControlModifier, -1)
	event.Simulate(window)
	if len(window.tabGistList) != currentLen+1 {
		t.Errorf("didn't open the tab: window.tabGistList = %d, want %d", len(window.tabGistList), currentLen+1)
	}
}

func TestNewGistSaveError(t *testing.T) { tRunner.Run(func() { testNewGistSaveError(t) }) }
func testNewGistSaveError(t *testing.T) {
	var (
		called      bool
		errorCalled bool
	)
	newGistTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
	}))
	defer newGistTS.Close()

	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.logger = &logger{
		errorFunc: func(str string) {
			errorCalled = true
		},
		warningFunc: func(str string) {},
	}
	window.gistService.API = newGistTS.URL
	window.newGist(true)
	tabWidget := tab.NewTabFromPointer(window.tabsWidget.CurrentWidget().Pointer())
	tabWidget.SaveButton().Click()
	if !called {
		t.Error("didn't call the server")
	}
	if !errorCalled {
		t.Error("didn't record the error")
	}
}

func TestNewGistSave(t *testing.T) { tRunner.Run(func() { testNewGistSave(t) }) }
func testNewGistSave(t *testing.T) {
	var (
		fileName    = "18I96IpY7p"
		description = "apoLIQrkbXEK5LpSWt"
		id          = "kUe"
		called      bool
		created     bool
	)
	newGistTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		var g gist.Gist
		data := new(bytes.Buffer)
		data.ReadFrom(r.Body)
		defer r.Body.Close()
		err := json.Unmarshal(data.Bytes(), &g)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if g.Description != description {
			t.Errorf("g.Description = %s, want %s", g.Description, description)
		}
		if _, ok := g.Files[fileName]; !ok {
			t.Errorf("%s not found in %v", fileName, g.Files)
		}

		g.ID = id
		b, err := json.Marshal(g)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(b)
	}))
	defer newGistTS.Close()
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.gistService.API = newGistTS.URL
	window.newGist(true)
	tab := tab.NewTabFromPointer(window.tabsWidget.CurrentWidget().Pointer())
	tab.SetDescription(description)
	tab.Files()[0].SetFileName(fileName)
	tab.ConnectGistCreated(func(g *gist.Gist) {
		created = true
		if g.Description != description {
			t.Errorf("not the same gist: g.Description = %s, want %s", g.Description, description)
		}
	})

	tab.SaveButton().Click()
	if !called {
		t.Error("didn't call the server")
	}
	if !created {
		t.Error("didn't send the created signal")
	}
	index := window.tabsWidget.IndexOf(tab)
	if window.tabsWidget.TabText(index) != fileName {
		t.Errorf("window.tabsWidget.TabText(%d) = %s, want %s", window.tabsWidget.TabText(index), index, fileName)
	}
	if !window.gistList.HasID(id) {
		t.Errorf("%s was not added to gistList", id)
	}
	if !window.searchbox.HasID(id) {
		t.Errorf("%s was not added to searchbox", id)
	}
}

func TestNewGistGlobalActions(t *testing.T) { tRunner.Run(func() { testNewGistGlobalActions(t) }) }
func testNewGistGlobalActions(t *testing.T) {
	var (
		fileName     = "Pv0LOCOHcuvPD"
		description  = "9JC4ExxYevl1znrd3H"
		called       bool
		clipboardTxt string
	)
	g := gist.Gist{
		ID:      "SdsEUebhhBSx",
		HTMLURL: "LISKELE1UyLThL6",
		Files: map[string]gist.File{
			fileName: gist.File{Content: "3NhvSGH"},
		},
	}
	newGistTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		b, err := json.Marshal(g)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(b)
	}))
	defer newGistTS.Close()
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.gistService.API = newGistTS.URL
	window.newGist(true)
	tab := tab.NewTabFromPointer(window.tabsWidget.CurrentWidget().Pointer())
	tab.SetDescription(description)
	tab.Files()[0].SetFileName(fileName)

	tab.SaveButton().Click()
	if !called {
		t.Error("didn't call the server")
	}
	window.clipboard = func() clipboard {
		return &fakeClipboard{
			textFunc: func(text string, mode gui.QClipboard__Mode) {
				clipboardTxt = text
			},
		}
	}
	c := window.menubar.Actions().CopyURL
	c.Trigger()
	if clipboardTxt == "" {
		t.Error("empty URL")
		return
	}
	if clipboardTxt != tab.Gist().HTMLURL {
		t.Errorf("clipboardTxt = `%s`, want `%s`", clipboardTxt, tab.Gist().HTMLURL)
	}
}

func TestDisplayNoTokenUsername(t *testing.T) { tRunner.Run(func() { testDisplayNoTokenUsername(t) }) }
func testDisplayNoTokenUsername(t *testing.T) {
	var err error
	_, window, cleanup, err := setup(t, appName, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.Display(app)
	tab := window.tabsWidget.CurrentWidget()
	w := conf.NewTabFromPointer(tab.Pointer())
	if w.Pointer() == nil {
		t.Errorf("didn't find conf tab")
	}
}

func TestDisplayTokenUsername(t *testing.T) { tRunner.Run(func() { testDisplayTokenUsername(t) }) }
func testDisplayTokenUsername(t *testing.T) {
	called := make(chan struct{})
	gres := gist.Gist{
		ID:          "WdeDv204lp",
		Description: "HmK2lZP9w4QC",
	}

	_, window, cleanup, err := setup(t, appName, []gist.Gist{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called <- struct{}{}
	}))
	defer gistTs.Close()

	window.gistService.API = gistTs.URL
	settings, cleanup2 := testSettings(appName)
	defer cleanup2()
	settings.SetValue(conf.AccessToken, core.NewQVariant17("KgeU5R2KAq9mHiZc0V"))
	settings.SetValue(conf.Username, core.NewQVariant17("SjnmG9dECJKUowzRVivpb76lcH"))
	settings.Sync()

	window.Display(app)
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Error("didn't call the server")
	}
}

func TestDisplayAfterConfig(t *testing.T) { tRunner.Run(func() { testDisplayAfterConfig(t) }) }
func testDisplayAfterConfig(t *testing.T) {
	called := make(chan struct{})
	gres := gist.Gist{
		ID:          "WdeDv204lp",
		Description: "HmK2lZP9w4QC",
	}
	_, window, cleanup, err := setup(t, appName, []gist.Gist{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called <- struct{}{}
	}))
	defer gistTs.Close()

	window.gistService.API = gistTs.URL
	window.Display(app)
	currentTab := window.tabsWidget.CurrentWidget()
	tab := conf.NewTabFromPointer(currentTab.Pointer())
	if tab.Pointer() == nil {
		t.Error("didn't find the config tab")
		return
	}

	// without entering it should not call the server
	tab.Close()
	select {
	case <-called:
		t.Error("didn't expect to call the server")
	case <-time.After(100 * time.Millisecond):
	}

	tab.UsernameInput.SetText("4xKmhkWG0WvIzPi4")
	tab.AccessTokenInput.SetText("H3iU3XdqlUzTuE2m")
	index := window.tabsWidget.IndexOf(tab)
	window.tabsWidget.TabCloseRequested(index)

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Error("didn't call the server")
	}
}

// test when the user changes the settings. The application should re-initialise
// the searchbox and gistlist.
func TestReconfigureSettings(t *testing.T) { tRunner.Run(func() { testReconfigureSettings(t) }) }
func testReconfigureSettings(t *testing.T) {
	t.Error("Not implemented yet")
}
