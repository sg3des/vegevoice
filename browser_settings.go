package main

import "github.com/sg3des/vegevoice/webkit"

var webkitSettings *webkit.WebSettings

func ReadSettings(settings map[string]interface{}) {
	webkitSettings = webkit.NewWebSettings()

	for name, value := range settings {
		webkitSettings.Set(name, value)
	}
}

func ApplySettings(webview *webkit.WebView) {
	webview.SetSettings(webkitSettings)

	if conf.Webview.Zoomlevel > 0 {
		webview.SetZoomLevel(conf.Webview.Zoomlevel)
	}
	webview.SetFullContentZoom(conf.Webview.FullContentZoom)
}

func GlobalSettings() {
	// w := &webkit.WebDatabase{}
	// log.Println(w.GetFilename())
	// webkit.SetWebDatabaseDirectoryPath("./testdata/webdatabase")
}
