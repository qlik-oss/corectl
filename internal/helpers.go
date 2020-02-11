package internal

import (
	"fmt"
	"math"
)

// FormatBytes takes a byte size integer and returns a string formatted with kilo, mega, giga prefixes.
func FormatBytes(bytes int) string {
	byteFloat := float64(bytes)
	unit := float64(1024)
	if byteFloat < unit {
		return fmt.Sprintf("%d", bytes)
	}
	exponent := (int)(math.Log(byteFloat) / math.Log(unit))
	prefix := string("kMGTPE"[exponent-1])
	return fmt.Sprintf("%.1f%s", byteFloat/math.Pow(unit, float64(exponent)), prefix)
}
