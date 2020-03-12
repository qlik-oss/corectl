package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"unicode"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
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
	in *bufio.Reader
}

func newShell() *shell {
	return &shell{
		cur: 0,
		buf: &bytes.Buffer{},
		in: bufio.NewReader(os.Stdin),
	}
}

func (s *shell) next() (rune, error) {
	r, _, err := s.in.ReadRune()
	return r, err
}

func (s *shell) run() {
	var err error
	for ;err == nil; {
		var r rune
		r, err = s.next()
		switch r {
		case '\x03', '\x04':
			tprint("exit")
			return
		case '\r', '\n':
			b, _ := ioutil.ReadAll(s.buf)
			line := string(b)
			clear()
			cb(len(line))
			tprint("\x1b[4m" + line + "\x1b[0m")
		case '\u007f', '\b':
			if s.cur > 0 {
				s.cur--
			}
			fmt.Print(string(r))
			s.buf.WriteRune(r)
		default:
			switch {
			case unicode.IsGraphic(r):
				s.cur++
				fmt.Print(string(r))
				s.buf.WriteRune(r)
			case unicode.IsControl(r):
				s.readCtrl()
			}
		}
	}
}

func (s *shell) readCtrl() {
	runes := make([]rune, 2)
	r, _ := s.next()
	runes[0] = r
	r, _ = s.next()
	runes[1] = r
	switch s := string(runes); s {
	case "[A":
	case "[B":
	case "[C":
		fmt.Print("\x1b" + s)
	case "[D":
		fmt.Print("\x1b" + s)
	default:
		fmt.Printf("%q", s)
	}
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
