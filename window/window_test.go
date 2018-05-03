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

	mw := &MainWindow{
		GistService: gist.Service{
			Username: "arsham",
			Token:    "token",
			API:      ts.URL,
			CacheDir: cacheDir,
			Logger:   l2,
		},
		app:      app,
		ConfName: name,
		logger:   l,
	}
	return ts, mw, func() {
		ts.Close()
		mw.window.Hide()
		s := getSettings(name)
		s.Clear()
		os.RemoveAll(cacheDir)
	}, nil
}

// testing vanilla setup
func testWindowStartupWidgets(t *testing.T) bool {
	name := "test"
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	oldLogger := mw.logger
	mw.logger = nil
	if err := mw.setupUI(); err != nil {
		t.Errorf("mw.setupUI() = %v, want nil", err)
		return false
	}

	if mw.logger == nil {
		t.Error("mw.logger = nil, want boxLogger")
	}
	mw.logger = oldLogger

	if mw.window == nil {
		t.Error("mw.window = nil, want *widgets.QMainWindow")
		return false
	}
	if mw.icon == nil {
		t.Error("mw.icon = nil, want *gui.QIcon")
		return false
	}
	if mw.menubar == nil {
		t.Error("mw.menubar = nil, want *widgets.QMenuBar")
		return false
	}
	if mw.menubar.quitAction == nil {
		t.Error("mw.menubar.quitAction = nil, want *widgets.QAction")
		return false
	}
	if mw.sysTray == nil {
		t.Error("mw.sysTray = nil, want *widgets.QSystemTrayIcon")
		return false
	}
	if mw.sysTray.Icon() == nil {
		t.Error("mw.sysTray.Icon() = nil, want *gui.QIcon")
	}
	if mw.sysTray.ContextMenu().Pointer() != mw.menubar.optionsMenu.Pointer() {
		t.Errorf("mw.sysTray.ContextMenu().Pointer() = %v, want %v",
			mw.sysTray.ContextMenu().Pointer(),
			mw.menubar.optionsMenu.Pointer(),
		)
	}

	if mw.statusbar == nil {
		t.Error("mw.statusbar = nil, want *widgets.QStatusBar")
		return false
	}
	if mw.tabWidget == nil {
		t.Error("mw.tabWidget = nil, want *widgets.QTabWidget")
		return false
	}
	if !mw.tabWidget.IsMovable() {
		t.Error("mw.tabWidget is not movable")
	}
	if mw.tabGistList == nil {
		t.Error("mw.tabGistList = nil, want []*tabGist")
		return false
	}
	if mw.listView == nil {
		t.Error("mw.listView = nil, want *widgets.QListView")
		return false
	}
	if mw.dockWidget == nil {
		t.Error("mw.dockWidget = nil, want *widgets.QDockWidget")
		return false
	}
	if mw.tabWidget.Count() < 1 {
		t.Errorf("mw.tabWidget.Count() = %d, want at least 1", mw.tabWidget.Count())
	}
	if mw.userInput == nil {
		t.Error("mw.userInput = nil, want *widgets.QDockWidget")
		return false
	}

	if !mw.userInput.IsClearButtonEnabled() {
		t.Error("userInput doesn't have a clear button")
	}
	return true
}

func testWindowModel(t *testing.T) bool {
	name := "test"
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	mw.setupUI()
	mw.setModel()
	if mw.model == nil {
		t.Error("mw.model = nil, want GistModel")
		return false
	}
	model := mw.model
	if mw.proxy == nil {
		t.Error("mw.proxy = nil, want *core.QSortFilterProxyModel")
		return false
	}
	if mw.proxy.SourceModel().Pointer() != model.Pointer() {
		t.Errorf("mw.proxy.SourceModel().Pointer() = %v, want %v", mw.proxy.SourceModel().Pointer(), model.Pointer())
	}
	if mw.proxy.FilterCaseSensitivity() != core.Qt__CaseInsensitive {
		t.Errorf("mw.proxy.FilterCaseSensitivity() = %d, want %d", mw.proxy.FilterCaseSensitivity(), core.Qt__CaseInsensitive)
	}
	if model.Pointer() != mw.model.Pointer() {
		t.Errorf("model.Pointer() = %v, want %v", model.Pointer(), mw.model.Pointer())
		return false
	}
	if mw.listView.Model().Pointer() != mw.proxy.Pointer() {
		t.Errorf("mw.listView.Model().Pointer() = %d, want %d", mw.listView.Model().Pointer(), mw.proxy.Pointer())
	}
	return true
}

func testPopulateError(t *testing.T) bool {
	var (
		name   = "test"
		called bool
	)
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	mw.GistService.Logger = nil
	defer cleanup()
	mw.setupUI()
	mw.setModel()

	mw.logger = &logger{
		errorFunc: func(str string) {
			called = true
		},
		warningFunc: func(str string) {},
	}
	mw.GistService.CacheDir = ""
	mw.populate()
	if mw.GistService.Logger == nil {
		t.Error("mw.GistService.Logger is not assigned")
		return false
	}
	if c := mw.model.RowCount(nil); c != 0 {
		t.Errorf("mw.model.RowCount() = %d, want 0", c)
	}
	if !called {
		t.Error("expected an error, didn't register the error")
	}
	if mw.GistService.CacheDir == "" {
		t.Error("mw.GistService.CacheDir is empty")
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
	ts, mw, cleanup, err := setup(t, name, []gist.Response{gres}, size)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	gres.URL = fmt.Sprintf("%s/gists/%s", ts.URL, gres.ID)

	mw.setupUI()
	mw.setModel()
	mw.populate()

	if c := mw.model.RowCount(nil); c != size {
		t.Errorf("mw.model.RowCount() = %d, want %d", c, size)
		return false
	}

	model := mw.listView.Model()
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
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	x, y, w, h := 400, 500, 600, 700
	tmpObj := widgets.NewQWidget(nil, 0)
	tmpObj.SetGeometry2(x, y, w, h)
	mw.settings = getSettings(name)
	size := tmpObj.SaveGeometry()
	mw.settings.SetValue(mainWindowGeometry, core.NewQVariant15(size))
	mw.settings.Sync()

	mw.loadSettings()
	geometry := mw.window.Geometry()
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
	mw.window.RestoreGeometry(newGeometry)
	check("to make sure: geometry.X()", mw.window.Geometry().X(), x)
	check("to make sure: geometry.Y()", mw.window.Geometry().Y(), y)
	check("to make sure: geometry.Width()", mw.window.Geometry().Width(), w)
	check("to make sure: geometry.Height()", mw.window.Geometry().Height(), h)

	mw.menubar.quitAction.Activate(widgets.QAction__Trigger)
	mw.menubar.quitAction.ConnectTriggered(func(bool) {
		tmp := widgets.NewQWidget(nil, 0)
		tmp.SetGeometry2(100, 100, 600, 600)
		defSize := tmp.SaveGeometry()
		sizeVar := mw.settings.Value(mainWindowGeometry, core.NewQVariant15(defSize))
		tmp.RestoreGeometry(sizeVar.ToByteArray())
		geometry := tmp.Geometry()

		check("after quiting: geometry.X()", geometry.X(), x)
		check("after quiting: geometry.Y()", geometry.Y(), y)
		check("after quiting: geometry.Width()", geometry.Width(), w)
		check("after quiting: geometry.Height()", geometry.Height(), h)
	})
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
	_, mw, cleanup, err := setup(t, name, res, 1)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	mw.setModel()
	mw.populate()
	mw.setupInteractions()

	mw.userInput.SetText("AAA")
	index := core.NewQModelIndex()
	if l := mw.proxy.RowCount(index); l != 1 {
		t.Errorf("listView row count = %d, want %d", l, 1)
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
	_, mw, cleanup, err := setup(t, name, res, 10)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	mw.setModel()
	mw.populate()
	mw.setupInteractions()

	app.SetActiveWindow(mw.window)
	mw.show()

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_Down, core.Qt__NoModifier, -1)
	event.Simulate(mw.userInput)

	if mw.userInput.HasFocus() {
		t.Error("userInput still in focus")
	}
	if !mw.listView.HasFocus() {
		t.Errorf("listView didn't get focused")
		return false
	}

	if i := mw.listView.CurrentIndex().Row(); i != 0 {
		t.Errorf("listView.CurrentIndex().Row() = %d, want 0", i)
	}
	event.Simulate(mw.listView)
	event.Simulate(mw.listView)
	if i := mw.listView.CurrentIndex().Row(); i != 2 {
		t.Errorf("listView.CurrentIndex().Row() = %d, want 2", i)
	}

	event = testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_Up, core.Qt__NoModifier, -1)
	event.Simulate(mw.userInput)
	if i := mw.listView.CurrentIndex().Row(); i != 2 {
		t.Errorf("listView.CurrentIndex().Row() = %d, want 2", i)
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

	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()
	mw.GistService.API = gistTs.URL

	mw.setupUI()
	if mw.tabWidget.Count() != 1 {
		t.Errorf("mw.tabWidget.Count() = %d, want 1", mw.tabWidget.Count())
		return false
	}

	if err := mw.openGist(badID); err == nil {
		t.Errorf("mw.openGist(%s) = nil, want error", badID)
		forward = false
	}

	if err := mw.openGist(id); err != nil {
		t.Errorf("mw.openGist(%s) = %s, want nil", id, err)
		forward = false
	}

	newIndex := 2
	if mw.tabWidget.Count() != newIndex {
		t.Errorf("mw.tabWidget.Count() = %d, want %d", mw.tabWidget.Count(), newIndex)
		return false
	}

	index := mw.tabWidget.CurrentIndex()
	tab := mw.tabWidget.Widget(index)
	guts := widgets.NewQPlainTextEditFromPointer(
		tab.FindChild("content", core.Qt__FindChildrenRecursively).Pointer(),
	)
	if guts.ToPlainText() != content {
		t.Errorf("content = %s, want %s", guts.ToPlainText(), content)
		forward = false
	}
	if mw.tabWidget.TabText(index) != fileName {
		t.Errorf("TabText(%d) = %s, want %s", index, mw.tabWidget.TabText(index), fileName)
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
	_, mw, cleanup, err := setup(t, name, []gist.Response{gres}, 10)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	mw.setModel()
	mw.populate()
	mw.GistService.API = gistTs.URL
	mw.setupInteractions()

	app.SetActiveWindow(mw.window)
	mw.show()

	var errCalled bool
	mw.logger = &logger{
		errorFunc:   func(str string) { errCalled = true },
		warningFunc: func(str string) { errCalled = true },
	}

	// with no selection, it should error because there is no item selectedisd,
	// hence no id.
	event := testlib.NewQTestEventList()
	event.AddKeyRelease(core.Qt__Key_Down, core.Qt__NoModifier, -1)
	event.AddKeyRelease(core.Qt__Key_Enter, core.Qt__NoModifier, -1)
	event.Simulate(mw.listView)
	if !errCalled {
		t.Error("didn't show error")
	}

	for _, key := range []core.Qt__Key{core.Qt__Key_Enter, core.Qt__Key_Return} {
		called = false
		mw.listView.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(core.Qt__Key_Down, core.Qt__NoModifier, -1)
		event.AddKeyRelease(key, core.Qt__NoModifier, -1)
		event.Simulate(mw.listView)

		if !called {
			t.Error("didn't call for gist")
		}
		delete(mw.tabGistList, gres.ID)
		os.RemoveAll(mw.GistService.CacheDir)
		os.Mkdir(mw.GistService.CacheDir, 0777)
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
	_, mw, cleanup, err := setup(t, name, []gist.Response{gres}, 10)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	mw.setupInteractions()
	app.SetActiveWindow(mw.window)
	mw.show()

	tcs := []core.Qt__Key{
		core.Qt__Key_A,
		core.Qt__Key_Colon,
		core.Qt__Key_Left,
		core.Qt__Key_Right,
		core.Qt__Key_Delete,
	}
	for _, tc := range tcs {
		mw.listView.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(tc, core.Qt__NoModifier, -1)
		event.Simulate(mw.listView)

		if mw.listView.HasFocus() {
			t.Errorf("%x: listView didn't loose focus", tc)
		}
		if !mw.userInput.HasFocus() {
			t.Errorf("%x: userInput didn't gain focus", tc)
		}
	}

	tcs = []core.Qt__Key{
		core.Qt__Key_Up,
		core.Qt__Key_Down,
	}
	for _, tc := range tcs {
		mw.listView.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(tc, core.Qt__NoModifier, -1)
		event.Simulate(mw.listView)

		if !mw.listView.HasFocus() {
			t.Errorf("%x: listView lost focus", tc)
		}
		if mw.userInput.HasFocus() {
			t.Errorf("%x: userInput gained focus", tc)
		}

		event.Simulate(mw.userInput)
		if !mw.listView.HasFocus() {
			t.Errorf("%x: listView didn't gain focus", tc)
		}
		if mw.userInput.HasFocus() {
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
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	mw.GistService.API = gistTs.URL

	startingSize := mw.tabWidget.Count()
	mw.openGist(id1)
	if mw.tabWidget.Count() != startingSize+1 {
		t.Errorf("mw.tabWidget.Count() = %d, want %d", mw.tabWidget.Count(), startingSize+1)
	}
	if mw.tabWidget.CurrentIndex() != 1 {
		t.Errorf("mw.tabWidget.CurrentIndex() = %d, want 1", mw.tabWidget.CurrentIndex())
	}

	mw.openGist(id2)
	if mw.tabWidget.Count() != startingSize+2 {
		t.Errorf("mw.tabWidget.Count() = %d, want %d", mw.tabWidget.Count(), startingSize+2)
	}
	if mw.tabWidget.CurrentIndex() != 2 {
		t.Errorf("mw.tabWidget.CurrentIndex() = %d, want 2", mw.tabWidget.CurrentIndex())
	}

	mw.openGist(id1)
	if mw.tabWidget.Count() != startingSize+2 {
		t.Errorf("mw.tabWidget.Count() = %d, want %d", mw.tabWidget.Count(), startingSize+2)
	}
	if mw.tabWidget.CurrentIndex() != 1 {
		t.Errorf("mw.tabWidget.CurrentIndex() = %d, want 1", mw.tabWidget.CurrentIndex())
	}

	return true
}

func testWindowStartupFocus(t *testing.T) bool {
	name := "test"
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	app.SetActiveWindow(mw.window)
	mw.show()

	if !mw.userInput.HasFocus() {
		t.Errorf("focus is on %s, want %s", app.FocusWidget().ObjectName(), mw.userInput.ObjectName())
	}

	return true
}

func testTypingOnListView(t *testing.T) bool {
	name := "test"
	_, mw, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return false
	}
	defer cleanup()

	mw.setupUI()
	mw.setupInteractions()

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
		mw.userInput.SetText(tc.prefix)
		mw.listView.SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease2(tc.input, core.Qt__NoModifier, -1)
		event.Simulate(mw.listView)
		if mw.userInput.Text() != tc.want {
			t.Errorf("mw.userInput.Text() = `%s`, want `%s`", mw.userInput.Text(), tc.want)
		}
	}

	return true
}
