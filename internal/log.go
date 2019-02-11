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
