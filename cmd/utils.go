package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/google/go-github/v26/github"
	ver "github.com/hashicorp/go-version"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		} else {
			fmt.Printf("corectl version: %s\n", version)
		}
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
		appName := viper.GetString("app")
		var state *internal.State
		if appName != "" {
			state = internal.PrepareEngineState(rootCtx, headers, false)
		} else {
			state = internal.PrepareEngineStateWithoutApp(rootCtx, headers)
		}
		printer.PrintStatus(state, viper.GetString("engine"))
	},
}

// Function for checking current version against latest released version on github
func checkLatestVersion() {
	client := github.NewClient(nil)
	rel, _, err := client.Repositories.GetLatestRelease(context.Background(), "qlik-oss", "corectl")
	if err != nil {
		// If we cannot connect to github just print the version
		fmt.Printf("corectl version: %s\n", version)
		return
	}
	latestVersion, outdated := isLatestVersion(version, *rel.TagName)
	if outdated {
		// Find absolute path of executable
		executable, _ := os.Executable()
		fmt.Println("--------------------------------------------------")
		fmt.Printf("corectl version: %s, latest version is %s\n", version, latestVersion)
		switch runtime.GOOS {
		case "darwin":
			fmt.Println("To update to the latest version using brew just run:")
			fmt.Print("\n  brew upgrade qlik-corectl\n\n")
			fmt.Println("If you don't already have the qlik-oss tap you have to first add the tap with:")
			fmt.Print("\n  brew tap qlik-oss/taps\n\n")
			fmt.Println("And after that, you have to install it with brew by running:")
			fmt.Print("\n  brew install qlik-corectl\n\n")
			fmt.Println("If you prefer curl, you can run:")
		case "linux":
			fmt.Println("To update to the latest version using snap just run:")
			fmt.Print("\n  snap refresh qlik-corectl\n\n")
			fmt.Println("If you prefer curl, you can run:")
		default:
			fmt.Println("To download the latest version using curl you can run:")
		}
		// Format a download string depending on OS
		var dwnl string
		if runtime.GOOS == "windows" {
			dwnl = fmt.Sprintf(`curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/v%s/corectl-windows-x86_64.zip" > corectl.zip && unzip ./corectl.zip -d "%s" && rm ./corectl.zip`, latestVersion, path.Dir(executable))
		} else {
			dwnl = fmt.Sprintf(`curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/v%s/corectl-%s-x86_64.tar.gz" | tar xz -C /tmp && mv /tmp/corectl %s`, latestVersion, runtime.GOOS, path.Dir(executable))
		}
		fmt.Printf("\n  %s\n\n", dwnl)
		fmt.Println("If you have any problems, questions or feedback you can find us on:")
		fmt.Println("  slack:\thttps://qlikbranch-slack-invite.herokuapp.com/")
		fmt.Println("  github:\thttps://github.com/qlik-oss/corectl")
		fmt.Println("---------------------------------------------------")
	} else {
		fmt.Printf("corectl version: %s\n", version)
	}
}

func isLatestVersion(currentTag string, latestTag string) (string, bool) {
	currentVersion, err := ver.NewVersion(currentTag)
	if err != nil {
		internal.FatalErrorf("Current version is not semantically versioned: %s", currentTag)
	}

	latestVersion, err := ver.NewVersion(latestTag[1:]) // Remove 'v' from the tag
	if err != nil {
		internal.FatalErrorf("Latest version is not semantically versioned: %s", latestVersion)
	}

	if currentVersion.LessThan(latestVersion) {
		return latestVersion.String(), true
	}
	return "", false
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
