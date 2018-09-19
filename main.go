package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var (
	state  *internal.State
	config string

	corectlCommand = &cobra.Command{
		Hidden: true,
		Use:    "corectl",
		Short:  "",
		Long:   `Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,

		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

			internal.QliVerbose = viper.GetBool("verbose")
			if config != "" {
				abs, err := filepath.Abs(config)
				if err != nil {
					fmt.Println("Error reading filepath: ", err.Error())
				}
				base := filepath.Base(abs)
				path := filepath.Dir(abs)
				viper.SetConfigName(strings.Split(base, ".")[0])
				viper.SetConfigType("yml")
				viper.AddConfigPath(path)

				if err := viper.ReadInConfig(); err == nil {
					internal.LogVerbose("Using config file: " + config)
				} else {
					fmt.Println(err)
				}
			} else {
				viper.SetConfigName("qli") // name of config file (without extension)
				viper.SetConfigType("yml")
				//viper.AddConfigPath("/etc/qli/") // paths to look for the config file
				//viper.AddConfigPath("$HOME/.qli")
				viper.AddConfigPath(".")

				if err := viper.ReadInConfig(); err == nil {
					internal.LogVerbose("Using config file in working directory")
				} else {
					internal.LogVerbose("No config file")
				}
			}

			internal.QliVerbose = viper.GetBool("verbose")
			engine := viper.GetString("engine")
			appID := viper.GetString("app")
			ttl := viper.GetString("ttl")
			sessionID := getSessionID(appID)
			state = internal.PrepareEngineState(ctx, engine, sessionID, appID, ttl)
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
			data := internal.GetModelMetadata(state.Ctx, state.Doc, state.MetaURL, false)
			printer.PrintFields(data, false)
		},
	}

	keysCommand = &cobra.Command{
		Use:   "keys",
		Short: "Print key-only field list",
		Long:  "Print key-only field list",

		Run: func(ccmd *cobra.Command, args []string) {
			data := internal.GetModelMetadata(state.Ctx, state.Doc, state.MetaURL, true)
			printer.PrintFields(data, true)

		},
	}

	associationsCommand = &cobra.Command{
		Use:   "assoc",
		Short: "Print table associations summary",
		Long:  "Print table associations summary",

		Run: func(ccmd *cobra.Command, args []string) {
			data := internal.GetModelMetadata(state.Ctx, state.Doc, state.MetaURL, false)
			printer.PrintAssociations(data)
		},
	}

	tablesCommand = &cobra.Command{
		Use:   "tables",
		Short: "Print tables summary",
		Long:  "Prints tables summary",

		Run: func(ccmd *cobra.Command, args []string) {
			data := internal.GetModelMetadata(state.Ctx, state.Doc, state.MetaURL, false)
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
			ctx := state.Ctx
			doc := state.Doc

			internal.Eval(ctx, doc, args)
		},
	}

	getScriptCmd = &cobra.Command{
		Use:   "script",
		Short: "Print the reload script",
		Long:  "Print the reload script",

		Run: func(ccmd *cobra.Command, args []string) {

			ctx := state.Ctx
			doc := state.Doc
			//global := state.Global

			script, err := doc.GetScript(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(script)
		},
	}

	metaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Shows metadata about the app",
		Long:  "Lists tables, fields, associations along with metadata like memory consumption, field cardinality etc",

		Run: func(ccmd *cobra.Command, args []string) {
			data := internal.GetModelMetadata(state.Ctx, state.Doc, state.MetaURL, false)
			printer.PrintMetadata(data)
		},
	}

	reloadCmd = &cobra.Command{
		Use:   "reload",
		Short: "Reloads the app",
		Long:  `Reloads the app. Example: corectl reload --connections ./myconnections.yml --script ./myscript.qvs`,

		Run: func(ccmd *cobra.Command, args []string) {

			ctx := state.Ctx
			doc := state.Doc
			global := state.Global

			connectionsFile := GetPathParameter(ccmd, "include-connections")
			if connectionsFile != "" {
				internal.SetupConnections(ctx, doc, connectionsFile, viper.ConfigFileUsed())
			} else {
				internal.SetupConnections(ctx, doc, "", viper.ConfigFileUsed())
			}

			scriptFile := GetPathParameter(ccmd, "script")
			if scriptFile != "" {
				internal.SetScript(ctx, doc, scriptFile)
			}

			internal.Reload(ctx, doc, global, true)
			if state.AppID != "" {
				internal.Save(ctx, doc, state.AppID)
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
			internal.PrintField(state.Ctx, state.Doc, args[0])
		},
	}

	statusCommand = &cobra.Command{
		Use:   "status",
		Short: "Prints status info about the connection to engine and current app",
		Long:  "Prints status info about the connection to engine and current app",

		Run: func(ccmd *cobra.Command, args []string) {
			printer.PrintStatus(state)
		},
	}
)

func getSessionID(appID string) string {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sessionID := base64.StdEncoding.EncodeToString([]byte("Corectl-" + currentUser.Username + "-" + hostName + "-" + appID))
	return sessionID
}

func init() {

	// Flags
	corectlCommand.PersistentFlags().StringVarP(&config, "config", "c", "", "path/to/config.yml where default parameters can be set")

	corectlCommand.PersistentFlags().StringP("engine", "e", "localhost", "URL to engine")
	viper.BindPFlag("engine", corectlCommand.PersistentFlags().Lookup("engine"))

	corectlCommand.PersistentFlags().String("ttl", "30", "Engine session time to live")
	viper.BindPFlag("ttl", corectlCommand.PersistentFlags().Lookup("ttl"))

	corectlCommand.PersistentFlags().String("engine-headers", "30", "HTTP headers to send to the engine")
	viper.BindPFlag("engine-headers", corectlCommand.PersistentFlags().Lookup("engine-headers"))

	corectlCommand.PersistentFlags().StringP("app", "a", "", "App name including .qvf file ending")
	viper.BindPFlag("app", corectlCommand.PersistentFlags().Lookup("app"))

	corectlCommand.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	viper.BindPFlag("verbose", corectlCommand.PersistentFlags().Lookup("verbose"))

	evalCmd.PersistentFlags().StringP("select", "s", "", "")
	viper.BindPFlag("select", evalCmd.PersistentFlags().Lookup("select"))

	fieldsCommand.PersistentFlags().StringP("select", "s", "", "")
	viper.BindPFlag("select", fieldsCommand.PersistentFlags().Lookup("select"))

	reloadCmd.PersistentFlags().String("include-connections", "", "path to connections file")
	//viper.BindPFlag("connections", reloadCmd.PersistentFlags().Lookup("connections"))

	reloadCmd.PersistentFlags().String("script", "", "Script file name")
	//viper.BindPFlag("script", reloadCmd.PersistentFlags().Lookup("script"))

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
