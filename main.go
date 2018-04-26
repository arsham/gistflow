// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/gisty/window"
)

func main() {
	token := "4fe4218d35fa707d9a964142bd120ea5d37428e3"
	gs := gist.Service{
		Username: "arsham",
		Token:    token,
	}
	ws := window.Service{
		GistService: gs,
	}
	err := ws.MainWindow()
	if err != nil {
		log.Fatal(err)
	}
}
