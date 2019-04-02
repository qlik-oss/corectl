package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var headersMap = make(map[string]string)
var explicitConfigFile = ""
var version = ""
var headers http.Header
var rootCtx = context.Background()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Hidden:                 true,
	Use:                    "corectl",
	Short:                  "",
	Long:                   `corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
	DisableAutoGenTag:      true,
	BashCompletionFunction: bashCompletionFunc,

	Annotations: map[string]string{
		"x-qlik-stability": "experimental",
	},

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		// if help, version or generate-docs command, no prerun is needed.
		if strings.Contains(ccmd.Use, "help") || ccmd.Use == "generate-docs" || ccmd.Use == "generate-spec" || ccmd.Use == "version" {
			return
		}
		internal.QliVerbose = viper.GetBool("verbose")
		internal.LogTraffic = viper.GetBool("traffic")
		if explicitConfigFile != "" {
			viper.SetConfigFile(strings.TrimSpace(explicitConfigFile))
			if err := viper.ReadInConfig(); err == nil {
				internal.LogVerbose("Using config file: " + explicitConfigFile)
			} else {
				fmt.Println(err)
			}
		} else {
			viper.SetConfigName("corectl") // name of config file (without extension)
			viper.SetConfigType("yml")
			viper.AddConfigPath(".")
			if err := viper.ReadInConfig(); err == nil {
				internal.LogVerbose("Using config file in working directory")
			} else {
				internal.LogVerbose("No config file")
			}
		}

		if len(headersMap) == 0 {
			headersMap = viper.GetStringMapString("headers")
		}
		headers = make(http.Header, 1)
		for key, value := range headersMap {
			headers.Set(key, value)
		}
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ccmd.HelpFunc()(ccmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(mainVersion string) {
	version = mainVersion
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Common commands
	rootCmd.AddCommand(getTablesCmd)
	rootCmd.AddCommand(getFieldsCmd)
	rootCmd.AddCommand(getAssociationsCmd)
	rootCmd.AddCommand(getKeysCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(evalCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(setAllCmd)
	rootCmd.AddCommand(getFieldCmd)
	rootCmd.AddCommand(getMetaCmd)

	// Subcommands
	rootCmd.AddCommand(measureCmd)
	rootCmd.AddCommand(dimensionCmd)
	rootCmd.AddCommand(objectCmd)
	rootCmd.AddCommand(connectionCmd)
	rootCmd.AddCommand(scriptCmd)
	rootCmd.AddCommand(appCmd)

	// Other
	rootCmd.AddCommand(catwalkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(getStatusCmd)

	rootCmd.PersistentFlags().StringVarP(&explicitConfigFile, "config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolP("traffic", "t", false, "Log JSON websocket traffic to stdout")
	viper.BindPFlag("traffic", rootCmd.PersistentFlags().Lookup("traffic"))

	rootCmd.PersistentFlags().StringP("engine", "e", "localhost:9076", "URL to the Qlik Associative Engine")
	viper.BindPFlag("engine", rootCmd.PersistentFlags().Lookup("engine"))

	rootCmd.PersistentFlags().String("ttl", "30", "Qlik Associative Engine session time to live in seconds")
	viper.BindPFlag("ttl", rootCmd.PersistentFlags().Lookup("ttl"))

	rootCmd.PersistentFlags().Bool("suppress", false, "Suppress all confirmation dialogues")
	viper.BindPFlag("suppress", rootCmd.PersistentFlags().Lookup("suppress"))

	rootCmd.PersistentFlags().Bool("no-data", false, "Open app without data")
	viper.BindPFlag("no-data", rootCmd.PersistentFlags().Lookup("no-data"))

	rootCmd.PersistentFlags().Bool("no-save", false, "Do not save the app")
	viper.BindPFlag("no-save", rootCmd.PersistentFlags().Lookup("no-save"))

	rootCmd.PersistentFlags().Bool("bash", false, "Bash flag used to adapt output to bash completion format")
	rootCmd.PersistentFlags().MarkHidden("bash")
	viper.BindPFlag("bash", rootCmd.PersistentFlags().Lookup("bash"))

	//not binding to viper since binding a map does not seem to work.
	rootCmd.PersistentFlags().StringToStringVar(&headersMap, "headers", nil, "Http headers to use when connecting to Qlik Associative Engine")

	rootCmd.PersistentFlags().StringP("app", "a", "", "App name, if no app is specified a session app is used instead.")
	viper.BindPFlag("app", rootCmd.PersistentFlags().Lookup("app"))
	// Set annotation to run bash completion function
	rootCmd.PersistentFlags().SetAnnotation("app", cobra.BashCompCustom, []string{"__corectl_get_apps"})

	for _, command := range []*cobra.Command{buildCmd, reloadCmd} {
		command.PersistentFlags().Bool("silent", false, "Do not log reload progress")
	}

	catwalkCmd.PersistentFlags().String("catwalk-url", "https://catwalk.core.qlik.com", "Url to an instance of catwalk, if not provided the qlik one will be used.")

	if runtime.GOOS != "windows" {
		// Do not add bash completion annotations for paths and files as they are not compatible with windows. On windows
		// we instead rely on the default bash behavior
		addFileRelatedBashAnnotations()
	}

	patchRootCommandUsageTemplate()
}

// GetRelativeParameter returns a parameter from the config file.
// It modifies the parameter to actually be relative to the config file and not the working directory
func GetRelativeParameter(paramName string) string {
	pathInConfigFile := viper.GetString(paramName)
	if pathInConfigFile != "" {
		return internal.RelativeToProject(viper.ConfigFileUsed(), pathInConfigFile)
	}
	return ""

}

func patchRootCommandUsageTemplate() {
	var originalUsageSnippet = `Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var rootSnippetMainSection = `Basic Commands:{{range .Commands}}{{if (and .IsAvailableCommand (eq (index .Annotations "command_category") ""))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Advanced Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (eq (index .Annotations "command_category") "sub"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Other Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (or (eq (index .Annotations "command_category") "other") (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var newUsageSnippet = `{{if (eq .Name "corectl")}}` + rootSnippetMainSection + `{{else}}` + originalUsageSnippet + "{{end}}"

	var patchedUsageTemplate = strings.Replace(rootCmd.UsageTemplate(), originalUsageSnippet, newUsageSnippet, 1)
	rootCmd.SetUsageTemplate(patchedUsageTemplate)
}
