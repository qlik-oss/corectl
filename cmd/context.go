package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setContextCmd = withLocalFlags(&cobra.Command{
	Use:   "set <context name>",
	Args:  cobra.ExactArgs(1),
	Short: "Set a context to the current configuration",
	Long: `Set a context to the current configuration

This command creates or updates a context by using the supplied flags and any
relevant config information found in the config file (if any).
The information stored will be engine url, headers and certificates (if present)
along with comment and the context-name.`,

	Example: `corectl context set local-engine
corectl context set rd-sense --engine localhost:9076 --comment "R&D Qlik Sense deployment"`,

	Run: func(ccmd *cobra.Command, args []string) {
		name := internal.SetContext(args[0], viper.GetString("comment"))
		printer.PrintCurrentContext(name)
	},
}, "comment")

var removeContextCmd = &cobra.Command{
	Use:   "rm <context name>...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Remove one or more contexts",
	Long:  "Remove one or more contexts",
	Example: `corectl context rm local-engine
corectl context rm ctx1 ctx2`,

	Run: func(ccmd *cobra.Command, args []string) {
		var removedCurrent bool
		for _, arg := range args {
			_, wasCurrent := internal.RemoveContext(arg)
			if wasCurrent {
				removedCurrent = true
			}
		}
		if removedCurrent {
			printer.PrintCurrentContext("")
		}
	},
}

var getContextCmd = &cobra.Command{
	Use:   "get [context name]",
	Args:  cobra.RangeArgs(0, 1),
	Short: "Get context, current context by default",
	Long:  "Get context, current context by default",
	Example: `corectl context get
corectl context get local-engine`,

	Run: func(ccmd *cobra.Command, args []string) {
		handler := internal.NewContextHandler()
		var name string

		if len(args) == 1 {
			name = args[0]
		}
		printer.PrintContext(name, handler)
	},
}

var listContextsCmd = &cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "List all contexts",
	Long:    "List all contexts",
	Example: "corectl context ls",

	Run: func(ccmd *cobra.Command, args []string) {
		handler := internal.NewContextHandler()
		printer.PrintContexts(handler, viper.GetBool("bash"))
	},
}

var useContextCmd = &cobra.Command{
	Use:     "use <context-name>",
	Args:    cobra.ExactArgs(1),
	Short:   "Specify what context to use",
	Long:    "Specify what context to use",
	Example: "corectl context use local-engine",

	Run: func(ccmd *cobra.Command, args []string) {
		name := internal.UseContext(args[0])
		printer.PrintCurrentContext(name)
	},
}

var clearContextCmd = &cobra.Command{
	Use:     "clear",
	Args:    cobra.ExactArgs(0),
	Short:   "Set the current context to none",
	Long:    "Set the current context to none",
	Example: "corectl context clear",

	Run: func(ccmd *cobra.Command, args []string) {
		previous := internal.ClearContext()
		if previous != "" {
			printer.PrintCurrentContext("")
		}
	},
}

var loginContextCmd = withLocalFlags(&cobra.Command{
	Use:   "login <context-name>",
	Args:  cobra.RangeArgs(0, 1),
	Short: "Login and set cookie for the named context",
	Long: `Login and set cookie for the named context
	
This is only applicable when connecting to 'Qlik Sense Enterprise for Windows' through its proxy using HTTPS.
If no 'context-name' is used as argument the 'current-context' defined in the config will be used instead.`,
	Example: `corectl context login
corectl context login context-name`,

	Run: func(ccmd *cobra.Command, args []string) {
		contextName := ""

		if len(args) > 0 {
			contextName = args[0]
		}

		internal.LoginContext(tlsClientConfig, contextName)
	},
}, "user", "password")

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Create, update and use contexts",
	Long: `Create, update and use contexts

Contexts store connection information such as engine url, certificates and headers,
similar to a config. The main difference between contexts and configs is that they
can be used globally. Use the context subcommands to configure contexts which
facilitate app development in environments where certificates and headers are needed.

The current context is the one that is being used. You can use "context get" to
display the contents of the current context and switch context with "context set"
or unset the current context with "context unset".

Note that contexts have the lowest precedence. This means that e.g. an --engine flag
(or an engine field in a config) will override the engine url in the current context.

Contexts are stored locally in your ~/.corectl/contexts.yml file.`,
	Annotations: map[string]string{
		"command_category": "other",
		"x-qlik-stability": "experimental",
	},
}

func init() {
	contextCmd.AddCommand(setContextCmd, removeContextCmd, listContextsCmd, useContextCmd, getContextCmd, clearContextCmd, loginContextCmd)
}
