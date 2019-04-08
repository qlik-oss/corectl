package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var setObjectsCmd = withLocalFlags(&cobra.Command{
	Use:   "set <glob-pattern-path-to-objects-files.json",
	Short: "Sets or updates the objects in the current app",
	Long: `Sets or updates the objects in the current app.
The JSON objects can be in either the GenericObjectProperties format or the GenericObjectEntry format`,
	Example: "corectl object set ./my-objects-glob-path.json",

	Run: func(ccmd *cobra.Command, args []string) {

		commandLineObjects := ""
		if len(args) > 0 {
			commandLineObjects = args[0]
		}

		state := internal.PrepareEngineState(rootCtx, headers, true)
		internal.SetupEntities(rootCtx, state.Doc, viper.ConfigFileUsed(), commandLineObjects, "object")
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var removeObjectCmd = withLocalFlags(&cobra.Command{
	Use:     "remove <object-id>...",
	Short:   "Remove one or many generic objects in the current app",
	Long:    "Remove one or many generic objects in the current app",
	Example: `corectl object remove ID-1 ID-2`,

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected atleast one object-id specify what object to remove from the app")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyObject(rootCtx, entity)
			if err != nil {
				internal.FatalError("Failed to remove generic object ", entity+" with error: "+err.Error())
			} else if !destroyed {
				internal.FatalError("Failed to remove generic object ", entity)
			}
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var listObjectsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Prints a list of all generic objects in the current app",
	Long:  "Prints a list of all generic objects in the current app in either plain text or JSON format",
	Example: `corectl object list
corectl get objects --json --app=myapp.qvf`,

	Run: func(ccmd *cobra.Command, args []string) {
		listEntities(ccmd, args, "object", !viper.GetBool("bash"))
	},
}

var getObjectPropertiesCmd = &cobra.Command{
	Use:     "properties <object-id>",
	Short:   "Prints the properties of the generic object",
	Long:    "Prints the properties of the generic object in JSON format",
	Example: "corectl object properties OBJECT-ID --app my-app.qvf",

	Run: func(ccmd *cobra.Command, args []string) {
		getEntityProperties(ccmd, args, "object")
	},
}

var getObjectLayoutCmd = &cobra.Command{
	Use:     "layout <object-id>",
	Short:   "Evaluates the hypercube layout of an generic object",
	Long:    "Evaluates the hypercube layout of an generic object in JSON format",
	Example: "corectl object layout OBJECT-ID --app my-app.qvf",

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected an object id specify what object to use as a parameter")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		printer.PrintGenericEntityLayout(state, args[0], "object")
	},
}

var getObjectDataCmd = &cobra.Command{
	Use:     "data <object-id>",
	Short:   "Evaluates the hypercube data of an generic object",
	Long:    "Evaluates the hypercube data of an generic object",
	Example: "corectl object data OBJECT-ID --app my-app.qvf",

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected an object id specify what object to use as a parameter")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		printer.EvalObject(rootCtx, state.Doc, args[0])
	},
}

var objectCmd = &cobra.Command{
	Use:   "object",
	Short: "Explore and manage generic objects",
	Long:  "Explore and manage generic objects",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	objectCmd.AddCommand(setObjectsCmd, listObjectsCmd, getObjectDataCmd, getObjectLayoutCmd, getObjectPropertiesCmd, removeObjectCmd)
}
