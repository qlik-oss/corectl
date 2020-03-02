package cmd

import (
	"context"
	"github.com/qlik-oss/corectl/pkg/boot"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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

	Run: func(ccmd *cobra.Command, args []string) {
		ccmd.HelpFunc()(ccmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
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
	//App Building Commands

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(unbuildCmd)

	// Common commands
	rootCmd.AddCommand(getAssociationsCmd)
	rootCmd.AddCommand(catwalkCmd)
	rootCmd.AddCommand(evalCmd)
	rootCmd.AddCommand(getFieldsCmd)
	rootCmd.AddCommand(getValuesCmd)
	rootCmd.AddCommand(getMetaCmd)
	rootCmd.AddCommand(getKeysCmd)
	rootCmd.AddCommand(getTablesCmd)

	// Subcommands
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(bookmarkCmd)
	rootCmd.AddCommand(connectionCmd)
	rootCmd.AddCommand(dimensionCmd)
	rootCmd.AddCommand(measureCmd)
	rootCmd.AddCommand(objectCmd)
	rootCmd.AddCommand(scriptCmd)
	rootCmd.AddCommand(alternateStateCmd)
	rootCmd.AddCommand(variableCmd)

	// Other
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)

	// Hidden administrative commands
	rootCmd.AddCommand(generateDocsCmd)
	rootCmd.AddCommand(generateSpecCmd)

	boot.InjectGlobalFlags(rootCmd)

	//initGlobalFlags(rootCmd.PersistentFlags())
	patchRootCommandUsageTemplate()

}

func ContextCommand() *cobra.Command {
	return contextCmd
}
func CompletionCommand() *cobra.Command {
	return completionCmd
}
func StatusCommand() *cobra.Command {
	return statusCmd
}

func GenerateDocsCommand() *cobra.Command {
	return generateDocsCmd
}
func GenerateSpecCommand() *cobra.Command {
	return generateDocsCmd
}

func EditAppSubCommand() *cobra.Command {
	var wsSommand = &cobra.Command{
		Use:   "ws",
		Short: "web socket things",
		Long:  `ws contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
	}

	wsSommand.AddCommand(buildCmd)
	wsSommand.AddCommand(reloadCmd)
	wsSommand.AddCommand(unbuildCmd)

	// Common commands
	wsSommand.AddCommand(getAssociationsCmd)
	wsSommand.AddCommand(catwalkCmd)
	wsSommand.AddCommand(evalCmd)
	wsSommand.AddCommand(getFieldsCmd)
	wsSommand.AddCommand(getValuesCmd)
	wsSommand.AddCommand(getMetaCmd)
	wsSommand.AddCommand(getKeysCmd)
	wsSommand.AddCommand(getTablesCmd)

	// Subcommands
	wsSommand.AddCommand(appCmd)
	wsSommand.AddCommand(bookmarkCmd)
	wsSommand.AddCommand(connectionCmd)
	wsSommand.AddCommand(dimensionCmd)
	wsSommand.AddCommand(measureCmd)
	wsSommand.AddCommand(objectCmd)
	wsSommand.AddCommand(scriptCmd)
	wsSommand.AddCommand(alternateStateCmd)
	wsSommand.AddCommand(variableCmd)

	return wsSommand
}
