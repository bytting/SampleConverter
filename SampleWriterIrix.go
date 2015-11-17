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
// Dag Robole (dag D0T robole AT gmail D0T com)

package main

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//const validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_().,"

// SampleWriterIrix Structure representing a sample writer
type SampleWriterIrix struct {
	KmlFile       string
	KmzFile       string
	UseScientific bool
	UseLabels     bool
	fd            *os.File
	fw            *bufio.Writer
}

// NewSampleWriterIrix Create a new sample writer
func NewSampleWriterIrix(kmzFile string, useScientific, useLabels bool) (SampleWriter, error) {

	// Initialize a sample writer
	sw := new(SampleWriterIrix)

	// Replace local characters from basename. Google Earth doesn't like them
	base := filepath.Base(kmzFile)
	newBase := ""
	for _, r := range base {
		if !strings.ContainsAny(validChars, string(r)) {
			newBase += "_"
		} else {
			newBase += string(r)
		}
	}
	sw.KmzFile = filepath.Join(filepath.Dir(kmzFile), string(os.PathSeparator), newBase)

	ext := filepath.Ext(sw.KmzFile)
	sw.KmlFile = strings.TrimSuffix(sw.KmzFile, ext) + ".kml"
	sw.UseScientific = useScientific
	sw.UseLabels = useLabels

	var err error
	sw.fd, err = os.Create(sw.KmlFile)
	if err != nil {
		return nil, err
	}

	sw.fw = bufio.NewWriter(sw.fd)

	// Add styles to the kml file
	var s Style
	sw.fw.WriteString(xml.Header)
	sw.fw.WriteString("\n<kml>\n  <Document>\n")

	colors := [...]string{"FFF0AA14", "FF78FFB4", "FF78FFF0", "FF14B4FF", "FF143CFF"}

	for i := 0; i < 5; i++ {
		s.ID = strconv.Itoa(i)
		s.IconStyle.Icon.Href = "files/donut.png"
		s.IconStyle.Scale = "0.5"
		s.IconStyle.Color = colors[i]
		s.LabelStyle.Scale = "0.5"
		b, err := xml.MarshalIndent(s, "    ", "    ")
		if err != nil {
			sw.fd.Close()
			os.Remove(sw.KmlFile)
			return nil, err
		}
		sw.fw.WriteString(string(b) + "\n")
	}

	return sw, nil
}

// Write Write a sample to the kml file
func (sw *SampleWriterIrix) Write(s *Sample) error {

	if strings.ToLower(s.Unit) != "sv/h" {
		return errors.New("Irix format requires unit to be \"Sv/h\"")
	}
	var p Placemark
	var styleID int

	// Calculate the style id for this sample
	if s.Value <= 1 {
		styleID = 0
	} else if s.Value <= 5 {
		styleID = 1
	} else if s.Value <= 10 {
		styleID = 2
	} else if s.Value <= 20 {
		styleID = 3
	} else {
		styleID = 4
	}

	// Set the number format
	mod := byte('f')
	if sw.UseScientific {
		mod = byte('E')
	}

	// Initialize a placemark structure
	if sw.UseLabels {
		p.Name = strconv.FormatFloat(s.Value, mod, -1, 64) + " " + s.Unit
	}
	p.StyleURL = "#" + strconv.Itoa(styleID)
	p.TimeStamp.When = s.Date.Format("2006-01-02T15:04:05")
	p.Point.Coordinates = strconv.FormatFloat(s.Longitude, 'f', -1, 64) + "," +
		strconv.FormatFloat(s.Latitude, 'f', -1, 64)
	p.Description = "Value: " + strconv.FormatFloat(s.Value, mod, -1, 64) + " Sv/h" +
		"\nLatitude: " + strconv.FormatFloat(s.Latitude, 'f', -1, 64) +
		"\nLongitude: " + strconv.FormatFloat(s.Longitude, 'f', -1, 64) +
		"\nTime: " + s.Date.String() + "\nFile: " + filepath.Base(sw.KmzFile)

	// Write placemark structure to the kml file
	b, err := xml.MarshalIndent(p, "    ", "    ")
	if err != nil {
		return err
	}
	sw.fw.WriteString(string(b) + "\n")

	return nil
}

// Close Finish the kml file and zip it to make a kmz file
func (sw *SampleWriterIrix) Close() error {

	sw.fw.WriteString("  </Document>\n</kml>")
	sw.fw.Flush()
	sw.fd.Close()

	err := sw.zipKml()
	if err != nil {
		return err
	}

	os.Remove(sw.KmlFile)

	return nil
}

// Zip the kml file
func (sw *SampleWriterIrix) zipKml() error {

	// Create kmz file
	zfout, err := os.Create(sw.KmzFile)
	if err != nil {
		return err
	}
	defer zfout.Close()

	zw := zip.NewWriter(zfout)
	defer zw.Close() // FIXME: error checking

	kmlFileBase := path.Base(filepath.ToSlash(sw.KmlFile))

	z, err := zw.Create(kmlFileBase)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(sw.KmlFile)
	if err != nil {
		return err
	}

	_, err = z.Write(b)
	if err != nil {
		return err
	}

	z, err = zw.Create("files/donut.png")
	if err != nil {
		return err
	}

	b, err = base64.StdEncoding.DecodeString(PngDonut)
	if err != nil {
		return err
	}

	_, err = z.Write(b)
	if err != nil {
		return err
	}

	return nil
}
