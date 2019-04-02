package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var setAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Sets the objects, measures, dimensions, connections and script in the current app",
	Long:  "Sets the objects, measures, dimensions, connections and script in the current app",
	Example: `corectl set all
corectl set all --app=my-app.qvf`,

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
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Reloads and saves the app after updating connections, dimensions, measures, objects and the script",
	Example: "corectl build --connections ./myconnections.yml --script ./myscript.qvs",
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
		viper.BindPFlag("silent", ccmd.PersistentFlags().Lookup("silent"))
	},
	Run: func(ccmd *cobra.Command, args []string) {
		ctx := rootCtx
		state := internal.PrepareEngineState(ctx, headers, true)

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

		doReload := true
		silent := viper.GetBool("silent")
		if doReload {
			internal.Reload(ctx, state.Doc, state.Global, silent, true)
		}

		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(ctx, state.Doc)
		}
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

var reloadCmd = &cobra.Command{
	Use:     "reload",
	Short:   "Reloads the app.",
	Long:    "Reloads the app.",
	Example: "corectl reload",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)

		viper.BindPFlag("silent", ccmd.PersistentFlags().Lookup("silent"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		silent := viper.GetBool("silent")

		internal.Reload(rootCtx, state.Doc, state.Global, silent, true)

		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

func init() {
	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("connections", "", "Path to a yml file containing the data connection definitions")
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("dimensions", "", "A list of generic dimension json paths")
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("measures", "", "A list of generic measures json paths")
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("objects", "", "A list of generic object json paths")
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("script", "", "path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded")
	}
}
