package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = withLocalFlags(&cobra.Command{
	Use:   "build",
	Args:  cobra.ExactArgs(0),
	Short: "Reload and save the app after updating connections, dimensions, measures, objects and the script",
	Example: `corectl build
corectl build --connections ./myconnections.yml --script ./myscript.qvs`,
	Annotations: map[string]string{
		"command_category": "build",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ctx := rootCtx
		state := internal.PrepareEngineState(ctx, headers, tlsClientConfig, true, false)

		separateConnectionsFile := ccmd.Flag("connections").Value.String()
		if separateConnectionsFile == "" {
			separateConnectionsFile = getPathFlagFromConfigFile("connections")
		}
		internal.SetupConnections(ctx, state.Doc, separateConnectionsFile)
		internal.SetDimensions(ctx, state.Doc, ccmd.Flag("dimensions").Value.String())
		internal.SetVariables(ctx, state.Doc, ccmd.Flag("variables").Value.String())
		internal.SetMeasures(ctx, state.Doc, ccmd.Flag("measures").Value.String())
		internal.SetObjects(ctx, state.Doc, ccmd.Flag("objects").Value.String())
		scriptFile := ccmd.Flag("script").Value.String()
		if scriptFile == "" {
			scriptFile = getPathFlagFromConfigFile("script")
		}
		if scriptFile != "" {
			internal.SetScript(ctx, state.Doc, scriptFile)
		}

		appProperties := ccmd.Flag("app-properties").Value.String()
		if appProperties == "" {
			appProperties = getPathFlagFromConfigFile("app-properties")
		}
		if appProperties != "" {
			internal.SetAppProperties(ctx, state.Doc, appProperties)
		}

		if !viper.GetBool("no-reload") {
			silent := viper.GetBool("silent")
			limit := viper.GetInt("limit")
			internal.Reload(ctx, state.Doc, state.Global, silent, limit)
		}

		if !viper.GetBool("no-save") {
			internal.Save(ctx, state.Doc)
		}
	},
}, "script", "app-properties", "connections", "dimensions", "measures", "variables", "bookmarks", "objects", "no-reload", "silent", "no-save", "limit")

var reloadCmd = withLocalFlags(&cobra.Command{
	Use:     "reload",
	Args:    cobra.ExactArgs(0),
	Short:   "Reload and save the app",
	Long:    "Reload and save the app",
	Example: "corectl reload",
	Annotations: map[string]string{
		"command_category": "build",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		silent := viper.GetBool("silent")
		limit := viper.GetInt("limit")

		internal.Reload(rootCtx, state.Doc, state.Global, silent, limit)

		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "silent", "no-save", "limit")
