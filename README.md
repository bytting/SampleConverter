# SampleConverter
Convert sample log files to standard formats like csv, json, xml, kmz etc.


# Plugins
Plugins for SampleConverter

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

For the plugin to be valid, it must define six variables: date, latitude, longitude, altitude, value and unit.
These variables should be set to their respective values in the body of the parseLine function.

- date (string)        => The date the sample was taken, in standard ISO format (yyyy-MM-ddThh:mm:ss)
- latitude (decimal)   => The latitude where the sample was taken (GPS format)
- longitude (decimal)  => The longitude where the sample was taken (GPS format)
- altitude (decimal)   => The altitude where the sample was taken (WGS84 format)
- value (decimal)      => The measurement value of the sample
- unit (string)        => The unit of the measurement value
