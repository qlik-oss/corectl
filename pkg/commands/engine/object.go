package engine

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"strings"
)

func CreateObjectCommand() *cobra.Command {
	var setObjectsCmd = WithLocalFlags(&cobra.Command{
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
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(true)
			internal.SetObjects(ctx, doc, dynconf.Glob(commandLineObjects))
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}, "no-save")

	var removeObjectCmd = WithLocalFlags(&cobra.Command{
		Use:               "rm <object-id>...",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: listValidObjectsForCompletion,
		Short:             "Remove one or many generic objects in the current app",
		Long:              "Remove one or many generic objects in the current app",
		Example:           "corectl object rm ID-1 ID-2",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			for _, entity := range args {
				destroyed, err := doc.DestroyObject(ctx, entity)
				if err != nil {
					log.Fatalf("could not remove generic object '%s': %s\n", entity, err)
				} else if !destroyed {
					log.Fatalf("could not remove generic object '%s'\n", entity)
				}
			}
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}, "no-save")

	var listObjectsCmd = WithLocalFlags(&cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(0),
		Short:   "Print a list of all generic objects in the current app",
		Long:    "Print a list of all generic objects in the current app",
		Example: "corectl object ls",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			items := internal.ListObjects(ctx, doc)
			printer.PrintNamedItemsListWithType(items, params.PrintMode())
		},
	}, "quiet")

	var getObjectPropertiesCmd = WithLocalFlags(&cobra.Command{
		Use:               "properties <object-id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: listValidObjectsForCompletion,
		Short:             "Print the properties of the generic object",
		Long:              "Print the properties of the generic object in JSON format",
		Example:           "corectl object properties OBJECT-ID",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			printer.PrintGenericEntityProperties(ctx, doc, args[0], "object", params.GetBool("minimum"), params.GetBool("full"))
		},
	}, "minimum", "full")

	var getObjectLayoutCmd = &cobra.Command{
		Use:               "layout <object-id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: listValidObjectsForCompletion,
		Short:             "Evaluate the hypercube layout of the generic object",
		Long:              "Evaluate the hypercube layout of the generic object",
		Example:           "corectl object layout OBJECT-ID",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			printer.PrintGenericEntityLayout(ctx, doc, args[0], "object")
		},
	}

	var getObjectDataCmd = &cobra.Command{
		Use:               "data <object-id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: listValidObjectsForCompletion,
		Short:             "Evaluate the hypercube data of a generic object",
		Long:              "Evaluate the hypercube data of a generic object",
		Example:           "corectl object data OBJECT-ID",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			printer.EvalObject(ctx, doc, args[0])
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

	objectCmd.AddCommand(listObjectsCmd, setObjectsCmd, getObjectPropertiesCmd, getObjectLayoutCmd, getObjectDataCmd, removeObjectCmd)
	return objectCmd
}

func listValidObjectsForCompletion(ccmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
	items := internal.ListObjects(ctx, doc)
	result := make([]string, 0)
	for _, item := range items {
		if strings.HasPrefix(item.ID, toComplete) {
			result = append(result, item.ID)
		}
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}
