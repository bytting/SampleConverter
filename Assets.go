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

// Howto for plugins
const TXT_Plugin_Howto = `
Plugins are standalone javascript files that are designed to parse a single line 
from a sampling log file. In order to be used as a plugin, it must implement a function
called "parseLine" with the following signature:

function parseLine(lineNumber, line)
{
        // Implementation goes here...
}

The input parameter "lineNumber" is the line number from the sampling file being parsed, starting with 1.
The input parameter "line" is the line of text from the sampling file.

The parseLine function shall return a boolean (true or false), indicating wether the current line should be
skipped or not. Returning false will instruct the converter to skip on to the next line.

For the plugin to be valid, it must define five variables: date, latitude, longitude, value and unit.
These variables should be set to their respective values in the body of the parseLine function.

- date (string)        => The date the sample was taken, in standard ISO format (yyyy-MM-ddThh:mm:ss)
- latitude (decimal)   => The latitude where the sample was taken (GPS format)
- longitude (decimal)  => The longitude where the sample was taken (GPS format)
- value (decimal)      => The measurement value of the sample
- unit (string)        => The unit of the measurement value
`

// Base64 encoded PNG image
const PNG_Donut = `
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
D49TzrLx8HXTAAAAAElFTkSuQmCC`
