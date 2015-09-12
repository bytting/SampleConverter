package main

import (
        "encoding/xml"
        "time"
)

// Structure representing a sample
type Sample struct {

        XMLName   xml.Name `xml:"sample" json:"-"`
        Date      time.Time `xml:"date" json:"date"`
	Latitude  float64 `xml:"latitude" json:"latitude"`
	Longitude float64 `xml:"longitude" json:"longitude"`
	Value     float64 `xml:"value" json:"value"`
	Unit      string `xml:"unit" json:"unit"`
}
