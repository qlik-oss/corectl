package main

import (
	"fmt"
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
	explicitConfigFile = ""
	version            = ""
	params             struct {
		engine    string
		appID     string
		sessionID string
		ttl       string
	}
	rootCtx = context.Background()

	corectlCommand = &cobra.Command{
		Hidden:            true,
		Use:               "corectl",
		Short:             "",
		Long:              `Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
		DisableAutoGenTag: true,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			// if help or generate-docs command, no prerun is needed.
			if strings.Contains(ccmd.Use, "help") || ccmd.Use == "generate-docs" {
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
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, false)
			printer.PrintFields(data, false)
		},
	}

	keysCommand = &cobra.Command{
		Use:   "keys",
		Short: "Print key-only field list",
		Long:  "Print key-only field list",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, true)
			printer.PrintFields(data, true)
		},
	}

	associationsCommand = &cobra.Command{
		Use:   "assoc",
		Short: "Print table associations summary",
		Long:  "Print table associations summary",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, false)
			printer.PrintAssociations(data)
		},
	}

	tablesCommand = &cobra.Command{
		Use:   "tables",
		Short: "Print tables summary",
		Long:  "Prints tables summary",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, false)
			printer.PrintTables(data)
		},
	}
	evalCmd = &cobra.Command{
		Use:   "eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...]",
		Short: "Evalutes a hypercube",
		Long:  `Evalutes a list of measures and dimensions. Meaures are separeted from dimensions by the "by" keyword. To omit dimensions and only use measures use "*" as dimension: eval <measures> by *`,

		Run: func(ccmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Expected at least one dimension or measure")
				ccmd.Usage()
				os.Exit(1)
			}
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
			internal.Eval(rootCtx, state.Doc, args)
		},
	}

	getScriptCmd = &cobra.Command{
		Use:   "script",
		Short: "Print the reload script",
		Long:  "Print the reload script",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
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
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
			data := internal.GetModelMetadata(rootCtx, state.Doc, state.MetaURL, false)
			printer.PrintMetadata(data)
		},
	}

	reloadCmd = &cobra.Command{
		Use:   "reload",
		Short: "Reloads the app",
		Long: `Reloads the app. Example: corectl reload --connections ./myconnections.yml --script ./myscript.qvs
			
`,

		Run: func(ccmd *cobra.Command, args []string) {
			ctx := rootCtx
			state := internal.PrepareEngineState(ctx, params.engine, params.appID, params.ttl, true)

			separateConnectionsFile := GetPathParameter(ccmd, "connections")
			internal.SetupConnections(ctx, state.Doc, separateConnectionsFile, viper.ConfigFileUsed())

			scriptFile := GetPathParameter(ccmd, "script")
			if scriptFile != "" {
				internal.SetScript(ctx, state.Doc, scriptFile)
			}

			silent := viper.GetBool("silent")

			internal.Reload(ctx, state.Doc, state.Global, silent, true)
			if state.AppID != "" {
				internal.Save(ctx, state.Doc, state.AppID)
			}
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
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
			internal.PrintField(rootCtx, state.Doc, args[0])
		},
	}

	statusCommand = &cobra.Command{
		Use:   "status",
		Short: "Prints status info about the connection to engine and current app",
		Long:  "Prints status info about the connection to engine and current app",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineState(rootCtx, params.engine, params.appID, params.ttl, false)
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
		Short: "Print app list",
		Long:  "Print app list",

		Run: func(ccmd *cobra.Command, args []string) {
			state := internal.PrepareEngineStateWithoutApp(rootCtx, params.engine, params.ttl)
			docList, err := state.Global.GetDocList(rootCtx)
			if err != nil {
				internal.FatalError(err)
			}
			printer.PrintApps(docList)
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Output the version of corectl",

		Run: func(_ *cobra.Command, args []string) {
			fmt.Printf("corectl version %s", version)
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

	corectlCommand.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	viper.BindPFlag("verbose", corectlCommand.PersistentFlags().Lookup("verbose"))

	reloadCmd.PersistentFlags().Bool("silent", false, "Do not log reload progress")
	viper.BindPFlag("silent", reloadCmd.PersistentFlags().Lookup("silent"))

	// Don't bind these to viper since paths are treated separately to support relative paths!
	reloadCmd.PersistentFlags().String("connections", "", "path/to/connections.yml that contains connections that are used in the reload. Note that when specifying connections in the config file they are specified inline, not as a file reference!")
	reloadCmd.PersistentFlags().String("script", "", "path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded")

	// commands
	corectlCommand.AddCommand(reloadCmd)
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
