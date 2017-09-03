package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sg3des/vegevoice/webkit"
)

//conf structure contains configuration
var conf struct {
	VegeVoice struct {
		StartPage       string `toml:"start-page"`
		HomogeneousTabs bool   `toml:"homogeneous-tabs"`
		HeightTabs      int    `toml:"height-tabs"`
	} `toml:"vegevoice"`

	Webkit map[string]interface{}

	Webview struct {
		Zoomlevel       float64 `toml:"zoom-level"`
		FullContentZoom bool    `toml:"full-content-zoom"`
	}
}

//ReadConfFile set default config values and parse config file
func ReadConfFile(filename string) {
	if _, err := os.Stat(filename); err != nil {
		log.Println(err)
		return
	}

	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		log.Println("failed decode config file", filename, "reason:", err)
		return
	}

	if conf.VegeVoice.HeightTabs == 0 {
		conf.VegeVoice.HeightTabs = -1
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
