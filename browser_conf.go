package main

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

//conf structure contains configuration
var conf struct {
	VegeVoice struct {
		StartPage string `toml:"start-page"`
	}
	Webkit  map[string]interface{}
	Webview struct {
		Zoomlevel float64 `toml:"zoom-level"`
	}
}

//ReadConf set default values for configuration and parse config file
func ReadConf() {
	//parse config files
	for _, configfile := range []string{
		path.Join(os.Getenv("XDG_CONFIG_HOME"), "vegevoice", "vegevoice.conf"),
		path.Join(os.Getenv("HOME"), ".config", "vegevoice", "vegevoice.conf"),
		"vegevoice.conf",
	} {
		if _, err := os.Stat(configfile); err != nil {
			continue
		}

		if _, err := toml.DecodeFile(configfile, &conf); err != nil {
			log.Println("failed decode config file", configfile, "reason:", err)
			continue
		}
		break
	}
}
