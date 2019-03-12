package internal

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger global logger
var Logger logrus.Logger

var defaultFields = &logrus.FieldMap{
	logrus.FieldKeyTime:  "timestamp",
	logrus.FieldKeyLevel: "logseverity",
	logrus.FieldKeyMsg:   "message",
}

var logToJSON bool
var disableTimestamps bool

// InitLogger function initializing the logger module
func InitLogger(logLevel string, logJSON bool, testFlag bool) {
	// Create a new logger
	Logger = *logrus.New()

	// If test flag is set we should skip adding timestamps to log entries
	disableTimestamps = testFlag

	// Set up logging formatters
	if logJSON {
		logToJSON = true
		Logger.SetFormatter(&logrus.JSONFormatter{
			FieldMap:         *defaultFields,
			DisableTimestamp: disableTimestamps,
		})
		Logger.Debug("Logging format set to JSON")
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "15:04:05",
			DisableTimestamp:       disableTimestamps,
			DisableSorting:         true,
			DisableLevelTruncation: true,
			FieldMap:               *defaultFields,
		})
	}

	// Set logging level
	var level, err = logrus.ParseLevel(logLevel)

	if err != nil {
		Logger.Warnf("Incorrect logging level set: %s", logLevel)
		level = logrus.InfoLevel
	}

	Logger.Debugf("Logging level set to %s", level)
	Logger.SetLevel(level)

	// Also log traffic if logging level is set to Debug or Trace
	if level == logrus.DebugLevel || level == logrus.TraceLevel {
		LogTraffic = true
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
	if logToJSON {
		fmt.Println(getTrafficLogEntry("sent", message))
	} else {
		Logger.WithFields(logrus.Fields{"traffic": "sent"}).Debug(string(message))
	}
}

// Received implements Received() method in enigma-go TrafficLogger interface
func (TrafficLogger) Received(message []byte) {
	if logToJSON {
		fmt.Println(getTrafficLogEntry("recv", message))
	} else {
		Logger.WithFields(logrus.Fields{"traffic": "recv"}).Debug(string(message))
	}
}

// Closed implements Closed() method in enigma-go TrafficLogger interface
func (TrafficLogger) Closed() {}

// Function for generating a traffic log entry in json
func getTrafficLogEntry(direction string, message []byte) string {
	var entry []byte
	if disableTimestamps {
		entry = []byte(fmt.Sprintf(`{"logseverity":"debug","message":"%s","traffic":%v}`, direction, string(message)))
	} else {
		entry = []byte(fmt.Sprintf(`{"logseverity":"debug","message":"%s","traffic":%v,"timestamp":"%s"}`, direction, string(message), time.Now().Format(time.RFC3339)))
	}

	return string(entry)
}
