package printer

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// PrintConnections prints a list of connections to standard out
func PrintConnections(connections []*enigma.Connection, mode PrintMode) {
	if mode.JsonMode() {
		log.PrintAsJSON(connections)
	} else if mode.QuietMode() || mode.BashMode() {
		for _, connection := range connections {
			PrintToBashComp(connection.Id)
		}
	} else {
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Id", "Name", "Type"})

		for _, connection := range connections {
			writer.Append([]string{connection.Id, connection.Name, connection.Type})
		}
		writer.Render()
	}
}

// PrintConnection prints a connection to standard out
func PrintConnection(connection *enigma.Connection) {
	log.PrintAsJSON(connection)
}
