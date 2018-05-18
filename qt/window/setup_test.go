// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package window

// this file contains all setup needed for testing.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arsham/gistflow/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// We setup an app and bootstrap the test with it. You should obtain an instance
// of the window through the setup() function and call the cleaup() function
// after each test is finished.
var (
	app     *widgets.QApplication
	appName = "windowTest"
)

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

// logger can be replaced by the messageBox of the MainWindow.
type logger struct {
	criticalFunc func(string) widgets.QMessageBox__StandardButton
	errorFunc    func(string)
	warningFunc  func(string)
}

func (l logger) Critical(msg string) widgets.QMessageBox__StandardButton { return l.criticalFunc(msg) }
func (l logger) Error(msg string)                                        { l.errorFunc(msg) }
func (l logger) Warning(msg string)                                      { l.warningFunc(msg) }
func (l logger) Warningf(format string, a ...interface{})                { l.Warning(fmt.Sprintf(format, a...)) }

// fakeClipboard is used in order to prevent the test's clipboard usage collide
// with normal clipboard system.
type fakeClipboard struct {
	textFunc func(string, gui.QClipboard__Mode)
}

func (f *fakeClipboard) SetText(text string, mode gui.QClipboard__Mode) { f.textFunc(text, mode) }

// testSettings creates a new QSettings and its file. The cleanup function
// should be called to clean up the settings, otherwise there might be
// unexpected results in other tests.
func testSettings(name string) (s *core.QSettings, cleanup func()) {
	s = core.NewQSettings3(
		core.QSettings__NativeFormat,
		core.QSettings__UserScope,
		"gistflow",
		name,
		nil,
	)
	cleanup = func() {
		os.Remove(s.FileName())
	}
	return
}

// The testserver (TS) will print up to `totalAmount` of `input` when is
// reached. After the totalAmount is exhausted, it will print `[\n]`.
func setup(t *testing.T, name string, input []gist.Gist, totalAmount int) (ts *httptest.Server, window *MainWindow, cleanup func(), err error) {
	var counter int
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if counter >= totalAmount {
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
	cacheDir, err := ioutil.TempDir("", "gistflow")
	if err != nil {
		return nil, nil, nil, err
	}

	window = NewMainWindow(nil, 0)
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
	window.clipboard = func() clipboard {
		return &fakeClipboard{textFunc: func(text string, mode gui.QClipboard__Mode) {}}
	}

	cleanup = func() {
		ts.Close()
		window.Hide()
		_, cleanup := testSettings(name)
		cleanup()
		os.RemoveAll(cacheDir)
	}

	return
}
