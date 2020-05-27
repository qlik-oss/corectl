package engine

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

func CreateDimensionCommand() *cobra.Command {
	var setDimensionsCmd = WithLocalFlags(&cobra.Command{
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
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(true)
			internal.SetDimensions(ctx, doc, dynconf.Glob(commandLineDimensions))
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}, "no-save")

	var removeDimensionCmd = WithLocalFlags(&cobra.Command{
		Use:               "rm <dimension-id>...",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: listValidDimensionsForCompletion,
		Short:             "Remove one or many dimensions in the current app",
		Long:              "Remove one or many dimensions in the current app",
		Example:           "corectl dimension rm ID-1",
		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			for _, entity := range args {
				destroyed, err := doc.DestroyDimension(ctx, entity)
				if err != nil {
					log.Fatalf("could not remove generic dimension '%s': %s\n", entity, err)
				} else if !destroyed {
					log.Fatalf("could not remove generic dimension '%s'\n", entity)
				}
			}
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}, "no-save")

	var listDimensionsCmd = WithLocalFlags(&cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(0),
		Short:   "Print a list of all generic dimensions in the current app",
		Long:    "Print a list of all generic dimensions in the current app",
		Example: "corectl dimension ls",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			items := internal.ListDimensions(ctx, doc)
			printer.PrintNamedItemsList(items, params.PrintMode(), false)
		},
	}, "quiet")

	var getDimensionPropertiesCmd = WithLocalFlags(&cobra.Command{
		Use:               "properties <dimension-id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: listValidDimensionsForCompletion,
		Short:             "Print the properties of the generic dimension",
		Long:              "Print the properties of the generic dimension",
		Example:           "corectl dimension properties DIMENSION-ID",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			printer.PrintGenericEntityProperties(ctx, doc, args[0], "dimension", params.GetBool("minimum"), false)
		},
	}, "minimum")

	var getDimensionLayoutCmd = &cobra.Command{
		Use:               "layout <dimension-id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: listValidDimensionsForCompletion,
		Short:             "Evaluate the layout of an generic dimension",
		Long:              "Evaluate the layout of an generic dimension",
		Example:           "corectl dimension layout DIMENSION-ID",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			printer.PrintGenericEntityLayout(ctx, doc, args[0], "dimension")
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

	dimensionCmd.AddCommand(listDimensionsCmd, setDimensionsCmd, getDimensionPropertiesCmd, getDimensionLayoutCmd, removeDimensionCmd)
	return dimensionCmd
}

func listValidDimensionsForCompletion(ccmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
	items := internal.ListDimensions(ctx, doc)
	result := make([]string, 0)
	for _, item := range items {
		result = append(result, item.ID)
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}
