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
// Copyright: (c) 2015 Norwegian Radiation Protection Authority
// Contributors: Dag Robøle (dag D0T robole AT gmail D0T com)

package main

import (
	"bufio"
	"encoding/xml"
	"os"
)

// SampleWriterXML Structure representing a sample writer
type SampleWriterXML struct {
	xmlFile string
	fd      *os.File
	fw      *bufio.Writer
}

// NewSampleWriterXML Create a new sample writer
func NewSampleWriterXML(xmlFile string) (SampleWriter, error) {

	// Initialize a sample writer
	sw := new(SampleWriterXML)
	sw.xmlFile = xmlFile

	var err error
	sw.fd, err = os.Create(sw.xmlFile)
	if err != nil {
		return nil, err
	}

	sw.fw = bufio.NewWriter(sw.fd)
	sw.fw.WriteString(xml.Header)
	sw.fw.WriteString("\n<samples>\n")

	return sw, nil
}

// Write Write a sample to the xml file
func (sw *SampleWriterXML) Write(s *Sample) error {

	b, err := xml.MarshalIndent(s, "  ", "    ")
	if err != nil {
		return err
	}
	sw.fw.WriteString(string(b) + "\n")

	return nil
}

// Close Finish the xml file
func (sw *SampleWriterXML) Close() error {

	sw.fw.WriteString("</samples>")
	sw.fw.Flush()
	sw.fd.Close()

	return nil
}
