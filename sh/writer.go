package main

import (
	"container/list"
	"fmt"
	"log"
)

type Writer struct {
	lines      []*list.List
	cursorLine int
	cursorCol  *list.Element
	prefix     string
}

type Direction int

const eol string = "EOL"

const (
	up Direction = iota
	down
	left
	right
)

func newLine(prefix string) *list.List {
	fmt.Print(prefix)
	l := list.New()
	l.PushBack(eol)
	return l
}

func isEOL(e *list.Element) bool {
	if e == nil {
		return false
	}
	return fmt.Sprintf("%v", e.Value) == eol
}

func newWriter(prefix string) *Writer {
	l := newLine(prefix)
	w := &Writer{
		lines:      make([]*list.List, 7),
		cursorLine: 0,
		cursorCol:  l.Back(),
		prefix:     prefix,
	}
	w.lines[0] = l
	return w
}

func (w *Writer) move(d Direction, seq []rune) {
	switch d {
	case up:
	case down:
	case left:
		log.Println("left")
		prev := w.cursorCol.Prev()
		if prev != nil {
			w.cursorCol = prev
			fmt.Print(string(seq))
		}
	case right:
		log.Println("right")
		next := w.cursorCol.Next()
		if next != nil {
			w.cursorCol = next
			fmt.Print(string(seq))
		}
	}
}

func (w *Writer) logLine() {
	for e, i := w.lines[0].Front(), 0; e != nil; e = e.Next() {
		log.Printf("%d: %v", i, e.Value)
		i++
	}
	log.Println("<--- end of line")
	if w.cursorCol != nil {
		log.Printf("%v<--- current pos\n\r", w.cursorCol.Value)
	} else {
		log.Printf("nil <--- current pos\n\r")
	}
}

func (w *Writer) writeCtrl(seq []rune) {
}

func (w *Writer) println() {
	fmt.Print("\r\n")
	w.cursorLine++
	w.lines[w.cursorLine] = newLine(w.prefix)
	w.cursorCol = w.lines[w.cursorLine].Back()
}

func (w *Writer) write(s string) {
	var curr *list.Element
	curr = w.lines[w.cursorLine].InsertBefore(s, w.cursorCol)
	w.print()
	w.cursorCol = curr.Next()
	w.logLine()
}

func (w *Writer) print() {
	i := 0
	for s := w.cursorCol.Prev(); !isEOL(s); s = s.Next() {
		fmt.Print(s.Value)
		i++
	}
	for k := 0; k < (i - 1); k++ {
		fmt.Print(string(leftKeySeq))
	}
}

func (w *Writer) getText() string {
	s := ""
	for _, v := range w.lines {
		if v != nil && v.Len() > 0 {
			for e := v.Front(); !isEOL(e); e = e.Next() {
				s += fmt.Sprintf("%v", e.Value)
			}
			s += "\r\n"
		}
	}
	return s
}
