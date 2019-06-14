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
	Short: "Creates a new context",
	Long:  "Creates a new context",
	Example: `corectl context create local-engine
corectl context create rd-sense --product "QSE" --comment "R&D Qlik Sense deployment"`,

	Run: func(ccmd *cobra.Command, args []string) {
		// Add validation of product
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

var listContextsCmd = &cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "List all contexts",
	Long:    "List all contexts",
	Example: "corectl context ls",

	Run: func(ccmd *cobra.Command, args []string) {
		currentContext, contexts := internal.GetContexts()
		printer.PrintContexts(contexts, currentContext, viper.GetBool("bash"))
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

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Explore and manage contexts",
	Long:  "Explore and manage contexts",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	contextCmd.AddCommand(createContextCmd, removeContextCmd, listContextsCmd, setContextCmd)
}
