package printer

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/qlik-oss/core-corectl/internal"
)

func PrintTables(data *internal.ModelMetadata) {

	tm.Println("---------- Tables ----------")

	tableList2 := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(tableList2, "Name\tRow count\tRAM\tFields\n")
	for _, table := range data.Metadata.Tables {
		if !table.Is_system {
			fmt.Fprintf(tableList2, "%s\t%d\t%s\t%s\n", table.Name, table.No_of_rows, formatBytes(table.Byte_size), data.FieldsInTable[table.Name])
		}
	}
	fmt.Fprintf(tableList2, "\t\t\n")
	fmt.Fprintf(tableList2, "Total RAM \t\t%s\n", formatBytes(data.Metadata.StaticByteSize))
	tm.Print(tableList2)
	tm.Flush()
}
