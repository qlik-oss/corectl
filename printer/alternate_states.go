package printer

import (
	"os"

	"github.com/qlik-oss/corectl/internal/log"

	"github.com/olekukonko/tablewriter"
)

// PrintStates prints a list of states to system out.
func PrintStates(statesList []string, mode log.PrintMode) {

	if mode.JsonMode() {
		log.PrintAsJSON(statesList)
	} else if mode.BashMode() || mode.QuietMode() {
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
