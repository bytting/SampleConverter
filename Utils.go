package main

import "os"

// Check if a file exists
func FileExists(filename string) bool {

	_, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}
