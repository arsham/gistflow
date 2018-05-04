// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/gisty/window"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func main() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	app := widgets.NewQApplication(len(os.Args), os.Args)
	token := os.Getenv("GISTY_TOKEN")
	if token == "" {
		log.Fatal("token cannot be empty")
	}
	g := gist.Service{
		Username: "arsham",
		Token:    token,
	}
	window := window.NewMainWindow(nil, 0)
	window.SetGistService(g)
	window.SetApp(app)
	err := window.Display()
	if err != nil {
		log.Fatal(err)
	}
}
