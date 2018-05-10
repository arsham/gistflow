// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"github.com/arsham/gisty/interface/messagebox"
	"github.com/therecipe/qt/widgets"
)

// File represents one file in a gist.
type File struct {
	widgets.QWidget

	_ func()       `constructor:"init"`
	_ func(string) `signal:"copyToClipboard"`
	_ func()       `signal:"updateGist"`
	_ func(string) `signal:"deleteFile"`

	fileName     *widgets.QLineEdit
	content      *widgets.QTextEdit
	copyButton   *widgets.QPushButton
	deleteButton *widgets.QPushButton
	messageBox   messagebox.Message
}

func (f *File) init() {
	f.SetObjectName("File")

	f.messageBox = messagebox.New(f)
	f.fileName = widgets.NewQLineEdit(f)
	f.fileName.SetPlaceholderText("Filename")

	f.deleteButton = widgets.NewQPushButton(f)
	f.deleteButton.SetText("Delete")
	f.deleteButton.SetToolTip("Deletes this file from this gist. This operations is irreversible")

	f.copyButton = widgets.NewQPushButton(f)
	f.copyButton.SetText("Copy Contents")
	f.copyButton.SetToolTip("Copy contents to system's clipboard")

	f.content = widgets.NewQTextEdit(f)
	f.content.SetObjectName("content")
	f.content.SetPlaceholderText("Contents")

	vLayout := widgets.NewQVBoxLayout2(f)
	hLayout := widgets.NewQHBoxLayout()
	hLayout.AddWidget(f.fileName, 0, 0)
	hSpacer := widgets.NewQSpacerItem(40, 20, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Minimum)
	hLayout.AddItem(hSpacer)
	hLayout.AddWidget(f.copyButton, 0, 0)
	hLayout.AddWidget(f.deleteButton, 0, 0)
	vLayout.AddLayout(hLayout, 0)
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

	f.deleteButton.ConnectClicked(func(bool) {
		b := f.messageBox.Critical("Are you sure you want to delete this file?")
		if b == widgets.QMessageBox__Ok {
			f.DeleteFile(f.fileName.Text())
		}
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
