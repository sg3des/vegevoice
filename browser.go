package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/addrs"
)

var ui *UserInterface

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	ReadConf()
	ReadSettings(conf.Webkit)
	// GlobalSettings()

	go addrs.ReadUrls(path.Join(os.Getenv("XDG_CONFIG_HOME"), "vegevoice", "addrs.list"))
	addrs.SetMaxItems(10)

	gtk.Init(nil)

	ui = CreateUi()
	ui.NewTab(conf.VegeVoice.StartPage)

	gtk.Main()
}

func downloadIcon(uri string) string {
	response, e := http.Get(uri)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	file, err := os.Create(path.Base(uri))
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	return path.Base(uri)
}
