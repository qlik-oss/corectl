package internal

import (
	"context"
	"fmt"
	"os"
	"io/ioutil"
	"strings"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// Eval builds a straight table  hypercube based on the supplied argument, evaluates it and prints the result to system out.
func Evil(ctx context.Context, doc *enigma.Doc) {

	ensureModelExists(ctx, doc)

	expr, err := shell()
	if err != nil {
		log.Fatal("EVIL COULD NOT SHELL!")
	}
	val, err := doc.EvaluateEx(ctx, string(expr))
	if err != nil {
		log.Fatal("EVIL FAILED: ", err.Error())
	}
	switch {
	case val.IsNumeric:
		fmt.Println("\x1b[38;5;48m", val.Number, "\x1b[38;0m")
	case val.Text == "-", strings.HasPrefix(val.Text, "Error:"):
		fmt.Println("\x1b[38;5;001m", val.Text, "\x1b[38;0m")
	default:
		fmt.Println("\x1b[38;5;126m", val.Text, "\x1b[38;0m")
	}
}

func shell() ([]byte, error) {
	fmt.Print("> ")
	return ioutil.ReadAll(os.Stdin)
}
