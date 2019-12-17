package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setObjectsCmd = withLocalFlags(&cobra.Command{
	Use:   "set <glob-pattern-path-to-objects-files.json",
	Args:  cobra.ExactArgs(1),
	Short: "Set or update the objects in the current app",
	Long: `Set or update the objects in the current app.
The JSON objects can be in either the GenericObjectProperties format or the GenericObjectEntry format`,
	Example: "corectl object set ./my-objects-glob-path.json",

	Run: func(ccmd *cobra.Command, args []string) {
		commandLineObjects := args[0]
		if commandLineObjects == "" {
			log.Fatalln("no objects specified")
		}
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, true, false)
		internal.SetObjects(rootCtx, state.Doc, commandLineObjects)
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var removeObjectCmd = withLocalFlags(&cobra.Command{
	Use:     "rm <object-id>...",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove one or many generic objects in the current app",
	Long:    "Remove one or many generic objects in the current app",
	Example: "corectl object rm ID-1 ID-2",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyObject(rootCtx, entity)
			if err != nil {
				log.Fatalf("could not remove generic object '%s': %s\n", entity, err)
			} else if !destroyed {
				log.Fatalf("could not remove generic object '%s'\n", entity)
			}
		}
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var listObjectsCmd = withLocalFlags(&cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all generic objects in the current app",
	Long:    "Print a list of all generic objects in the current app",
	Example: "corectl object ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		items := internal.ListObjects(state.Ctx, state.Doc)
		printer.PrintNamedItemsListWithType(items, viper.GetBool("bash"))
	},
}, "quiet")

var getObjectPropertiesCmd = withLocalFlags(&cobra.Command{
	Use:     "properties <object-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Print the properties of the generic object",
	Long:    "Print the properties of the generic object in JSON format",
	Example: "corectl object properties OBJECT-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityProperties(state, args[0], "object", viper.GetBool("minimum"))
	},
}, "minimum", "full")

var getObjectLayoutCmd = &cobra.Command{
	Use:     "layout <object-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Evaluate the hypercube layout of the generic object",
	Long:    "Evaluate the hypercube layout of the generic object",
	Example: "corectl object layout OBJECT-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityLayout(state, args[0], "object")
	},
}

var getObjectDataCmd = &cobra.Command{
	Use:     "data <object-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Evaluate the hypercube data of a generic object",
	Long:    "Evaluate the hypercube data of a generic object",
	Example: "corectl object data OBJECT-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
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
	objectCmd.AddCommand(listObjectsCmd, setObjectsCmd, getObjectPropertiesCmd, getObjectLayoutCmd, getObjectDataCmd, removeObjectCmd)
}
