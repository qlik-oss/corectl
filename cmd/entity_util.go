package cmd

import (
	"fmt"
	"os"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

func getEntityProperties(ccmd *cobra.Command, args []string, entityType string, minimum bool) {
	if len(args) < 1 {
		fmt.Println("Expected an " + entityType + " id to specify what " + entityType + " to use as a parameter")
		ccmd.Usage()
		os.Exit(1)
	}
	state := internal.PrepareEngineState(rootCtx, headers, false)
	printer.PrintGenericEntityProperties(state, args[0], entityType, minimum)
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
