package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = withLocalFlags(&cobra.Command{
	Use:   "build",
	Args: cobra.ExactArgs(0),
	Short: "Reload and save the app after updating connections, dimensions, measures, objects and the script",
	Example: `corectl build
corectl build --connections ./myconnections.yml --script ./myscript.qvs`,
	Annotations: map[string]string{
		"command_category": "build",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ctx := rootCtx
		state := internal.PrepareEngineState(ctx, headers, true)

		separateConnectionsFile := ccmd.Flag("connections").Value.String()
		if separateConnectionsFile == "" {
			separateConnectionsFile = getPathFlagFromConfigFile("connections")
		}
		internal.SetupConnections(ctx, state.Doc, separateConnectionsFile)
		internal.SetupEntities(ctx, state.Doc, ccmd.Flag("dimensions").Value.String(), "dimension")
		internal.SetupEntities(ctx, state.Doc, ccmd.Flag("measures").Value.String(), "measure")
		internal.SetupEntities(ctx, state.Doc, ccmd.Flag("objects").Value.String(), "object")
		scriptFile := ccmd.Flag("script").Value.String()
		if scriptFile == "" {
			scriptFile = getPathFlagFromConfigFile("script")
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
}, "script", "connections", "dimensions", "measures", "objects", "no-reload", "silent", "no-save")

var reloadCmd = withLocalFlags(&cobra.Command{
	Use:     "reload",
	Args: cobra.ExactArgs(0),
	Short:   "Reload and save the app",
	Long:    "Reload and save the app",
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
