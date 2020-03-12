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
	cexpr, err := check(ctx, doc, string(expr))
	if err != nil {
		if cexpr != "" {
			fmt.Print(cexpr)
		}
		log.Fatalln(err)
	}
	val, err := doc.EvaluateEx(ctx, cexpr)
	if err != nil {
		log.Fatalln(err)
	}
	switch {
	case val.IsNumeric:
		fmt.Println("\x1b[38;5;48m", val.Number, "\x1b[38;0m")
	case val.Text == "-", strings.HasPrefix(val.Text, "Error:"):
		fmt.Println("\x1b[38;5;1m", val.Text, "\x1b[38;0m")
	default:
		fmt.Println("\x1b[38;5;126m", val.Text, "\x1b[38;0m")
	}
}

func shell() ([]byte, error) {
	fmt.Print("> ")
	return ioutil.ReadAll(os.Stdin)
}

func check(ctx context.Context, doc *enigma.Doc, expr string) (string, error) {
	var err error
	str, bads, err := doc.CheckNumberOrExpression(ctx, expr)
	if err != nil {
		return "", err
	}
	if str != "" {
		return "", fmt.Errorf(str)
	}
	var feedback string
	var last int
	for _, bad := range bads {
		i, j := bad.From, bad.From + bad.Count
		feedback += expr[:i]
		feedback += "\x1b[48;5;1m\x1b[5m" + expr[i:j] + "\x1b[0m"
		last = j
	}
	if last != 0 {
		err = errBadField
	}
	feedback += expr[last:]
	return feedback, err
}

var errBadField = fmt.Errorf("found bad field names in expression")
