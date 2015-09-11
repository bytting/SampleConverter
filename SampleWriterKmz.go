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
)

type SampleWriterKmz struct {
	SampleFile    string
	KmlFile       string
	KmzFile       string
	MinValue      float64
	MaxValue      float64
	UseScientific bool
	UseLabels     bool
	fd            *os.File
	fw            *bufio.Writer
}

type Style struct {
	Id        string `xml:"id,attr"`
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

type Placemark struct {
	Name      string `xml:"name"`
	TimeStamp struct {
		When string `xml:"when"`
	}
	Description string `xml:"description"`
	Point       struct {
		Coordinates string `xml:"coordinates"`
	}
	StyleUrl string `xml:"styleUrl"`
}

func NewSampleWriterKmz(sampleFile string, useScientific, useLabels bool, minValue, maxValue float64) (SampleWriter, error) {
	sw := new(SampleWriterKmz)
	sw.SampleFile = sampleFile
	sw.KmlFile = sw.SampleFile + ".kml"
	sw.KmzFile = sw.SampleFile + ".kmz"
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

	err = sw.addStylesToKml()
	if err != nil {
		return nil, err
	}

	return sw, nil
}

func (sw *SampleWriterKmz) Write(s *Sample) error {

	err := sw.addSampleToKml(s)
	if err != nil {
		return err
	}

	return nil
}

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

func (sw *SampleWriterKmz) addStylesToKml() error {

	var s Style
	sw.fw.WriteString(xml.Header)
	sw.fw.WriteString("\n<kml>\n  <Document>\n")

	colors := [...]string{"FFF0FF14", "FF78FFF0", "FF14B4FF", "FF1400FF"}

	for i := 0; i < 4; i++ {
		s.Id = strconv.Itoa(i)
		s.IconStyle.Icon.Href = "files/donut.png"
		s.IconStyle.Scale = "0.5"
		s.IconStyle.Color = colors[i]
		s.LabelStyle.Scale = "0.5"
		b, err := xml.MarshalIndent(s, "    ", "    ")
		if err != nil {
			return err
		}
		sw.fw.WriteString(string(b) + "\n")
	}

	return nil
}

func (sw *SampleWriterKmz) addSampleToKml(s *Sample) error {

	var p Placemark
	var styleId int

	sector := (sw.MaxValue - sw.MinValue) / 4.0

	if s.Value <= sw.MinValue+sector {
		styleId = 0
	} else if s.Value <= sw.MinValue+sector*2 {
		styleId = 1
	} else if s.Value <= sw.MinValue+sector*3 {
		styleId = 2
	} else {
		styleId = 3
	}

	mod := byte('f')
	if sw.UseScientific {
		mod = byte('E')
	}

	if sw.UseLabels {
		p.Name = strconv.FormatFloat(s.Value, mod, -1, 64) + " " + s.Unit
	}
	p.StyleUrl = "#" + strconv.Itoa(styleId)
	p.TimeStamp.When = s.Date
	p.Point.Coordinates = strconv.FormatFloat(s.Longitude, 'f', -1, 64) + "," +
		strconv.FormatFloat(s.Latitude, 'f', -1, 64)
	p.Description = "Value: " + strconv.FormatFloat(s.Value, mod, -1, 64) + " " + s.Unit +
		"\nLatitude: " + strconv.FormatFloat(s.Latitude, 'f', -1, 64) +
		"\nLongitude: " + strconv.FormatFloat(s.Longitude, 'f', -1, 64) +
		"\nTime: " + s.Date + "\nFile: " + filepath.Base(sw.SampleFile)

	b, err := xml.MarshalIndent(p, "    ", "    ")
	if err != nil {
		return err
	}

	sw.fw.WriteString(string(b) + "\n")

	return nil
}

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

	b, err = base64.StdEncoding.DecodeString(PNG_Donut)
	if err != nil {
		return err
	}

	_, err = z.Write(b)
	if err != nil {
		return err
	}

	return nil
}
