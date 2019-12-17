package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var headersMap = make(map[string]string)
var explicitConfigFile = ""
var explicitCertificatePath = ""
var version = ""
var commit = ""
var branch = ""
var headers http.Header
var tlsClientConfig *tls.Config
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
		"x-qlik-stability": "stable",
	},

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		// For some commands we don't want to do a prerun.
		if skipPreRun(ccmd) {
			return
		}
		// Depending on the command, we might not want to use context when loading config.
		withContext := shouldUseContext(ccmd)
		internal.ReadConfig(explicitConfigFile, explicitCertificatePath, withContext)

		tlsClientConfig = &tls.Config{}

		if certPath := viper.GetString("certificates"); certPath != "" {
			tlsClientConfig = internal.ReadCertificates(tlsClientConfig, certPath)
		}

		if viper.GetBool("insecure") {
			tlsClientConfig.InsecureSkipVerify = true
		}

		if len(headersMap) == 0 {
			headersMap = viper.GetStringMapString("headers")
		}
		headers = make(http.Header, 1)
		for key, value := range headersMap {
			headers.Set(key, value)
		}

		headers.Set("User-Agent", fmt.Sprintf("corectl/%s (%s)", version, runtime.GOOS))

		// Initiate the printers mode
		printer.Init()
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ccmd.HelpFunc()(ccmd, args)
	},
}

func skipPreRun(ccmd *cobra.Command) bool {
	// Depending on what command we are using, we might not want to do the prerun,
	// i.e. reading config and such.
	// Note: 'path' is the complete path for the command for example: 'corectl app build'.
	// Using path instead of 'ccmd.Use' since it only gives you the last word of the command.
	path := ccmd.CommandPath()
	switch {
	case strings.Contains(path, "help"):
		return true
	case strings.Contains(path, "generate-docs"):
		return true
	case strings.Contains(path, "generate-spec"):
		return true
	case strings.Contains(path, "version"):
		return true
	// For contexts we only want to do a prerun for context set.
	case strings.Contains(path, "context"):
		if strings.Contains(path, "context set") || strings.Contains(path, "context login") {
			return false
		}
		return true
	}
	return false
}

func shouldUseContext(ccmd *cobra.Command) bool {
	path := ccmd.CommandPath()
	// Switch so cases can be easily added.
	// Note that this is only import for commands that
	// actually do the complete prerun. That is: not the
	// ones included in the skipPreRun function.
	switch {
	case strings.Contains(path, "context set"):
		return false
	}
	return true
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(mainVersion, branchName, commitSha string) {
	version = mainVersion
	branch = branchName
	commit = commitSha
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints an error message so we just want to exit
		os.Exit(1)
	}
}

func patchRootCommandUsageTemplate() {
	var originalUsageSnippet = `Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var rootSnippetMainSection = `App Building Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (eq (index .Annotations "command_category") "build"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

App Analysis Commands:{{range .Commands}}{{if (and .IsAvailableCommand (eq (index .Annotations "command_category") ""))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Advanced Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (eq (index .Annotations "command_category") "sub"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Other Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (or (eq (index .Annotations "command_category") "other") (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var newUsageSnippet = `{{if (eq .Name "corectl")}}` + rootSnippetMainSection + `{{else}}` + originalUsageSnippet + "{{end}}"

	var patchedUsageTemplate = strings.Replace(rootCmd.UsageTemplate(), originalUsageSnippet, newUsageSnippet, 1)
	rootCmd.SetUsageTemplate(patchedUsageTemplate)
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
	rootCmd.AddCommand(getValuesCmd)
	rootCmd.AddCommand(getMetaCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(unbuildCmd)

	// Subcommands
	rootCmd.AddCommand(alternateStateCmd)
	rootCmd.AddCommand(measureCmd)
	rootCmd.AddCommand(dimensionCmd)
	rootCmd.AddCommand(objectCmd)
	rootCmd.AddCommand(variableCmd)
	rootCmd.AddCommand(bookmarkCmd)
	rootCmd.AddCommand(connectionCmd)
	rootCmd.AddCommand(scriptCmd)
	rootCmd.AddCommand(appCmd)

	// Other
	rootCmd.AddCommand(catwalkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(generateDocsCmd)
	rootCmd.AddCommand(generateSpecCmd)

	initGlobalFlags(rootCmd.PersistentFlags())
	patchRootCommandUsageTemplate()

}
