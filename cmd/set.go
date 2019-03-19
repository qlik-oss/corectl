package cmd

import (
	"fmt"
	"os"

	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Set commands
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets one or several resources",
	Long:  "Sets one or several resources",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
		viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
		viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
		viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
	},
}

var setAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Sets the objects, measures, dimensions, connections and script in the current app",
	Long:  "Sets the objects, measures, dimensions, connections and script in the current app",
	Example: `corectl set all
corectl set all --app=my-app.qvf`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		setCmd.PersistentPreRun(setCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {

		state := internal.PrepareEngineState(rootCtx, headers, true)
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

var setConnectionsCmd = &cobra.Command{
	Use:     "connections <path-to-connections-file.yml>",
	Short:   "Sets or updates the connections in the current app",
	Long:    "Sets or updates the connections in the current app",
	Example: "corectl set connections ./my-connections.yml",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		setCmd.PersistentPreRun(setCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {

		state := internal.PrepareEngineState(rootCtx, headers, true)
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

var setDimensionsCmd = &cobra.Command{
	Use:     "dimensions <glob-pattern-path-to-dimensions-files.json>",
	Short:   "Sets or updates the dimensions in the current app",
	Long:    "Sets or updates the dimensions in the current app",
	Example: "corectl set dimensions ./my-dimensions-glob-path.json",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		setCmd.PersistentPreRun(setCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {

		commandLineDimensions := ""
		if len(args) > 0 {
			commandLineDimensions = args[0]
		}
		state := internal.PrepareEngineState(rootCtx, headers, true)
		internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineDimensions, "dimension")
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

var setMeasuresCmd = &cobra.Command{
	Use:     "measures <glob-pattern-path-to-measures-files.json>",
	Short:   "Sets or updates the measures in the current app",
	Long:    "Sets or updates the measures in the current app",
	Example: "corectl set measures ./my-measures-glob-path.json",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		setCmd.PersistentPreRun(setCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {

		commandLineMeasures := ""
		if len(args) > 0 {
			commandLineMeasures = args[0]
		}
		state := internal.PrepareEngineState(rootCtx, headers, true)
		internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineMeasures, "measure")
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

var setObjectsCmd = &cobra.Command{
	Use:   "objects <glob-pattern-path-to-objects-files.json",
	Short: "Sets or updates the objects in the current app",
	Long: `Sets or updates the objects in the current app.
The JSON objects can be in either the GenericObjectProperties format or the GenericObjectEntry format`,
	Example: "corectl set objects ./my-objects-glob-path.json",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		setCmd.PersistentPreRun(setCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {

		commandLineObjects := ""
		if len(args) > 0 {
			commandLineObjects = args[0]
		}

		state := internal.PrepareEngineState(rootCtx, headers, true)
		internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineObjects, "object")
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

var setScriptCmd = &cobra.Command{
	Use:     "script <path-to-script-file.yml>",
	Short:   "Sets the script in the current app",
	Long:    "Sets the script in the current app",
	Example: "corectl set script ./my-script-file",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		setCmd.PersistentPreRun(setCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {

		state := internal.PrepareEngineState(rootCtx, headers, true)
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

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.AddCommand(setAllCmd)
	setCmd.AddCommand(setConnectionsCmd)
	setCmd.AddCommand(setDimensionsCmd)
	setCmd.AddCommand(setMeasuresCmd)
	setCmd.AddCommand(setObjectsCmd)
	setCmd.AddCommand(setScriptCmd)
}
