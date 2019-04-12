package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var setConnectionsCmd = &cobra.Command{
	Use:   "set <path-to-connections-file.yml>",
	Short: "Set or update the connections in the current app",
	Long:  "Set or update the connections in the current app",
	Example: `corectl connection set
corectl connection set ./my-connections.yml`,

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, true)
		separateConnectionsFile := ""
		if len(args) > 0 {
			separateConnectionsFile = args[0]
		}
		if separateConnectionsFile == "" {
			separateConnectionsFile = getPathFlagFromConfigFile("connections")
		}
		internal.SetupConnections(rootCtx, state.Doc, separateConnectionsFile)
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}

var removeConnectionCmd = &cobra.Command{
	Use:   "remove <connection-id>...",
	Short: "Remove the specified connection(s)",
	Long:  "Remove one or many connections from the app",
	Example: `corectl connection remove
corectl connection remove ID-1
corectl connection remove ID-1 ID-2`,

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected at least one identifier of a connection to delete.")
			ccmd.Usage()
			os.Exit(1)
		}

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

var getConnectionsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "Print a list of all connections in the current app",
	Long:    "Print a list of all connections in the current app",
	Example: `corectl connection ls`,

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		connections, err := state.Doc.GetConnections(rootCtx)
		if err != nil {
			internal.FatalError(err)
		}
		printer.PrintConnections(connections, !viper.GetBool("bash"), viper.GetBool("bash"))
	},
}

var getConnectionCmd = &cobra.Command{
	Use:     "get",
	Short:   "Show the properties for a specific connection",
	Long:    "Show the properties for a specific connection",
	Example: "corectl connection get CONNECTION-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Expected a connection name as parameter")
			ccmd.Usage()
			os.Exit(1)
		}
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
	connectionCmd.AddCommand(setConnectionsCmd, getConnectionCmd, getConnectionsCmd, removeConnectionCmd)
}
