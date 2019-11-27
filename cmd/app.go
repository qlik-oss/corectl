package cmd

import (
	"fmt"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/internal/rest"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listAppsCmd = withLocalFlags(&cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all apps available in the current engine",
	Long:    "Print a list of all apps available in the current engine",
	Example: "corectl app ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, true)
		docList, err := state.Global.GetDocList(rootCtx)
		if err != nil {
			log.Fatalf("could not retrieve app list: %s\n", err)
		}
		printer.PrintApps(docList, viper.GetBool("bash"))
	},
}, "quiet")

var removeAppCmd = withLocalFlags(&cobra.Command{
	Use:     "rm <app-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Remove the specified app",
	Long:    "Remove the specified app",
	Example: "corectl app rm APP-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		app := args[0]

		if ok, err := internal.AppExists(rootCtx, viper.GetString("engine"), app, headers, tlsClientConfig); !ok {
			log.Fatalln(err)
		}
		confirmed := askForConfirmation(fmt.Sprintf("Do you really want to delete the app: %s?", app))

		if confirmed {
			internal.DeleteApp(rootCtx, viper.GetString("engine"), app, headers, tlsClientConfig)
		}
	},
}, "suppress")

var importAppCmd = withLocalFlags(&cobra.Command{
	Use:     "import",
	Args:    cobra.ExactArgs(1),
	Short:   "Import the specified app into the engine, returns the ID of the created app",
	Long:    "Import the specified app into the engine, returns the ID of the created app",
	Example: "corectl import <path-to-app.qvf>",
	Annotations: map[string]string{
		"x-qlik-stability": "experimental",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		appPath := args[0]
		engine := internal.GetEngineURL()
		appID, appName, err := rest.ImportApp(appPath, engine, headers, tlsClientConfig)
		if err != nil {
			log.Fatalln(err)
		}
		internal.SetAppIDToKnownApps(appName, appID, false)
		log.Info("Imported app with new ID: ")
		log.Quiet(appID)
	},
}, "quiet")

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Explore and manage apps",
	Long:  "Explore and manage apps",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	appCmd.AddCommand(listAppsCmd, removeAppCmd, importAppCmd)
}
