package main

import (
	"log"
	"os"
	"path"

	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/addrs"
)

var UI *UserInterface

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	ReadConf()
	ReadSettings(conf.Webkit)
	GlobalSettings()

	go addrs.ReadUrls(path.Join(os.Getenv("XDG_CONFIG_HOME"), "vegevoice", "addrs.list"))
	addrs.SetMaxItems(10)

	gtk.Init(nil)

	UI = CreateUi()
	UI.NewTab(conf.VegeVoice.StartPage)

	gtk.Main()
}
