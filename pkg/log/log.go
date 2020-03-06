package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// QUIET, VERBOSE, TRAFFIC, JSON

type logLevel int

func (l logLevel) String() string {
	switch l {
	case quiet:
		return ""
	case fatal:
		return "Error"
	case err:
		return "Error"
	case warn:
		return ""
	case info:
		return ""
	case verbose:
		return ""
	}
	return ""
}

const (
	quiet logLevel = iota
	fatal
	err
	warn
	info
	verbose
)

type logBuffer struct {
	levels   []logLevel
	messages []string
}

func newBuffer() *logBuffer {
	buf := &logBuffer{
		levels:   []logLevel{},
		messages: []string{},
	}
	return buf
}

func (b *logBuffer) add(lvl logLevel, a ...interface{}) {
	b.levels = append(b.levels, lvl)
	b.messages = append(b.messages, fmt.Sprint(a...))
	// If it's fatal, we just want to dump everything in the buffer
	if lvl == fatal {
		buffering = false
		b.flush()
	}
}

func (b *logBuffer) flush() {
	for i, lvl := range b.levels {
		print(lvl, b.messages[i])
	}
	b.levels = []logLevel{}
	b.messages = []string{}
}

var level logLevel

// printJSON represents whether all output should be in JSON format or not.
var printJSON bool

// Traffic represents wether traffic should be printed or not.
var Traffic bool

var buffering bool
var buffer *logBuffer

func init() {
	level = info
	buffer = newBuffer()
	buffering = true
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

func Verboseln(a ...interface{}) {
	println(verbose, a...)
}

func Verbosef(format string, a ...interface{}) {
	printf(verbose, format, a...)
}

func Verbose(a ...interface{}) {
	print(verbose, a...)
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
			fmt.Println(os.Stderr, FormatAsJSON(msg))
		}
	} else {
		prefix := lvl.String()
		if prefix != "" {
			prefix += ": "
		}
		if buffering {
			buffer.add(lvl, a...)
		} else {
			str := prefix + fmt.Sprint(a...)
			fmt.Fprintln(os.Stderr, str)
		}
	}
}

// PrintAsJSON prints data as JSON to standard out. If already encoded as []byte or json.RawMessage it will be reformated with readable indentation
func PrintAsJSON(data interface{}) {
	fmt.Fprintln(os.Stdout, FormatAsJSON(data))
}

// FormatAsJSON is a utility method that formats the supplied data in a readable way without printing anything. If already encoded as []byte or json.RawMessage it will be reformated with readable indentation
func FormatAsJSON(data interface{}) string {
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
	return buffer.String()
}
