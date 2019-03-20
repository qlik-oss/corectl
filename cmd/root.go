package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var headersMap = make(map[string]string)
var explicitConfigFile = ""
var version = ""
var headers http.Header
var rootCtx = context.Background()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Hidden:            true,
	Use:               "corectl",
	Short:             "",
	Long:              `Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
	DisableAutoGenTag: true,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		// if help, version or generate-docs command, no prerun is needed.
		if strings.Contains(ccmd.Use, "help") || ccmd.Use == "generate-docs" || ccmd.Use == "generate-API-spec" || ccmd.Use == "version" {
			return
		}
		internal.QliVerbose = viper.GetBool("verbose")
		internal.LogTraffic = viper.GetBool("traffic")
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolP("traffic", "t", false, "Log JSON websocket traffic to stdout")
	viper.BindPFlag("traffic", rootCmd.PersistentFlags().Lookup("traffic"))

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

	for _, command := range []*cobra.Command{buildCmd, evalCmd, getCmd, reloadCmd, setCmd, removeCmd} {
		//not binding to viper since binding a map does not seem to work.
		command.PersistentFlags().StringToStringVar(&headersMap, "headers", nil, "Headers to use when connecting to qix engine")
	}

	for _, command := range []*cobra.Command{buildCmd, catwalkCmd, evalCmd, getAssociationsCmd, getConnectionsCmd, getConnectionCmd, getDimensionsCmd, getDimensionCmd, getFieldsCmd, getKeysCmd, getFieldCmd, getMeasuresCmd, getMeasureCmd, getMetaCmd, getObjectsCmd, getObjectCmd, getScriptCmd, getStatusCmd, getTablesCmd, reloadCmd, removeCmd, setCmd} {
		command.PersistentFlags().StringP("app", "a", "", "App name, if no app is specified a session app is used instead.")
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

	for _, command := range []*cobra.Command{removeCmd} {
		command.PersistentFlags().Bool("suppress", false, "Suppress all confirmation dialogues")
	}

	catwalkCmd.PersistentFlags().String("catwalk-url", "https://catwalk.core.qlik.com", "Url to an instance of catwalk, if not provided the qlik one will be used.")
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
