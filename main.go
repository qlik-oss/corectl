package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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
	params             struct {
		engine    string
		appID     string
		sessionID string
		ttl       string
		headers   http.Header
	}
	rootCtx = context.Background()

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
			params.engine = viper.GetString("engine")
			params.appID = viper.GetString("app")
			params.ttl = viper.GetString("ttl")

			if len(headersMap) == 0 {
				headersMap = viper.GetStringMapString("headers")
			}
			params.headers = make(http.Header, 1)
			for key, value := range headersMap {
				params.headers.Set(key, value)
			}
		},

		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	fieldsCommand = &cobra.Command{
		Use:   "fields",
		Short: "Print field list",
		Long:  "Print field list",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, params.headers, false)
			printer.PrintFields(data, false)
		},
	}

	keysCommand = &cobra.Command{
		Use:   "keys",
		Short: "Print key-only field list",
		Long:  "Print key-only field list",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, params.headers, true)
			printer.PrintFields(data, true)
		},
	}

	associationsCommand = &cobra.Command{
		Use:   "assoc",
		Short: "Print table associations summary",
		Long:  "Print table associations summary",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, params.headers, false)
			printer.PrintAssociations(data)
		},
	}

	tablesCommand = &cobra.Command{
		Use:   "tables",
		Short: "Print tables summary",
		Long:  "Prints tables summary",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, params.headers, false)
			printer.PrintTables(data)
		},
	}
	evalCmd = &cobra.Command{
		Use:   "eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...]",
		Short: "Evalutes a list of measures and dimensions",
		Long:  `Evalutes a list of measures and dimensions. Meaures are separeted from dimensions by the "by" keyword. To omit dimensions and only use measures use "*" as dimension: eval <measures> by *`,

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Expected at least one dimension or measure")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			internal.Eval(rootCtx, state.Doc, args)
		},
	}

	getScriptCmd = &cobra.Command{
		Use:   "script",
		Short: "Print the reload script",
		Long:  "Print the reload script",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			script, err := state.Doc.GetScript(rootCtx)
			if err != nil {
				internal.FatalError(err)
			}
			fmt.Println(script)
		},
	}

	metaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Shows metadata about the app",
		Long:  "Lists tables, fields, associations along with metadata like memory consumption, field cardinality etc",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, params.headers, false)
			printer.PrintMetadata(data)
		},
	}

	reloadCmd = &cobra.Command{
		Use:   "reload",
		Short: "Reloads and saves the app after updating connections, objects and the script",
		Long: `Reloads the app. Example: corectl reload --connections ./myconnections.yml --script ./myscript.qvs
			
`,

		Run: func(ccmd *cobra.Command, args []string) {
			updateOrReload(ccmd, args, true)
		},
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates connections, objects and script and saves the app",
		Long: `Updates connections, objects and script in the app. Example: corectl update	

`,

		Run: func(ccmd *cobra.Command, args []string) {
			updateOrReload(ccmd, args, false)
		},
	}

	fieldCmd = &cobra.Command{
		Use:   "field <field name>",
		Short: "Shows content of a field",
		Long:  ``,

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Expected a field name as parameter")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			internal.PrintField(rootCtx, state.Doc, args[0])
		},
	}

	statusCommand = &cobra.Command{
		Use:   "status",
		Short: "Prints status info about the connection to engine and current app",
		Long:  "Prints status info about the connection to engine and current app",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			engine := params.engine
			if engine == "" {
				engine = "localhost:9076"
			}
			printer.PrintStatus(state, engine)
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

	listAppsCmd = &cobra.Command{
		Use:   "apps",
		Short: "Prints a list of all apps available in the current engine",
		Long:  "Prints a list of all apps available in the current engine",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineStateWithoutApp(rootCtx, params.engine, params.ttl, params.headers)
			docList, err := state.Global.GetDocList(rootCtx)
			if err != nil {
				internal.FatalError(err)
			}
			printer.PrintApps(docList, viper.GetBool("json"))
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of corectl",

		Run: func(_ *cobra.Command, args []string) {
			fmt.Printf("corectl version %s\n", version)
		},
	}

	listObjectsCmd = &cobra.Command{
		Use:   "objects",
		Short: "Prints a list of all objects in the current app",
		Long:  "Prints a list of all objects in the current app",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			internal.SetupObjects(state.Ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String())
			allInfos, err := state.Doc.GetAllInfos(rootCtx)
			if err != nil {
				internal.FatalError(err)
			}
			printer.PrintObjects(allInfos)
		},
	}

	getObjectPropertiesCmd = &cobra.Command{
		Use:   "properties",
		Short: "Prints the properties of the object identified by the --object flag",
		Long:  "Prints the properties of the object identified by the --object flag",

		Run: func(ccmd *cobra.Command, args []string) {
			objectID := ccmd.Flag("object").Value.String()
			if objectID == "" {
				fmt.Println("Expected a --object flag to specify what object to use")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			internal.SetupObjects(state.Ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String())
			printer.PrintObject(state, objectID)
		},
	}

	getObjectLayoutCmd = &cobra.Command{
		Use:   "layout",
		Short: "Evalutes the hypercube layout of an object defined by the --object parameter",
		Long:  `Evalutes the hypercube layout of an object defined by the --object parameter`,

		Run: func(ccmd *cobra.Command, args []string) {
			objectID := ccmd.Flag("object").Value.String()
			if objectID == "" {
				fmt.Println("Expected a --object flag to specify what object to use")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			internal.SetupObjects(state.Ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String())
			printer.PrintObjectLayout(state, objectID)
		},
	}

	getObjectDataCmd = &cobra.Command{
		Use:   "data",
		Short: "Evalutes the hypercube data of an object defined by the --object parameter. Note that only basic hypercubes like straight tables are supported",
		Long:  `Evalutes the hypercube data of an object defined by the --object parameter. Note that only basic hypercubes like straight tables are supported`,

		Run: func(ccmd *cobra.Command, args []string) {
			objectID := ccmd.Flag("object").Value.String()
			if objectID == "" {
				fmt.Println("Expected a --object flag to specify what object to use")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, params.headers, false)
			internal.SetupObjects(state.Ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String())
			printer.EvalObject(rootCtx, state.Doc, objectID)
		},
	}
)

func init() {
	corectlCommand.PersistentFlags().StringVarP(&explicitConfigFile, "config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")

	corectlCommand.PersistentFlags().StringP("engine", "e", "", "URL to engine (default \"localhost:9076\")")
	viper.BindPFlag("engine", corectlCommand.PersistentFlags().Lookup("engine"))

	corectlCommand.PersistentFlags().String("ttl", "30", "Engine session time to live")
	viper.BindPFlag("ttl", corectlCommand.PersistentFlags().Lookup("ttl"))

	corectlCommand.PersistentFlags().StringP("app", "a", "", "App name including .qvf file ending. If no app is specified a session app is used instead.")
	viper.BindPFlag("app", corectlCommand.PersistentFlags().Lookup("app"))

	//not binding to viper since binding a map does not seem to work.
	corectlCommand.PersistentFlags().StringToStringVar(&headersMap, "headers", nil, "Headers to use when connecting to qix engine")

	corectlCommand.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	viper.BindPFlag("verbose", corectlCommand.PersistentFlags().Lookup("verbose"))

	reloadCmd.PersistentFlags().Bool("silent", false, "Do not log reload progress")
	viper.BindPFlag("silent", reloadCmd.PersistentFlags().Lookup("silent"))

	listAppsCmd.PersistentFlags().Bool("json", false, "Prints the apps in json format")
	viper.BindPFlag("json", listAppsCmd.PersistentFlags().Lookup("json"))

	for _, command := range []*cobra.Command{reloadCmd, updateCmd, getObjectPropertiesCmd, getObjectLayoutCmd, getObjectDataCmd, listObjectsCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("objects", "", "A list of object json paths")
	}
	for _, command := range []*cobra.Command{getObjectPropertiesCmd, getObjectLayoutCmd, getObjectDataCmd} {
		//Don't bind to vibe rsince this parameter is purely interactive
		command.PersistentFlags().StringP("object", "o", "", "ID of a generic object")
	}
	for _, command := range []*cobra.Command{reloadCmd, updateCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("script", "", "path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded")
		command.PersistentFlags().String("connections", "", "path/to/connections.yml that contains connections that are used in the reload. Note that when specifying connections in the config file they are specified inline, not as a file reference!")
	}

	// commands
	corectlCommand.AddCommand(reloadCmd)
	corectlCommand.AddCommand(updateCmd)
	corectlCommand.AddCommand(evalCmd)
	corectlCommand.AddCommand(metaCmd)
	corectlCommand.AddCommand(getScriptCmd)
	corectlCommand.AddCommand(fieldsCommand)
	corectlCommand.AddCommand(keysCommand)
	corectlCommand.AddCommand(tablesCommand)
	corectlCommand.AddCommand(fieldCmd)
	corectlCommand.AddCommand(associationsCommand)
	corectlCommand.AddCommand(statusCommand)
	corectlCommand.AddCommand(generateDocsCommand)
	corectlCommand.AddCommand(listAppsCmd)
	corectlCommand.AddCommand(versionCmd)
	corectlCommand.AddCommand(listObjectsCmd)
	corectlCommand.AddCommand(getObjectPropertiesCmd)
	corectlCommand.AddCommand(getObjectLayoutCmd)
	corectlCommand.AddCommand(getObjectDataCmd)
}

// GetPathParameter returns a parameter from either the command line or the config file.
// Compared to using BindPFlag this function modifies relative paths in the config file
// to actually be relative to the config file and not the working directory
func GetPathParameter(ccmd *cobra.Command, paramName string) string {
	if pathOnCommandLine := ccmd.Flag(paramName).Value.String(); pathOnCommandLine != "" {
		return pathOnCommandLine
	}
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

func updateOrReload(ccmd *cobra.Command, args []string, reload bool) {
	ctx := rootCtx
	state := internal.PrepareEngineState(ctx, params.engine, params.appID, params.ttl, params.headers, true)

	separateConnectionsFile := GetPathParameter(ccmd, "connections")
	internal.SetupConnections(ctx, state.Doc, separateConnectionsFile, viper.ConfigFileUsed())
	internal.SetupObjects(ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String())
	scriptFile := GetPathParameter(ccmd, "script")
	if scriptFile != "" {
		internal.SetScript(ctx, state.Doc, scriptFile)
	}

	silent := viper.GetBool("silent")

	if reload {
		internal.Reload(ctx, state.Doc, state.Global, silent, true)
	}
	if state.AppID != "" {
		internal.Save(ctx, state.Doc, state.AppID)
	}
}
