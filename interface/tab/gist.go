// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"github.com/therecipe/qt/widgets"
)

// File represents one file in a gist.
type File struct {
	widgets.QWidget

	_ func()       `constructor:"init"`
	_ func(string) `signal:"copyToClipboard"`
	_ func()       `signal:"updateGist"`

	fileName   *widgets.QLineEdit
	content    *widgets.QTextEdit
	copyButton *widgets.QPushButton
}

func (f *File) init() {
	f.SetObjectName("File")
	vLayout := widgets.NewQVBoxLayout2(f)
	hLayout := widgets.NewQHBoxLayout()
	f.fileName = widgets.NewQLineEdit(f)
	f.fileName.SetPlaceholderText("Filename")
	hLayout.AddWidget(f.fileName, 0, 0)
	hSpacer := widgets.NewQSpacerItem(40, 20, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Minimum)
	hLayout.AddItem(hSpacer)
	f.copyButton = widgets.NewQPushButton(f)
	f.copyButton.SetText("Copy Contents")
	hLayout.AddWidget(f.copyButton, 0, 0)
	vLayout.AddLayout(hLayout, 0)
	f.content = widgets.NewQTextEdit(f)
	f.content.SetObjectName("content")
	f.content.SetPlaceholderText("Contents")
	vLayout.AddWidget(f.content, 0, 0)

	f.copyButton.ConnectClicked(func(bool) {
		f.CopyToClipboard(f.content.ToPlainText())
	})
	f.content.ConnectTextChanged(func() {
		f.UpdateGist()
	})
	f.fileName.ConnectTextChanged(func(text string) {
		f.UpdateGist()
	})
}

// FileName returns the fileName.
func (f *File) FileName() string { return f.fileName.Text() }

// SetFileName returns the fileName.
func (f *File) SetFileName(fileName string) { f.fileName.SetText(fileName) }

// Content returns the content.
func (f *File) Content() *widgets.QTextEdit { return f.content }

// CopyButton returns the copyButton.
func (f *File) CopyButton() *widgets.QPushButton { return f.copyButton }
