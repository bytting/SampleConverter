

package main

import (
	"bufio"
	"encoding/json"
	"os"
)

// Structure representing a sample writer
type SampleWriterJson struct {

	JsonFile      string
        UseScientific bool
	fd            *os.File
	fw            *bufio.Writer
        sep string
}

// Create a new sample writer
func NewSampleWriterJson(jsonFile string, useScientific bool) (SampleWriter, error) {

        // Initialize a sample writer
	sw := new(SampleWriterJson)
	sw.JsonFile = jsonFile
        sw.UseScientific = useScientific
        sw.sep = ""

	var err error
	sw.fd, err = os.Create(sw.JsonFile)
	if err != nil {
		return nil, err
	}

	sw.fw = bufio.NewWriter(sw.fd)
        sw.fw.WriteString("[\n")

	return sw, nil
}

// Write a sample to the json file
func (sw *SampleWriterJson) Write(s *Sample) error {

        // Set the number format
	/*mod := byte('f')
	if sw.UseScientific {
		mod = byte('E')
	}*/

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	sw.fw.WriteString(sw.sep + string(b))

        if len(sw.sep) == 0 {
                sw.sep = ","
        }

	return nil
}

// Finish the json file
func (sw *SampleWriterJson) Close() error {

        sw.fw.WriteString("\n]")
	sw.fw.Flush()
	sw.fd.Close()

	return nil
}
