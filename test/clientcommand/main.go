package main

import (
	"fmt"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/commands/standard"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/spf13/cobra"
	"strings"
)

func createCommandTree() *cobra.Command {

	var subCommand = &cobra.Command{
		Use:   "raw",
		Short: "command used to invoke Qlik Sense for kubernetes using rest",
		Long:  "command used to invoke Qlik Sense for kubernetes using rest",

		Run: func(ccmd *cobra.Command, args []string) {
			comm := boot.NewCommunicator(ccmd)
			rest := comm.RestCaller()
			var result []byte
			err := rest.CallStd("GET", comm.GetString("url"), comm.GetStringMap("query"), nil, &result)
			if err != nil {
				fmt.Println(err)
			}
			log.PrintAsJSON(result)
		},
	}
	subCommand.Flags().StringP("url", "u", "", "v1/tenants/me")
	subCommand.Flags().StringToStringP("query", "q", nil, "\"x=firstvalue,y=secondvalue\"")

	var rootCommand = &cobra.Command{
		Use:   "clientcommand",
		Args:  cobra.ExactArgs(0),
		Short: "test",
		Long:  "test test ",

		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}
	rootCommand.AddCommand(subCommand)

	rootCommand.AddCommand(standard.CompletionCommand())
	rootCommand.AddCommand(standard.GenerateSpecCommand())
	rootCommand.AddCommand(standard.GenerateDocsCommand())
	rootCommand.AddCommand(standard.ContextCommand())
	rootCommand.AddCommand(standard.StatusCommand())

	boot.InjectGlobalFlags(rootCommand)

	patchRootCommandUsageTemplate(rootCommand)

	return rootCommand
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

var rootCommand = createCommandTree()

func main() {
	rootCommand.Execute()
}
