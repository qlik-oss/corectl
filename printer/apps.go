package printer

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/qlik-oss/enigma-go"
)

// PrintApps prints a list of apps and some meta to system out.
func PrintApps(docList []*enigma.DocListEntry) {

	docTable := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(docTable, "Id\tName\tLast-Reloaded\tReadOnly\tTitle\n")
	for _, doc := range docList {
		fmt.Fprintf(docTable, "%s\t%s\t%s\t%t\t%s\n", doc.DocId, doc.DocName, doc.LastReloadTime, doc.ReadOnly, doc.Title)
	}
	fmt.Print(docTable)
}
