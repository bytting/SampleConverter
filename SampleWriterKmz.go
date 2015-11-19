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
	"archive/zip"
	"bufio"
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_().,"

// SampleWriterKmz Structure representing a sample writer
type SampleWriterKmz struct {
	KmlFile       string
	KmzFile       string
	MinValue      float64
	MaxValue      float64
	UseScientific bool
	UseLabels     bool
	fd            *os.File
	fw            *bufio.Writer
}

// Style Structure representing a kml style
type Style struct {
	ID        string `xml:"id,attr"`
	IconStyle struct {
		Icon struct {
			Href string `xml:"href"`
		}
		Scale string `xml:"scale"`
		Color string `xml:"color"`
	}
	LabelStyle struct {
		Scale string `xml:"scale"`
	}
}

// Placemark Structure representing a kml placemark
type Placemark struct {
	Name      string `xml:"name"`
	TimeStamp struct {
		When string `xml:"when"`
	}
	Description string `xml:"description"`
	Point       struct {
		Coordinates string `xml:"coordinates"`
	}
	StyleURL string `xml:"styleUrl"`
}

// NewSampleWriterKmz Create a new sample writer
func NewSampleWriterKmz(kmzFile string, useScientific, useLabels bool, minValue, maxValue float64) (SampleWriter, error) {

	// Initialize a sample writer
	sw := new(SampleWriterKmz)

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
	sw.MinValue = minValue
	sw.MaxValue = maxValue
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

	colors := [...]string{"FFF0FF14", "FF78FFF0", "FF14B4FF", "FF1400FF"}

	for i := 0; i < 4; i++ {
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
func (sw *SampleWriterKmz) Write(s *Sample) error {

	var p Placemark
	var styleID int

	// Calculate the style id for this sample
	sector := (sw.MaxValue - sw.MinValue) / 4.0

	if s.Value <= sw.MinValue+sector {
		styleID = 0
	} else if s.Value <= sw.MinValue+sector*2 {
		styleID = 1
	} else if s.Value <= sw.MinValue+sector*3 {
		styleID = 2
	} else {
		styleID = 3
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
	p.Description = "Value: " + strconv.FormatFloat(s.Value, mod, -1, 64) + " " + s.Unit +
		"\nLatitude: " + strconv.FormatFloat(s.Latitude, 'f', -1, 64) +
		"\nLongitude: " + strconv.FormatFloat(s.Longitude, 'f', -1, 64) +
		"\nAltitude: " + strconv.FormatFloat(s.Altitude, 'f', -1, 64) +
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
func (sw *SampleWriterKmz) Close() error {

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
func (sw *SampleWriterKmz) zipKml() error {

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
