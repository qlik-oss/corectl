package printer

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/qlik-oss/corectl/internal"
)

// PrintTables prints a list of tables along with meta data to system out.
func PrintTables(data *internal.ModelMetadata) {
	tableList2 := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(tableList2, "Name\tRow count\tRAM\tFields\n")
	for _, table := range data.Tables {

		fmt.Fprintf(tableList2, "%s\t%d\t%s\t%s\n", table.Name, table.NoOfRows, table.MemUsage(), data.FieldsInTableTexts[table.Name])
	}
	fmt.Fprintf(tableList2, "\t\t\n")
	fmt.Fprintf(tableList2, "Total RAM \t\t%s\n", data.MemUsage())
	tm.Print(tableList2)
	tm.Flush()
}
