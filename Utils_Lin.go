// +build linux

package main

import (
	"os"
	"path/filepath"
)

func ConfigFile() string {

	p := filepath.Join(HomeDir(), ".config", "sampleconverter")
	os.MkdirAll(p, 0777)
	return filepath.Join(p, "settings.json")
}

func PluginDir() string {

	p := filepath.Join(HomeDir(), "sampleconverter", "plugins")
	os.MkdirAll(p, 0777)
	return p
}
