package printer

import (
	"os"

	"github.com/qlik-oss/corectl/internal"

	"github.com/olekukonko/tablewriter"
)

// PrintStates prints a list of states to system out.
func PrintStates(statesList []string, printAsBash bool) {
	if internal.PrintJSON {
		internal.PrintAsJSON(statesList)
	} else if printAsBash {
		for _, state := range statesList {
			PrintToBashComp(state)
		}
	} else {
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Name"})

		for _, state := range statesList {
			writer.Append([]string{state})
		}
		writer.Render()
	}
}
