package printer

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal"
)

// PrintAssociations prints a list of associations to system out.
func PrintAssociations(data *internal.ModelMetadata) {

	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader([]string{"Field(s)", "Linked tables"})

	for _, key := range data.SourceKeys {
		fieldInfo := ""
		for f, field := range key.KeyFields {
			if f > 0 {
				fieldInfo += " + "
			}
			fieldInfo += field
		}
		tableInfo := ""
		for f, table := range key.Tables {
			if f > 0 {
				tableInfo += " <--> "
			}
			tableInfo += table
		}
		writer.Append([]string{fieldInfo, tableInfo})
	}
	writer.Render()
}
