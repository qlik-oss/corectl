package printer

import (
	"fmt"
	"github.com/fatih/color"
	tm "github.com/buger/goterm"
	"github.com/qlik-oss/corectl/internal"
)

// PrintTables prints a list of tables along with meta data to system out.
func PrintTables(data *internal.ModelMetadata) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	tableList2 := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(tableList2, "%s\tRow count\tRAM\t%s\n", yellow("Name"), green("Fields"))
	for _, table := range data.Tables {

		fmt.Fprintf(tableList2, "%s\t%d\t%s\t%s\n", yellow(table.Name), table.NoOfRows, table.MemUsage(), green(data.FieldsInTableTexts[table.Name]))
	}
	fmt.Fprintf(tableList2, "\t\t\n")
	fmt.Fprintf(tableList2, "Total RAM \t\t%s\n", data.MemUsage())
	fmt.Print(tableList2)
}
