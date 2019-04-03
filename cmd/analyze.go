package cmd

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
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

var evalCmd = &cobra.Command{
	Use:   "eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...]",
	Short: "Evaluates a list of measures and dimensions",
	Long:  `Evaluates a list of measures and dimensions. To evaluate a measure for a specific dimension use the <measure> by <dimension> notation. If dimensions are omitted then the eval will be evaluated over all dimensions.`,
	Example: `corectl eval "Count(a)" // returns the number of values in field "a"
corectl eval "1+1" // returns the calculated value for 1+1
corectl eval "Avg(Sales)" by "Region" // returns the average of measure "Sales" for dimension "Region"
corectl eval by "Region" // Returns the values for dimension "Region"`,

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Expected at least one dimension or measure")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		internal.Eval(rootCtx, state.Doc, args)
	},
}

var catwalkCmd = withCommonLocalFlags(&cobra.Command{
	Use:   "catwalk",
	Short: "Opens the specified app in catwalk",
	Long:  `Opens the specified app in catwalk. If no app is specified the catwalk hub will be opened.`,
	Example: `corectl catwalk --app my-app.qvf
corectl catwalk --app my-app.qvf --catwalk-url http://localhost:8080`,

	Run: func(ccmd *cobra.Command, args []string) {
		catwalkURL := viper.GetString("catwalk-url") + "?engine_url=" + internal.TidyUpEngineURL(viper.GetString("engine")) + "/apps/" + viper.GetString("app")
		if !strings.HasPrefix(catwalkURL, "www") && !strings.HasPrefix(catwalkURL, "https://") && !strings.HasPrefix(catwalkURL, "http://") {
			fmt.Println("Please provide a valid URL starting with 'https://', 'http://' or 'www'")
			os.Exit(1)
		}
		err := browser.OpenURL(catwalkURL)
		if err != nil {
			fmt.Println("Could not open URL", err)
			os.Exit(1)
		}
	},
}, "catwalk-url")
