package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
)

var unbuildCmd = withLocalFlags(&cobra.Command{
	Use:     "unbuild",
	Args:    cobra.ExactArgs(0),
	Short:   "Split upp an existing app into separate entities",
	Long:    "Split upp an existing app into separate entities",
	Example: `corectl unbuild`,
	Annotations: map[string]string{
		"command_category": "build",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ctx := rootCtx
		outdir := ccmd.Flag("dir").Value.String()
		state := internal.PrepareEngineState(ctx, headers, false)

		if outdir == DefaultUnbuild {
			appLayout, _ := state.Doc.GetAppLayout(ctx)
			if appLayout.Title != "" {
				outdir = internal.BuildRootFolderFromTitle(appLayout.Title)
			} else {
				outdir = "unknown-app-unbuild"
			}
		}

		internal.Unbuild(ctx, state.Doc, state.Global, outdir)
	},
}, "dir")
