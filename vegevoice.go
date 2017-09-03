package main

import (
	"log"
	"os"
	"path"

	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/argum"
	"github.com/sg3des/vegevoice/urlstorage"
)

var args struct {
	URLs   []string
	Config string `argum:"-c,--config"`
}

var ui *UserInterface

func init() {
	log.SetFlags(log.Lshortfile)

	argum.MustParse(&args)
}

func resolveWD() (dirConf, dirStrg string) {
	envConfig := os.Getenv("XDG_CONFIG_HOME")
	if len(envConfig) == 0 {
		envConfig = path.Join(os.Getenv("HOME"), ".config")
	}
	dirConf = path.Join(envConfig, "vegevoice")
	os.MkdirAll(dirConf, 0755)

	dirStrg = path.Join(os.Getenv("HOME"), ".local", "share", "vegevoice")
	os.MkdirAll(dirConf, 0755)

	return
}

func main() {
	dirConf, dirStrg := resolveWD()

	urlstorage.Initialize(dirStrg)
	urlstorage.SetMaxItems(10)

	if args.Config == "" {
		args.Config = path.Join(dirConf, "vegevoice.conf")
	}

	ReadConfFile(args.Config)
	SetCacheDir(dirStrg)

	gtk.Init(nil)

	ui = CreateUi()

	for _, u := range urlstorage.GetPinnedTabs() {
		t := NewTab(u)
		ui.AppendTab(t)
		t.Pin()
	}

	for _, url := range args.URLs {
		ui.AppendTab(NewTab(url))
	}

	if len(args.URLs) == 0 {
		if conf.VegeVoice.StartPage != "" {
			ui.AppendTab(NewTab(conf.VegeVoice.StartPage))
		} else {
			ui.AppendTab(NewTab("https://google.com/"))
		}
	}

	gtk.Main()
}
