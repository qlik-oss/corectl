package cmd

import (
	"fmt"
	"os"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Get commands
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Lists one or several resources",
	Long:  "Lists one or several resources",
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
	},
}

var getAppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Prints a list of all apps available in the current engine",
	Long:  "Prints a list of all apps available in the current engine",
	Example: `corectl get apps
corectl get apps --engine=localhost:9276`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
		viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineStateWithoutApp(rootCtx, headers)
		docList, err := state.Global.GetDocList(rootCtx)
		if err != nil {
			internal.FatalError(err)
		}
		printer.PrintApps(docList, viper.GetBool("json"))
	},
}

var getAssociationsCmd = &cobra.Command{
	Use:     "assoc",
	Aliases: []string{"associations"},
	Short:   "Print table associations summary",
	Long:    "Print table associations summary",
	Example: `corectl get assoc
corectl get associations`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
		printer.PrintAssociations(data)
	},
}

var getConnectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Prints a list of all connections in the specified app",
	Long:  "Prints a list of all connections in the specified app",
	Example: `corectl get connections
corectl get connections --json`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
		viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, true)
		connections, err := state.Doc.GetConnections(rootCtx)
		if err != nil {
			internal.FatalError(err)
		}
		printer.PrintConnections(connections, viper.GetBool("json"))
	},
}

var getConnectionCmd = &cobra.Command{
	Use:     "connection",
	Short:   "Shows the properties for a specific connection",
	Long:    "Shows the properties for a specific connection",
	Example: "corectl get connection CONNECTION-ID",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Expected a connection name as parameter")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, true)
		connection, err := state.Doc.GetConnection(rootCtx, args[0])
		if err != nil {
			internal.FatalError(err)
		}
		printer.PrintConnection(connection)
	},
}

var getDimensionsCmd = &cobra.Command{
	Use:     "dimensions",
	Short:   "Prints a list of all generic dimensions in the current app",
	Long:    "Prints a list of all generic dimensions in the current app",
	Example: "corectl get dimensions",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
		viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntities(ccmd, args, "dimension", viper.GetBool("json"))
	},
}

var getDimensionCmd = &cobra.Command{
	Use:     "dimension <dimension-id>",
	Short:   "Shows content of an generic dimension",
	Long:    "Shows content of an generic dimension. If no subcommand is specified the properties will be shown.",
	Example: "corectl get dimension DIMENSION-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	//if no specific subcommand is used show the Dimensions properties
	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "dimension")
	},
}

var getDimensionPropertiesCmd = &cobra.Command{
	Use:     "properties <dimension-id>",
	Short:   "Prints the properties of the generic dimension",
	Long:    "Prints the properties of the generic dimension",
	Example: "corectl get dimension properties DIMENSION-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getDimensionCmd.PersistentPreRun(getDimensionCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "dimension")
	},
}

var getDimensionLayoutCmd = &cobra.Command{
	Use:     "layout <dimension-id>",
	Short:   "Evaluates the layout of an generic dimension",
	Long:    "Evaluates the layout of an generic dimension",
	Example: "corectl get dimension layout DIMENSION-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getDimensionCmd.PersistentPreRun(getDimensionCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityLayout(ccmd, args, "dimension")
	},
}

var getFieldCmd = &cobra.Command{
	Use:     "field <field name>",
	Short:   "Shows content of a field",
	Long:    "Prints all the values for a specific field in your data model",
	Example: "corectl get field FIELD",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

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

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

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

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, true)
		printer.PrintFields(data, true)
	},
}

var getMeasuresCmd = &cobra.Command{
	Use:   "measures",
	Short: "Prints a list of all generic measures in the current app",
	Long:  "Prints a list of all generic measures in the current app in either plain text or json",
	Example: `corectl get measures
corectl get measures --json`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
		viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntities(ccmd, args, "measure", viper.GetBool("json"))
	},
}

var getMeasureCmd = &cobra.Command{
	Use:     "measure <measure-id>",
	Short:   "Shows content of an generic measure",
	Long:    "Shows content of an generic measure. If no subcommand is specified the properties will be shown.",
	Example: "corectl get measure MEASURE-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	//if no specific subcommand is used show the Measure properties
	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "measure")
	},
}

var getMeasurePropertiesCmd = &cobra.Command{
	Use:   "properties <measure-id>",
	Short: "Prints the properties of the generic measure",
	Long:  "Prints the properties of the generic measure in JSON format",
	Example: `corectl get measure properties MEASURE-ID
corectl get measure properties MEASURE-ID --app my-app.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getMeasureCmd.PersistentPreRun(getMeasureCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "measure")
	},
}

var getMeasureLayoutCmd = &cobra.Command{
	Use:   "layout <measure-id>",
	Short: "Evaluates the layout of an generic measure",
	Long:  "Evaluates the layout of an generic measure and prints in JSON format",
	Example: `corectl get measure layout MEASURE-ID
corectl get measure layout MEASURE-ID --app my-app.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getMeasureCmd.PersistentPreRun(getMeasureCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityLayout(ccmd, args, "measure")
	},
}

var getMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Shows metadata about the app",
	Long:  "Lists tables, fields, associations along with metadata like memory consumption, field cardinality etc",
	Example: `corectl get meta
corectl get meta --app my-app.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
		printer.PrintMetadata(data)
	},
}

var getObjectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "Prints a list of all generic objects in the current app",
	Long:  "Prints a list of all generic objects in the current app in either plain text or JSON format",
	Example: `corectl get objects
corectl get objects --json --app=myapp.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
		viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntities(ccmd, args, "object", viper.GetBool("json"))
	},
}

var getObjectCmd = &cobra.Command{
	Use:     "object <object-id>",
	Short:   "Shows content of an generic object",
	Long:    "Shows content of an generic object. If no subcommand is specified the properties will be shown.",
	Example: "corectl get object OBJECT-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	//if no specific subcommand is used show the objects properties
	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "object")
	},
}

var getObjectPropertiesCmd = &cobra.Command{
	Use:     "properties <object-id>",
	Short:   "Prints the properties of the generic object",
	Long:    "Prints the properties of the generic object in JSON format",
	Example: "corectl get object properties OBJECT-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getObjectCmd.PersistentPreRun(getObjectCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "object")
	},
}

var getObjectLayoutCmd = &cobra.Command{
	Use:     "layout <object-id>",
	Short:   "Evaluates the hypercube layout of an generic object",
	Long:    "Evaluates the hypercube layout of an generic object in JSON format",
	Example: "corectl get object layout OBJECT-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getObjectCmd.PersistentPreRun(getObjectCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected an object id specify what object to use as a parameter")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		printer.PrintGenericEntityLayout(state, args[0], "object")
	},
}

var getObjectDataCmd = &cobra.Command{
	Use:     "data <object-id>",
	Short:   "Evaluates the hypercube data of an generic object",
	Long:    "Evaluates the hypercube data of an generic object",
	Example: "corectl get object data OBJECT-ID --app my-app.qvf",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getObjectCmd.PersistentPreRun(getObjectCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected an object id specify what object to use as a parameter")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		printer.EvalObject(rootCtx, state.Doc, args[0])
	},
}

var getScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Print the reload script",
	Long:  "Fetches the script currently set in the app and prints it in plain text.",
	Example: `corectl get script
corectl get script --app=my-app.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		script, err := state.Doc.GetScript(rootCtx)
		if err != nil {
			internal.FatalError(err)
		}
		fmt.Println(script)
	},
}

var getStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints status info about the connection to engine and current app",
	Long:  "Prints status info about the connection to engine and current app, and also the status of the data model",
	Example: `corectl get status
corectl get status --app=my-app.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		printer.PrintStatus(state, viper.GetString("engine"))
	},
}

var getTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Print tables summary",
	Long:  "Prints tables summary for the data model in an app",
	Example: `corectl get tables
corectl get tables --app=my-app.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		getCmd.PersistentPreRun(getCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
		printer.PrintTables(data)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(getAppsCmd)
	getCmd.AddCommand(getAssociationsCmd)
	getCmd.AddCommand(getConnectionsCmd)
	getCmd.AddCommand(getConnectionCmd)
	getCmd.AddCommand(getDimensionCmd)
	getCmd.AddCommand(getDimensionsCmd)
	getCmd.AddCommand(getFieldCmd)
	getCmd.AddCommand(getFieldsCmd)
	getCmd.AddCommand(getKeysCmd)
	getCmd.AddCommand(getMeasureCmd)
	getCmd.AddCommand(getMeasuresCmd)
	getCmd.AddCommand(getMetaCmd)
	getCmd.AddCommand(getObjectCmd)
	getCmd.AddCommand(getObjectsCmd)
	getCmd.AddCommand(getScriptCmd)
	getCmd.AddCommand(getStatusCmd)
	getCmd.AddCommand(getTablesCmd)

	getObjectCmd.AddCommand(getObjectPropertiesCmd)
	getObjectCmd.AddCommand(getObjectLayoutCmd)
	getObjectCmd.AddCommand(getObjectDataCmd)

	getDimensionCmd.AddCommand(getDimensionPropertiesCmd)
	getDimensionCmd.AddCommand(getDimensionLayoutCmd)

	getMeasureCmd.AddCommand(getMeasurePropertiesCmd)
	getMeasureCmd.AddCommand(getMeasureLayoutCmd)
}
