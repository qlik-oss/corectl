package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var setMeasuresCmd = withLocalFlags(&cobra.Command{
	Use:   "set <glob-pattern-path-to-measures-files.json>",
	Short: "Set or update the measures in the current app",
	Long:  "Set or update the measures in the current app",
	Example: `corectl measure set
corectl measure set ./my-measures-glob-path.json`,

	Run: func(ccmd *cobra.Command, args []string) {

		commandLineMeasures := ""
		if len(args) > 0 {
			commandLineMeasures = args[0]
		}
		state := internal.PrepareEngineState(rootCtx, headers, true)
		internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineMeasures, "measure")
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var removeMeasureCmd = withLocalFlags(&cobra.Command{
	Use:     "remove <measure-id>...",
	Short:   "Remove one or many generic measures in the current app",
	Long:    "Remove one or many generic measures in the current app",
	Example: `corectl measure remove ID-1 ID-2`,

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected atleast one measure-id specify what measure to remove from the app")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyMeasure(rootCtx, entity)
			if err != nil {
				internal.FatalError("Failed to remove generic measure ", entity+" with error: "+err.Error())
			} else if !destroyed {
				internal.FatalError("Failed to remove generic measure ", entity)
			}
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var listMeasuresCmd = &cobra.Command{
	Use:     "ls",
	Short:   "Print a list of all generic measures in the current app",
	Long:    "Print a list of all generic measures in the current app",
	Example: `corectl measure ls`,

	Run: func(ccmd *cobra.Command, args []string) {
		listEntities(ccmd, args, "measure", !viper.GetBool("bash"))
	},
}

var getMeasurePropertiesCmd = &cobra.Command{
	Use:     "properties <measure-id>",
	Short:   "Print the properties of the generic measure",
	Long:    "Print the properties of the generic measure",
	Example: `corectl measure properties MEASURE-ID`,

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "measure")
	},
}

var getMeasureLayoutCmd = &cobra.Command{
	Use:     "layout <measure-id>",
	Short:   "Evaluate the layout of an generic measure",
	Long:    "Evaluate the layout of an generic measure and prints in JSON format",
	Example: `corectl measure layout MEASURE-ID`,

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityLayout(ccmd, args, "measure")
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
