package cmd

import (
	"fmt"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getAppsCmd = &cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all apps available in the current engine",
	Long:    "Print a list of all apps available in the current engine",
	Example: `corectl app ls`,

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineStateWithoutApp(rootCtx, headers)
		docList, err := state.Global.GetDocList(rootCtx)
		if err != nil {
			internal.FatalError(err)
		}
		printer.PrintApps(docList, !viper.GetBool("bash"), viper.GetBool("bash"))
	},
}

var removeAppCmd = withLocalFlags(&cobra.Command{
	Use:     "rm <app-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Remove the specified app",
	Long:    `Remove the specified app`,
	Example: `corectl app rm APP-ID`,

	Run: func(ccmd *cobra.Command, args []string) {
		app := viper.GetString("app")
		app = args[0]

		confirmed := askForConfirmation(fmt.Sprintf("Do you really want to delete the app: %s?", app))

		if confirmed {
			internal.DeleteApp(rootCtx, viper.GetString("engine"), app, viper.GetString("ttl"), headers)
		}
	},
}, "suppress")

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Explore and manage apps",
	Long:  "Explore and manage apps",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	appCmd.AddCommand(getAppsCmd, removeAppCmd)
}
