package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var setDimensionsCmd = &cobra.Command{
	Use:     "set <glob-pattern-path-to-dimensions-files.json>",
	Short:   "Sets or updates the dimensions in the current app",
	Long:    "Sets or updates the dimensions in the current app",
	Example: "corectl set dimensions ./my-dimensions-glob-path.json",

	Run: func(ccmd *cobra.Command, args []string) {

		commandLineDimensions := ""
		if len(args) > 0 {
			commandLineDimensions = args[0]
		}
		state := internal.PrepareEngineState(rootCtx, headers, true)
		internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineDimensions, "dimension")
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var removeDimensionCmd = &cobra.Command{
	Use:     "remove <dimension-id>...",
	Short:   "Remove one or many dimensions in the current app",
	Long:    "Remove one or many dimensions in the current app",
	Example: `corectl dimension remove ID-1`,

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected atleast one dimension-id specify what dimension to remove from the app")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyDimension(rootCtx, entity)
			if err != nil {
				internal.FatalError("Failed to remove generic dimension ", entity+" with error: "+err.Error())
			} else if !destroyed {
				internal.FatalError("Failed to remove generic dimension ", entity)
			}
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var getDimensionsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "Prints a list of all generic dimensions in the current app",
	Long:    "Prints a list of all generic dimensions in the current app",
	Example: "corectl dimension list",

	Run: func(ccmd *cobra.Command, args []string) {
		getEntities(ccmd, args, "dimension", !viper.GetBool("bash"))
	},
}

var getDimensionPropertiesCmd = &cobra.Command{
	Use:     "properties <dimension-id>",
	Short:   "Prints the properties of the generic dimension",
	Long:    "Prints the properties of the generic dimension",
	Example: "corectl dimension properties DIMENSION-ID --app my-app.qvf",

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "dimension")
	},
}

var getDimensionLayoutCmd = &cobra.Command{
	Use:     "layout <dimension-id>",
	Short:   "Evaluates the layout of an generic dimension",
	Long:    "Evaluates the layout of an generic dimension",
	Example: "corectl dimension layout DIMENSION-ID --app my-app.qvf",

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityLayout(ccmd, args, "dimension")
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
	dimensionCmd.AddCommand(setDimensionsCmd)
	dimensionCmd.AddCommand(getDimensionsCmd)
	dimensionCmd.AddCommand(getDimensionLayoutCmd)
	dimensionCmd.AddCommand(getDimensionPropertiesCmd)
	dimensionCmd.AddCommand(removeDimensionCmd)
}
