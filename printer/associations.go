package printer

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/qlik-oss/core-cli/internal"
)

func PrintAssociations(data *internal.ModelMetadata) {

	tm.Println("---------- Associations ----------")
	keyList := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(keyList, "Field(s)\tLinked tables\n")
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
		fmt.Fprintf(keyList, "%s\t%s\n", fieldInfo, tableInfo)
	}
	tm.Println(keyList)
	tm.Flush()
}
