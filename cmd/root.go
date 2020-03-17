package cmd

import (
	"context"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/commands/engine"
	"github.com/qlik-oss/corectl/pkg/commands/standard"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCtx = context.Background()

func CreateRootCommand(version, branch, commit string) *cobra.Command {
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

	//App Building Commands
	appBuildCommands := []*cobra.Command{
		engine.CreateBuildCmd(),
		engine.CreateReloadCmd(),
		engine.CreateUnbuildCmd(),
	}
	Annotate("command_category", "build", appBuildCommands...)
	rootCmd.AddCommand(appBuildCommands...)

	// Common commands
	commonCommands := []*cobra.Command{
		engine.CreategetAssociationsCmd(),
		engine.CreateCatwalkCmd(),
		engine.CreateEvalCmd(),
		engine.CreateGetFieldsCmd(),
		engine.CreateGetValuesCmd(),
		engine.CreateGetMetaCmd(),
		engine.CreateGetKeysCmd(),
		engine.CreategetTablesCmd(),
	}
	//Annotate("command_category", "common", getAssociationsCmd, catwalkCmd, evalCmd,
	//		getFieldsCmd, getValuesCmd, getMetaCmd, getKeysCmd, getTablesCmd)
	rootCmd.AddCommand(commonCommands...)

	// Subcommands
	subCommands := []*cobra.Command{
		engine.CreateAppCmd(),
		engine.CreateAppCmd(),
		engine.CreateBookmarkCommand(),
		engine.CreateConnectionCommand(),
		engine.CreateDimensionCommand(),
		engine.CreateMeasureCommand(),
		engine.CreateObjectCommand(),
		engine.CreateScriptCommand(),
		engine.CreateAlternateStateCommand(),
		engine.CreateVariableCommand(),
	}
	Annotate("command_category", "sub", subCommands...)
	rootCmd.AddCommand(subCommands...)

	// Other
	otherCommands := []*cobra.Command{
		standard.CreateCompletionCommand("corectl"),
		getContextCmd,
		standard.CreateStatusCmd(),
		standard.CreateVersionCmd(version, branch, commit),
	}
	Annotate("command_category", "other", otherCommands...)
	rootCmd.AddCommand(otherCommands...)

	// Hidden administrative commands
	rootCmd.AddCommand(standard.CreateGenarateDocsCommand())
	rootCmd.AddCommand(standard.CreateGenerateSpecCommand(version))

	boot.InjectGlobalFlags(rootCmd, false)

	//initGlobalFlags(rootCmd.PersistentFlags())
	patchRootCommandUsageTemplate(rootCmd)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, branch, commit string) {
	rootCmd := CreateRootCommand(version, branch, commit)
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints an error message so we just want to exit
		os.Exit(1)
	}
}

func patchRootCommandUsageTemplate(rootCmd *cobra.Command) {
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

func ContextCommand() *cobra.Command {
	return contextCmd
}

func Annotate(key, value string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		if cmd.Annotations == nil {
			cmd.Annotations = map[string]string{}
		}
		cmd.Annotations[key] = value
	}
}
