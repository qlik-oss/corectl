package printer

import (
	"context"
	"fmt"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/viper"
)

//PrintStatus prints the name of the app and the engine corectl is connected to.
// It also prints if the data model is empty or not
func PrintStatus(state *internal.State) {
	if state.AppID != "" {
		fmt.Println("Connected to " + state.AppID + " @ " + viper.GetString("engine"))
	} else {
		fmt.Println("Connected to session app @ " + viper.GetString("engine"))
	}
	tableCount := dataModelTableCount(state.Ctx, state.Doc)
	if tableCount == 0 {
		fmt.Println("The data model is empty.")
	} else if tableCount == 1 {
		fmt.Println("The data model has 1 table.")
	} else {
		fmt.Printf("The data model has %d tables.", tableCount)
	}

}

// Returns the number of tables in the data model
func dataModelTableCount(ctx context.Context, doc *enigma.Doc) int {
	tables, _, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		return 0
	}
	return len(tables)
}
