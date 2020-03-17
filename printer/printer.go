package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/qlik-oss/corectl/pkg/log"
)

func Quietf(format string, a ...interface{}) {
	qprint(fmt.Sprintf(format, a...))
}

func Quiet(a ...interface{}) {
	qprint(a...)
}

func qprint(a ...interface{}) {
	s := log.Appendln(fmt.Sprint(a...))
	fmt.Fprint(os.Stdout, s)
}

// PrintToBashComp handles strings that should be included as options when using auto completion
func PrintToBashComp(str string) {
	if strings.Contains(str, " ") {
		// If string includes whitespaces we need to add quotes
		Quietf("%q", str)
	} else {
		Quiet(str)
	}
}

// PrintAsJSON prints data as JSON to standard out.
func PrintAsJSON(data interface{}) {
	s := log.FormatAsJSON(data)
	fmt.Fprint(os.Stdout, s)
}
