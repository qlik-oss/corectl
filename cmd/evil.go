package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
)

var evilCmd = &cobra.Command{
	Use:   "evil",
	Args:  cobra.ExactArgs(0),
	Short: "",
	Long:  "",
	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		internal.Evil(rootCtx, state.Doc)
	},
}
