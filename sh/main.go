package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	state, err := terminal.MakeRaw(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer terminal.Restore(0, state)
	sh := newShell()
	sh.run()
}

type shell struct {
	cur int
	buf *bytes.Buffer
	in  *bufio.Reader
}

func newShell() *shell {
	return &shell{
		cur: 0,
		buf: &bytes.Buffer{},
		in:  bufio.NewReader(os.Stdin),
	}
}

func (s *shell) next() (rune, error) {
	r, _, err := s.in.ReadRune()
	return r, err
}

func (s *shell) run() {
	var err error
	w := newWriter("> ")
	for err == nil {
		var r rune
		r, err = s.next()
		k := Key(r)
		switch {
		case k.IsCtrlChar():
			switch k {
			case Esc:
				seq, err := s.readCtrlSeq()
				if err != nil {
					fmt.Println("ERROR: ", err)
					os.Exit(1)
				}
				switch {
				case IsUpArrow(seq):
					w.move(up, seq)
				case IsDownArrow(seq):
					w.move(down, seq)
				case IsLeftArrow(seq):
					w.move(left, seq)
				case IsRightArrow(seq):
					w.move(right, seq)
				case IsDelete(seq):
					fmt.Printf(string(seq))
				case IsBkspc(seq):
					w.writeCtrl(seq)
					// case IsHome(seq):
					// case IsEnd(seq):
					// case IsPageUp(seq):
					// case IsPageDown(seq):
				}
			case Enter:
				w.println()
				// evaluate expression on enter?
				// expr := w.getText()
				// evaluate(expr)
			default:
				fmt.Printf("%s", string(r))
			case CtrlC:
				return
			}
		default:
			w.write(string(r))
		}
	}
}

func (s *shell) readCtrlSeq() ([]rune, error) {
	runes := make([]rune, 0, 2)
	runes = append(runes, rune(Esc))

	for {
		if s.in.Buffered() == 0 {
			break
		}
		r, _, err := s.in.ReadRune()
		if err != nil {
			return []rune{}, err
		}
		runes = append(runes, r)
	}
	return runes, nil
}

func tprint(a ...interface{}) {
	s := fmt.Sprintln(a...)
	fmt.Print(s)
	fmt.Printf("\x1b[%dD", len(s))
}

func clear() {
	fmt.Print("\x1b[2K")
}

func cb(n int) {
	fmt.Printf("\x1b[%dD", n)
}

const (
	CL = "[2k"
)
