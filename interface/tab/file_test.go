// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"testing"

	"github.com/therecipe/qt/widgets"
)

func TestFile(t *testing.T) { tRunner.Run(func() { testFile(t) }) }
func testFile(t *testing.T) {
	file := NewFile(widgets.NewQWidget(nil, 0), 0)
	if file.fileName == nil {
		t.Error("file.fileName = nil, want *widgets.QLineEdit")
	}
	if file.content == nil {
		t.Error("file.content = nil, want *widgets.QTextEdit")
	}
	if file.copyButton == nil {
		t.Error("file.copyButton = nil, want *widgets.QPushButton")
	}
	if file.deleteButton == nil {
		t.Error("file.deleteButton = nil, want *widgets.QPushButton")
	}
	if file.messageBox == nil {
		t.Error("file.messageBox = nil, want messagebox.Message")
	}
}

func TestFileSignals(t *testing.T) { tRunner.Run(func() { testFileSignals(t) }) }
func testFileSignals(t *testing.T) {
	var (
		called  bool
		content = "I2FmBmMNdtRBypCxGYq"
	)
	file := NewFile(widgets.NewQWidget(nil, 0), 0)
	file.ConnectCopyToClipboard(func(text string) {
		called = true
		if text != content {
			t.Errorf("text = %s, want %s", text, content)
		}
	})

	file.ConnectUpdateGist(func() {
		called = true
	})
	file.content.SetText(content)
	file.CopyButton().Click()
	if !called {
		t.Error("didn't trigger the signal")
	}

	called = false
	file.SetFileName("new fileName")
	if !called {
		t.Error("didn't trigger the signal")
	}
}

func TestFileDeleteConfirmation(t *testing.T) { tRunner.Run(func() { testFileDeleteConfirmation(t) }) }
func testFileDeleteConfirmation(t *testing.T) {
	var (
		called       bool
		deleteCalled bool
		fileName     = "rSYfamuuIDvXsrjUgl"
		button       widgets.QMessageBox__StandardButton
	)
	w := widgets.NewQWidget(nil, 0)
	file := NewFile(w, 0)
	file.fileName.SetText(fileName)
	file.messageBox = logger{
		criticalFunc: func(string) widgets.QMessageBox__StandardButton {
			called = true
			return button
		},
	}
	file.ConnectDeleteFile(func(name string) {
		deleteCalled = true
		if name != fileName {
			t.Errorf("name = %s, want %s", name, fileName)
		}
	})

	button = widgets.QMessageBox__Cancel
	file.deleteButton.Click()
	if !called {
		t.Error("didn't trigger the signal")
	}
	if deleteCalled {
		t.Error("didn't expect to send deletion signal")
	}

	button = widgets.QMessageBox__Ok
	called = false
	file.deleteButton.Click()
	if !called {
		t.Error("didn't trigger the signal")
	}
	if !deleteCalled {
		t.Error("didn't send deletion signal")
	}
}
