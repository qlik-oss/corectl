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

var LogTraffic bool

type TrafficLogger struct{}

func (l *TrafficLogger) Opened() {}

func (l *TrafficLogger) Sent(message []byte) {
	fmt.Println("-->", string(message[:]))
}

func (l *TrafficLogger) Received(message []byte) {
	fmt.Println("<--", string(message[:]))
}

func (l *TrafficLogger) Closed() {}
