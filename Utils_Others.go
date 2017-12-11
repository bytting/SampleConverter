// +build !linux

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
	"os"
	"path/filepath"
)

func ConfigFile() string {

	p := filepath.Join(HomeDir(), "sampleconverter")
	os.MkdirAll(p, 0777)
	return filepath.Join(p, "settings.json")
}

func PluginDir() string {

	p := filepath.Join(HomeDir(), "sampleconverter", "plugins")
	os.MkdirAll(p, 0777)
	return p
}
