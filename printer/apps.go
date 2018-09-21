package printer

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/qlik-oss/enigma-go"
)

// PrintApps prints a list of apps and some meta to system out.
func PrintApps(docList []*enigma.DocListEntry) {

	docTable := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(docTable, "Name\tLast Modified\tSize\tLast Reloaded\tReadOnly\tTitle\n")
	for _, doc := range docList {
		fmt.Fprintf(docTable, "%s\t%f\t%f\t%s\t%t\t%s\n", doc.DocName, doc.FileTime, doc.FileSize, doc.LastReloadTime, doc.ReadOnly, doc.Title)
	}
	fmt.Print(docTable)
}
