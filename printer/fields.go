package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/enigma-go"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
)

func uniqueAndTotal(field *internal.FieldModel) string {
	total := ""
	if field.Cardinal < field.TotalCount {
		total = fmt.Sprintf("%d/%d", field.Cardinal, field.TotalCount)
	} else {
		total = fmt.Sprintf("%d", field.TotalCount)
	}
	return total
}

// PrintFields prints a tables of fields along with various metadata to system out.
func PrintFields(data *internal.ModelMetadata, keyOnly bool) {
	if len(data.Fields) == 0 {
		if keyOnly {
			log.Infoln("No key fields found.")
			return
		}
		log.Infoln("No fields found.")
		return
	}

	if mode == bashMode || mode == quietMode {
		for _, field := range data.Fields {
			if field != nil && !field.IsSystem {
				PrintToBashComp(field.Name)
			}
		}
		return
	}

	writer := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Field", "Uniq/Tot", "RAM", "Tags", "Tables"}

	// Add a sample content header if samples exists
	if data.SampleContentByFieldName != nil {
		headers = append(headers, "Sample content")
	}
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader(headers)

	// Add rows
	writer.SetRowLine(true)
	for _, field := range data.Fields {
		if field != nil && !field.IsSystem {
			total := uniqueAndTotal(field)
			cells := []string{field.Name, total, field.MemUsage(), strings.Join(field.Tags, ", ")}
			tablesPresence := ""
			for tableIndex, fieldInTable := range field.FieldInTable {
				if fieldInTable != nil {
					tablesPresence += data.Tables[tableIndex].Name + ": " + fieldInTableToText(fieldInTable) + "\n"
				}
			}
			cells = append(cells, tablesPresence)
			if data.SampleContentByFieldName != nil {
				cells = append(cells, data.SampleContentByFieldName[field.Name])
			}
			writer.Append(cells)
		}
	}

	// Add total ram as footer
	footers := []string{"Total RAM", "-", data.MemUsage(), "-", "-", "-"}
	writer.SetFooter(footers)
	writer.Render()
}

func fieldInTableToText(fieldInTable *enigma.FieldInTableData) string {
	if fieldInTable != nil {

		info := fmt.Sprintf("%d/%d", fieldInTable.NTotalDistinctValues, fieldInTable.NNonNulls)
		if fieldInTable.NRows > fieldInTable.NNonNulls {
			info += fmt.Sprintf("+%d", fieldInTable.NRows-fieldInTable.NNonNulls)
		}
		if fieldInTable.KeyType == "NOT_KEY" {
			info += ""
		} else if fieldInTable.KeyType == "ANY_KEY" {
			info += "*"
		} else if fieldInTable.KeyType == "PRIMARY_KEY" {
			info += "**"
		} else if fieldInTable.KeyType == "PERFECT_KEY" {
			info += "***"
		} else {
			info += "?"
		}
		return info
	}
	return ""
}
