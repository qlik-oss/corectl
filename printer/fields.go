package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/enigma-go"

	"github.com/qlik-oss/corectl/internal"
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

// PrintFields prints a table sof fields along with various metadata to system out.
func PrintFields(data *internal.ModelMetadata, keyOnly bool) {
	writer := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Field", "Uniq/Tot", "RAM", "Tags"}
	// For each table we should also add a header
	for _, table := range data.Tables {
		headers = append(headers, table.Name)
	}

	// Add a sample content header if samples exists
	if data.SampleContentByFieldName != nil {
		headers = append(headers, "Sample content")
	}
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader(headers)

	// Add rows
	writer.SetRowLine(true)
	for _, field := range data.Fields {
		if field != nil && !field.IsSystem && (!keyOnly || isKey(field)) {
			total := uniqueAndTotal(field)
			rows := []string{field.Name, total, field.MemUsage(), strings.Join(field.Tags, ", ")}
			for _, fieldInTable := range field.FieldInTable {
				rows = append(rows, fieldInTableToText(fieldInTable))
			}
			if data.SampleContentByFieldName != nil {
				rows = append(rows, data.SampleContentByFieldName[field.Name])
			}
			writer.Append(rows)
		}
	}

	// Add total ram as footer
	footers := []string{"Total RAM", "", data.MemUsage(), ""}
	for _, table := range data.Tables {
		footers = append(footers, table.MemUsage())
	}
	if data.SampleContentByFieldName != nil {
		footers = append(footers, "")
	}
	writer.SetFooter(footers)

	writer.Render()
}

func isKey(field *internal.FieldModel) bool {
	for _, tag := range field.Tags {
		if tag == "$key" {
			return true
		}
	}
	return false
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
