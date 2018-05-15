// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package conf

import (
	"os"
	"testing"

	"github.com/therecipe/qt/core"

	"github.com/therecipe/qt/widgets"
)

var (
	app         *widgets.QApplication
	appName     = "configTest"
	userName    = "Rbrmk6ZOP2VRn"
	accessToken = "pGfGUZroCva"
)

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func testSettings(name string) (*core.QSettings, func()) {
	s := core.NewQSettings3(
		core.QSettings__NativeFormat,
		core.QSettings__UserScope,
		"gistflow",
		name,
		nil,
	)
	s.SetValue(Username, core.NewQVariant17(userName))
	s.SetValue(AccessToken, core.NewQVariant17(accessToken))
	s.Sync()
	return s, func() {
		os.Remove(s.FileName())
	}
}

func TestFirstTimeSettings(t *testing.T) { tRunner.Run(func() { testFirstTimeSettings(t) }) }
func testFirstTimeSettings(t *testing.T) {
	settings, err := New(appName)
	if err == nil {
		t.Error("err = nil, want error")
	}
	if settings == nil {
		t.Error("settings = nil, want *Settings")
	}

	ts, cleanup := testSettings(appName)
	defer cleanup()

	ts.Remove(Username)
	ts.Sync()

	settings, err = New(appName)
	if err == nil {
		t.Error("err = nil, want error")
	}
	if settings == nil {
		t.Error("settings = nil, want *Settings")
	}

	ts.Remove(AccessToken)
	ts.SetValue(Username, core.NewQVariant17("zBXyfBu"))
	ts.Sync()
	settings, err = New(appName)
	if err == nil {
		t.Error("err = nil, want error")
	}
	if settings == nil {
		t.Error("settings = nil, want *Settings")
	}

	ts.SetValue(Username, core.NewQVariant17("Sj5Qziu6ylvJ8KOK25"))
	ts.SetValue(AccessToken, core.NewQVariant17("VEteiRh3"))
	ts.Sync()
	settings, err = New(appName)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if settings == nil {
		t.Error("settings = nil, want *Settings")
	}
}

func TestTab(t *testing.T) { tRunner.Run(func() { testTab(t) }) }
func testTab(t *testing.T) {
	tab := NewTab(nil)
	if tab.GridLayout == nil {
		t.Error("tab.GridLayout = nil, want *widgets.QGridLayout")
	}
	if tab.UsernameInput == nil {
		t.Error("tab.UsernameInput = nil, want *widgets.QLineEdit")
	}
	if tab.AccessTokenInput == nil {
		t.Error("tab.AccessTokenInput = nil, want *widgets.QLineEdit")
	}

	_, cleanup := testSettings(appName)
	defer cleanup()

	settings, err := New(appName)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	tab.SetSettings(settings)
	var (
		username = "JsQZlcBUreOQkQy"
		token    = "JnaUZhtza6jN"
	)
	tab.UsernameInput.SetText(username)

	v := settings.Value(Username, core.NewQVariant17(""))
	if v.ToString() != username {
		t.Errorf("written username = %s, want %s", v.ToString(), username)
	}
	if settings.Username != username {
		t.Errorf("settings.Username = %s, want %s", settings.Username, username)
	}

	tab.AccessTokenInput.SetText(token)
	v = settings.Value(AccessToken, core.NewQVariant17(""))
	if v.ToString() != token {
		t.Errorf("written token = %s, want %s", v.ToString(), token)
	}
	if settings.Token != token {
		t.Errorf("settings.Token = %s, want %s", settings.Token, token)
	}
}

func TestTabPrePopulate(t *testing.T) { tRunner.Run(func() { testTabPrePopulate(t) }) }
func testTabPrePopulate(t *testing.T) {
	_, cleanup := testSettings(appName)
	defer cleanup()
	settings, err := New(appName)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}

	tab := NewTab(nil)
	tab.SetSettings(settings)
	if tab.UsernameInput.Text() != userName {
		t.Errorf("tab.UsernameInput.Text() = %s, want %s", tab.UsernameInput.Text(), userName)
	}
	if tab.AccessTokenInput.Text() != accessToken {
		t.Errorf("tab.AccessTokenInput.Text() = %s, want %s", tab.AccessTokenInput.Text(), accessToken)
	}
}
