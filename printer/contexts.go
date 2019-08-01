package printer

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal"
)

// PrintContext prints all information in a context
func PrintContext(name string, context *internal.Context) {
	if context == nil {
		fmt.Println("No current context")
		return
	}
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("  Product: %s\n", context.Product)
	fmt.Printf("  Comment: %s\n", context.Comment)
	fmt.Printf("  Engine: %s\n", context.Engine)
	fmt.Printf("  Certificates: %s\n", context.Certificates)
	fmt.Println("  Headers:")
	for k, v := range context.Headers {
		fmt.Printf("    %s: %s\n", k, v)
	}
}

func PrintCurrentContext(name string) {
	if name == "" {
		fmt.Println("Context: <NONE>")
	} else {
		fmt.Printf("Context: %s\n", name)
	}
}

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
		header := []string{"Name", "Product", "Engine", "Current", "Comment"}
		writer.SetHeader(header)

		for _, k := range sortedContextKeys {
			context := handler.Get(k)
			row := []string{k, context.Product, context.Engine, "", context.Comment}
			if k == handler.Current {
				// In case we change header order
				for i, h := range header {
					if h == "Current" {
						row[i] = "*"
						break
					}
				}
			}
			writer.Append(row)
		}
		writer.Render()
	}
}
