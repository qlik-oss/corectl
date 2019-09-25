package log

import (
	"fmt"
)

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
