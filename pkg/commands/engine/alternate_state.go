package engine

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

func CreateAlternateStateCommand() *cobra.Command {
	var listAlternateStatesCmd = WithLocalFlags(&cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(0),
		Short:   "Print a list of all alternate states in the current app",
		Long:    "Print a list of all alternate states in the current app",
		Example: "corectl state ls",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			items := internal.ListAlternateStates(ctx, doc)
			printer.PrintStates(items, params.PrintMode())
		},
	}, "quiet") //TODO is quiet really supported?

	var addAlternateStateCmd = &cobra.Command{
		Use:     "add <alternate-state-name>",
		Args:    cobra.ExactArgs(1),
		Short:   "Add an alternate states in the current app",
		Long:    "Add an alternate states in the current app",
		Example: "corectl state add NAME-1",

		Run: func(ccmd *cobra.Command, args []string) {
			stateName := args[0]
			if stateName == "" {
				log.Fatalln("no state name specified")
			}
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			internal.AddAlternateState(ctx, doc, stateName)
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}

	var removeAlternateStateCmd = &cobra.Command{
		Use:               "rm <alternate-state-name>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: listValidAlternateStatesForCompletion,
		Short:             "Removes an alternate state in the current app",
		Long:              "Removes an alternate state in the current app",
		Example:           "corectl state rm NAME-1",

		Run: func(ccmd *cobra.Command, args []string) {
			stateName := args[0]
			if stateName == "" {
				log.Fatalln("no state name specified")
			}
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			internal.RemoveAlternateState(ctx, doc, stateName)
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}

	var alternateStateCmd = &cobra.Command{
		Use:     "state",
		Short:   "Explore and manage alternate states",
		Long:    "Explore and manage alternate states",
		Aliases: []string{"alternatestate"},
		Annotations: map[string]string{
			"command_category": "sub",
			"x-qlik-stability": "experimental",
		},
	}

	alternateStateCmd.AddCommand(listAlternateStatesCmd, addAlternateStateCmd, removeAlternateStateCmd)
	return alternateStateCmd
}

func listValidAlternateStatesForCompletion(ccmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
	items := internal.ListAlternateStates(ctx, doc)
	return items, cobra.ShellCompDirectiveNoFileComp
}
