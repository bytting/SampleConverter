
package main

import (
	"encoding/csv"
	"os"
        "strconv"
)

// Structure representing a sample writer
type SampleWriterCsv struct {

	CsvFile       string
        UseScientific bool
	fd            *os.File
	fw            *csv.Writer
}

// Create a new sample writer
func NewSampleWriterCsv(csvFile string, useScientific bool) (SampleWriter, error) {

        // Initialize a sample writer
	sw := new(SampleWriterCsv)
	sw.CsvFile = csvFile
        sw.UseScientific = useScientific

	var err error
	sw.fd, err = os.Create(sw.CsvFile)
	if err != nil {
		return nil, err
	}

	sw.fw = csv.NewWriter(sw.fd)
	sw.fw.Write([]string {"Date", "Latitude", "Longitude", "Value", "Unit"})

	return sw, nil
}

// Write a sample to the csv file
func (sw *SampleWriterCsv) Write(s *Sample) error {

        // Set the number format
	mod := byte('f')
	if sw.UseScientific {
		mod = byte('E')
	}

        lat := strconv.FormatFloat(s.Latitude, 'f', 8, 64)
        lon := strconv.FormatFloat(s.Longitude, 'f', 8, 64)
        val := strconv.FormatFloat(s.Value, mod, 8, 64)

	sw.fw.Write([]string {s.Date.String(), lat, lon, val, s.Unit})

	return nil
}

// Finish the csv file
func (sw *SampleWriterCsv) Close() error {

	sw.fw.Flush()
	sw.fd.Close()

	return nil
}
