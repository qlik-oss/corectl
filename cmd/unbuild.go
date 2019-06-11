package cmd

import (
	"context"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var unbuildCmd = withLocalFlags(&cobra.Command{
	Use:     "unbuild",
	Args:    cobra.ExactArgs(0),
	Short:   "Split upp an existing app into separate entities",
	Long:    "Split upp an existing app into separate entities",
	Example: `corectl unbuild`,
	Annotations: map[string]string{
		"command_category": "build",
		"x-qlik-stability": "experimental",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ctx := rootCtx
		viper.Set("no-data", "true") // Force no-data since we only use metadata
		outdir := ccmd.Flag("dir").Value.String()
		state := internal.PrepareEngineState(ctx, headers, false)
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
