package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setMeasuresCmd = withLocalFlags(&cobra.Command{
	Use:     "set <glob-pattern-path-to-measures-files.json>",
	Args:    cobra.ExactArgs(1),
	Short:   "Set or update the measures in the current app",
	Long:    "Set or update the measures in the current app",
	Example: "corectl measure set ./my-measures-glob-path.json",

	Run: func(ccmd *cobra.Command, args []string) {
		commandLineMeasures := args[0]
		if commandLineMeasures == "" {
			log.Fatalln("no measures specified")
		}
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, true, false)
		internal.SetMeasures(rootCtx, state.Doc, commandLineMeasures)
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var removeMeasureCmd = withLocalFlags(&cobra.Command{
	Use:     "rm <measure-id>...",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove one or many generic measures in the current app",
	Long:    "Remove one or many generic measures in the current app",
	Example: "corectl measure rm ID-1 ID-2",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyMeasure(rootCtx, entity)
			if err != nil {
				log.Fatalf("could not remove generic measure '%s': %s\n", entity, err)
			} else if !destroyed {
				log.Fatalf("could not remove generic measure '%s'\n", entity)
			}
		}
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var listMeasuresCmd = withLocalFlags(&cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all generic measures in the current app",
	Long:    "Print a list of all generic measures in the current app",
	Example: "corectl measure ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		items := internal.ListMeasures(state.Ctx, state.Doc)
		printer.PrintNamedItemsList(items, viper.GetBool("bash"), false)
	},
}, "quiet")

var getMeasurePropertiesCmd = withLocalFlags(&cobra.Command{
	Use:     "properties <measure-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Print the properties of the generic measure",
	Long:    "Print the properties of the generic measure",
	Example: "corectl measure properties MEASURE-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityProperties(state, args[0], "measure", viper.GetBool("minimum"))
	},
}, "minimum")

var getMeasureLayoutCmd = &cobra.Command{
	Use:     "layout <measure-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Evaluate the layout of an generic measure",
	Long:    "Evaluate the layout of an generic measure and prints in JSON format",
	Example: "corectl measure layout MEASURE-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityLayout(state, args[0], "measure")
	},
}

var measureCmd = &cobra.Command{
	Use:   "measure",
	Short: "Explore and manage measures",
	Long:  "Explore and manage measures",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	measureCmd.AddCommand(listMeasuresCmd, setMeasuresCmd, getMeasurePropertiesCmd, getMeasureLayoutCmd, removeMeasureCmd)
}
