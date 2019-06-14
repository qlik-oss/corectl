package printer

import (
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

// PrintContexts prints a list of contexts to standard out
func PrintContexts(contexts map[interface{}]interface{}, currentContext string, printAsBash bool) {
	var sortedContextKeys []string
	for k := range contexts {
		sortedContextKeys = append(sortedContextKeys, k.(string))
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

		for _, v := range sortedContextKeys {
			context := map[interface{}]interface{}{}
			context = contexts[v].(map[interface{}]interface{})
			if v == currentContext {
				writer.Append([]string{v, context["product"].(string), "*", context["comment"].(string)})
			} else {
				writer.Append([]string{v, context["product"].(string), "", context["comment"].(string)})
			}
		}
		writer.Render()
	}
}
