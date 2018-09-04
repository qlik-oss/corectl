package printer

import (
	"fmt"
	"math"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/qlik-oss/core-corectl/internal"
)

func PrintFields(data *internal.ModelMetadata, keyOnly bool) {
	fieldList := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(fieldList, "Field\tRows\tRAM\tTags\t")
	for _, table := range data.TableNames {
		fmt.Fprintf(fieldList, "%s\t", table)
	}
	if data.SampleContentByFieldName != nil {
		fmt.Fprintf(fieldList, "Sample content")
	}
	fmt.Fprintf(fieldList, "\n")
	for _, fieldName := range data.FieldNames {
		field := data.FieldMetadataByName[fieldName]
		if field != nil && !field.Is_system {
			total := uniqueAndTotal(field)
			fmt.Fprintf(fieldList, "%s\t%s\t%s\t%s\t", field.Name, total, formatBytes(field.Byte_size), strings.Join(field.Tags, ", "))
			fieldInfo := data.FieldSourceTableInfoByName[field.Name]
			for _, ff := range fieldInfo {
				fmt.Fprintf(fieldList, "%s\t", ff.RowCount)
			}
			if data.SampleContentByFieldName != nil {
				fmt.Fprintf(fieldList, "%s\t", data.SampleContentByFieldName[fieldName])
			}
			fmt.Fprintf(fieldList, "\n")
		}
	}

	fmt.Fprintf(fieldList, "\t\t\t")
	for range data.TableNames {
		fmt.Fprintf(fieldList, "\t")
	}
	fmt.Fprintf(fieldList, "\n")
	fmt.Fprintf(fieldList, "Total RAM \t\t%s\t", formatBytes(data.Metadata.StaticByteSize))
	for _, tableName := range data.TableNames {
		table := data.TableMetadataByName[tableName]
		if table != nil {
			fmt.Fprintf(fieldList, "\t%s", formatBytes(table.Byte_size))
		}
	}
	fmt.Print(fieldList, "\n\n")
}

func uniqueAndTotal(field *internal.FieldMetadata) string {
	total := ""
	if field.Cardinal < field.Total_count {
		total = fmt.Sprintf("%d/%d", field.Cardinal, field.Total_count)
	} else {
		total = fmt.Sprintf("%d", field.Total_count)
	}
	return total
}

func formatBytes(bytes int) string {
	byteFloat := float64(bytes)
	unit := float64(1024)
	if byteFloat < unit {
		return fmt.Sprintf("%d", bytes)
	}
	exponent := (int)(math.Log(byteFloat) / math.Log(unit))
	prefix := string("kMGTPE"[exponent-1])
	return fmt.Sprintf("%.1f%s", byteFloat/math.Pow(unit, float64(exponent)), prefix)
}
