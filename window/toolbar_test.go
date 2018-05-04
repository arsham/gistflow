// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package window

import (
	"testing"

	"github.com/therecipe/qt/widgets"
)

func TestToolbarActions(t *testing.T) {
	action := NewAppAction(nil)
	toolbar := NewAppToolbar2(nil)
	toolbar.SetAction(action)

	tcs := map[string]*widgets.QAction{
		"actionClipboard": action.actionClipboard,
		"actionSync":      action.actionSync,
		"actionCopyURL":   action.actionCopyURL,
	}
	for name, a := range tcs {
		if !isIn(toolbar.Actions(), a) {
			t.Errorf("%s was not found in actions", name)
		}
	}
}

func isIn(actions []*widgets.QAction, action *widgets.QAction) bool {
	for _, a := range actions {
		if a.Pointer() == action.Pointer() {
			return true
		}
	}
	return false
}
