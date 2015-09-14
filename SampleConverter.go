/*
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2015:
// Dag Robøle (dag D0T robole AT gmail D0T com)

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
        "errors"
)

var progName string
var version string = "0.2"

// Flag variables
var (
	usePlugin           string
        useFormat           string
	listPlugins         bool
	setPluginDirectory  string
	showPluginDirectory bool
        showVersion         bool
	useLabels           bool
	useScientific       bool
        showHowto           bool
)

// Settings structure
type Settings struct {

	PluginDirectory string `json:"PluginDirectory"`
}

func init() {

        // Save program name
        progName = filepath.Base(os.Args[0])

	// Load flags
	flag.StringVar(&usePlugin, "use-plugin", "", "Convert one or more sample files using the given plugin")
	flag.StringVar(&useFormat, "use-format", "kmz", "Use the given output format")
	flag.BoolVar(&listPlugins, "list-plugins", false, "List all available plugins")
        flag.StringVar(&setPluginDirectory, "set-plugin-directory", "", "Set the directory where " + progName + " looks for plugins")
	flag.BoolVar(&showPluginDirectory, "show-plugin-directory", false, "Show the directory where " + progName + " looks for plugins")
	flag.BoolVar(&showVersion, "version", false, "Show " + progName + " version")
	flag.BoolVar(&useLabels, "use-labels", false, "Use labels for markers")
	flag.BoolVar(&useScientific, "use-scientific", false, "Use scientific notation for decimal values")
	flag.BoolVar(&showHowto, "show-plugin-howto", false, "Show the plugin howto")
}

func main() {

	flag.Parse()

        useFormat = strings.ToLower(useFormat)

	// Load settings
	var settings Settings
	settingsFile := filepath.Join(ExecutableDir(), "settings.json")

	if !FileExists(settingsFile) {

		settings = Settings{PluginDirectory: filepath.Join(ExecutableDir(), "plugins")}
		os.MkdirAll(settings.PluginDirectory, 0777)
		sbytes, _ := json.Marshal(&settings)
		ioutil.WriteFile(settingsFile, sbytes, 0644)

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

	} else if showPluginDirectory {

		// Print current plugin directory to stdout
		fmt.Println(settings.PluginDirectory)

	} else if showVersion {

		// Print version
                fmt.Println(version)

	} else if showHowto {

                // Show plugin howto
                fmt.Println(TXT_Plugin_Howto)

	} else if len(setPluginDirectory) > 0 {

		// Set plugin directory
		dir := filepath.Clean(setPluginDirectory)
		os.MkdirAll(dir, 0777)
		settings.PluginDirectory = dir
		sbytes, _ := json.Marshal(&settings)
		ioutil.WriteFile(settingsFile, sbytes, 0644)

	} else if len(usePlugin) > 0 {

		// Convert sample files
		if flag.NArg() < 1 {
                        log.Fatalln("ERROR: No input files given")
		}

		pluginFile := filepath.Join(settings.PluginDirectory, usePlugin+".js")
		if !FileExists(pluginFile) {
                        log.Fatalf("ERROR: Plugin %s does not exist", pluginFile)
		}

                sampleFiles := ArgumentFiles()
                if len(sampleFiles) == 0 {
                        log.Fatalln("ERROR: No valid input files given")
                }

	        for _, sampleFile := range ArgumentFiles() {

                        if !FileExists(sampleFile) {
                                fmt.Errorf("ERROR: Sampling file %s does not exist", sampleFile)
                                continue
                        }

                        err := convertSampleFile(pluginFile, sampleFile)
                        if err != nil {
                                log.Fatalln(err.Error())
                        }
                }

	} else {
                log.Fatalf("ERROR: Missing arguments.\nUse \"%s -h\" for a description of possible arguments", progName)
	}
}

// Convert a single sample file
func convertSampleFile(pluginFile, sampleFile string) error {

	fmt.Printf("Converting file '%s' with plugin '%s' using format '%s'\n", filepath.Base(sampleFile), filepath.Base(pluginFile), useFormat)

	sr, err := NewSampleReader(pluginFile, sampleFile)
	if err != nil {
                return err
	}
	defer sr.Close()

        sw, err := createSampleWriter(sampleFile, sr.MinValue, sr.MaxValue)
	if err != nil {
                return err
	}
	defer sw.Close()

	for {
		s, more, err := sr.Read()
		if err != nil {
                        return err
		}

		if !more {
			break
		}

		err = sw.Write(s)
		if err != nil {
                        return err
		}
	}

        return nil
}

// Create the correct sample writer based on the useFormat flag
func createSampleWriter(sampleFile string, minValue, maxValue float64) (SampleWriter, error) {

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
