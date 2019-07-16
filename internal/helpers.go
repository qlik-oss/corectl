package internal

import (
	"fmt"
	"math"
	"path/filepath"
)

// RelativeToProject transforms a path to be relative to a base path of the project file
func RelativeToProject(path string) string {
	if ConfigDir != "" && !filepath.IsAbs(path) {
		fullpath := filepath.Join(ConfigDir, path)
		return fullpath
	}
	return path
}

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

// Contains will return true if slice contains string
func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
