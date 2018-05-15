// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package conf

import (
	"errors"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// Variable names in settings.
const (
	AccessToken = "access_token"
	Username    = "username"
)

// Tab is a tab shown in the tabWidget area that contains the application's
// settings.
type Tab struct {
	widgets.QTabWidget

	_ func() `constructor:"init"`

	Name             string // This should be the application name to reflect the right entry.
	settings         *Settings
	GridLayout       *widgets.QGridLayout
	UsernameInput    *widgets.QLineEdit
	AccessTokenInput *widgets.QLineEdit
}

func (t *Tab) init() {
	t.SetObjectName("SettingsTab")
	groupBox := widgets.NewQGroupBox2("Essentials", t)
	groupBox.SetGeometry(core.NewQRect4(10, 20, 511, 161))
	t.GridLayout = widgets.NewQGridLayout(groupBox)
	t.GridLayout.SetObjectName("gridLayout")
	t.GridLayout.SetContentsMargins(0, 0, 0, 0)
	t.GridLayout.SetSpacing(0)
	label := widgets.NewQLabel2("Username", groupBox, core.Qt__Widget)
	t.GridLayout.AddWidget3(label, 0, 0, 1, 1, 0)
	t.UsernameInput = widgets.NewQLineEdit(groupBox)
	t.UsernameInput.SetClearButtonEnabled(true)
	t.GridLayout.AddWidget3(t.UsernameInput, 0, 1, 1, 1, 0)
	label2 := widgets.NewQLabel2("Access token", t, core.Qt__Widget)
	t.GridLayout.AddWidget3(label2, 1, 0, 1, 1, 0)
	t.AccessTokenInput = widgets.NewQLineEdit(groupBox)
	t.AccessTokenInput.SetClearButtonEnabled(true)
	t.GridLayout.AddWidget3(t.AccessTokenInput, 1, 1, 1, 1, 0)

	labelText := "Click <a href='https://github.com/settings/tokens'>here</a> to create a new access token. This will take you to a take where you can generate a new token. Copy the token and leave it in the box above."
	label3 := widgets.NewQLabel2(labelText, groupBox, core.Qt__Widget)
	label3.SetTextInteractionFlags(core.Qt__TextBrowserInteraction)
	label3.SetOpenExternalLinks(true)
	t.GridLayout.AddWidget3(label3, 2, 0, 1, 2, 0)
}

// SetSettings assigns the Settings instance and updates it when the values are
// changed.
func (t *Tab) SetSettings(s *Settings) {
	t.settings = s
	t.UsernameInput.ConnectTextChanged(func(text string) {
		s.SetValue(Username, core.NewQVariant17(text))
		s.Username = text
		s.Sync()
	})
	t.AccessTokenInput.ConnectTextChanged(func(text string) {
		s.SetValue(AccessToken, core.NewQVariant17(text))
		s.Token = text
		s.Sync()
	})

	v := s.Value(Username, core.NewQVariant17(""))
	if v.ToString() != "" {
		t.UsernameInput.SetText(v.ToString())
	}
	v = s.Value(AccessToken, core.NewQVariant17(""))
	if v.ToString() != "" {
		t.AccessTokenInput.SetText(v.ToString())
	}
}

// Settings holds the written settings loaded from system.
type Settings struct {
	*core.QSettings
	Token    string
	Username string
}

// New returns an instance of Settings. name is the application name, which is
// used to identify the settings. It returns an error if the token and the
// username has not been set yet.
func New(name string) (*Settings, error) {
	var err error
	s := core.NewQSettings3(
		core.QSettings__NativeFormat,
		core.QSettings__UserScope,
		"gistflow",
		name,
		nil,
	)
	token := s.Value(AccessToken, core.NewQVariant17(""))
	if token.ToString() == "" {
		err = errors.New("empty token")
	}
	username := s.Value(Username, core.NewQVariant17(""))
	if username.ToString() == "" {
		err = errors.New("empty username")
	}
	return &Settings{
		Token:     token.ToString(),
		Username:  username.ToString(),
		QSettings: s,
	}, err
}
