package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = withCommonLocalFlags(&cobra.Command{
	Use:     "build",
	Short:   "Reloads and saves the app after updating connections, dimensions, measures, objects and the script",
	Example: "corectl build --connections ./myconnections.yml --script ./myscript.qvs",
	Annotations: map[string]string{
		"command_category": "build",
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

		if !viper.GetBool("no-reload") {
			silent := viper.GetBool("silent")
			internal.Reload(ctx, state.Doc, state.Global, silent, true)
		}

		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(ctx, state.Doc)
		}
	},
}, "no-reload", "silent", "no-save")

var reloadCmd = withCommonLocalFlags(&cobra.Command{
	Use:     "reload",
	Short:   "Reloads the app.",
	Long:    "Reloads the app.",
	Example: "corectl reload",
	Annotations: map[string]string{
		"command_category": "build",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		silent := viper.GetBool("silent")

		internal.Reload(rootCtx, state.Doc, state.Global, silent, true)

		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "silent", "no-save")

func init() {
	// Don't bind these to viper since paths are treated separately to support relative paths!
	buildCmd.PersistentFlags().String("connections", "", "Path to a yml file containing the data connection definitions")
	// Don't bind these to viper since paths are treated separately to support relative paths!
	buildCmd.PersistentFlags().String("dimensions", "", "A list of generic dimension json paths")
	// Don't bind these to viper since paths are treated separately to support relative paths!
	buildCmd.PersistentFlags().String("measures", "", "A list of generic measures json paths")
	// Don't bind these to viper since paths are treated separately to support relative paths!
	buildCmd.PersistentFlags().String("objects", "", "A list of generic object json paths")
	// Don't bind these to viper since paths are treated separately to support relative paths!
	buildCmd.PersistentFlags().String("script", "", "path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded")

}
