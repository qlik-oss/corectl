package printer

import (
	"github.com/spf13/viper"
)

// Since the printing mode is mostly relevant for the printer package
// this seems like a good place to put it. But, internal.log uses some
// of these flags as well so maybe it should be in state.

type pmode int

func (m pmode) String() string {
	switch m {
	case jsonMode:
		return "json"
	case bashMode:
		return "bash"
	case quietMode:
		return "quiet"
	default:
		return "standard"
	}
}

const (
	_ pmode = iota
	quietMode
	bashMode
	jsonMode
)

var mode pmode

// SetMode sets the printing mode based on the flags provided to viper.
// The flags have the following precedence: json > bash > quiet
func SetMode() {
	switch {
	case viper.GetBool("json"):
		mode = jsonMode
	case viper.GetBool("bash"):
		mode = bashMode
	case viper.GetBool("quiet"):
		mode = quietMode
	}
	// Example implementation if we would like to log precedence information
	/*
		for m := jsonMode; m > 0; m-- {
			if mode > m {
				fmt.Printf("flag '--%s' overriden by '--%s'\n", m.String(), mode.String())
			}
		}
	*/
}
