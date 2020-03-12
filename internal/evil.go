package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
	"golang.org/x/crypto/ssh/terminal"
)

// Eval builds a straight table  hypercube based on the supplied argument, evaluates it and prints the result to system out.
func Evil(ctx context.Context, doc *enigma.Doc) {
	ensureModelExists(ctx, doc)
	shell(ctx, doc)
}

func evil(ctx context.Context, doc *enigma.Doc, expr string) {
	val, err := doc.EvaluateEx(ctx, expr)
	if err != nil {
		return
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

func shell(ctx context.Context, doc *enigma.Doc) {
	// fd 0 is stdin
	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}

	defer func() {
		if err := terminal.Restore(0, state); err != nil {
			fmt.Println("warning, failed to restore terminal:", err)
		}
	}()

	in := bufio.NewReader(os.Stdin)
	var buf bytes.Buffer
	for {
		r, _, err := in.ReadRune()
		if err != nil {
			fmt.Println("stdin:", err)
			break
		}
		exit := false

		switch r {
		case '\x03':
			exit = true
			fmt.Println("exit")
		case '\r':
			expr, err := check(ctx, doc, buf.String())
			if err == errBadField {
				fmt.Print("\x1b[2K", expr)
			}
			break
			evil(ctx, doc, expr)
		default:
			fmt.Printf("%s", string(r))
			buf.WriteRune(r)
		}

		if exit {
			break
		}
	}
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
		feedback += "\x1b[4;38;5;1m" + expr[i:j] + "\x1b[0m"
		last = j
	}
	if last != 0 {
		err = errBadField
	}
	feedback += expr[last:]
	return feedback, err
}

var errBadField = fmt.Errorf("found bad field names in expression")
