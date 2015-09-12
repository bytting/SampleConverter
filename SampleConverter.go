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
        "errors"
)

// Flag variables
var (
	usePlugin           string
        useFormat           string
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

        progName := filepath.Base(os.Args[0])

	// Load flags
	flag.StringVar(&usePlugin, "use-plugin", "", "Convert one or more sample files using the given plugin")
	flag.StringVar(&useFormat, "use-format", "kmz", "Use the given output format")
	flag.BoolVar(&listPlugins, "list-plugins", false, "List all available plugins")
        flag.StringVar(&setPluginDirectory, "set-plugin-directory", "", "Set the directory where " + progName + " looks for plugins")
	flag.BoolVar(&showPluginDirectory, "show-plugin-directory", false, "Show the directory where " + progName + " looks for plugins")
	flag.BoolVar(&useLabels, "use-labels", false, "Use labels for markers")
	flag.BoolVar(&useScientific, "use-scientific", false, "Use scientific notation for decimal values")
	flag.Parse()

        useFormat = strings.ToLower(useFormat)

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

			sr, err := NewSampleReader(pluginFile, sampFile)
			if err != nil {
				log.Fatalln(err.Error())
			}

                        sw, err := createSampleWriter(sampFile, useScientific, useLabels, sr.MinValue, sr.MaxValue)
			if err != nil {
				log.Fatalln(err.Error())
			}

			log.Printf("Converting file %s with plugin %s using format \"%s\"\n", sampFile, pluginFile, useFormat)

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

func createSampleWriter(sampleFile string, useScientific, useLabels bool, minValue, maxValue float64) (SampleWriter, error) {

        switch useFormat {
        case "xml":
                return NewSampleWriterXml(sampleFile + ".xml")
        case "kmz":
                return NewSampleWriterKmz(sampleFile + ".kmz", useScientific, useLabels, minValue, maxValue)
        case "json":
                return NewSampleWriterJson(sampleFile + ".json")
        case "csv":
                return NewSampleWriterCsv(sampleFile + ".csv", useScientific)
        }

        return nil, errors.New("Output format not supported: " + useFormat)
}
