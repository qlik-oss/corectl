package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setConnectionsCmd = &cobra.Command{
	Use:     "set <path-to-connections-file.yml>",
	Args:    cobra.ExactArgs(1),
	Short:   "Set or update the connections in the current app",
	Long:    "Set or update the connections in the current app",
	Example: "corectl connection set ./my-connections.yml",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, true)
		separateConnectionsFile := args[0]
		if separateConnectionsFile == "" {
			internal.FatalError("Error: No connection file specified.")
		}
		internal.SetupConnections(rootCtx, state.Doc, separateConnectionsFile)
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var removeConnectionCmd = &cobra.Command{
	Use:   "rm <connection-id>...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Remove the specified connection(s)",
	Long:  "Remove one or many connections from the app",
	Example: `corectl connection rm
corectl connection rm ID-1
corectl connection rm ID-1 ID-2`,

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		for _, connection := range args {
			err := state.Doc.DeleteConnection(rootCtx, connection)
			if err != nil {
				internal.FatalError("Failed to remove connection: ", connection, " with error: ", err.Error())
			}
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var listConnectionsCmd = &cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all connections in the current app",
	Long:    "Print a list of all connections in the current app",
	Example: "corectl connection ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		connections, err := state.Doc.GetConnections(rootCtx)
		if err != nil {
			internal.FatalError(err)
		}
		printer.PrintConnections(connections, viper.GetBool("bash"))
	},
}

var getConnectionCmd = &cobra.Command{
	Use:     "get <connection-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Show the properties for a specific connection",
	Long:    "Show the properties for a specific connection",
	Example: "corectl connection get CONNECTION-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		connection, err := state.Doc.GetConnection(rootCtx, args[0])
		if err != nil {
			internal.FatalError(err)
		}
		printer.PrintConnection(connection)
	},
}

var connectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "Explore and manage connections",
	Long:  "Explore and manage connections",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	connectionCmd.AddCommand(setConnectionsCmd, getConnectionCmd, listConnectionsCmd, removeConnectionCmd)
}
