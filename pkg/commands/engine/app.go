package engine

import (
	"bufio"
	"fmt"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/rest"
	"os"
	"strings"

	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

type ListAppResponse struct {
	Data []ListItem `json:"data"`
}

type ListItem struct {
	Name string `json:"name"`
	Id   string `json:"resourceID"`
}

func CreateListAppsCmd() *cobra.Command {
	return withLocalFlags(&cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(0),
		Short:   "Print a list of all apps available in the current engine",
		Long:    "Print a list of all apps available in the current engine",
		Example: "corectl app ls",

		Run: func(ccmd *cobra.Command, args []string) {

			comm := boot.NewCommunicator(ccmd)
			if comm.IsSenseForKubernetes() {
				docList, err := comm.RestCaller().ListApps()
				if err != nil {
					log.Fatalf("could not retrieve app list: %s\n", err)
				}
				rest.PrintApps(docList, comm.PrintMode())
			} else {
				ctx, global, params := comm.OpenGlobalSocket()
				docList, err := global.GetDocList(ctx)
				if err != nil {
					log.Fatalf("could not retrieve app list: %s\n", err)
				}
				printer.PrintApps(docList, params.PrintMode())
			}
		},
	}, "quiet")
}
func CreateRemoveAppCmd() *cobra.Command {
	return withLocalFlags(&cobra.Command{
		Use:     "rm <app-id>",
		Args:    cobra.ExactArgs(1),
		Short:   "Remove the specified app",
		Long:    "Remove the specified app",
		Example: "corectl app rm APP-ID",

		Run: func(ccmd *cobra.Command, args []string) {
			app := args[0]
			comm := boot.NewCommunicator(ccmd)
			comm.OverrideSetting("app", app)
			if ok, err := comm.AppExists(); !ok {
				log.Fatalln(err)
			}

			confirmed := comm.GetBool("suppress") || askForConfirmation(fmt.Sprintf("Do you really want to delete the app: %s?", app))

			if confirmed {
				comm.DeleteApp(app)
			}
		},
	}, "suppress")
}
func CreateImportAppCmd() *cobra.Command {
	return withLocalFlags(&cobra.Command{
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
			comm := boot.NewCommunicator(ccmd)
			rest := comm.RestCaller()

			appID, appName, err := rest.ImportApp(appPath)
			if err != nil {
				log.Fatalln(err)
			}
			boot.SetAppIDToKnownApps(comm.AppIdMappingNamespace(), appName, appID, false)
			log.Info("Imported app with new ID: ")
			log.Quiet(appID)
		},
	}, "quiet")
}
func CreateAppCmd() *cobra.Command {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Explore and manage apps",
		Long:  "Explore and manage apps",
		Annotations: map[string]string{
			"command_category": "sub",
		},
	}
	appCmd.AddCommand(CreateListAppsCmd(), CreateRemoveAppCmd(), CreateImportAppCmd())
	return appCmd
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", s)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
