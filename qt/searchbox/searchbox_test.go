// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package searchbox

import (
	"os"
	"testing"

	"github.com/arsham/gistflow/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/testlib"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func TestDialog(t *testing.T) { tRunner.Run(func() { testDialog(t) }) }
func testDialog(t *testing.T) {
	g := NewDialog(nil, 0)
	if g.results == nil {
		t.Error("g.results = nil, want *widgets.QListView")
	}
	if g.input == nil {
		t.Error("g.input = nil, want *widgets.QLineEdit")
	}
	if g.model == nil {
		t.Error("g.model = nil, want *ListModel")
	}
	if g.proxy == nil {
		t.Error("g.proxy = nil, want *ListModel")
	}
	if g.proxy.SourceModel().Pointer() != g.model.Pointer() {
		t.Errorf("g.proxy.SourceModel() = %v, want %v", g.proxy.SourceModel().Pointer(), g.model.Pointer())
	}
	if g.results.Model().Pointer() != g.proxy.Pointer() {
		t.Errorf("Results().Model() = %v, want %v", g.results.Model().Pointer(), g.proxy.Pointer())
	}
}

func TestEscape(t *testing.T) { tRunner.Run(func() { testEscape(t) }) }
func testEscape(t *testing.T) {
	var called bool
	m := NewDialog(nil, 0)
	app.SetActiveWindow(m)
	m.Show()
	defer m.Hide()

	m.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		// checking the test's logic here.
		if event.Key() == int(core.Qt__Key_Escape) {
			called = true
		}
	})

	event := testlib.NewQTestEventList()
	event.AddKeyClick(core.Qt__Key_Escape, core.Qt__NoModifier, -1)
	event.Simulate(m)
	if !called {
		t.Error("Escape didn't trigger the event")
	}
	if m.IsVisible() {
		t.Error("Widget is still visible")
	}
}

func TestID(t *testing.T) { tRunner.Run(func() { testID(t) }) }
func testID(t *testing.T) {
	var (
		id1   = "blUdwZYJyG1PZ1hlQ"
		id2   = "FMxVi90zkbjU"
		d     = NewDialog(widgets.NewQWidget(nil, 0), 0)
		g1    = NewListItem(nil)
		g2    = NewListItem(nil)
		model = NewListModel(nil)
	)
	d.results.SetModel(model)

	g1.GistID = id1
	g2.GistID = id2

	model.AddGist(g1)
	model.AddGist(g2)

	if d.ID(0) != id1 {
		t.Errorf("d.ID(0) = %s, want %s", d.ID(0), id1)
	}

	if d.ID(1) != id2 {
		t.Errorf("d.ID(1) = %s, want %s", d.ID(1), id2)
	}
}

func TestDescription(t *testing.T) { tRunner.Run(func() { testDescription(t) }) }
func testDescription(t *testing.T) {
	var (
		description1 = "q5a"
		description2 = "G23teuAJT"
		d            = NewDialog(widgets.NewQWidget(nil, 0), 0)
		g1           = NewListItem(nil)
		g2           = NewListItem(nil)
		model        = NewListModel(nil)
	)
	d.results.SetModel(model)

	g1.Description = description1
	g2.Description = description2

	model.AddGist(g1)
	model.AddGist(g2)

	if d.Description(0) != description1 {
		t.Errorf("d.Description(:0) = %s, want %s", d.Description(0), description1)
	}

	if d.Description(1) != description2 {
		t.Errorf("d.Description(:1) = %s, want %s", d.Description(1), description2)
	}
}

func TestFiltering(t *testing.T) { tRunner.Run(func() { testFiltering(t) }) }
func testFiltering(t *testing.T) {
	var (
		description1 = "adfasdfA A Aasdfakfj"
		description2 = "klsjdhfB B Bsdfklsjhf"
		d            = NewDialog(widgets.NewQWidget(nil, 0), 0)
		g1           = NewListItem(nil)
		g2           = NewListItem(nil)
	)

	g1.Description = description1
	g2.Description = description2
	d.Add(gist.Gist{Description: description1})
	d.Add(gist.Gist{Description: description2})

	if d.Model().RowCount(core.NewQModelIndex()) != 2 {
		t.Errorf("RowCount() = %d, want 2", d.Model().RowCount(core.NewQModelIndex()))
		return
	}

	d.input.SetText("A*A*A")

	if d.Model().RowCount(core.NewQModelIndex()) != 1 {
		t.Errorf("RowCount() = %d, want 1", d.Model().RowCount(core.NewQModelIndex()))
	}
	if d.Description(0) != description1 {
		t.Errorf("d.Description(0) = %s, want %s", d.Description(0), description1)
	}

	d.input.SetText("B*B*B")

	if d.Model().RowCount(core.NewQModelIndex()) != 1 {
		t.Errorf("RowCount() = %d, want 1", d.Model().RowCount(core.NewQModelIndex()))
	}
	if d.Description(0) != description2 {
		t.Errorf("d.Description(0) = %s, want %s", d.Description(0), description2)
	}
}

func TestKeepTopMostIndexOnResults(t *testing.T) {
	tRunner.Run(func() { testKeepTopMostIndexOnResults(t) })
}
func testKeepTopMostIndexOnResults(t *testing.T) {
	var (
		description1 = "rhd1f2OTPqJ"
		description2 = "xDmPggqRJSKP8R"
		parent       = widgets.NewQWidget(nil, 0)
		d            = NewDialog(parent, 0)
		g1           = NewListItem(nil)
		g2           = NewListItem(nil)
	)

	g1.Description = description1
	g2.Description = description2
	d.Add(gist.Gist{Description: description1})
	d.Add(gist.Gist{Description: description2})

	app.SetActiveWindow(d)
	d.View(parent.Geometry())
	if d.results.CurrentIndex().Row() < 0 {
		t.Errorf("CurrentIndex().Row() = %d, want >= 0", d.results.CurrentIndex().Row())
	}

	tcs := []string{description1, description2, description2[2:4], ""}
	for _, tc := range tcs {
		d.input.SetText(tc)
		if d.results.CurrentIndex().Row() < 0 {
			t.Errorf("`%s`: CurrentIndex().Row() = %d, want >= 0", tc, d.results.CurrentIndex().Row())
		}
	}
}

func TestNagigatingSearchBox(t *testing.T) {
	up := testlib.NewQTestEventList()
	up.AddKeyPress(core.Qt__Key_Up, core.Qt__NoModifier, -1)
	down := testlib.NewQTestEventList()
	down.AddKeyPress(core.Qt__Key_Down, core.Qt__NoModifier, -1)

	tRunner.Run(func() { testNagigatingSearchBox(t, "up", up) })
	tRunner.Run(func() { testNagigatingSearchBox(t, "down", down) })
}
func testNagigatingSearchBox(t *testing.T, name string, dir *testlib.QTestEventList) {
	var (
		description = "kIdCJuTsv7N2R2"
		parent      = widgets.NewQWidget(nil, 0)
		d           = NewDialog(parent, 0)
		g           = NewListItem(nil)
	)
	g.Description = description
	d.Add(gist.Gist{Description: description})

	app.SetActiveWindow(d)
	d.View(parent.Geometry())
	dir.Simulate(d)
	if d.results.CurrentIndex().Row() != 0 {
		t.Errorf("%s: CurrentIndex().Row() = %d, want 0", name, d.results.CurrentIndex().Row())
	}
	if !d.results.HasFocus() {
		t.Errorf("%s: results didn't get focused", name)
	}
}

func TestOpenGistSlot(t *testing.T) { tRunner.Run(func() { testOpenGistSlot(t) }) }
func testOpenGistSlot(t *testing.T) {
	var (
		id     = "gB4MGfquZxDUhNJas"
		d      = NewDialog(widgets.NewQWidget(nil, 0), 0)
		called bool
	)
	app.SetActiveWindow(d)
	d.Show()
	defer d.Hide()
	d.Add(gist.Gist{ID: id})

	d.ConnectOpenGist(func(text string) {
		called = true
		if text != id {
			t.Errorf("id = %s, want %s", text, id)
		}
	})

	keys := []core.Qt__Key{core.Qt__Key_Enter, core.Qt__Key_Return}
	for _, key := range keys {
		event := testlib.NewQTestEventList()
		event.AddKeyClick(core.Qt__Key_Down, core.Qt__NoModifier, -1)
		event.Simulate(d)

		event = testlib.NewQTestEventList()
		event.AddKeyClick(key, core.Qt__NoModifier, -1)
		event.Simulate(d.results)
		if !called {
			t.Error("signal wasn't received")
		}
		called = false
	}
}

// TODO: typing in results should add to the input
