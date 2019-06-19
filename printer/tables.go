package printer

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal"
)

// PrintTables prints a list of tables along with meta data to system out.
func PrintTables(data *internal.ModelMetadata) {
	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetHeader([]string{"Name", "Row count", "RAM", "Fields"})
	writer.SetAutoFormatHeaders(false)
	writer.SetRowLine(true)
	for _, table := range data.Tables {
		writer.Append([]string{table.Name, strconv.Itoa(table.NoOfRows), table.MemUsage(), data.FieldsInTableTexts[table.Name]})
	}
	writer.SetFooter([]string{" ", "Total RAM", data.MemUsage(), " "})
	writer.Render()
}
