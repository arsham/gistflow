// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

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
	"time"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/testlib"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/gisty/interface/tab"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

type logger struct {
	errorFunc   func(string)
	warningFunc func(string)
}

func (l logger) Error(msg string)                         { l.errorFunc(msg) }
func (l logger) Warning(msg string)                       { l.warningFunc(msg) }
func (l logger) Warningf(format string, a ...interface{}) { l.Warning(fmt.Sprintf(format, a...)) }

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

	if window.StatusArea() == nil {
		t.Error("window.StatusArea() = nil, want *widgets.QStatusBar")
		return
	}
	if window.TabsWidget() == nil {
		t.Error("window.TabsWidget() = nil, want *widgets.QTabWidget")
		return
	}
	if !window.TabsWidget().IsMovable() {
		t.Error("window.TabsWidget() is not movable")
	}
	if window.tabGistList == nil {
		t.Error("window.tabGistList = nil, want []*tabGist")
		return
	}
	if window.GistList() == nil {
		t.Error("window.GistList() = nil, want *widgets.QListView")
		return
	}
	if window.dockWidget == nil {
		t.Error("window.dockWidget = nil, want *widgets.QDockWidget")
		return
	}
	if window.TabsWidget().Count() < 1 {
		t.Errorf("window.TabsWidget().Count() = %d, want at least 1", window.TabsWidget().Count())
	}
	if window.userInput == nil {
		t.Error("window.userInput = nil, want *widgets.QDockWidget")
		return
	}

	if !window.userInput.IsClearButtonEnabled() {
		t.Error("userInput doesn't have a clear button")
	}

	if window.clipboard() == nil {
		t.Error("window.clipboard() = nil, want *gui.QClipboard")
		return
	}
}

func TestWindowModel(t *testing.T) { tRunner.Run(func() { testWindowModel(t) }) }
func testWindowModel(t *testing.T) {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	window.setModel()
	if window.model == nil {
		t.Error("window.model = nil, want listGistModel")
		return
	}
	model := window.model
	if window.proxy == nil {
		t.Error("window.proxy = nil, want *core.QSortFilterProxyModel")
		return
	}
	if window.proxy.SourceModel().Pointer() != model.Pointer() {
		t.Errorf("window.proxy.SourceModel().Pointer() = %v, want %v", window.proxy.SourceModel().Pointer(), model.Pointer())
	}
	if window.proxy.FilterCaseSensitivity() != core.Qt__CaseInsensitive {
		t.Errorf("window.proxy.FilterCaseSensitivity() = %d, want %d", window.proxy.FilterCaseSensitivity(), core.Qt__CaseInsensitive)
	}
	if model.Pointer() != window.model.Pointer() {
		t.Errorf("model.Pointer() = %v, want %v", model.Pointer(), window.model.Pointer())
		return
	}
	if window.GistList().Model().Pointer() != window.proxy.Pointer() {
		t.Errorf("window.GistList().Model().Pointer() = %d, want %d", window.GistList().Model().Pointer(), window.proxy.Pointer())
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
	window.setModel()

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
	if c := window.model.RowCount(nil); c != 0 {
		t.Errorf("window.model.RowCount() = %d, want 0", c)
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

	window.setModel()
	window.populate()

	if c := window.model.RowCount(nil); c != size {
		t.Errorf("window.model.RowCount() = %d, want %d", c, size)
		return
	}

	model := window.GistList().Model()
	item := model.Index(0, 0, core.NewQModelIndex())
	desc := item.Data(tab.Description).ToString()
	id := item.Data(tab.GistID).ToString()
	if desc != gres.Description {
		t.Errorf("Display = %s, want %s", desc, gres.Description)
	}
	if id != gres.ID {
		t.Errorf("Display = %s, want %s", id, gres.ID)
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
	window.SetSettings(getSettings(name))
	size := tmpObj.SaveGeometry()
	window.Settings().SetValue(mainWindowGeometry, core.NewQVariant15(size))
	window.Settings().Sync()

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

func TestFilteringGists(t *testing.T) { tRunner.Run(func() { testFilteringGists(t) }) }
func testFilteringGists(t *testing.T) {
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
		return
	}
	defer cleanup()

	window.setModel()
	window.populate()

	window.userInput.SetText("AAA")
	index := core.NewQModelIndex()
	if l := window.proxy.RowCount(index); l != 1 {
		t.Errorf("proxy row count = %d, want %d", l, 1)
		return
	}
}

func TestListViewKeys(t *testing.T) { tRunner.Run(func() { testListViewKeys(t) }) }
func testListViewKeys(t *testing.T) {
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
		return
	}
	defer cleanup()

	window.setModel()
	window.populate()

	app.SetActiveWindow(window)
	window.show()

	event := testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_Down, core.Qt__NoModifier, -1)
	event.Simulate(window.userInput)

	if window.userInput.HasFocus() {
		t.Error("userInput still in focus")
	}
	if !window.GistList().HasFocus() {
		t.Errorf("window.GistList() didn't get focused")
		return
	}

	if i := window.GistList().CurrentIndex().Row(); i != 0 {
		t.Errorf("window.GistList().CurrentIndex().Row() = %d, want 0", i)
	}
	event.Simulate(window.GistList())
	event.Simulate(window.GistList())
	if i := window.GistList().CurrentIndex().Row(); i != 2 {
		t.Errorf("window.GistList().CurrentIndex().Row() = %d, want 2", i)
	}

	event = testlib.NewQTestEventList()
	event.AddKeyPress(core.Qt__Key_Up, core.Qt__NoModifier, -1)
	event.Simulate(window.userInput)
	if i := window.GistList().CurrentIndex().Row(); i != 2 {
		t.Errorf("window.GistList().CurrentIndex().Row() = %d, want 2", i)
	}
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

	if window.TabsWidget().Count() != 1 {
		t.Errorf("window.TabsWidget().Count() = %d, want 1", window.TabsWidget().Count())
		return
	}

	if err := window.openGist(badID); err == nil {
		t.Errorf("window.openGist(%s) = nil, want error", badID)
	}

	if err := window.openGist(id); err != nil {
		t.Errorf("window.openGist(%s) = %s, want nil", id, err)
	}

	newIndex := 2
	if window.TabsWidget().Count() != newIndex {
		t.Errorf("window.TabsWidget().Count() = %d, want %d", window.TabsWidget().Count(), newIndex)
		return
	}

	index := window.TabsWidget().CurrentIndex()
	tab := tab.NewTabFromPointer(window.TabsWidget().Widget(index).Pointer())
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
	if window.TabsWidget().TabText(index) != fileName {
		t.Errorf("TabText(%d) = %s, want %s", index, window.TabsWidget().TabText(index), fileName)
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

	window.setModel()
	window.populate()
	window.gistService.API = gistTs.URL

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
	event.Simulate(window.GistList())
	if !errCalled {
		t.Error("didn't show error")
	}

	for _, key := range []core.Qt__Key{core.Qt__Key_Enter, core.Qt__Key_Return} {
		called = false
		window.GistList().SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(core.Qt__Key_Down, core.Qt__NoModifier, -1)
		event.AddKeyRelease(key, core.Qt__NoModifier, -1)
		event.Simulate(window.GistList())

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

func TestExchangingFocus(t *testing.T) { tRunner.Run(func() { testExchangingFocus(t) }) }
func testExchangingFocus(t *testing.T) {
	name := "test"
	gres := gist.Response{
		ID:          "QXhJNchXAK",
		Description: "kfxLTwoCOkqEuPlp",
	}
	_, window, cleanup, err := setup(t, name, []gist.Response{gres}, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

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
		window.GistList().SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(tc, core.Qt__NoModifier, -1)
		event.Simulate(window.GistList())

		if window.GistList().HasFocus() {
			t.Errorf("%x: GistList() didn't loose focus", tc)
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
		window.GistList().SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease(tc, core.Qt__NoModifier, -1)
		event.Simulate(window.GistList())

		if !window.GistList().HasFocus() {
			t.Errorf("%x: GistList() lost focus", tc)
		}
		if window.userInput.HasFocus() {
			t.Errorf("%x: userInput gained focus", tc)
		}

		event.Simulate(window.userInput)
		if !window.GistList().HasFocus() {
			t.Errorf("%x: GistList() didn't gain focus", tc)
		}
		if window.userInput.HasFocus() {
			t.Errorf("%x: userInput didn't loose focus", tc)
		}
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

	startingSize := window.TabsWidget().Count()
	window.openGist(id1)
	if window.TabsWidget().Count() != startingSize+1 {
		t.Errorf("window.TabsWidget().Count() = %d, want %d", window.TabsWidget().Count(), startingSize+1)
	}
	if window.TabsWidget().CurrentIndex() != 1 {
		t.Errorf("window.TabsWidget().CurrentIndex() = %d, want 1", window.TabsWidget().CurrentIndex())
	}

	window.openGist(id2)
	if window.TabsWidget().Count() != startingSize+2 {
		t.Errorf("window.TabsWidget().Count() = %d, want %d", window.TabsWidget().Count(), startingSize+2)
	}
	if window.TabsWidget().CurrentIndex() != 2 {
		t.Errorf("window.TabsWidget().CurrentIndex() = %d, want 2", window.TabsWidget().CurrentIndex())
	}

	window.openGist(id1)
	if window.TabsWidget().Count() != startingSize+2 {
		t.Errorf("window.TabsWidget().Count() = %d, want %d", window.TabsWidget().Count(), startingSize+2)
	}
	if window.TabsWidget().CurrentIndex() != 1 {
		t.Errorf("window.TabsWidget().CurrentIndex() = %d, want 1", window.TabsWidget().CurrentIndex())
	}
}

func TestWindowStartupFocus(t *testing.T) { tRunner.Run(func() { testWindowStartupFocus(t) }) }
func testWindowStartupFocus(t *testing.T) {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	app.SetActiveWindow(window)
	window.show()

	if !window.userInput.HasFocus() {
		t.Errorf("focus is on %s, want %s", app.FocusWidget().ObjectName(), window.userInput.ObjectName())
	}
}

func TestTypingOnListView(t *testing.T) { tRunner.Run(func() { testTypingOnListView(t) }) }
func testTypingOnListView(t *testing.T) {
	name := "test"
	_, window, cleanup, err := setup(t, name, nil, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

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
		window.GistList().SetFocus2()
		event := testlib.NewQTestEventList()
		event.AddKeyRelease2(tc.input, core.Qt__NoModifier, -1)
		event.Simulate(window.GistList())
		if window.userInput.Text() != tc.want {
			t.Errorf("window.userInput.Text() = `%s`, want `%s`", window.userInput.Text(), tc.want)
		}
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
	window.TabsWidget().SetCurrentWidget(tab1)

	c.Trigger()
	if clpText != url1 {
		t.Errorf("clpText = `%s`, want `%s`", clpText, url1)
	}

	window.TabsWidget().SetCurrentWidget(tab2)
	c.Trigger()
	if clpText != url2 {
		t.Errorf("clpText = `%s`, want `%s`", clpText, url2)
	}
}
