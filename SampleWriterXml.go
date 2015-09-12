
package main

import (
	"bufio"
	"encoding/xml"
	"os"
)

// Structure representing a sample writer
type SampleWriterXml struct {

	XmlFile       string
        UseScientific bool
	fd            *os.File
	fw            *bufio.Writer
}

// Create a new sample writer
func NewSampleWriterXml(xmlFile string, useScientific bool) (SampleWriter, error) {

        // Initialize a sample writer
	sw := new(SampleWriterXml)
	sw.XmlFile = xmlFile
        sw.UseScientific = useScientific

	var err error
	sw.fd, err = os.Create(sw.XmlFile)
	if err != nil {
		return nil, err
	}

	sw.fw = bufio.NewWriter(sw.fd)
	sw.fw.WriteString(xml.Header)
	sw.fw.WriteString("\n<samples>\n")

	return sw, nil
}

// Write a sample to the xml file
func (sw *SampleWriterXml) Write(s *Sample) error {

        // Set the number format
	/*mod := byte('f')
	if sw.UseScientific {
		mod = byte('E')
	}*/

        // Write placemark structure to the kml file
	b, err := xml.MarshalIndent(s, "  ", "    ")
	if err != nil {
		return err
	}
	sw.fw.WriteString(string(b) + "\n")

	return nil
}

// Finish the xml file
func (sw *SampleWriterXml) Close() error {

	sw.fw.WriteString("</samples>")
	sw.fw.Flush()
	sw.fd.Close()

	return nil
}
