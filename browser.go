package main

import (
	"log"

	"github.com/mattn/go-gtk/gtk"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	gtk.Init(nil)
	ui := CreateUi()
	ui.NewTab("google.com")

	gtk.Main()
}
