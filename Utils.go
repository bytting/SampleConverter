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
	"flag"
	"os"
	"os/user"
	"path/filepath"
)

// Check if a file exists
func FileExists(filename string) bool {

	_, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

func ArgumentFiles() []string {

	var allFiles []string
	for _, pattern := range flag.Args() {
		files, _ := filepath.Glob(pattern)
		allFiles = append(allFiles, files...)
	}

	return allFiles
}

func ExecutableFile() string {

	exe, _ := filepath.Abs(os.Args[0])
	return exe
}

func ExecutableDir() string {

	return filepath.Dir(ExecutableFile())
}

func HomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
