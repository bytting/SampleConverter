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
	"errors"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"time"
)

// SampleReader Structure representing a sample reader
type SampleReader struct {
	pluginFile string
	sampleFile string
	MinValue   float64
	MaxValue   float64
	fd         *os.File
	scanner    *bufio.Scanner
	lineNum    int
	vm         *otto.Otto
}

// NewSampleReader Create a new sample reader
func NewSampleReader(pluginFile, sampleFile string) (*SampleReader, error) {

	// Initialize a sample reader structure
	sr := new(SampleReader)
	sr.pluginFile = pluginFile
	sr.sampleFile = sampleFile
	sr.MinValue = 0.0
	sr.MaxValue = 0.0

	var err error

	sr.fd, err = os.Open(sr.sampleFile)
	if err != nil {
		return nil, err
	}

	sr.scanner = bufio.NewScanner(sr.fd)
	sr.lineNum = 0

	// Create a otto javascript runtime
	sr.vm, err = sr.createPluginRuntime()
	if err != nil {
		return nil, err
	}

	// Scan the sample file to find the min and max measurement values (God, make it stop)
	// The kmz sample writer need these values to calculate the correct placemark colors
	initialized := false

	for sr.scanner.Scan() {

		sr.lineNum++
		samp, err := sr.execPlugin(sr.scanner.Text(), sr.lineNum)
		if err != nil {
			return nil, err
		}

		if samp == nil {
			continue
		}

		if !initialized {
			initialized = true
			sr.MinValue = samp.Value
			sr.MaxValue = samp.Value
		} else {
			if samp.Value < sr.MinValue {
				sr.MinValue = samp.Value
			}
			if samp.Value > sr.MaxValue {
				sr.MaxValue = samp.Value
			}
		}
	}

	err = sr.scanner.Err()
	if err != nil {
		return nil, err
	}

	// Reset the scanner for later use
	sr.fd.Seek(0, 0)
	sr.scanner = bufio.NewScanner(sr.fd)
	sr.lineNum = 0

	return sr, nil
}

// Read the next line from the sample file using a javascript plugin and make a sample structure from it
func (sr *SampleReader) Read() (*Sample, bool, error) {

	for {
		b := sr.scanner.Scan()
		if !b {
			break
		}

		err := sr.scanner.Err()
		if err != nil {
			return nil, false, err
		}

		sr.lineNum++

		sample, err := sr.execPlugin(sr.scanner.Text(), sr.lineNum)
		if err != nil {
			return nil, false, err
		}

		if sample == nil {
			continue
		}

		return sample, true, nil
	}

	return nil, false, nil
}

// Close the sample reader and clean up
func (sr *SampleReader) Close() error {

	sr.fd.Close()
	return nil
}

// Create a javascript runtime
func (sr *SampleReader) createPluginRuntime() (*otto.Otto, error) {

	// Read plugin file
	b, err := ioutil.ReadFile(sr.pluginFile)
	if err != nil {
		return nil, err
	}

	// Create runtime and load plugin file
	vm := otto.New()
	_, err = vm.Run(string(b))
	if err != nil {
		return nil, err
	}

	return vm, nil
}

// Execute plugin and extract a sample
func (sr *SampleReader) execPlugin(line string, lineNum int) (*Sample, error) {

	// Prepare arguments
	argLineNum, err := sr.vm.ToValue(lineNum)
	if err != nil {
		return nil, err
	}

	argLine, err := sr.vm.ToValue(line)
	if err != nil {
		return nil, err
	}

	// Execute plugin
	retVal, err := sr.vm.Call("parseLine", nil, argLineNum, argLine)
	if err != nil {
		return nil, err
	}

	// Extract and evaluate return value
	ret, err := retVal.ToBoolean()
	if err != nil {
		return nil, err
	}

	if !ret {
		return nil, nil
	}

	// Extract a full sample from javascript runtime
	sample, err := sr.getSample()
	if err != nil {
		return nil, err
	}

	return sample, nil
}

// Helper function to populate a sample structure with a single sample
func (sr *SampleReader) getSample() (*Sample, error) {

	var err error
	var v otto.Value

	s := new(Sample)

	// Extract date field from javascript runtime
	v, err = sr.vm.Get("date")
	if err != nil {
		return nil, err
	}

	if !v.IsDefined() {
		return nil, errors.New("date not defined")
	}

	ds, err := v.ToString()
	if err != nil {
		return nil, err
	}

	s.Date, err = time.Parse("2006-01-02T15:04:05", ds)
	if err != nil {
		return nil, err
	}

	// Extract latitude field from javascript runtime
	v, err = sr.vm.Get("latitude")
	if err != nil {
		return nil, err
	}

	if !v.IsDefined() {
		return nil, errors.New("latitude not defined")
	}

	s.Latitude, err = v.ToFloat()
	if err != nil {
		return nil, err
	}

	// Extract longitude field from javascript runtime
	v, err = sr.vm.Get("longitude")
	if err != nil {
		return nil, err
	}

	if !v.IsDefined() {
		return nil, errors.New("longitude not defined")
	}

	s.Longitude, err = v.ToFloat()
	if err != nil {
		return nil, err
	}

	// Extract altitude field from javascript runtime
	v, err = sr.vm.Get("altitude")
	if err != nil {
		return nil, err
	}

	if !v.IsDefined() {
		return nil, errors.New("altitude not defined")
	}

	s.Altitude, err = v.ToFloat()
	if err != nil {
		return nil, err
	}

	// Extract value field from javascript runtime
	v, err = sr.vm.Get("value")
	if err != nil {
		return nil, err
	}

	if !v.IsDefined() {
		return nil, errors.New("value not defined")
	}

	s.Value, err = v.ToFloat()
	if err != nil {
		return nil, err
	}

	// Extract unit field from javascript runtime
	v, err = sr.vm.Get("unit")
	if err != nil {
		return nil, err
	}

	if !v.IsDefined() {
		return nil, errors.New("unit not defined")
	}

	s.Unit, err = v.ToString()
	if err != nil {
		return nil, err
	}

	return s, nil
}
