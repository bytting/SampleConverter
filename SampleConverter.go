package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Flag variables
var (
	usePlugin           string
	listPlugins         bool
	setPluginDirectory  string
	showPluginDirectory bool
	useLabels           bool
	useScientific       bool
)

// Settings structure
type Settings struct {
	PluginDirectory string `json:"PluginDirectory"`
}

func main() {

	// Load flags
	flag.StringVar(&usePlugin, "use-plugin", "", "Generate one or more kmz files using the given plugin")
	flag.BoolVar(&listPlugins, "list-plugins", false, "List all available plugins")
	flag.StringVar(&setPluginDirectory, "set-plugin-directory", "", "Set the directory where makekmz looks for plugins")
	flag.BoolVar(&showPluginDirectory, "show-plugin-directory", false, "Show the directory where makekmz looks for plugins")
	flag.BoolVar(&useLabels, "use-labels", false, "Use labels for markers")
	flag.BoolVar(&useScientific, "use-scientific", false, "Use scientific notation for decimal values")
	flag.Parse()

	usr, _ := user.Current()

	// Initialize log
	/*logFile := usr.HomeDir + "/makekmz.log"
	logfd, err := os.Create(logFile)
	if err != nil {
		panic(err)
	}
	defer logfd.Close()
	log.SetOutput(logfd)*/

	// Load settings
	var settings Settings
	settingsFile := filepath.Join(usr.HomeDir, "makekmz.json")

	if !FileExists(settingsFile) {
		settings = Settings{PluginDirectory: filepath.Join(usr.HomeDir, "makekmz-plugins")}
		os.MkdirAll(settings.PluginDirectory, 0777)

	} else {
		sbytes, _ := ioutil.ReadFile(settingsFile)
		json.Unmarshal(sbytes, &settings)
	}

	// Execute operation based on flags
	if listPlugins {

		// Print plugin names to stdout
		files, _ := ioutil.ReadDir(settings.PluginDirectory)
		for _, f := range files {
			ext := filepath.Ext(f.Name())
			if strings.ToLower(ext) == ".js" {
				fmt.Printf("%s\n", strings.TrimSuffix(f.Name(), ext))
			}
		}
		os.Exit(0)

	} else if showPluginDirectory {

		// Print current plugin directory to stdout
		fmt.Println(settings.PluginDirectory)
		os.Exit(0)

	} else if len(setPluginDirectory) > 0 {

		// Set plugin directory
		dir := filepath.Clean(setPluginDirectory)
		os.MkdirAll(dir, 0777)
		settings.PluginDirectory = dir
		sbytes, _ := json.Marshal(&settings)
		ioutil.WriteFile(settingsFile, sbytes, 0644)
		os.Exit(0)

	} else if len(usePlugin) > 0 {

		// Generate kmz files
		if flag.NArg() < 1 {
			log.Fatalln("No input files given")
		}

		pluginFile := filepath.Join(settings.PluginDirectory, usePlugin+".js")
		if !FileExists(pluginFile) {
			log.Fatalf("Plugin %s does not exist", pluginFile)
		}

		for _, sampFile := range flag.Args() {
			if !FileExists(sampFile) {
				log.Printf("Sampling file %s does not exist", sampFile)
				continue
			}

			log.Printf("Generating file %s with plugin %s\n", sampFile, pluginFile)

			sr, err := NewSampleReader(pluginFile, sampFile)
			if err != nil {
				log.Fatalln(err.Error())
			}
			//defer sr.Close()

			sw, err := NewSampleWriterKmz(sampFile, useScientific, useLabels, sr.MinValue, sr.MaxValue)
			if err != nil {
				log.Fatalln(err.Error())
			}
			//defer sw.Close()

			for {
				s, ok, err := sr.Read()
				if err != nil {
					log.Fatalln(err.Error())
				}

				if !ok {
					break
				}

				err = sw.Write(s)
				if err != nil {
					log.Fatalln(err.Error())
				}

				//log.Printf("%s %f %f %f %s\n", s.Date, s.Latitude, s.Longitude, s.Value, s.Unit)
			}
			sr.Close()
			sw.Close()
		}
		os.Exit(0)

	} else {
		log.Fatalln("Missing arguments")
		os.Exit(1)
	}
}
