package main

import (
	"bufio"
	"errors"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
)

// Structure representing a sample reader
type SampleReader struct {

	PluginFile string
	SampleFile string
	MinValue   float64
	MaxValue   float64
	fd         *os.File
	scanner    *bufio.Scanner
	lineNum    int
	vm         *otto.Otto
}

// Create a new sample reader
func NewSampleReader(pluginFile, sampleFile string) (*SampleReader, error) {

        // Initialize a sample reader structure
	sr := new(SampleReader)
	sr.PluginFile = pluginFile
	sr.SampleFile = sampleFile
	sr.MinValue = 0.0
	sr.MaxValue = 0.0

	var err error

	sr.fd, err = os.Open(sr.SampleFile)
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
			return nil, b, nil
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
	b, err := ioutil.ReadFile(sr.PluginFile)
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

	argLineNum, err := sr.vm.ToValue(lineNum)
	if err != nil {
		return nil, err
	}

	argLine, err := sr.vm.ToValue(line)
	if err != nil {
		return nil, err
	}

	retVal, err := sr.vm.Call("parseLine", nil, argLineNum, argLine)
	if err != nil {
		return nil, err
	}

	ret, err := retVal.ToBoolean()
	if err != nil {
		return nil, err
	}

	if !ret {
		return nil, nil
	}

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

	v, err = sr.vm.Get("date")
	if err != nil {
		return nil, err
	}

	if !v.IsDefined() {
		return nil, errors.New("date not defined")
	}

	s.Date, err = v.ToString()
	if err != nil {
		return nil, err
	}

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
