// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package searchbox

import (
	"testing"

	"github.com/therecipe/qt/widgets"
)

func TestRemove(t *testing.T) { tRunner.Run(func() { testRemove(t) }) }
func testRemove(t *testing.T) {
	var (
		id1   = "VFN"
		id2   = "ruec9sg"
		l     = NewListModel(widgets.NewQWidget(nil, 0))
		item1 = NewListItem(nil)
		item2 = NewListItem(nil)
	)
	item1.GistID = id1
	item2.GistID = id2

	l.AddGist(item1)
	l.AddGist(item2)

	isIn := func(id string) bool {
		for _, m := range l.gists {
			if m.GistID == id {
				return true
			}
		}
		return false
	}

	l.remove(id1)
	if isIn(id1) {
		t.Errorf("%s was not removed", id1)
	}

	l.remove(id2)
	if isIn(id2) {
		t.Errorf("%s was not removed", id2)
	}
}
