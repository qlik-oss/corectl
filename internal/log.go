package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

// PrintVerbose is set to true to indicate that verbose logging should be enabled.
var PrintVerbose bool

// PrintJSON represents whether all output should be in JSON format or not.
var PrintJSON bool

// InitLogOutput reads the viper flags json, verbose and traffic and sets the
// internal loggin variables PrintJSON, PrintVerbose and LogTraffic accordingly.
func InitLogOutput() {
	PrintJSON = viper.GetBool("json")
	if !PrintJSON {
		PrintVerbose = viper.GetBool("verbose")
		LogTraffic = viper.GetBool("traffic")
	}
}

// LogVerbose prints the supplied message to system out if verbose logging is enabled.
func LogVerbose(message string) {
	if PrintVerbose {
		fmt.Println(message)
	}
}

// LogTraffic is set to true when websocket traffic should be printed to stdout
var LogTraffic bool

// TrafficLogger is a struct implementing the TrafficLogger interface in enigma-go
type TrafficLogger struct{}

// Opened implements Opened() method in enigma-go TrafficLogger interface
func (TrafficLogger) Opened() {}

// Sent implements Sent() method in enigma-go TrafficLogger interface
func (TrafficLogger) Sent(message []byte) {
	fmt.Println("-->", string(message))
}

// Received implements Received() method in enigma-go TrafficLogger interface
func (TrafficLogger) Received(message []byte) {
	fmt.Println("<--", string(message))
}

// Closed implements Closed() method in enigma-go TrafficLogger interface
func (TrafficLogger) Closed() {}
