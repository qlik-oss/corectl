package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var getAppsCmd = &cobra.Command{
	Use:     "ls",
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
	Use:   "remove <app-id>",
	Short: "Remove the specified app",
	Long:  `Remove the specified app`,
	Example: `corectl app remove
corectl app remove APP-ID`,

	Run: func(ccmd *cobra.Command, args []string) {
		app := viper.GetString("app")

		if len(args) != 1 && app == "" {
			fmt.Println("Expected an identifier of the app to delete.")
			ccmd.Usage()
			os.Exit(1)
		}

		if len(args) == 1 {
			app = args[0]
		}

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
