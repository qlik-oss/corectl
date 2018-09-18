package printer

import (
	"fmt"
	"github.com/qlik-oss/enigma-go"
	"strings"

	tm "github.com/buger/goterm"
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
	fieldList := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(fieldList, "Field\tUniq/Tot\tRAM\tTags\t")
	for _, table := range data.Tables {
		fmt.Fprintf(fieldList, "%s\t", table.Name)
	}
	if data.SampleContentByFieldName != nil {
		fmt.Fprintf(fieldList, "Sample content")
	}
	fmt.Fprintf(fieldList, "\n")
	for _, field := range data.Fields {
		if field != nil && !field.IsSystem {
			total := uniqueAndTotal(field)
			fmt.Fprintf(fieldList, "%s\t%s\t%s\t%s\t", field.Name, total, field.MemUsage(), strings.Join(field.Tags, ", "))
			//fieldInfo := data.FieldSourceTableInfoByName[field.Name]
			for _, fieldInTable := range field.FieldInTable {
				fmt.Fprintf(fieldList, "%s\t", fieldInTableToText(fieldInTable))
			}
			if data.SampleContentByFieldName != nil {
				fmt.Fprintf(fieldList, "%s\t", data.SampleContentByFieldName[field.Name])
			}
			fmt.Fprintf(fieldList, "\n")
		}
	}

	fmt.Fprintf(fieldList, "\t\t\t")
	for range data.Tables {
		fmt.Fprintf(fieldList, "\t")
	}

	fmt.Fprintf(fieldList, "\n")
	fmt.Fprintf(fieldList, "Total RAM \t\t%s\t", data.MemUsage())
	for _, table := range data.Tables {
		fmt.Fprintf(fieldList, "\t%s", table.MemUsage())
	}

	fmt.Print(fieldList, "\n\n")
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
