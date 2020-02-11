package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/pkg/urtag"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

var setConnectionsCmd = &cobra.Command{
	Use:     "set <path-to-connections-file.yml>",
	Args:    cobra.ExactArgs(1),
	Short:   "Set or update the connections in the current app",
	Long:    "Set or update the connections in the current app",
	Example: "corectl connection set ./my-connections.yml",

	Run: func(ccmd *cobra.Command, args []string) {
		ctx, _, doc, params := urtag.NewCommunicator(ccmd).OpenAppSocket(true)
		separateConnectionsFile := args[0]
		if separateConnectionsFile == "" {
			log.Fatalln("no connections config file specified")
		}
		internal.SetupConnections(ctx, doc, urtag.ReadConnectionsFile(separateConnectionsFile))
		if !params.NoSave() {
			internal.Save(ctx, doc, params.NoData())
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
		ctx, _, doc, params := urtag.NewCommunicator(ccmd).OpenAppSocket(false)
		for _, connection := range args {
			err := doc.DeleteConnection(ctx, connection)
			if err != nil {
				log.Fatalf("could not remove connection '%s': %s\n", connection, err)
			}
		}
		if !params.NoSave() {
			internal.Save(ctx, doc, params.NoData())
		}
	},
}

var listConnectionsCmd = withLocalFlags(&cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all connections in the current app",
	Long:    "Print a list of all connections in the current app",
	Example: "corectl connection ls",

	Run: func(ccmd *cobra.Command, args []string) {
		ctx, _, doc, params := urtag.NewCommunicator(ccmd).OpenAppSocket(false)
		connections, err := doc.GetConnections(ctx)
		if err != nil {
			log.Fatalf("could not retrieve list of connections: %s\n", err)
		}
		printer.PrintConnections(connections, params.PrintMode())
	},
}, "quiet")

var getConnectionCmd = &cobra.Command{
	Use:     "get <connection-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Show the properties for a specific connection",
	Long:    "Show the properties for a specific connection",
	Example: "corectl connection get CONNECTION-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		ctx, _, doc, _ := urtag.NewCommunicator(ccmd).OpenAppSocket(false)
		connection, err := doc.GetConnection(ctx, args[0])
		if err != nil {
			log.Fatalf("could not retrieve connection '%s': %s\n", args[0], err)
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
