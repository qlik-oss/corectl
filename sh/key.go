package main

type Key int

const (
	CtrlC Key = 3
	CtrlD     = 4
	Enter     = 13
	Esc       = 27
	Del       = 0x7f
)

func (k Key) IsCtrlChar() bool {
	return k < ' ' || k == 0x7f
}

func runesEq(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

var rightKeySeq = []rune{rune(Esc), rune(91), rune(67)}
var leftKeySeq = []rune{rune(Esc), rune(91), rune(68)}

func IsUpArrow(seq []rune) bool {
	upArrow := []rune{rune(Esc), rune(91), rune(65)}
	return runesEq(seq, upArrow)
}

func IsDownArrow(seq []rune) bool {
	downArrow := []rune{rune(Esc), rune(91), rune(66)}
	return runesEq(seq, downArrow)
}

func IsRightArrow(seq []rune) bool {
	return runesEq(seq, rightKeySeq)
}

func IsLeftArrow(seq []rune) bool {
	return runesEq(seq, leftKeySeq)
}

func IsDelete(seq []rune) bool {
	delete := []rune{rune(Esc), rune(91), rune(51), rune(126)}
	return runesEq(seq, delete)
}

func IsBkspc(seq []rune) bool {
	return false
}
