package printer

import (
	"os"

	"github.com/qlik-oss/corectl/internal/log"

	"github.com/olekukonko/tablewriter"
)

// PrintStates prints a list of states to system out.
func PrintStates(statesList []string, printAsBash bool) {
	switch mode {
	case jsonMode:
		log.PrintAsJSON(statesList)
	case bashMode:
		fallthrough
	case quietMode:
		for _, state := range statesList {
			PrintToBashComp(state)
		}
	default:
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Name"})

		for _, state := range statesList {
			writer.Append([]string{state})
		}
		writer.Render()
	}
}
