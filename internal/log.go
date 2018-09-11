package internal

import "fmt"

// QliVerbose is set to true to indicate that verbose logging should be enabled.
var QliVerbose bool

// LogVerbose prints the supplied message to system out if verbose logging is enabled.
func LogVerbose(message string) {
	if QliVerbose {
		fmt.Println(message)
	}
}
