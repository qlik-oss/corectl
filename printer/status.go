package printer

import (
	"context"
	"fmt"

	"github.com/qlik-oss/enigma-go"
)

//PrintStatus prints the name of the app and the engine corectl is connected to.
// It also prints if the data model is empty or not
func PrintStatus(ctx context.Context, doc *enigma.Doc, engine string, appName string) {
	if doc != nil {
		fmt.Println("Connected to " + appName + " @ " + engine)
		tableCount := dataModelTableCount(ctx, doc)
		if tableCount == 0 {
			fmt.Println("The data model is empty.")
		} else if tableCount == 1 {
			fmt.Printf("The data model has %d table.\n", tableCount)
		} else {
			fmt.Printf("The data model has %d tables.\n", tableCount)
		}
	} else {
		fmt.Println("Connected without app to " + engine)
	}
}
func PrintStatusRest(url string) {
	fmt.Println("Connected without app to ", url)
}

// Returns the number of tables in the data model
func dataModelTableCount(ctx context.Context, doc *enigma.Doc) int {
	tables, _, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		return 0
	}
	return len(tables)
}
