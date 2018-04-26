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
