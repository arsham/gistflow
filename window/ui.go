package window

// func (m *MainWindow) listViewWidget() *widgets.QListView {
// 	openGist := func(index *core.QModelIndex) {
// 		err := m.gistDialog(index)
// 		if err != nil {
// 			id := index.Data(GistID).ToString()
// 			m.logger.warningf("listViewWidget: %s: %s", err, id)
// 		}
// 	}

// 	listView := widgets.NewQListView(m.dialog)
// 	listView.SetObjectName("listView")
// 	})

// 	return listView
// }

// type gistDialog struct {
// 	gist      *widgets.QPlainTextEdit
// 	ok        *widgets.QPushButton
// 	clipboard *widgets.QPushButton
// 	browser   *widgets.QPushButton
// }

// func (g *gistDialog) setupUI(parent *widgets.QDialog) {
// 	vLayout := widgets.NewQVBoxLayout2(nil)
// 	vLayout.SetObjectName("vLayout")
// 	hLayout := widgets.NewQHBoxLayout2(nil)
// 	hLayout.SetObjectName("hLayout")
// 	hLayout2 := widgets.NewQHBoxLayout2(nil)
// 	hLayout2.SetObjectName("hLayout2")

// 	g.gist = widgets.NewQPlainTextEdit(parent)
// 	g.gist.SetObjectName("gist")
// 	vLayout.AddWidget(g.gist, 0, 0)

// 	spacer := widgets.NewQSpacerItem(479, 20, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Minimum)
// 	hLayout2.AddItem(spacer)

// 	g.ok = widgets.NewQPushButton2("Ok", parent)
// 	g.ok.SetObjectName("ok")

// 	hLayout2.AddWidget(g.ok, 0, 0)
// 	g.clipboard = widgets.NewQPushButton2("Copy to Clipboard", parent)
// 	g.clipboard.SetObjectName("clipboard")
// 	hLayout2.AddWidget(g.clipboard, 0, 0)
// 	g.browser = widgets.NewQPushButton2("In Browser", parent)
// 	g.browser.SetObjectName("browser")
// 	hLayout2.AddWidget(g.browser, 0, 0)
// 	vLayout.AddLayout(hLayout2, 0)
// 	hLayout.AddLayout(vLayout, 0)
// 	parent.SetLayout(hLayout)
// }
