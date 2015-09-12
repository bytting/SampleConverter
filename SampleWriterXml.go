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
	"bufio"
	"encoding/xml"
	"os"
)

// Structure representing a sample writer
type sampleWriterXml struct {

	XmlFile       string
	fd            *os.File
	fw            *bufio.Writer
}

// Create a new sample writer
func NewSampleWriterXml(xmlFile string) (SampleWriter, error) {

        // Initialize a sample writer
	sw := new(sampleWriterXml)
	sw.XmlFile = xmlFile

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
func (sw *sampleWriterXml) Write(s *Sample) error {

        // Write placemark structure to the kml file
	b, err := xml.MarshalIndent(s, "  ", "    ")
	if err != nil {
		return err
	}
	sw.fw.WriteString(string(b) + "\n")

	return nil
}

// Finish the xml file
func (sw *sampleWriterXml) Close() error {

	sw.fw.WriteString("</samples>")
	sw.fw.Flush()
	sw.fd.Close()

	return nil
}
