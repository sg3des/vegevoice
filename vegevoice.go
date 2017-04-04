package main

import (
	"log"
	"os"
	"path"

	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/urlstorage"
)

var (
	ui *UserInterface

	dirConf string
	dirStrg string
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func resolveWD() {
	envConfig := os.Getenv("XDG_CONFIG_HOME")
	if len(envConfig) == 0 {
		envConfig = path.Join(os.Getenv("HOME"), ".config")
	}
	dirConf = path.Join(envConfig, "vegevoice")
	os.MkdirAll(dirConf, 0755)

	dirStrg = path.Join(os.Getenv("HOME"), ".local", "share", "vegevoice")
	os.MkdirAll(dirConf, 0755)
}

func main() {
	resolveWD()
	ReadConf(dirConf)
	SetCacheDir(dirStrg)

	go urlstorage.Initialize(dirStrg)
	urlstorage.SetMaxItems(10)

	gtk.Init(nil)

	ui = CreateUi()
	for _, u := range urlstorage.GetPinnedTabs() {
		ui.NewTab(u).Pinned = true
	}

	if len(ui.tabs) == 0 {
		if conf.VegeVoice.StartPage != "" {
			ui.NewTab(conf.VegeVoice.StartPage)
		} else {
			ui.NewTab("http://golang.org")
		}
	}

	gtk.Main()
}
