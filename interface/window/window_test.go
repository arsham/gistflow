// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/gisty/interface/gistlist"
	"github.com/arsham/gisty/interface/searchbox"
	"github.com/arsham/gisty/interface/tab"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/testlib"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

type logger struct {
	criticalFunc func(string) widgets.QMessageBox__StandardButton
	errorFunc    func(string)
	warningFunc  func(string)
}

func (l logger) Critical(msg string) widgets.QMessageBox__StandardButton { return l.criticalFunc(msg) }
func (l logger) Error(msg string)                                        { l.errorFunc(msg) }
func (l logger) Warning(msg string)                                      { l.warningFunc(msg) }
func (l logger) Warningf(format string, a ...interface{})                { l.Warning(fmt.Sprintf(format, a...)) }

type fakeClipboard struct {
	textFunc func(string, gui.QClipboard__Mode)
}

func (f *fakeClipboard) SetText(text string, mode gui.QClipboard__Mode) { f.textFunc(text, mode) }

func setup(t *testing.T, name string, input []gist.Response, answers int) (*httptest.Server, *MainWindow, func(), error) {
	var counter int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter >= answers {
			w.Write([]byte("[\n]"))
			return
		}
		counter++
		b, err := json.Marshal(input)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	}))
	l := &logger{
		errorFunc:   func(string) {},
		warningFunc: func(string) {},
	}
	l2 := &logger{
		errorFunc:   func(string) {},
		warningFunc: func(string) {},
	}
	cacheDir, err := ioutil.TempDir("", "gisty")
	if err != nil {
		return nil, nil, nil, err
	}

	window := NewMainWindow(nil, 0)
	window.gistService = gist.Service{
		Username: "arsham",
		Token:    "token",
		API:      ts.URL,
		CacheDir: cacheDir,
		Logger:   l2,
	}
	window.SetApp(app)
	window.name = name
	window.logger = l
	window.clipboard = func() clipboard {
		return &fakeClipboard{textFunc: func(text string, mode gui.QClipboard__Mode) {}}
	}

	return ts, window, func() {
		ts.Close()
		window.Hide()
		s := getSettings(name)
		s.Clear()
		os.RemoveAll(cacheDir)
	}, nil
}

// testing vanilla setup
func TestWindowStartupWidgets(t *testing.T) { tRunner.Run(func() { testWindowStartupWidgets(t) }) }
func testWindowStartupWidgets(t *testing.T) {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	oldLogger := window.logger
	window.logger = nil
	window.setupUI()
	window.clipboard = func() clipboard {
		return &fakeClipboard{textFunc: func(text string, mode gui.QClipboard__Mode) {}}
	}

	if window.logger == nil {
		t.Error("window.logger = nil, want boxLogger")
	}
	window.logger = oldLogger

	if window == nil {
		t.Error("window = nil, want *widgets.QMainWindow")
		return
	}
	if window.icon == nil {
		t.Error("window.icon = nil, want *gui.QIcon")
		return
	}
	if window.menubar == nil {
		t.Error("window.menubar = nil, want *widgets.QMenuBar")
		return
	}
	if window.sysTray == nil {
		t.Error("window.sysTray = nil, want *widgets.QSystemTrayIcon")
		return
	}
	if window.sysTray.Icon() == nil {
		t.Error("window.sysTray.Icon() = nil, want *gui.QIcon")
	}
	if window.sysTray.ContextMenu().Pointer() != window.menubar.Options().Pointer() {
		t.Errorf("window.sysTray.ContextMenu().Pointer() = %v, want %v",
			window.sysTray.ContextMenu().Pointer(),
			window.menubar.Options().Pointer(),
		)
	}

	if window.statusArea == nil {
		t.Error("window.statusArea = nil, want *widgets.QStatusBar")
		return
	}
	if window.tabsWidget == nil {
		t.Error("window.tabsWidget = nil, want *widgets.QTabWidget")
		return
	}
	if !window.tabsWidget.IsMovable() {
		t.Error("window.tabsWidget is not movable")
	}
	if window.tabGistList == nil {
		t.Error("window.tabGistList = nil, want []*tabGist")
		return
	}
	if window.gistList == nil {
		t.Error("window.gistList = nil, want *widgets.QListView")
		return
	}
	if window.dockWidget == nil {
		t.Error("window.dockWidget = nil, want *widgets.QDockWidget")
		return
	}
	if window.tabsWidget.Count() < 1 {
		t.Errorf("window.tabsWidget.Count() = %d, want at least 1", window.tabsWidget.Count())
	}

	if window.clipboard() == nil {
		t.Error("window.clipboard() = nil, want *gui.QClipboard")
		return
	}

	if window.gistList == nil {
		t.Error("window.gistList = nil, want *gistlist.Container")
		return
	}
}

func TestPopulateError(t *testing.T) { tRunner.Run(func() { testPopulateError(t) }) }
func testPopulateError(t *testing.T) {
	name := "test"
	called := make(chan struct{})
	_, window, cleanup, err := setup(t, name, nil, 0)
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
		warningFunc: func(str string) {},
	}
	window.gistService.CacheDir = ""
	window.populate()
	if window.gistService.Logger == nil {
		t.Error("window.gistService.Logger is not assigned")
		return
	}
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

	if window.gistService.CacheDir == "" {
		t.Error("window.gistService.CacheDir is empty")
	}
}

func TestPopulate(t *testing.T) { tRunner.Run(func() { testPopulate(t) }) }
func testPopulate(t *testing.T) {
	name := "test"
	size := 5
	gres := gist.Response{
		ID:          "QXhJNchXAK",
		Description: "kfxLTwoCOkqEuPlp",
	}
	ts, window, cleanup, err := setup(t, name, []gist.Response{gres}, size)
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
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	x, y, w, h := 400, 500, 600, 700
	tmpObj := widgets.NewQWidget(nil, 0)
	tmpObj.SetGeometry2(x, y, w, h)
	window.settings = getSettings(name)
	size := tmpObj.SaveGeometry()
	window.settings.SetValue(mainWindowGeometry, core.NewQVariant15(size))
	window.settings.Sync()

	window.loadSettings()
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
		called   bool
		name     = "test"
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
		called = true
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

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.gistService.API = gistTs.URL

	if window.tabsWidget.Count() != 1 {
		t.Errorf("window.tabsWidget.Count() = %d, want 1", window.tabsWidget.Count())
		return
	}

	if err := window.openGist(badID); err == nil {
		t.Errorf("window.openGist(%s) = nil, want error", badID)
	}

	if err := window.openGist(id); err != nil {
		t.Errorf("window.openGist(%s) = %s, want nil", id, err)
	}

	newIndex := 2
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
	name := "test"
	var called bool
	gres := gist.Response{
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
	_, window, cleanup, err := setup(t, name, []gist.Response{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.populate()
	window.gistService.API = gistTs.URL

	app.SetActiveWindow(window)

	var errCalled bool
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

	if !called {
		return
	}
}

func TestOpeningGistTwice(t *testing.T) { tRunner.Run(func() { testOpeningGistTwice(t) }) }
func testOpeningGistTwice(t *testing.T) {
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

	window.gistService.API = gistTs.URL

	startingSize := window.tabsWidget.Count()
	window.openGist(id1)
	if window.tabsWidget.Count() != startingSize+1 {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), startingSize+1)
	}
	if window.tabsWidget.CurrentIndex() != 1 {
		t.Errorf("window.tabsWidget.CurrentIndex() = %d, want 1", window.tabsWidget.CurrentIndex())
	}

	window.openGist(id2)
	if window.tabsWidget.Count() != startingSize+2 {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), startingSize+2)
	}
	if window.tabsWidget.CurrentIndex() != 2 {
		t.Errorf("window.tabsWidget.CurrentIndex() = %d, want 2", window.tabsWidget.CurrentIndex())
	}

	window.openGist(id1)
	if window.tabsWidget.Count() != startingSize+2 {
		t.Errorf("window.tabsWidget.Count() = %d, want %d", window.tabsWidget.Count(), startingSize+2)
	}
	if window.tabsWidget.CurrentIndex() != 1 {
		t.Errorf("window.tabsWidget.CurrentIndex() = %d, want 1", window.tabsWidget.CurrentIndex())
	}
}

func TestToggle(t *testing.T) { tRunner.Run(func() { testToggle(t) }) }
func testToggle(t *testing.T) {
	name := "test"
	window := NewMainWindow(nil, 0)
	window.name = name
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
		name     = "test"
		id1      = "wqWKsfoQevEbGjhmz"
		content  = "AFMQydAKTiJLa"
		id2      = "yuaosJCTsGUqEldvigi"
		api, url string
		clpText  string
	)

	gistTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = fmt.Sprintf("%s%s", api, r.URL.Path)
		gres := gist.Gist{
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
		name     = "test"
		content  = "CNF5EmQJxiGvzwedbmTME3p0Y"
		fileName = "84nkJJG0"
	)

	gres := gist.Response{
		ID: "QXhJNchXAK",
		Files: map[string]gist.File{
			fileName: gist.File{Content: content},
		},
	}
	_, window, cleanup, err := setup(t, name, []gist.Response{gres}, 10)
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
	_, window, cleanup, err := setup(t, "test", nil, 0)
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
	name := "test"
	var called bool
	gres := gist.Response{
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
	_, window, cleanup, err := setup(t, name, []gist.Response{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	window.populate()
	window.gistService.API = gistTs.URL

	app.SetActiveWindow(window)
	window.Show()

	var errCalled bool
	window.logger = &logger{
		errorFunc:   func(str string) { errCalled = true },
		warningFunc: func(str string) { errCalled = true },
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

	if !called {
		return
	}
}

func TestUpdateGistError(t *testing.T) { tRunner.Run(func() { testUpdateGistError(t) }) }
func testUpdateGistError(t *testing.T) {
	var (
		name       = "test"
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
	_, window, cleanup, err := setup(t, name, nil, 0)
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
		name       = "test"
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
	_, window, cleanup, err := setup(t, name, nil, 0)
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
	var name = "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
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
		return
	}
}

func TestNewGistSaveError(t *testing.T) { tRunner.Run(func() { testNewGistSaveError(t) }) }
func testNewGistSaveError(t *testing.T) {
	var (
		name        = "test"
		called      bool
		errorCalled bool
	)
	newGistTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
	}))
	defer newGistTS.Close()

	_, window, cleanup, err := setup(t, name, nil, 0)
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
		name        = "test"
		fileName    = "18I96IpY7p"
		description = "apoLIQrkbXEK5LpSWt"
		called      bool
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer newGistTS.Close()
	_, window, cleanup, err := setup(t, name, nil, 0)
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
}

// Test closing dirty gists
