package printer

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal"
	enigma "github.com/qlik-oss/enigma-go"
)

// PrintConnections prints a list of connections to standard out
func PrintConnections(connections []*enigma.Connection, printAsJSON bool, printAsBash bool) {
	if internal.PrintJSON {
		internal.PrintAsJSON(connections)
	} else if printAsBash {
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
	if internal.PrintJSON {
		internal.PrintAsJSON(connection)
	}
}
