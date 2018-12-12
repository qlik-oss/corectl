package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/browser"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var (
	headersMap         = make(map[string]string)
	explicitConfigFile = ""
	version            = ""
	headers            http.Header
	rootCtx            = context.Background()

	corectlCommand = &cobra.Command{
		Hidden:            true,
		Use:               "corectl",
		Short:             "",
		Long:              `Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
		DisableAutoGenTag: true,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			// if help, version or generate-docs command, no prerun is needed.
			if strings.Contains(ccmd.Use, "help") || ccmd.Use == "generate-docs" || ccmd.Use == "version" {
				return
			}
			internal.QliVerbose = viper.GetBool("verbose")
			if explicitConfigFile != "" {
				viper.SetConfigFile(strings.TrimSpace(explicitConfigFile))
				if err := viper.ReadInConfig(); err == nil {
					internal.LogVerbose("Using config file: " + explicitConfigFile)
				} else {
					fmt.Println(err)
				}
			} else {
				viper.SetConfigName("corectl") // name of config file (without extension)
				viper.SetConfigType("yml")
				viper.AddConfigPath(".")
				if err := viper.ReadInConfig(); err == nil {
					internal.LogVerbose("Using config file in working directory")
				} else {
					internal.LogVerbose("No config file")
				}
			}
			internal.QliVerbose = viper.GetBool("verbose")

			if len(headersMap) == 0 {
				headersMap = viper.GetStringMapString("headers")
			}
			headers = make(http.Header, 1)
			for key, value := range headersMap {
				headers.Set(key, value)
			}
		},

		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Reloads and saves the app after updating connections, dimensions, measures, objects and the script",
		Long: `Builds the app. Example: corectl build --connections ./myconnections.yml --script ./myscript.qvs
			
`,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			corectlCommand.PersistentPreRun(corectlCommand, args)
			viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
			viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("silent", ccmd.PersistentFlags().Lookup("silent"))
		},
		Run: func(ccmd *cobra.Command, args []string) {
			build(ccmd, args)
		},
	}

	catwalkCmd = &cobra.Command{
		Use:   "catwalk",
		Short: "Opens the specified app in catwalk",
		Long:  `Opens the specified app in catwalk. If no app is specified the catwalk hub will be opened.`,
		Example: `corectl catwalk --app my-app.qvf
corectl catwalk --app my-app.qvf --catwalk-url http://localhost:8080`,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			corectlCommand.PersistentPreRun(corectlCommand, args)
			viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("catwalk-url", ccmd.PersistentFlags().Lookup("catwalk-url"))
		},
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
	}

	evalCmd = &cobra.Command{
		Use:   "eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...]",
		Short: "Evaluates a list of measures and dimensions",
		Long:  `Evaluates a list of measures and dimensions. To evaluate a measure for a specific dimension use the <measure> by <dimension> notation. If dimensions are omitted then the eval will be evaluated over all dimensions.`,
		Example: `corectl eval "Count(a)" // returns the number of values in field "a"
corectl eval "1+1" // returns the calculated value for 1+1
corectl eval "Avg(Sales)" by "Region" // returns the average of measure "Sales" for dimension "Region"
corectl eval by "Region" // Returns the values for dimension "Region"`,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			corectlCommand.PersistentPreRun(corectlCommand, args)
			viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
			viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Expected at least one dimension or measure")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			internal.Eval(rootCtx, state.Doc, args)
		},
	}

	//Get commands

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Lists one or several resources",
		Long:  "Lists one or several resources",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			corectlCommand.PersistentPreRun(corectlCommand, args)
			viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
			viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
		},
	}

	getAppsCmd = &cobra.Command{
		Use:   "apps",
		Short: "Prints a list of all apps available in the current engine",
		Long:  "Prints a list of all apps available in the current engine",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineStateWithoutApp(rootCtx, viper.GetString("engine"), viper.GetString("ttl"), headers)
			docList, err := state.Global.GetDocList(rootCtx)
			if err != nil {
				internal.FatalError(err)
			}
			printer.PrintApps(docList, viper.GetBool("json"))
		},
	}

	getAssociationsCmd = &cobra.Command{
		Use:   "assoc",
		Short: "Print table associations summary",
		Long:  "Print table associations summary",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
			printer.PrintAssociations(data)
		},
	}

	getConnectionsCmd = &cobra.Command{
		Use:     "connections",
		Short:   "Prints a list of all connections in the current app",
		Long:    "Prints a list of all connections in the current app",
		Example: "corectl get connections",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			connections, err := state.Doc.GetConnections(rootCtx)
			if err != nil {
				internal.FatalError(err)
			}
			printer.PrintConnections(connections, viper.GetBool("json"))
		},
	}

	getConnectionCmd = &cobra.Command{
		Use:     "connection",
		Short:   "Shows the content of a specific connector",
		Long:    "Shows the content of a specific connector",
		Example: "corectl get connection CONNECTION-ID",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected a connection name as parameter")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			connection, err := state.Doc.GetConnection(rootCtx, args[0])
			if err != nil {
				internal.FatalError(err)
			}
			printer.PrintConnection(connection)
		},
	}

	getDimensionsCmd = &cobra.Command{
		Use:   "dimensions",
		Short: "Prints a list of all generic dimensions in the current app",
		Long:  "Prints a list of all generic dimensions in the current app",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntities(ccmd, args, "dimension", viper.GetBool("json"))
		},
	}

	getDimensionCmd = &cobra.Command{
		Use:   "dimension <dimension-id>",
		Short: "Shows content of an generic dimension",
		Long:  "Shows content of an generic dimension. If no subcommand is specified the properties will be shown. Example: corectl get dimension DIMENSION-ID --app my-app.qvf",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		//if no specific subcommand is used show the Dimensions properties
		Run: func(ccmd *cobra.Command, args []string) {
			getEntityProperties(ccmd, args, "dimension")
		},
	}

	getDimensionPropertiesCmd = &cobra.Command{
		Use:   "properties <dimension-id>",
		Short: "Prints the properties of the generic dimension",
		Long:  "Prints the properties of the generic dimension. Example: corectl get dimension properties DIMENSION-ID --app my-app.qvf",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getDimensionCmd.PersistentPreRun(getDimensionCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntityProperties(ccmd, args, "dimension")
		},
	}

	getDimensionLayoutCmd = &cobra.Command{
		Use:   "layout <dimension-id>",
		Short: "Evaluates the layout of an generic dimension",
		Long:  `Evaluates the layout of an generic dimension. Example: corectl get dimension layout DIMENSION-ID --app my-app.qvf`,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getDimensionCmd.PersistentPreRun(getDimensionCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntityLayout(ccmd, args, "dimension")
		},
	}

	getFieldCmd = &cobra.Command{
		Use:   "field <field name>",
		Short: "Shows content of a field",
		Long:  ``,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected a field name as parameter")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			internal.PrintField(rootCtx, state.Doc, args[0])
		},
	}

	getFieldsCmd = &cobra.Command{
		Use:   "fields",
		Short: "Print field list",
		Long:  "Print field list",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
			printer.PrintFields(data, false)
		},
	}

	getKeysCmd = &cobra.Command{
		Use:   "keys",
		Short: "Print key-only field list",
		Long:  "Print key-only field list",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, true)
			printer.PrintFields(data, true)
		},
	}

	getMeasuresCmd = &cobra.Command{
		Use:   "measures",
		Short: "Prints a list of all generic measures in the current app",
		Long:  "Prints a list of all generic measures in the current app",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntities(ccmd, args, "measure", viper.GetBool("json"))
		},
	}

	getMeasureCmd = &cobra.Command{
		Use:   "measure <measure-id>",
		Short: "Shows content of an generic measure",
		Long:  "Shows content of an generic measure. If no subcommand is specified the properties will be shown. Example: corectl get measure MEASURE-ID --app my-app.qvf",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		//if no specific subcommand is used show the Measure properties
		Run: func(ccmd *cobra.Command, args []string) {
			getEntityProperties(ccmd, args, "measure")
		},
	}

	getMeasurePropertiesCmd = &cobra.Command{
		Use:   "properties <measure-id>",
		Short: "Prints the properties of the generic measure",
		Long:  "Prints the properties of the generic measure. Example: corectl get measure properties MEASURE-ID --app my-app.qvf",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getMeasureCmd.PersistentPreRun(getMeasureCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntityProperties(ccmd, args, "measure")
		},
	}

	getMeasureLayoutCmd = &cobra.Command{
		Use:   "layout <measure-id>",
		Short: "Evaluates the layout of an generic measure",
		Long:  `Evaluates the layout of an generic measure. Example: corectl get measure layout MEASURE-ID --app my-app.qvf`,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getMeasureCmd.PersistentPreRun(getMeasureCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntityLayout(ccmd, args, "measure")
		},
	}

	getMetaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Shows metadata about the app",
		Long:  "Lists tables, fields, associations along with metadata like memory consumption, field cardinality etc",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
			printer.PrintMetadata(data)
		},
	}

	getObjectsCmd = &cobra.Command{
		Use:   "objects",
		Short: "Prints a list of all generic objects in the current app",
		Long:  "Prints a list of all generic objects in the current app",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("json", ccmd.PersistentFlags().Lookup("json"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntities(ccmd, args, "object", viper.GetBool("json"))
		},
	}

	getObjectCmd = &cobra.Command{
		Use:   "object <object-id>",
		Short: "Shows content of an generic object",
		Long:  "Shows content of an generic object. If no subcommand is specified the properties will be shown. Example: corectl get object OBJECT-ID --app my-app.qvf",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		//if no specific subcommand is used show the objects properties
		Run: func(ccmd *cobra.Command, args []string) {
			getEntityProperties(ccmd, args, "object")
		},
	}

	getObjectPropertiesCmd = &cobra.Command{
		Use:   "properties <object-id>",
		Short: "Prints the properties of the generic object",
		Long:  "Prints the properties of the generic object. Example: corectl get object properties OBJECT-ID --app my-app.qvf",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getObjectCmd.PersistentPreRun(getObjectCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {
			getEntityProperties(ccmd, args, "object")
		},
	}

	getObjectLayoutCmd = &cobra.Command{
		Use:   "layout <object-id>",
		Short: "Evaluates the hypercube layout of an generic object",
		Long:  `Evaluates the hypercube layout of an generic object. Example: corectl get object layout OBJECT-ID --app my-app.qvf`,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getObjectCmd.PersistentPreRun(getObjectCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected an object id specify what object to use as a parameter")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			printer.PrintGenericEntityLayout(state, args[0], "object")
		},
	}

	getObjectDataCmd = &cobra.Command{
		Use:   "data <object-id>",
		Short: "Evaluates the hypercube data of an generic object",
		Long:  `Evaluates the hypercube data of an generic object. Example: corectl get object data OBJECT-ID --app my-app.qvf`,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getObjectCmd.PersistentPreRun(getObjectCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected an object id specify what object to use as a parameter")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			printer.EvalObject(rootCtx, state.Doc, args[0])
		},
	}

	getScriptCmd = &cobra.Command{
		Use:   "script",
		Short: "Print the reload script",
		Long:  "Print the reload script",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			script, err := state.Doc.GetScript(rootCtx)
			if err != nil {
				internal.FatalError(err)
			}
			fmt.Println(script)
		},
	}

	getStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Prints status info about the connection to engine and current app",
		Long:  "Prints status info about the connection to engine and current app",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			engine := viper.GetString("engine")
			if engine == "" {
				engine = "localhost:9076"
			}
			state := internal.PrepareEngineState(rootCtx, engine, viper.GetString("app"), viper.GetString("ttl"), headers, false)
			printer.PrintStatus(state, engine)
		},
	}

	getTablesCmd = &cobra.Command{
		Use:   "tables",
		Short: "Print tables summary",
		Long:  "Prints tables summary",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			getCmd.PersistentPreRun(getCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, headers, false)
			printer.PrintTables(data)
		},
	}

	//Reload command

	reloadCmd = &cobra.Command{
		Use:   "reload",
		Short: "Reloads the app.",
		Long:  "Reloads the app. Example: corectl reload",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			corectlCommand.PersistentPreRun(corectlCommand, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
			viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
			viper.BindPFlag("silent", ccmd.PersistentFlags().Lookup("silent"))
			viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			silent := viper.GetBool("silent")

			internal.Reload(rootCtx, state.Doc, state.Global, silent, true)

			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	//Remove commands

	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove entities (connections, dimensions, measures, objects) in the app or the app itself",
		Long:  "Remove one or mores generic entities (connections, dimensions, measures, objects) in the app",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			corectlCommand.PersistentPreRun(corectlCommand, args)
			viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
			viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
		},
	}

	removeAppCmd = &cobra.Command{
		Use:     "app <app-id>",
		Short:   "removes the specified app.",
		Long:    `removes the specified app.`,
		Example: "corectl remove app APP-ID",
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			removeCmd.PersistentPreRun(removeCmd, args)
		},
		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected an identifier of the app to delete.")
				ccmd.Usage()
				os.Exit(1)
			}
			internal.DeleteApp(rootCtx, viper.GetString("engine"), args[0], viper.GetString("ttl"), headers)
		},
	}

	removeConnectionCmd = &cobra.Command{
		Use:     "connection <connection-id>",
		Short:   "removes the specified connection.",
		Long:    `removes the specified connection.`,
		Example: "corectl remove connection CONNECTION-ID",
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			removeCmd.PersistentPreRun(removeCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
		},
		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected an identifier of the connection to delete.")
				ccmd.Usage()
				os.Exit(1)
			}

			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			err := state.Doc.DeleteConnection(rootCtx, args[0])
			if err != nil {
				internal.FatalError("Failed to remove connection", args[0])
			}
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}

		},
	}

	removeDimensionsCmd = &cobra.Command{
		Use:   "dimensions <dimension-id>...",
		Short: "Removes the specified generic dimensions in the current app",
		Long:  "Removes the specified generic dimensions in the current app. Example: corectl remove dimension ID-1 ID-2",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			removeCmd.PersistentPreRun(removeCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected atleast one entity id specify what entity to remove from the app")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			for _, entity := range args {
				destroyed, err := state.Doc.DestroyDimension(rootCtx, entity)
				if err != nil {
					internal.FatalError("Failed to remove generic dimension ", entity+" with error: "+err.Error())
				} else if !destroyed {
					internal.FatalError("Failed to remove generic dimension ", entity)
				}
			}
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	removeMeasuresCmd = &cobra.Command{
		Use:   "measures <measure-id>...",
		Short: "Removes the specified generic measures in the current app",
		Long:  "Removes the specified generic measures in the current app. Example: corectl remove measures ID-1 ID-2",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			removeCmd.PersistentPreRun(removeCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected atleast one entity id specify what entity to remove from the app")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			for _, entity := range args {
				destroyed, err := state.Doc.DestroyMeasure(rootCtx, entity)
				if err != nil {
					internal.FatalError("Failed to remove generic measure ", entity+" with error: "+err.Error())
				} else if !destroyed {
					internal.FatalError("Failed to remove generic measure ", entity)
				}
			}
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	removeObjectsCmd = &cobra.Command{
		Use:   "objects <object-id>...",
		Short: "Removes the specified generic objects in the current app",
		Long:  "Removes the specified generic objects in the current app. Example: corectl remove objects ID-1 ID-2",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			removeCmd.PersistentPreRun(removeCmd, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Expected atleast one entity id specify what entity to remove from the app")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
			for _, entity := range args {
				destroyed, err := state.Doc.DestroyObject(rootCtx, entity)
				if err != nil {
					internal.FatalError("Failed to remove generic object ", entity+" with error: "+err.Error())
				} else if !destroyed {
					internal.FatalError("Failed to remove generic object ", entity)
				}
			}
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	//Set commands

	setCmd = &cobra.Command{
		Use:   "set",
		Short: "Sets one or several resources",
		Long:  "Sets one or several resources",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			corectlCommand.PersistentPreRun(corectlCommand, args)
			viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
			viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
			viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
			viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))

		},
	}

	setAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Sets the objects, measures, dimensions, connections and script in the current app",
		Long:  "Sets the objects, measures, dimensions, connections and script in the current app",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			setCmd.PersistentPreRun(setCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {

			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			separateConnectionsFile := ccmd.Flag("connections").Value.String()
			if separateConnectionsFile == "" {
				separateConnectionsFile = GetRelativeParameter("connections")
			}
			internal.SetupConnections(rootCtx, state.Doc, separateConnectionsFile, viper.ConfigFileUsed())
			internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("dimensions").Value.String(), "dimension")
			internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("measures").Value.String(), "measure")
			internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String(), "object")
			scriptFile := ccmd.Flag("script").Value.String()
			if scriptFile == "" {
				scriptFile = GetRelativeParameter("script")
			}
			if scriptFile != "" {
				internal.SetScript(rootCtx, state.Doc, scriptFile)
			}

			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	setConnectionsCmd = &cobra.Command{
		Use:   "connections <path-to-connections-file.yml>",
		Short: "Sets or updates the connections in the current app",
		Long:  "Sets or updates the connections in the current app. Example corectl set connections ./my-connections.yml",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			setCmd.PersistentPreRun(setCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {

			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			separateConnectionsFile := ""
			if len(args) > 0 {
				separateConnectionsFile = args[0]
			}
			if separateConnectionsFile == "" {
				separateConnectionsFile = GetRelativeParameter("connections")
			}
			internal.SetupConnections(rootCtx, state.Doc, separateConnectionsFile, viper.ConfigFileUsed())
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	setDimensionsCmd = &cobra.Command{
		Use:   "dimensions <glob-pattern-path-to-dimensions-files.json>",
		Short: "Sets or updates the dimensions in the current app",
		Long:  "Sets or updates the dimensions in the current app. Example corectl set dimensions ./my-dimensions-glob-path.json",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			setCmd.PersistentPreRun(setCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {

			commandLineDimensions := ""
			if len(args) > 0 {
				commandLineDimensions = args[0]
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineDimensions, "dimension")
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	setMeasuresCmd = &cobra.Command{
		Use:   "measures <glob-pattern-path-to-measures-files.json>",
		Short: "Sets or updates the measures in the current app",
		Long:  "Sets or updates the measures in the current app. Example corectl set measures ./my-measures-glob-path.json",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			setCmd.PersistentPreRun(setCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {

			commandLineMeasures := ""
			if len(args) > 0 {
				commandLineMeasures = args[0]
			}
			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineMeasures, "measure")
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	setObjectsCmd = &cobra.Command{
		Use:   "objects <glob-pattern-path-to-objects-files.json",
		Short: "Sets or updates the objects in the current app",
		Long:  "Sets or updates the objects in the current app Example corectl set objects ./my-objects-glob-path.json",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			setCmd.PersistentPreRun(setCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {

			commandLineObjects := ""
			if len(args) > 0 {
				commandLineObjects = args[0]
			}

			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineObjects, "object")
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	setScriptCmd = &cobra.Command{
		Use:   "script <path-to-script-file.yml>",
		Short: "Sets the script in the current app",
		Long:  "Sets the script in the current app. Example: corectl set script ./my-script-file",

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			setCmd.PersistentPreRun(setCmd, args)
		},

		Run: func(ccmd *cobra.Command, args []string) {

			state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
			scriptFile := ""
			if len(args) > 0 {
				scriptFile = args[0]
			}
			if scriptFile == "" {
				scriptFile = GetRelativeParameter("script")
			}
			if scriptFile != "" {
				internal.SetScript(rootCtx, state.Doc, scriptFile)
			} else {
				fmt.Println("Expected the path to a file containing the qlik script")
				ccmd.Usage()
				os.Exit(1)
			}
			if state.AppID != "" && !viper.GetBool("no-save") {
				internal.Save(rootCtx, state.Doc, state.AppID)
			}
		},
	}

	//Other commands

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of corectl",

		Run: func(_ *cobra.Command, args []string) {
			fmt.Printf("corectl version %s\n", version)
		},
	}

	generateDocsCommand = &cobra.Command{
		Use:    "generate-docs",
		Short:  "Generate markdown docs based on cobra commands",
		Long:   "Generate markdown docs based on cobra commands",
		Hidden: true,

		Run: func(ccmd *cobra.Command, args []string) {
			fmt.Println("Generating documentation")
			doc.GenMarkdownTree(corectlCommand, "./docs")
		},
	}
)

func init() {

	corectlCommand.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	viper.BindPFlag("verbose", corectlCommand.PersistentFlags().Lookup("verbose"))

	//Is it nicer to have one loop per argument or group the commands together if they all are used in the same commands?
	for _, command := range []*cobra.Command{buildCmd, catwalkCmd, evalCmd, getCmd, reloadCmd, removeCmd, setCmd} {
		command.PersistentFlags().StringVarP(&explicitConfigFile, "config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")
	}

	//since several commands are using the same flag, the viper binding has to be done in the commands prerun function, otherwise they overwrite.

	for _, command := range []*cobra.Command{buildCmd, catwalkCmd, evalCmd, getCmd, reloadCmd, removeCmd, setCmd} {
		command.PersistentFlags().StringP("engine", "e", "", "URL to engine (default \"localhost:9076\")")
	}

	for _, command := range []*cobra.Command{buildCmd, evalCmd, getCmd, reloadCmd, removeCmd, setCmd} {
		command.PersistentFlags().String("ttl", "30", "Engine session time to live in seconds")
	}

	for _, command := range []*cobra.Command{buildCmd, evalCmd, getCmd, reloadCmd, removeCmd, setCmd} {
		//not binding to viper since binding a map does not seem to work.
		command.PersistentFlags().StringToStringVar(&headersMap, "headers", nil, "Headers to use when connecting to qix engine")
	}

	for _, command := range []*cobra.Command{buildCmd, catwalkCmd, evalCmd, getAssociationsCmd, getConnectionsCmd, getConnectionCmd, getDimensionsCmd, getDimensionCmd, getFieldsCmd, getKeysCmd, getFieldCmd, getMeasuresCmd, getMeasureCmd, getMetaCmd, getObjectsCmd, getObjectCmd, getScriptCmd, getStatusCmd, getTablesCmd, reloadCmd, removeConnectionCmd, removeDimensionsCmd, removeMeasuresCmd, removeObjectsCmd, setCmd} {
		command.PersistentFlags().StringP("app", "a", "", "App name including .qvf file ending. If no app is specified a session app is used instead.")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd, setConnectionsCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("connections", "", "path/to/connections.yml that contains connections that are used in the reload. Note that when specifying connections in the config file they are specified inline, not as a file reference!")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("dimensions", "", "A list of generic dimension json paths")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("measures", "", "A list of generic measures json paths")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("objects", "", "A list of generic object json paths")
	}

	for _, command := range []*cobra.Command{buildCmd, reloadCmd} {
		command.PersistentFlags().Bool("silent", false, "Do not log reload progress")
	}

	for _, command := range []*cobra.Command{reloadCmd, removeConnectionCmd, removeDimensionsCmd, removeMeasuresCmd, removeObjectsCmd, setCmd} {
		command.PersistentFlags().Bool("no-save", false, "Do not save the app")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("script", "", "path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded")
	}

	for _, command := range []*cobra.Command{getAppsCmd, getConnectionsCmd, getDimensionsCmd, getMeasuresCmd, getObjectsCmd} {
		command.PersistentFlags().Bool("json", false, "Prints the information in json format")
	}

	catwalkCmd.PersistentFlags().String("catwalk-url", "https://catwalk.core.qlik.com", "Url to an instance of catwalk, if not provided the qlik one will be used.")

	// commands
	corectlCommand.AddCommand(buildCmd)
	corectlCommand.AddCommand(catwalkCmd)
	corectlCommand.AddCommand(evalCmd)
	corectlCommand.AddCommand(generateDocsCommand)
	corectlCommand.AddCommand(getCmd)
	corectlCommand.AddCommand(reloadCmd)
	corectlCommand.AddCommand(removeCmd)
	corectlCommand.AddCommand(setCmd)
	corectlCommand.AddCommand(versionCmd)

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

	removeCmd.AddCommand(removeAppCmd)
	removeCmd.AddCommand(removeConnectionCmd)
	removeCmd.AddCommand(removeDimensionsCmd)
	removeCmd.AddCommand(removeMeasuresCmd)
	removeCmd.AddCommand(removeObjectsCmd)

	setCmd.AddCommand(setAllCmd)
	setCmd.AddCommand(setConnectionsCmd)
	setCmd.AddCommand(setDimensionsCmd)
	setCmd.AddCommand(setMeasuresCmd)
	setCmd.AddCommand(setObjectsCmd)
	setCmd.AddCommand(setScriptCmd)

}

// GetRelativeParameter returns a parameter from the config file.
// It modifies the parameter to actually be relative to the config file and not the working directory
func GetRelativeParameter(paramName string) string {
	pathInConfigFile := viper.GetString(paramName)
	if pathInConfigFile != "" {
		return internal.RelativeToProject(viper.ConfigFileUsed(), pathInConfigFile)
	}
	return ""
}

func main() {
	if err := corectlCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getEntityProperties(ccmd *cobra.Command, args []string, entityType string) {
	if len(args) < 1 {
		fmt.Println("Expected an " + entityType + " id to specify what " + entityType + " to use as a parameter")
		ccmd.Usage()
		os.Exit(1)
	}
	state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
	printer.PrintGenericEntityProperties(state, args[0], entityType)
}

func getEntityLayout(ccmd *cobra.Command, args []string, entityType string) {
	if len(args) < 1 {
		fmt.Println("Expected an " + entityType + " id to specify what " + entityType + " to use as a parameter")
		ccmd.Usage()
		os.Exit(1)
	}
	state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
	printer.PrintGenericEntityLayout(state, args[0], entityType)
}

func getEntities(ccmd *cobra.Command, args []string, entityType string, printAsJSON bool) {
	state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
	allInfos, err := state.Doc.GetAllInfos(rootCtx)
	if err != nil {
		internal.FatalError(err)
	}
	printer.PrintGenericEntities(allInfos, entityType, printAsJSON)
}

func build(ccmd *cobra.Command, args []string) {
	ctx := rootCtx
	state := internal.PrepareEngineState(ctx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)

	separateConnectionsFile := ccmd.Flag("connections").Value.String()
	if separateConnectionsFile == "" {
		separateConnectionsFile = GetRelativeParameter("connections")
	}
	internal.SetupConnections(ctx, state.Doc, separateConnectionsFile, viper.ConfigFileUsed())
	internal.SetupEntities(ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("dimensions").Value.String(), "dimension")
	internal.SetupEntities(ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("measures").Value.String(), "measure")
	internal.SetupEntities(ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String(), "object")
	scriptFile := ccmd.Flag("script").Value.String()
	if scriptFile == "" {
		scriptFile = GetRelativeParameter("script")
	}
	if scriptFile != "" {
		internal.SetScript(ctx, state.Doc, scriptFile)
	}

	silent := viper.GetBool("silent")

	internal.Reload(ctx, state.Doc, state.Global, silent, true)

	if state.AppID != "" {
		internal.Save(ctx, state.Doc, state.AppID)
	}
}
