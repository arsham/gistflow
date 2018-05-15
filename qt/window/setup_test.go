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

var (
	app     *widgets.QApplication
	appName = "windowTest"
)

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

func testSettings(name string) (*core.QSettings, func()) {
	s := core.NewQSettings3(
		core.QSettings__NativeFormat,
		core.QSettings__UserScope,
		"gistflow",
		name,
		nil,
	)
	return s, func() {
		os.Remove(s.FileName())
	}
}

func setup(t *testing.T, name string, input []gist.Gist, answers int) (*httptest.Server, *MainWindow, func(), error) {
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
	cacheDir, err := ioutil.TempDir("", "gistflow")
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
		_, cleanup := testSettings(name)
		cleanup()
		os.RemoveAll(cacheDir)
	}, nil
}
