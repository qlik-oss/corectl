package engine

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

func CreateVariableCommand() *cobra.Command {
	var setVariablesCmd = WithLocalFlags(&cobra.Command{
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
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(true)
			internal.SetVariables(ctx, doc, dynconf.Glob(commandLineVariables))
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}, "no-save")

	var removeVariableCmd = WithLocalFlags(&cobra.Command{
		Use:     "rm <variable-name>...",
		Args:    cobra.MinimumNArgs(1),
		Short:   "Remove one or many variables in the current app",
		Long:    "Remove one or many variables in the current app",
		Example: "corectl variable rm NAME-1",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			for _, entity := range args {
				destroyed, err := doc.DestroyVariableByName(ctx, entity)
				if err != nil {
					log.Fatalf("could not remove generic variable '%s': %s\n", entity, err)
				} else if !destroyed {
					log.Fatalf("could not remove generic variable '%s'\n", entity)
				}
			}
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}, "no-save")

	var listVariablesCmd = WithLocalFlags(&cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(0),
		Short:   "Print a list of all generic variables in the current app",
		Long:    "Print a list of all generic variables in the current app",
		Example: "corectl variable ls",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			items := internal.ListVariables(ctx, doc)
			printer.PrintNamedItemsList(items, params.PrintMode(), !params.GetBool("quiet"))
		},
	}, "quiet")

	var getVariablePropertiesCmd = WithLocalFlags(&cobra.Command{
		Use:     "properties <variable-name>",
		Args:    cobra.ExactArgs(1),
		Short:   "Print the properties of the generic variable",
		Long:    "Print the properties of the generic variable",
		Example: "corectl variable properties VARIABLE-NAME",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			printer.PrintGenericEntityProperties(ctx, doc, args[0], "variable", params.GetBool("minimum"), false)
		},
	}, "minimum")

	var getVariableLayoutCmd = &cobra.Command{
		Use:     "layout <variable-id>",
		Args:    cobra.ExactArgs(1),
		Short:   "Evaluate the layout of an generic variable",
		Long:    "Evaluate the layout of an generic variable",
		Example: "corectl variable layout VARIABLE-NAME",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			printer.PrintGenericEntityLayout(ctx, doc, args[0], "variable")
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

	variableCmd.AddCommand(setVariablesCmd, removeVariableCmd, listVariablesCmd, getVariablePropertiesCmd, getVariableLayoutCmd)
	return variableCmd
}
