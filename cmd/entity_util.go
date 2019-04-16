package cmd

import (
	"fmt"
	"os"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getEntityProperties(ccmd *cobra.Command, args []string, entityType string) {
	if len(args) < 1 {
		fmt.Println("Expected an " + entityType + " id to specify what " + entityType + " to use as a parameter")
		ccmd.Usage()
		os.Exit(1)
	}
	state := internal.PrepareEngineState(rootCtx, headers, false)
	printer.PrintGenericEntityProperties(state, args[0], entityType)
}

func getEntityLayout(ccmd *cobra.Command, args []string, entityType string) {
	if len(args) < 1 {
		fmt.Println("Expected an " + entityType + " id to specify what " + entityType + " to use as a parameter")
		ccmd.Usage()
		os.Exit(1)
	}
	state := internal.PrepareEngineState(rootCtx, headers, false)
	printer.PrintGenericEntityLayout(state, args[0], entityType)
}

func listEntities(ccmd *cobra.Command, args []string, entityType string, printAsJSON bool) {
	state := internal.PrepareEngineState(rootCtx, headers, false)
	allInfos, err := state.Doc.GetAllInfos(rootCtx)
	if err != nil {
		internal.FatalError(err)
	}
	printer.PrintGenericEntities(allInfos, entityType, printAsJSON, viper.GetBool("bash"))
}
