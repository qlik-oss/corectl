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
		return "QUIET"
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
	err
	warn
	info
	debug
	trace // not used
)

var level logLevel

// printJSON represents whether all output should be in JSON format or not.
var printJSON bool

// Traffic represents wether traffic should be printed or not.
var Traffic bool

// Init reads the log-related viper flags json, verbose, traffic and quiet and sets the
// internal (log.) level, printJSON and traffic variables accordingly.
func Init() {
	printJSON = viper.GetBool("json")
	if !printJSON {
		Traffic = viper.GetBool("traffic")
		switch {
		case viper.GetBool("quiet"):
			level = quiet
		case viper.GetBool("verbose"):
			level = debug
		default:
			level = info
		}
	} else {
		level = info
	}
}

func init() {
	level = info
}

func Quietln(a ...interface{}) {
	println(quiet, a...)
}

func Quietf(format string, a ...interface{}) {
	printf(quiet, format, a...)
}

func Quiet(a ...interface{}) {
	print(quiet, a...)
}

func Fatalf(format string, a ...interface{}) {
	printf(fatal, format, a...)
	os.Exit(1)
}

func Fatalln(a ...interface{}) {
	println(fatal, a...)
	os.Exit(1)
}

func Fatal(a ...interface{}) {
	print(fatal, a...)
	os.Exit(1)
}

func Errorln(a ...interface{}) {
	println(err, a...)
}

func Errorf(format string, a ...interface{}) {
	printf(err, format, a...)
}

func Error(a ...interface{}) {
	print(err, a...)
}

func Warnln(a ...interface{}) {
	println(warn, a...)
}

func Warnf(format string, a ...interface{}) {
	printf(warn, format, a...)
}

func Warn(a ...interface{}) {
	print(warn, a...)
}

func Infoln(a ...interface{}) {
	println(info, a...)
}

func Infof(format string, a ...interface{}) {
	printf(info, format, a...)
}

func Info(a ...interface{}) {
	print(info, a...)
}

func Debugln(a ...interface{}) {
	println(debug, a...)
}

func Debugf(format string, a ...interface{}) {
	printf(debug, format, a...)
}

func Debug(a ...interface{}) {
	print(debug, a...)
}

func printf(lvl logLevel, format string, a ...interface{}) {
	print(lvl, fmt.Sprintf(format, a...))
}

func println(lvl logLevel, a ...interface{}) {
	print(lvl, fmt.Sprintln(a...))
}

// print handles all the printing.
// If it is called for a higher log level than what is set, it will not print.
// If printJSON is set it will print as json instead, using the log level as the key.
func print(lvl logLevel, a ...interface{}) {
	if lvl > level { //If level supplied is larger than the current log level, don't print.
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
		str := fmt.Sprint(a...)
		// If the string has a carriage return ('\r') as its first character
		// we need to remove it before adding prefixes and prepend it to the formatted string.
		cr := false
		if str[0] == '\r' {
			cr = true
			str = str[1:]
		}
		prefix := ""
		if lvl != quiet {
			prefix = lvl.String() + ": "
		}
		formatted := format(prefix, str)
		if cr {
			formatted = "\r" + formatted
		}
		fmt.Print(formatted)
	}
}

// format adds the specified prefix to each line contained in the specified string
// returns the resulting string
func format(prefix, str string) string {
	lines := strings.Split(str, "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = prefix + l
		}
	}
	return strings.Join(lines, "\n")
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
