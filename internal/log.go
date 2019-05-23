package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

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

// FatalError prints the supplied message and exists the process with code 1
func FatalError(fatalMessage ...interface{}) {
	if PrintJSON {
		errMsg := map[string]string{
			"error": fmt.Sprint(fatalMessage...),
		}
		PrintAsJSON(errMsg)
	} else {
		fmt.Println("ERROR", fmt.Sprint(fatalMessage...))
	}
	os.Exit(1)
}

// PrintAsJSON prints data as JSON. If already encoded as []byte or json.RawMessage it will be reformated with proper indentation
func PrintAsJSON(data interface{}) {
	var jsonBytes json.RawMessage
	var err error
	switch v := data.(type) {
	case json.RawMessage:
		jsonBytes = v
	case []byte:
		jsonBytes = json.RawMessage(v)
	default:
		jsonBytes, err = json.Marshal(data)
	}
	if err != nil {
		FatalError(err)
	}
	var buffer bytes.Buffer
	json.Indent(&buffer, jsonBytes, "", "  ")
	fmt.Println(buffer.String())
}
