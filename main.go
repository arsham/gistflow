// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights
// reserved. Use of this source code is governed by the LGPL-v3 License that can
// be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/arsham/gistflow/qt/window"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func main() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	app := widgets.NewQApplication(len(os.Args), os.Args)
	window := window.NewMainWindow(nil, 0)
	err := window.Display(app)
	if err != nil {
		log.Fatal(err)
	}
	widgets.QApplication_Exec()
}
