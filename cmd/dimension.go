package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setDimensionsCmd = withLocalFlags(&cobra.Command{
	Use:     "set <glob-pattern-path-to-dimensions-files.json>",
	Args:    cobra.ExactArgs(1),
	Short:   "Set or update the dimensions in the current app",
	Long:    "Set or update the dimensions in the current app",
	Example: "corectl dimension set ./my-dimensions-glob-path.json",

	Run: func(ccmd *cobra.Command, args []string) {
		commandLineDimensions := args[0]
		if commandLineDimensions == "" {
			log.Fatalln("no dimensions specified")
		}
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, true, false)
		internal.SetDimensions(rootCtx, state.Doc, commandLineDimensions)
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var removeDimensionCmd = withLocalFlags(&cobra.Command{
	Use:     "rm <dimension-id>...",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove one or many dimensions in the current app",
	Long:    "Remove one or many dimensions in the current app",
	Example: "corectl dimension rm ID-1",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyDimension(rootCtx, entity)
			if err != nil {
				log.Fatalf("could not remove generic dimension '%s': %s\n", entity, err)
			} else if !destroyed {
				log.Fatalf("could not remove generic dimension '%s'\n", entity)
			}
		}
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var listDimensionsCmd = withLocalFlags(&cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all generic dimensions in the current app",
	Long:    "Print a list of all generic dimensions in the current app",
	Example: "corectl dimension ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		items := internal.ListDimensions(state.Ctx, state.Doc)
		printer.PrintNamedItemsList(items, viper.GetBool("bash"), false)
	},
}, "quiet")

var getDimensionPropertiesCmd = withLocalFlags(&cobra.Command{
	Use:     "properties <dimension-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Print the properties of the generic dimension",
	Long:    "Print the properties of the generic dimension",
	Example: "corectl dimension properties DIMENSION-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityProperties(state, args[0], "dimension", viper.GetBool("minimum"))
	},
}, "minimum")

var getDimensionLayoutCmd = &cobra.Command{
	Use:     "layout <dimension-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Evaluate the layout of an generic dimension",
	Long:    "Evaluate the layout of an generic dimension",
	Example: "corectl dimension layout DIMENSION-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityLayout(state, args[0], "dimension")
	},
}

var dimensionCmd = &cobra.Command{
	Use:   "dimension",
	Short: "Explore and manage dimensions",
	Long:  "Explore and manage dimensions",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	dimensionCmd.AddCommand(listDimensionsCmd, setDimensionsCmd, getDimensionPropertiesCmd, getDimensionLayoutCmd, removeDimensionCmd)
}
