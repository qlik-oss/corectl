package cmd

import (
	"context"

	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var unbuildCmd = withLocalFlags(&cobra.Command{
	Use:   "unbuild",
	Args:  cobra.ExactArgs(0),
	Short: "Split up an existing app into separate json and yaml files",
	Long: `Extracts generic objects, dimensions, measures, variables, reload script and connections from an app in an engine into separate json and yaml files.
In addition to the resources from the app a corectl.yml configuration file is generated that binds them all together.
Passwords in the connection definitions can not be exported from the app and hence need to be handled manually.
Generic Object trees (e.g. Qlik Sense sheets) are exported as a full property tree which means that child objects are found inside the parent´s json (the qChildren array).
`,
	Example: `corectl unbuild
corectl unbuild --app APP-ID`,
	Annotations: map[string]string{
		"command_category": "build",
		"x-qlik-stability": "experimental",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ctx := rootCtx
		viper.Set("no-data", "true") // Force no-data since we only use metadata
		outdir := ccmd.Flag("dir").Value.String()
		state := internal.PrepareEngineState(ctx, headers, tlsClientConfig, false, false)
		if outdir == DefaultUnbuildFolder {
			outdir = getDefaultOutDir(ctx, state)
		}
		internal.Unbuild(ctx, state.Doc, state.Global, outdir)
	},
}, "dir")

func getDefaultOutDir(ctx context.Context, state *internal.State) string {
	appLayout, _ := state.Doc.GetAppLayout(ctx)
	var defaultFolder string
	if appLayout.Title != "" {
		defaultFolder = appLayout.Title
	} else if state.AppName != "" {
		defaultFolder = state.AppName
	} else {
		defaultFolder = state.AppID
	}
	return internal.BuildRootFolderFromTitle(defaultFolder)
}
