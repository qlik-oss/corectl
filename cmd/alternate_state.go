package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listAlternateStatesCmd = withLocalFlags(&cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all alternate states in the current app",
	Long:    "Print a list of all alternate states in the current app",
	Example: "corectl state ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		items := internal.ListAlternateStates(state.Ctx, state.Doc)
		printer.PrintStates(items, viper.GetBool("bash"))
	},
}, "quiet")

var addAlternateStateCmd = &cobra.Command{
	Use:     "add <alternate-state-name>",
	Args:    cobra.ExactArgs(1),
	Short:   "Add an alternate states in the current app",
	Long:    "Add an alternate states in the current app",
	Example: "corectl state add NAME-1",

	Run: func(ccmd *cobra.Command, args []string) {
		stateName := args[0]
		if stateName == "" {
			log.Fatalln("no state name specified")
		}
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		internal.AddAlternateState(state.Ctx, state.Doc, stateName)
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var removeAlternateStateCmd = &cobra.Command{
	Use:     "rm <alternate-state-name>",
	Args:    cobra.ExactArgs(1),
	Short:   "Removes an alternate state in the current app",
	Long:    "Removes an alternate state in the current app",
	Example: "corectl state rm NAME-1",

	Run: func(ccmd *cobra.Command, args []string) {
		stateName := args[0]
		if stateName == "" {
			log.Fatalln("no state name specified")
		}
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		internal.RemoveAlternateState(state.Ctx, state.Doc, stateName)
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var alternateStateCmd = &cobra.Command{
	Use:     "state",
	Short:   "Explore and manage alternate states",
	Long:    "Explore and manage alternate states",
	Aliases: []string{"alternatestate"},
	Annotations: map[string]string{
		"command_category": "sub",
		"x-qlik-stability": "experimental",
	},
}

func init() {
	alternateStateCmd.AddCommand(listAlternateStatesCmd, addAlternateStateCmd, removeAlternateStateCmd)
}
