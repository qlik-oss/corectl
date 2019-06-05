package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
)

var unbuildCmd = withLocalFlags(&cobra.Command{
	Use:     "unbuild",
	Args:    cobra.ExactArgs(0),
	Short:   "Split upp an existing qvf into separate entities",
	Example: `corectl unbuild`,
	Annotations: map[string]string{
		"command_category": "build",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ctx := rootCtx
		state := internal.PrepareEngineState(ctx, headers, true)
		internal.Unbuild(ctx, state.Doc, state.Global)

	},
})
