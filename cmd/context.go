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
	Long:  "Create a new context",
	Example: `corectl context create local-engine
corectl context create rd-sense --product "QSE" --comment "R&D Qlik Sense deployment"`,

	Run: func(ccmd *cobra.Command, args []string) {
		internal.CreateContext(args[0], viper.GetString("product"), viper.GetString("comment"))
	},
}, "product", "comment")

var removeContextCmd = &cobra.Command{
	Use:     "rm <context name>",
	Args:    cobra.ExactArgs(1),
	Short:   "Removes a context",
	Long:    "Removes a context",
	Example: "corectl context rm local-engine",

	Run: func(ccmd *cobra.Command, args []string) {
		internal.RemoveContext(args[0])
	},
}

var getContextCmd = &cobra.Command{
	Use:     "get",
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Get context, current context by default",
	Long:    "Get context, current context by default",
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
		internal.SetCurrentContext(args[0])
	},
}

var unsetContextCmd = &cobra.Command{
	Use:	"unset",
	Args:	cobra.ExactArgs(0),
	Short: "Unset current context",
	Long: "Unset current context",
	Example: "corectl context unset",

	Run: func(ccmd *cobra.Command, args []string) {
		internal.UnsetCurrentContext()
	},
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Explore and manage contexts",
	Long:  "Explore and manage contexts",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	contextCmd.AddCommand(createContextCmd, removeContextCmd, listContextsCmd, setContextCmd, getContextCmd, unsetContextCmd)
}
