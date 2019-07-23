package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// QUIET, VERBOSE, TRAFFIC, JSON

type logLevel int

func (l logLevel) String() string {
	switch l {
	case quiet:
		return ""
	case fatal:
		return "FATAL"
	case err:
		return "ERROR"
	case warn:
		return "WARN"
	case info:
		return "INFO"
	case debug:
		return "DEBUG"
	case trace:
		return "TRACE"
	}
	return ""
}

const (
	quiet logLevel = iota
	fatal
	err //not used
	warn
	info
	debug
	trace // not used
)

var level logLevel

// printJSON represents whether all output should be in JSON format or not.
var printJSON bool

// traffic represents wheter traffic should be printed or not.
var traffic bool

// Init reads the log-related viper flags json, verbose, traffic and quiet and sets the
// internal (log.) level, printJSON and traffic variables accordingly.
func Init() {
	printJSON = viper.GetBool("json")
	if !printJSON {
		traffic = viper.GetBool("traffic")
		switch {
		case viper.GetBool("quiet"):
			level = quiet
		case viper.GetBool("verbose"):
			level = debug
		default:
			level = info
		}
	}
}

func Quiet(a ...interface{}) {
	print(quiet, a...)
}

func Quietf(format string, a ...interface{}) {
	printf(quiet, format, a...)
}

func Debug(a ...interface{}) {
	print(debug, a...)
}

func Debugf(format string, a ...interface{}) {
	printf(debug, format, a...)
}

func Info(a ...interface{}) {
	print(info, a...)
}

func Infof(format string, a ...interface{}) {
	printf(info, format, a...)
}

func Warn(a ...interface{}) {
	print(warn, a...)
}

func Warnf(format string, a ...interface{}) {
	printf(warn, format, a...)
}

// FatalError prints the supplied message and exits the process with code 1
func Fatal(a ...interface{}) {
	print(err, a...)
	os.Exit(1)
}

// FatalErrorf formats and prints the message according to the applied format and exits the process with code 1
func Fatalf(format string, a...interface{}) {
	printf(err, format, a...)
	os.Exit(1)
}


func printf(lvl logLevel, format string, a ...interface{}) {
	print(lvl, fmt.Sprintf(format, a...))
}

func print(lvl logLevel, a ...interface{}) {
	if lvl > level { //If level supplied is larger than the current log level, don't print
		return
	}
	if printJSON {
		if lvl != quiet {
			msg := map[string]string{
				strings.ToLower(lvl.String()): fmt.Sprint(a...),
			}
			PrintAsJSON(msg)
		}
	} else {
		str := ""
		if lvl != quiet {
			str += lvl.String() + ": "
		}
		str += fmt.Sprint(a...)
		fmt.Println(str)
	}
}


// PrintAsJSON prints data as JSON. If already encoded as []byte or json.RawMessage it will be reformated with readable indentation
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
		Fatal(err)
	}
	var buffer bytes.Buffer
	json.Indent(&buffer, jsonBytes, "", "  ")
	fmt.Println(buffer.String())
}
