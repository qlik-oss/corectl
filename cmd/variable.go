package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setVariablesCmd = withLocalFlags(&cobra.Command{
	Use:     "set <glob-pattern-path-to-variables-files.json>",
	Args:    cobra.ExactArgs(1),
	Short:   "Set or update the variables in the current app",
	Long:    "Set or update the variables in the current app",
	Example: "corectl variable set ./my-variables-glob-path.json",

	Run: func(ccmd *cobra.Command, args []string) {
		commandLineVariables := args[0]
		if commandLineVariables == "" {
			log.Fatalln("no variables specified")
		}
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, true, false)
		internal.SetVariables(rootCtx, state.Doc, commandLineVariables)
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var removeVariableCmd = withLocalFlags(&cobra.Command{
	Use:     "rm <variable-name>...",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove one or many variables in the current app",
	Long:    "Remove one or many variables in the current app",
	Example: "corectl variable rm NAME-1",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyVariableByName(rootCtx, entity)
			if err != nil {
				log.Fatalf("could not remove generic variable '%s': %s\n", entity, err)
			} else if !destroyed {
				log.Fatalf("could not remove generic variable '%s'\n", entity)
			}
		}
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var listVariablesCmd = withLocalFlags(&cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all generic variables in the current app",
	Long:    "Print a list of all generic variables in the current app",
	Example: "corectl variable ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		items := internal.ListVariables(state.Ctx, state.Doc)
		printer.PrintNamedItemsList(items, viper.GetBool("bash"), true)
	},
}, "quiet")

var getVariablePropertiesCmd = withLocalFlags(&cobra.Command{
	Use:     "properties <variable-name>",
	Args:    cobra.ExactArgs(1),
	Short:   "Print the properties of the generic variable",
	Long:    "Print the properties of the generic variable",
	Example: "corectl variable properties VARIABLE-NAME",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityProperties(state, args[0], "variable", viper.GetBool("minimum"))
	},
}, "minimum")

var getVariableLayoutCmd = &cobra.Command{
	Use:     "layout <variable-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Evaluate the layout of an generic variable",
	Long:    "Evaluate the layout of an generic variable",
	Example: "corectl variable layout VARIABLE-NAME",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityLayout(state, args[0], "variable")
	},
}

var variableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Explore and manage variables",
	Long:  "Explore and manage variables",
	Annotations: map[string]string{
		"command_category": "sub",
		"x-qlik-stability": "experimental",
	},
}

func init() {
	variableCmd.AddCommand(setVariablesCmd, removeVariableCmd, listVariablesCmd, getVariablePropertiesCmd, getVariableLayoutCmd)
}
