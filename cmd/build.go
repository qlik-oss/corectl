package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/spf13/cobra"
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
		ctx, global, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(true)

		var config *boot.ConnectionsConfig

		if params.IsString("connections") {
			connectionsPath := params.GetPath("connections")
			if connectionsPath != "" {
				config = boot.ReadConnectionsFile(connectionsPath)
			} else {
				config = nil
			}
		} else {
			config = boot.ReadConnectionsFile(params.ConfigFilePath())
		}
		internal.SetupConnections(ctx, doc, config)
		internal.SetDimensions(ctx, doc, params.GetGlobFiles("dimensions"))
		internal.SetVariables(ctx, doc, params.GetGlobFiles("variables"))
		internal.SetMeasures(ctx, doc, params.GetGlobFiles("measures"))
		internal.SetObjects(ctx, doc, params.GetGlobFiles("objects"))

		if params.GetPath("script") != "" {
			internal.SetScript(ctx, doc, params.GetPath("script"))
		}

		if params.GetPath("app-properties") != "" {
			internal.SetAppProperties(ctx, doc, params.GetPath("app-properties"))
		}

		if !params.GetBool("no-reload") {
			silent := params.GetBool("silent")
			limit := params.GetInt("limit")
			internal.Reload(ctx, doc, global, silent, limit)
		}

		if !params.NoSave() {
			internal.Save(ctx, doc, params.NoData())
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
		ctx, global, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)

		internal.Reload(ctx, doc, global, params.GetBool("silent"), params.GetInt("limit"))

		if !params.NoSave() {
			internal.Save(ctx, doc, params.NoData())
		}
	},
}, "silent", "no-save", "limit")
