package printer

import (
	"encoding/json"
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/qlik-oss/corectl/internal"
	enigma "github.com/qlik-oss/enigma-go"
)

// PrintConnections prints a list of connections to standard out
func PrintConnections(connections []*enigma.Connection, printAsJSON bool, printAsBash bool) {
	if printAsJSON {
		jsonPrinter(connections)
	} else if printAsBash {
		for _, connection := range connections {
			PrintToBashComp(connection.Id)
		}
	} else {
		connectionsTable := tm.NewTable(0, 10, 3, ' ', 0)
		fmt.Fprintf(connectionsTable, "Id\tName\n")
		for _, connection := range connections {
			fmt.Fprintf(connectionsTable, "%s\t%s\n", connection.Id, connection.Name)
		}
		fmt.Print(connectionsTable)
	}
}

// PrintConnection prints a connection to standard out
func PrintConnection(connection *enigma.Connection) {
	jsonPrinter(connection)
}

func jsonPrinter(v interface{}) {
	buffer, err := json.Marshal(v)
	if err != nil {
		internal.FatalError(err)
	}
	fmt.Println(prettyJSON(buffer))
}
