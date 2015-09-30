// +build linux

package main

import (
	"os"
	"path/filepath"
)

// ConfigFile Get the config directory for Linux systems
func ConfigFile() string {

	p := filepath.Join(HomeDir(), ".config", "sampleconverter")
	os.MkdirAll(p, 0777)
	return filepath.Join(p, "settings.json")
}

// PluginDir Get the default plugin directory
func PluginDir() string {

	p := filepath.Join(HomeDir(), "sampleconverter", "plugins")
	os.MkdirAll(p, 0777)
	return p
}
