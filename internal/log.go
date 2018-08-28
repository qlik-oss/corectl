package internal

import "fmt"

var QliVerbose bool

func LogVerbose(message string) {
	if QliVerbose {
		fmt.Println(message)
	}
}
