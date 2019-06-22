package build

import (
	"fmt"
)

// Details represents known data for a given build
type Details struct {
	Version string `json:"version,omitempty"`
	Date    string `json:"date,omitempty"`
}

var version, buildDate string

// String returns build details as a string with formatting
// suitable for console output.
//
// i.e.
// Build Details:
//         Version:        v0.1.0-155-g1a20f8b
//         Date:           2018-11-05-14:33:14-UTC
func String() string {
	return fmt.Sprintf("Build Details:\n\tVersion:\t%s\n\tDate:\t\t%s", version, buildDate)
}

// Data returns build details as a struct
func Data() Details {
	return Details{
		Version: version,
		Date:    buildDate,
	}
}
