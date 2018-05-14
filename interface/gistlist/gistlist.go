// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package gistlist

import (
	"github.com/arsham/gisty/gist"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const (
	maxLen      = 40
	truncateStr = "..."
)

// Container holds a list of gists. This is used for the left side gist list.
type Container struct {
	widgets.QListWidget

	_ func()              `constructor:"init"`
	_ func(gist.Response) `signal:"add"`

	items map[string]*widgets.QListWidgetItem
}

func init() {
	Container_QRegisterMetaType()
}

func (c *Container) init() {
	c.ConnectAdd(c.add)
	c.items = make(map[string]*widgets.QListWidgetItem, 10)
}

func (c *Container) add(r gist.Response) {
	item := widgets.NewQListWidgetItem(c, 0)
	description := r.Description
	if description == "" {
		for n := range r.Files {
			description = n
			break
		}
	}
	if len(description) > maxLen {
		description = description[:maxLen-len(truncateStr)] + truncateStr
	}
	item.SetText(description)
	item.SetData(int(core.Qt__UserRole), core.NewQVariant14(r.ID))
	c.AddItem2(item)
	c.items[r.ID] = item
}

// ID returns the ID of the gist associated with the index.
func (c *Container) ID(row int) string {
	return c.Item(row).Data(int(core.Qt__UserRole)).ToString()
}

// Description returns the Description of the gist associated with the index.
func (c *Container) Description(row int) string {
	return c.Item(row).Text()
}

// IndexID returns the ID of the gist associated with the index.
func (c *Container) IndexID(index *core.QModelIndex) string {
	return index.Data(int(core.Qt__UserRole)).ToString()
}

// IndexDescription returns the Description of the gist associated with the index.
func (c *Container) IndexDescription(index *core.QModelIndex) string {
	return index.Data(int(core.Qt__DisplayRole)).ToString()
}

// Remove removes an item from the list by its gist ID.
func (c *Container) Remove(id string) {
	if item, ok := c.items[id]; ok {
		row := c.Row(item)
		c.TakeItem(row)
		delete(c.items, id)
	}
}

// HasID returns true if the id is found in the items
func (c *Container) HasID(id string) bool {
	if _, ok := c.items[id]; ok {
		return true
	}
	return false
}
