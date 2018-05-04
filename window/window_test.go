// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/therecipe/qt/testlib"

	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func TestMainWindow(t *testing.T) {
	type testCase struct {
		name string
		f    func(t *testing.T) bool
	}
	tRunner.Run(func() {
		tcs := []testCase{
			{"testWindowStartupWidgets", testWindowStartupWidgets},
			{"testWindowModel", testWindowModel},
			{"testPopulateError", testPopulateError},
			{"testPopulate", testPopulate},
			{"testLoadingGeometry", testLoadingGeometry},
			{"testFilteringGists", testFilteringGists},
			{"testListViewKeys", testListViewKeys},
			{"testViewGist", testViewGist},
			{"testClickViewGist", testClickViewGist},
			{"testExchangingFocus", testExchangingFocus},
			{"testOpeningGistTwice", testOpeningGistTwice},
			{"testWindowStartupFocus", testWindowStartupFocus},
			{"testTypingOnListView", testTypingOnListView},
		}
		for _, tc := range tcs {
			if !tc.f(t) {
				t.Errorf("stopped at %s", tc.name)
				return
			}
		}
	})
}

type logger struct {
	errorFunc   func(string)
	warningFunc func(string)
}

func (l logger) Error(msg string)                         { l.errorFunc(msg) }
func (l logger) Warning(msg string)                       { l.warningFunc(msg) }
func (l logger) Warningf(format string, a ...interface{}) { l.Warning(fmt.Sprintf(format, a...)) }

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
	window.app = app
	window.name = name
	window.logger = l

	return ts, window, func() {
		ts.Close()
		window.Hide()
		s := getSettings(name)
		s.Clear()
		os.RemoveAll(cacheDir)
	}, nil
}

// testing vanilla setup
func testWindowStartupWidgets(t *testing.T) bool {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	oldLogger := window.logger
	window.logger = nil
	window.setupUI()

	if window.logger == nil {
		t.Error("window.logger = nil, want boxLogger")
	}
	window.logger = oldLogger

	if window == nil {
		t.Error("window = nil, want *widgets.QMainWindow")
		return false
	}
	if window.icon == nil {
		t.Error("window.icon = nil, want *gui.QIcon")
		return false
	}
	if window.menubar == nil {
		t.Error("window.menubar = nil, want *widgets.QMenuBar")
		return false
	}
	if window.sysTray == nil {
		t.Error("window.sysTray = nil, want *widgets.QSystemTrayIcon")
		return false
	}
	if window.sysTray.Icon() == nil {
		t.Error("window.sysTray.Icon() = nil, want *gui.QIcon")
	}
	if window.sysTray.ContextMenu().Pointer() != window.menubar.menuOptions.Pointer() {
		t.Errorf("window.sysTray.ContextMenu().Pointer() = %v, want %v",
			window.sysTray.ContextMenu().Pointer(),
			window.menubar.menuOptions.Pointer(),
		)
	}

	if window.statusbar == nil {
		t.Error("window.statusbar = nil, want *widgets.QStatusBar")
		return false
	}
	if window.tabWidget == nil {
		t.Error("window.tabWidget = nil, want *widgets.QTabWidget")
		return false
	}
	if !window.tabWidget.IsMovable() {
		t.Error("window.tabWidget is not movable")
	}
	if window.tabGistList == nil {
		t.Error("window.tabGistList = nil, want []*tabGist")
		return false
	}
	if window.gistList == nil {
		t.Error("window.gistList = nil, want *widgets.QListView")
		return false
	}
	if window.dockWidget == nil {
		t.Error("window.dockWidget = nil, want *widgets.QDockWidget")
		return false
	}
	if window.tabWidget.Count() < 1 {
		t.Errorf("window.tabWidget.Count() = %d, want at least 1", window.tabWidget.Count())
	}
	if window.userInput == nil {
		t.Error("window.userInput = nil, want *widgets.QDockWidget")
		return false
	}

	if !window.userInput.IsClearButtonEnabled() {
		t.Error("userInput doesn't have a clear button")
	}
	return true
}

func testWindowModel(t *testing.T) bool {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	window.setupUI()
	window.setModel()
	if window.model == nil {
		t.Error("window.model = nil, want GistModel")
		return false
	}
	model := window.model
	if window.proxy == nil {
		t.Error("window.proxy = nil, want *core.QSortFilterProxyModel")
		return false
	}
	if window.proxy.SourceModel().Pointer() != model.Pointer() {
		t.Errorf("window.proxy.SourceModel().Pointer() = %v, want %v", window.proxy.SourceModel().Pointer(), model.Pointer())
	}
	if window.proxy.FilterCaseSensitivity() != core.Qt__CaseInsensitive {
		t.Errorf("window.proxy.FilterCaseSensitivity() = %d, want %d", window.proxy.FilterCaseSensitivity(), core.Qt__CaseInsensitive)
	}
	if model.Pointer() != window.model.Pointer() {
		t.Errorf("model.Pointer() = %v, want %v", model.Pointer(), window.model.Pointer())
		return false
	}
	if window.gistList.Model().Pointer() != window.proxy.Pointer() {
		t.Errorf("window.gistList.Model().Pointer() = %d, want %d", window.gistList.Model().Pointer(), window.proxy.Pointer())
	}
	return true
}

func testPopulateError(t *testing.T) bool {
	var (
		name   = "test"
		called bool
	)
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	window.gistService.Logger = nil
	defer cleanup()
	window.setupUI()
	window.setModel()

	window.logger = &logger{
		errorFunc: func(str string) {
			called = true
		},
		warningFunc: func(str string) {},
	}
	window.gistService.CacheDir = ""
	window.populate()
	if window.gistService.Logger == nil {
		t.Error("window.gistService.Logger is not assigned")
		return false
	}
	if c := window.model.RowCount(nil); c != 0 {
		t.Errorf("window.model.RowCount() = %d, want 0", c)
	}
	if !called {
		t.Error("expected an error, didn't register the error")
	}
	if window.gistService.CacheDir == "" {
		t.Error("window.gistService.CacheDir is empty")
	}

	return true
}

func testPopulate(t *testing.T) bool {
	name := "test"
	size := 5
	gres := gist.Response{
		ID:          "QXhJNchXAK",
		Description: "kfxLTwoCOkqEuPlp",
	}
	ts, window, cleanup, err := setup(t, name, []gist.Response{gres}, size)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	gres.URL = fmt.Sprintf("%s/gists/%s", ts.URL, gres.ID)

	window.setupUI()
	window.setModel()
	window.populate()

	if c := window.model.RowCount(nil); c != size {
		t.Errorf("window.model.RowCount() = %d, want %d", c, size)
		return false
	}

	model := window.gistList.Model()
	item := model.Index(0, 0, core.NewQModelIndex())
	desc := item.Data(Description).ToString()
	id := item.Data(GistID).ToString()
	if desc != gres.Description {
		t.Errorf("Display = %s, want %s", desc, gres.Description)
	}
	if id != gres.ID {
		t.Errorf("Display = %s, want %s", id, gres.ID)
	}
	return true
}

func testLoadingGeometry(t *testing.T) bool {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	window.setupUI()
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

	return true
}

func testFilteringGists(t *testing.T) bool {
	name := "test"
	res := []gist.Response{
		gist.Response{
			ID:          "QXhJNchXAK",
			Description: "666666 A 6 A 6 A 66666",
		},
		gist.Response{
			ID:          "mLRBtHGeAKRkfENd",
			Description: "666666 BBB 66666",
		},
	}
	_, window, cleanup, err := setup(t, name, res, 1)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	window.setupUI()
	window.setModel()
	window.populate()
	window.setupInteractions()

	window.userInput.SetText("AAA")
	index := core.NewQModelIndex()
	if l := window.proxy.RowCount(index); l != 1 {
		t.Errorf("gistList row count = %d, want %d", l, 1)
		return false
	}
	return true
}

func testListViewKeys(t *testing.T) bool {
	name := "test"
	res := []gist.Response{
		gist.Response{
			ID:          "QXhJNchXAK",
			Description: "666666AAA66666",
		},
	}
	_, window, cleanup, err := setup(t, name, res, 10)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	window.setupUI()
	window.setModel()
	window.populate()
	window.setupInteractions()

	app.SetActiveWindow(window)
	window.show()

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_Down, core.Qt__NoModifier, -1)
	event.Simulate(window.userInput)

	if window.userInput.HasFocus() {
		t.Error("userInput still in focus")
	}
	if !window.gistList.HasFocus() {
		t.Errorf("gistList didn't get focused")
		return false
	}

	if i := window.gistList.CurrentIndex().Row(); i != 0 {
		t.Errorf("gistList.CurrentIndex().Row() = %d, want 0", i)
	}
	event.Simulate(window.gistList)
	event.Simulate(window.gistList)
	if i := window.gistList.CurrentIndex().Row(); i != 2 {
		t.Errorf("gistList.CurrentIndex().Row() = %d, want 2", i)
	}

	event = testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_Up, core.Qt__NoModifier, -1)
	event.Simulate(window.userInput)
	if i := window.gistList.CurrentIndex().Row(); i != 2 {
		t.Errorf("gistList.CurrentIndex().Row() = %d, want 2", i)
	}

	return true
}

func testViewGist(t *testing.T) (forward bool) {
	var (
		called   bool
		name     = "test"
		id       = "uWIkJYdkFuVwYcyy"
		badID    = "kJuZxkDCBp"
		fileName = "LpqrRCgBBYY"
		content  = "fLGLysiOuxReut\nASUonvyd"
	)
	forward = true

	files := map[string]gist.ResponseFile{
		fileName: gist.ResponseFile{Content: content},
	}
	gres := gist.ResponseGist{
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
			forward = false
			return
		}
		w.Write(b)
	}))
	defer gistTs.Close()

	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	window.gistService.API = gistTs.URL

	window.setupUI()
	if window.tabWidget.Count() != 1 {
		t.Errorf("window.tabWidget.Count() = %d, want 1", window.tabWidget.Count())
		return false
	}

	if err := window.openGist(badID); err == nil {
		t.Errorf("window.openGist(%s) = nil, want error", badID)
		forward = false
	}

	if err := window.openGist(id); err != nil {
		t.Errorf("window.openGist(%s) = %s, want nil", id, err)
		forward = false
	}

	newIndex := 2
	if window.tabWidget.Count() != newIndex {
		t.Errorf("window.tabWidget.Count() = %d, want %d", window.tabWidget.Count(), newIndex)
		return false
	}

	index := window.tabWidget.CurrentIndex()
	tab := window.tabWidget.Widget(index)
	guts := widgets.NewQPlainTextEditFromPointer(
		tab.FindChild("content", core.Qt__FindChildrenRecursively).Pointer(),
	)
	if guts.ToPlainText() != content {
		t.Errorf("content = %s, want %s", guts.ToPlainText(), content)
		forward = false
	}
	if window.tabWidget.TabText(index) != fileName {
		t.Errorf("TabText(%d) = %s, want %s", index, window.tabWidget.TabText(index), fileName)
		forward = false
	}
	return forward
}

func testClickViewGist(t *testing.T) bool {
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
		return false
	}
	defer cleanup()

	window.setupUI()
	window.setModel()
	window.populate()
	window.gistService.API = gistTs.URL
	window.setupInteractions()

	app.SetActiveWindow(window)
	window.show()

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
		return false
	}

	return true
}

func testExchangingFocus(t *testing.T) bool {
	name := "test"
	gres := gist.Response{
		ID:          "QXhJNchXAK",
		Description: "kfxLTwoCOkqEuPlp",
	}
	_, window, cleanup, err := setup(t, name, []gist.Response{gres}, 10)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	window.setupUI()
	window.setupInteractions()
	app.SetActiveWindow(window)
	window.show()

	tcs := []core.Qt__Key{
		core.Qt__Key_A,
		core.Qt__Key_Colon,
		core.Qt__Key_Left,
		core.Qt__Key_Right,
		core.Qt__Key_Delete,
	}
	for _, tc := range tcs {
		window.gistList.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(tc, core.Qt__NoModifier, -1)
		event.Simulate(window.gistList)

		if window.gistList.HasFocus() {
			t.Errorf("%x: gistList didn't loose focus", tc)
		}
		if !window.userInput.HasFocus() {
			t.Errorf("%x: userInput didn't gain focus", tc)
		}
	}

	tcs = []core.Qt__Key{
		core.Qt__Key_Up,
		core.Qt__Key_Down,
	}
	for _, tc := range tcs {
		window.gistList.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(tc, core.Qt__NoModifier, -1)
		event.Simulate(window.gistList)

		if !window.gistList.HasFocus() {
			t.Errorf("%x: gistList lost focus", tc)
		}
		if window.userInput.HasFocus() {
			t.Errorf("%x: userInput gained focus", tc)
		}

		event.Simulate(window.userInput)
		if !window.gistList.HasFocus() {
			t.Errorf("%x: gistList didn't gain focus", tc)
		}
		if window.userInput.HasFocus() {
			t.Errorf("%x: userInput didn't loose focus", tc)
		}
	}

	return true
}

func testOpeningGistTwice(t *testing.T) bool {
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
		return false
	}
	defer cleanup()

	window.setupUI()
	window.gistService.API = gistTs.URL

	startingSize := window.tabWidget.Count()
	window.openGist(id1)
	if window.tabWidget.Count() != startingSize+1 {
		t.Errorf("window.tabWidget.Count() = %d, want %d", window.tabWidget.Count(), startingSize+1)
	}
	if window.tabWidget.CurrentIndex() != 1 {
		t.Errorf("window.tabWidget.CurrentIndex() = %d, want 1", window.tabWidget.CurrentIndex())
	}

	window.openGist(id2)
	if window.tabWidget.Count() != startingSize+2 {
		t.Errorf("window.tabWidget.Count() = %d, want %d", window.tabWidget.Count(), startingSize+2)
	}
	if window.tabWidget.CurrentIndex() != 2 {
		t.Errorf("window.tabWidget.CurrentIndex() = %d, want 2", window.tabWidget.CurrentIndex())
	}

	window.openGist(id1)
	if window.tabWidget.Count() != startingSize+2 {
		t.Errorf("window.tabWidget.Count() = %d, want %d", window.tabWidget.Count(), startingSize+2)
	}
	if window.tabWidget.CurrentIndex() != 1 {
		t.Errorf("window.tabWidget.CurrentIndex() = %d, want 1", window.tabWidget.CurrentIndex())
	}

	return true
}

func testWindowStartupFocus(t *testing.T) bool {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	window.setupUI()
	app.SetActiveWindow(window)
	window.show()

	if !window.userInput.HasFocus() {
		t.Errorf("focus is on %s, want %s", app.FocusWidget().ObjectName(), window.userInput.ObjectName())
	}

	return true
}

func testTypingOnListView(t *testing.T) bool {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	window.setupUI()
	window.setupInteractions()

	tcs := []struct {
		prefix string
		input  string
		want   string
	}{
		{"", "a", "a"},
		{"", ":", ":"},
		{"", "-", "-"},
		{"a", "a", "aa"},
		{"a", ":", "a:"},
		{"a", "-", "a-"},
		{"a ", "a", "a a"},
		{"a ", ":", "a :"},
		{"a ", "-", "a -"},
	}
	for _, tc := range tcs {
		window.userInput.SetText(tc.prefix)
		window.gistList.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease2(tc.input, core.Qt__NoModifier, -1)
		event.Simulate(window.gistList)
		if window.userInput.Text() != tc.want {
			t.Errorf("window.userInput.Text() = `%s`, want `%s`", window.userInput.Text(), tc.want)
		}
	}

	return true
}
