package main

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/sg3des/vegevoice/webkit"
)

//conf structure contains configuration
var conf struct {
	VegeVoice struct {
		StartPage string `toml:"start-page"`
	} `toml:"vegevoice"`

	Webkit map[string]interface{}

	Webview struct {
		Zoomlevel       float64 `toml:"zoom-level"`
		FullContentZoom bool    `toml:"full-content-zoom"`
	}
}

//ReadConf set default values for configuration and parse config file
func ReadConf(dir string) {
	configfile := path.Join(dir, "vegevoice.conf")

	if _, err := os.Stat(configfile); err != nil {
		log.Println(err)
		return
	}

	if _, err := toml.DecodeFile(configfile, &conf); err != nil {
		log.Println("failed decode config file", configfile, "reason:", err)
		return
	}

	ReadSettings(conf.Webkit)
}

var webkitSettings *webkit.WebSettings

func ReadSettings(settings map[string]interface{}) {
	webkitSettings = webkit.NewWebSettings()

	for name, value := range settings {
		webkitSettings.Set(name, value)
	}
}

func ApplySettings(webview *webkit.WebView) {
	if webkitSettings == nil {
		return
	}

	webview.SetSettings(webkitSettings)

	if conf.Webview.Zoomlevel > 0 {
		webview.SetZoomLevel(conf.Webview.Zoomlevel)
	}
	webview.SetFullContentZoom(conf.Webview.FullContentZoom)
}
