package main

import (
	"log"

	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/addrs"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	ReadConf()
	ReadSettings(conf.Webkit)
	GlobalSettings()

	go addrs.ReadUrls("addrs/addrs.list")
	addrs.SetMaxItems(10)

	gtk.Init(nil)

	ui := CreateUi()
	ui.NewTab("google.com")

	gtk.Main()
}
