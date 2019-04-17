package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-latest"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Args:    cobra.ExactArgs(0),
	Short:   "Print the version of corectl",
	Example: "corectl version",
	Annotations: map[string]string{
		"command_category": "other",
	},

	Run: func(_ *cobra.Command, args []string) {

		if version != "development build" {
			checkLatestVersion()
		}

		fmt.Printf("corectl version: %s\n", version)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Args:  cobra.ExactArgs(0),
	Short: "Print status info about the connection to the engine and current app",
	Long:  "Print status info about the connection to the engine and current app, and also the status of the data model",
	Example: `corectl status
corectl status --app=my-app.qvf`,
	Annotations: map[string]string{
		"command_category": "other",
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		printer.PrintStatus(state, viper.GetString("engine"))
	},
}

// Function for checking current version against latest released version on github
func checkLatestVersion() {
	githubTag := &latest.GithubTag{
		Owner:      "qlik-oss",
		Repository: "corectl",
	}

	res, err := latest.Check(githubTag, version)

	if err == nil && res.Outdated {

		// Find absolute path of executable
		executable, _ := os.Executable()

		// Format a download string depending on OS
		var dwnl string
		if runtime.GOOS == "windows" {
			dwnl = fmt.Sprintf(`curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/v%s/corectl-windows-x86_64.zip" > corectl.zip && unzip ./corectl.zip -d "%s" && rm ./corectl.zip`, res.Current, path.Dir(executable))
		} else {
			dwnl = fmt.Sprintf(`curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/v%s/corectl-%s-x86_64.tar.gz" | tar xz -C /tmp && mv /tmp/corectl %s`, res.Current, runtime.GOOS, path.Dir(executable))
		}

		fmt.Println("-------------------------------------------------")
		fmt.Printf("There is a new version available! Please upgrade for the latest features and bug fixes. You are on %s, latest version is %s. \n", version, res.Current)
		fmt.Printf("To download the latest version you can use this command: \n")
		fmt.Printf(`'%s'`, dwnl)
		fmt.Println("\n-------------------------------------------------")
	}
}

func askForConfirmation(s string) bool {
	if viper.GetString("suppress") == "true" {
		return true
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
