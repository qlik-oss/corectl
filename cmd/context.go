package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createContextCmd = withLocalFlags(&cobra.Command{
	Use:   "create <context name>",
	Args:  cobra.ExactArgs(1),
	Short: "Create a new context",
	Long: `Create a new context

This command creates a new context using the supplied flags and any relevant
config information found in the config file (if any). The information stored
will be engine url, headers and certificates (if present) along with comment
and the context-name. The current context will become the newly created one.`,

	Example: `corectl context create local-engine
corectl context create rd-sense --engine localhost:9076 --comment "R&D Qlik Sense deployment"`,

	Run: func(ccmd *cobra.Command, args []string) {
		name := internal.CreateContext(args[0], viper.GetString("comment"))
		printer.PrintCurrentContext(name)
	},
}, "comment")

var removeContextCmd = &cobra.Command{
	Use:     "rm <context name>",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Removes a context",
	Long:    "Removes a context",
	Example: "corectl context rm local-engine",

	Run: func(ccmd *cobra.Command, args []string) {
		_, wasCurrent := internal.RemoveContext(args[0])
		if wasCurrent {
			printer.PrintCurrentContext("")
		}
	},
}

var getContextCmd = &cobra.Command{
	Use:   "get",
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
		} else {
			name = handler.Current
		}
		context := handler.Get(name)
		printer.PrintContext(name, context)
	},
}

var updateContextCmd = withLocalFlags(&cobra.Command{
	Use:   "update",
	Args:  cobra.RangeArgs(0, 1),
	Short: "Update context, current context by default",
	Long:  "Update context, current context by default",
	Example: `corectl context update
corectl context update local-engine`,

	Run: func(ccmd *cobra.Command, args []string) {
		var name string
		if len(args) == 1 {
			name = args[0]
		}
		internal.UpdateContext(name, viper.GetString("comment"))
	},
}, "comment")

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

var setContextCmd = &cobra.Command{
	Use:     "set",
	Args:    cobra.ExactArgs(1),
	Short:   "Set a current context",
	Long:    "Set a current context",
	Example: "corectl context set local-engine",

	Run: func(ccmd *cobra.Command, args []string) {
		name := internal.SetCurrentContext(args[0])
		printer.PrintCurrentContext(name)
	},
}

var unsetContextCmd = &cobra.Command{
	Use:     "unset",
	Args:    cobra.ExactArgs(0),
	Short:   "Unset current context",
	Long:    "Unset current context",
	Example: "corectl context unset",

	Run: func(ccmd *cobra.Command, args []string) {
		previousContext := internal.UnsetCurrentContext()
		if previousContext != "" {
			printer.PrintCurrentContext("")
		}
	},
}

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

Note that contexts have the lowest precedence. This means that a e.g. an --engine flag
(or an engine field in a config) will override the engine url in the current context.

Contexts are stored locally in your ~/.corectl/contexts.yml file.`,
	Annotations: map[string]string{
		"command_category": "other",
		"x-qlik-stability": "experimental",
	},
}

func init() {
	contextCmd.AddCommand(createContextCmd, removeContextCmd, updateContextCmd, listContextsCmd, setContextCmd, getContextCmd, unsetContextCmd)
}
