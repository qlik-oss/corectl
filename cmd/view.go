package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"os"
)

var getAssociationsCmd = &cobra.Command{
	Use:     "assoc",
	Aliases: []string{"associations"},
	Short:   "Print table associations summary",
	Long:    "Print table associations summary",
	Example: `corectl get assoc
corectl get associations`,

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
		printer.PrintAssociations(data)
	},
}

var getTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Print tables summary",
	Long:  "Prints tables summary for the data model in an app",
	Example: `corectl get tables
corectl get tables --app=my-app.qvf`,

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
		printer.PrintTables(data)
	},
}

var getMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Shows metadata about the app",
	Long:  "Lists tables, fields, associations along with metadata like memory consumption, field cardinality etc",
	Example: `corectl meta
corectl get meta --app my-app.qvf`,

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
		printer.PrintMetadata(data)
	},
}

var getFieldCmd = &cobra.Command{
	Use:     "field <field name>",
	Short:   "Shows content of a field",
	Long:    "Prints all the values for a specific field in your data model",
	Example: "corectl get field FIELD",

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Expected a field name as parameter")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		internal.PrintField(rootCtx, state.Doc, args[0])
	},
}

var getFieldsCmd = &cobra.Command{
	Use:     "fields",
	Short:   "Print field list",
	Long:    "Prints all the fields in an app, and for each field also some sample content, tags and and number of values",
	Example: "corectl get fields",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
		printer.PrintFields(data, false)
	},
}

var getKeysCmd = &cobra.Command{
	Use:     "keys",
	Short:   "Print key-only field list",
	Long:    "Prints a fields list containing key-only fields",
	Example: "corectl get keys",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, true)
		printer.PrintFields(data, true)
	},
}
