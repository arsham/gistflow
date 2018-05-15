// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package gistlist

import (
	"os"
	"strings"
	"testing"

	"github.com/arsham/gistflow/gist"
	"github.com/therecipe/qt/widgets"
)

var app *widgets.QApplication

func TestMain(m *testing.M) {
	app = widgets.NewQApplication(len(os.Args), os.Args)
	go func() { app.Exit(m.Run()) }()
	app.Exec()
}

func TestTruncateDescription(t *testing.T) { tRunner.Run(func() { testTruncateDescription(t) }) }
func testTruncateDescription(t *testing.T) {
	tcs := []struct {
		name     string
		text     string
		truncate bool
	}{
		{"short", strings.Repeat("A", maxLen-1), false},
		{"exact", strings.Repeat("A", maxLen), false},
		{"long", strings.Repeat("A", maxLen+1), true},
	}

	for _, tc := range tcs {
		c := NewContainer(widgets.NewQWidget(nil, 0))
		r := gist.Gist{
			Description: tc.text,
		}
		c.Add(r)
		ret := c.Item(0).Text()
		switch tc.truncate {
		case true:
			if len(ret) > maxLen {
				t.Errorf("%s: len(%s) = %d, want at most %d", tc.name, ret, len(ret), maxLen)
			}
			if !strings.HasSuffix(ret, truncateStr) {
				t.Errorf("%s: ret = %s, want %s at the end", tc.name, ret, truncateStr)
			}
		case false:
			if ret != tc.text {
				t.Errorf("%s: ret = %s, want %s", tc.name, ret, tc.text)
			}
		}
	}
}

func TestID(t *testing.T) { tRunner.Run(func() { testID(t) }) }
func testID(t *testing.T) {
	var (
		id1 = "8QyPFAPTJ"
		id2 = "lezFNeVL7cZUWnM"
	)
	c := NewContainer(widgets.NewQWidget(nil, 0))
	c.Add(gist.Gist{ID: id1})
	c.Add(gist.Gist{ID: id2})
	if c.ID(0) != id1 {
		t.Errorf("c.ID(0) = %s, want %s", c.ID(0), id1)
	}

	if c.ID(1) != id2 {
		t.Errorf("c.ID(1) = %s, want %s", c.ID(1), id2)
	}
}

func TestDescription(t *testing.T) { tRunner.Run(func() { testDescription(t) }) }
func testDescription(t *testing.T) {
	var (
		desciption1 = "82CEM7bn"
		desciption2 = "JcuFLi"
	)
	c := NewContainer(widgets.NewQWidget(nil, 0))
	c.Add(gist.Gist{Description: desciption1})
	c.Add(gist.Gist{Description: desciption2})
	if c.Description(0) != desciption1 {
		t.Errorf("c.Description(0) = %s, want %s", c.Description(0), desciption1)
	}

	if c.Description(1) != desciption2 {
		t.Errorf("c.Description(1) = %s, want %s", c.Description(1), desciption2)
	}
}

func TestIndexID(t *testing.T) { tRunner.Run(func() { testIndexID(t) }) }
func testIndexID(t *testing.T) {
	var (
		id1 = "8rV9Kg30Mc5YByEl7j"
		id2 = "bQFmUmPsgpZRNo8Y"
	)
	c := NewContainer(widgets.NewQWidget(nil, 0))
	c.Add(gist.Gist{ID: id1})
	c.Add(gist.Gist{ID: id2})
	item := c.Item(0)
	index := c.IndexFromItem(item)
	if c.IndexID(index) != id1 {
		t.Errorf("c.IndexID(:0) = %s, want %s", c.IndexID(index), id1)
	}

	item = c.Item(1)
	index = c.IndexFromItem(item)
	if c.IndexID(index) != id2 {
		t.Errorf("c.IndexID(:1) = %s, want %s", c.IndexID(index), id2)
	}
}

func TestIndexDescription(t *testing.T) { tRunner.Run(func() { testIndexDescription(t) }) }
func testIndexDescription(t *testing.T) {
	var (
		desciption1 = "IAHOOv4pz9zk"
		desciption2 = "kNsiedEcfF"
		c           = NewContainer(widgets.NewQWidget(nil, 0))
	)
	c.Add(gist.Gist{Description: desciption1})
	c.Add(gist.Gist{Description: desciption2})
	item := c.Item(0)
	index := c.IndexFromItem(item)
	if c.IndexDescription(index) != desciption1 {
		t.Errorf("c.IndexDescription(:0) = %s, want %s", c.IndexDescription(index), desciption1)
	}

	item = c.Item(1)
	index = c.IndexFromItem(item)
	if c.IndexDescription(index) != desciption2 {
		t.Errorf("c.IndexDescription(:1) = %s, want %s", c.IndexDescription(index), desciption2)
	}
}

func TestEmptyDescription(t *testing.T) { tRunner.Run(func() { testEmptyDescription(t) }) }
func testEmptyDescription(t *testing.T) {
	var (
		id       = "kqapzK9iRVppLxHbS"
		fileName = "LC7BVBWAuCXY"
		c        = NewContainer(widgets.NewQWidget(nil, 0))
	)
	c.Add(gist.Gist{
		ID: id,
		Files: map[string]gist.File{
			fileName: gist.File{},
		},
	})
	if c.Description(0) != fileName {
		t.Errorf("c.Description(:0) = %s, want %s", c.Description(0), fileName)
	}
}

func TestRemoveItem(t *testing.T) { tRunner.Run(func() { testRemoveItem(t) }) }
func testRemoveItem(t *testing.T) {
	var (
		id1 = "QJn1eTU5bHOzUPc"
		id2 = "not found"
		id3 = "086vmZLyedK"
		c   = NewContainer(widgets.NewQWidget(nil, 0))
	)
	c.Add(gist.Gist{ID: id1})
	c.Add(gist.Gist{ID: id3})

	currentLen := len(c.items)
	if len(c.items) != currentLen {
		t.Errorf("len(c.items) = %d, want %d", len(c.items), currentLen)
		return
	}

	item := c.items[id1]
	c.Remove(id1)
	if len(c.items) != currentLen-1 {
		t.Errorf("len(c.items) = %d, want %d", len(c.items), currentLen-1)
		return
	}
	row := c.IndexFromItem(item).Row()
	if row != -1 {
		t.Errorf("row  %d, want -1", row)
	}

	c.Remove(id2)
	if len(c.items) != currentLen-1 {
		t.Errorf("len(c.items) = %d, want %d", len(c.items), currentLen-1)
		return
	}

	item = c.items[id3]
	c.Remove(id3)
	if len(c.items) != currentLen-2 {
		t.Errorf("len(c.items) = %d, want %d", len(c.items), currentLen-2)
		return
	}
	row = c.IndexFromItem(item).Row()
	if row != -1 {
		t.Errorf("row  %d, want -1", row)
	}
}
