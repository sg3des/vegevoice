package main

import (
	"log"
	"os"
	"path"

	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/addrs"
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

	go addrs.ReadUrls(dirConf)
	addrs.SetMaxItems(10)

	gtk.Init(nil)

	ui = CreateUi()
	ui.NewTab(conf.VegeVoice.StartPage)

	gtk.Main()
}
