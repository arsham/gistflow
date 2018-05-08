// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package tab

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// File represents one file in a gist.
type File struct {
	widgets.QWidget

	_ func()       `constructor:"init"`
	_ func(string) `signal:"copyToClipboard"`

	_ *widgets.QLabel      `property:"information"`
	_ *widgets.QTextEdit   `property:"content"`
	_ *widgets.QPushButton `property:"copy"`
}

func (f *File) init() {
	f.SetObjectName("File")
	vLayout := widgets.NewQVBoxLayout2(f)
	hLayout := widgets.NewQHBoxLayout()
	f.SetInformation(widgets.NewQLabel(f, core.Qt__Widget))
	hLayout.AddWidget(f.Information(), 0, 0)
	hSpacer := widgets.NewQSpacerItem(40, 20, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Minimum)
	hLayout.AddItem(hSpacer)
	f.SetCopy(widgets.NewQPushButton(f))
	f.Copy().SetText("Copy")
	hLayout.AddWidget(f.Copy(), 0, 0)
	vLayout.AddLayout(hLayout, 0)
	f.SetContent(widgets.NewQTextEdit(f))
	f.Content().SetObjectName("content")
	vLayout.AddWidget(f.Content(), 0, 0)

	f.Copy().ConnectClicked(func(bool) {
		f.CopyToClipboard(f.Content().ToPlainText())
	})
}
