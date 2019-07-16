package printer

import (
	"os"
	"sort"

	"github.com/qlik-oss/corectl/internal"
	"github.com/olekukonko/tablewriter"
)

// PrintContexts prints a list of contexts to standard out
func PrintContexts(handler *internal.ContextHandler, printAsBash bool) {
	var sortedContextKeys []string
	for k := range handler.Contexts {
		sortedContextKeys = append(sortedContextKeys, k)
	}

	sort.Strings(sortedContextKeys)

	if printAsBash {
		for _, v := range sortedContextKeys {
			PrintToBashComp(v)
		}
	} else {
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetRowLine(true)
		writer.SetHeader([]string{"Name", "Product", "Current", "Comment"})

		for _, k := range sortedContextKeys {
			context := handler.Get(k)
			row := []string{k, context.Product, "", context.Comment}
			if k == handler.Current {
				row[2] = "*"
			}
			writer.Append(row)
		}
		writer.Render()
	}
}
