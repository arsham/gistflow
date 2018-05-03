// Copyright 2018 Arsham Shirvani <arshamshirvani@gmail.com>. All rights reserved.
// Use of this source code is governed by the MIT license
// License that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/arsham/gisty/gist"
	"github.com/arsham/gisty/window"
)

func main() {
	token := os.Getenv("GISTY_TOKEN")
	if token == "" {
		log.Fatal("token cannot be empty")
	}
	gs := gist.Service{
		Username: "arsham",
		Token:    token,
	}
	ws := window.MainWindow{
		GistService: gs,
	}
	err := ws.Display()
	if err != nil {
		log.Fatal(err)
	}
}
