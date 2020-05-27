package engine

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

func CreateConnectionCommand() *cobra.Command {
	var setConnectionsCmd = &cobra.Command{
		Use:     "set <path-to-connections-file.yml>",
		Args:    cobra.ExactArgs(1),
		Short:   "Set or update the connections in the current app",
		Long:    "Set or update the connections in the current app",
		Example: "corectl connection set ./my-connections.yml",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(true)
			separateConnectionsFile := args[0]
			if separateConnectionsFile == "" {
				log.Fatalln("no connections config file specified")
			}
			internal.SetupConnections(ctx, doc, boot.ReadConnectionsFile(separateConnectionsFile))
			if !params.NoSave() {
				internal.Save(ctx, doc, params.NoData())
			}
		},
	}

	var removeConnectionCmd = &cobra.Command{
		Use:               "rm <connection-id>...",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: listValidConnectionsForCompletion,
		Short:             "Remove the specified connection(s)",
		Long:              "Remove one or many connections from the app",
		Example: `corectl connection rm
corectl connection rm ID-1
corectl connection rm ID-1 ID-2`,

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
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

	var listConnectionsCmd = WithLocalFlags(&cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(0),
		Short:   "Print a list of all connections in the current app",
		Long:    "Print a list of all connections in the current app",
		Example: "corectl connection ls",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, params := boot.NewCommunicator(ccmd).OpenAppSocket(false)
			connections, err := doc.GetConnections(ctx)
			if err != nil {
				log.Fatalf("could not retrieve list of connections: %s\n", err)
			}
			printer.PrintConnections(connections, params.PrintMode())
		},
	}, "quiet")

	var getConnectionCmd = &cobra.Command{
		Use:               "get <connection-id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: listValidConnectionsForCompletion,
		Short:             "Show the properties for a specific connection",
		Long:              "Show the properties for a specific connection",
		Example:           "corectl connection get CONNECTION-ID",

		Run: func(ccmd *cobra.Command, args []string) {
			ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
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

	connectionCmd.AddCommand(setConnectionsCmd, getConnectionCmd, listConnectionsCmd, removeConnectionCmd)
	return connectionCmd
}

func listValidConnectionsForCompletion(ccmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ctx, _, doc, _ := boot.NewCommunicator(ccmd).OpenAppSocket(false)
	items, err := doc.GetConnections(ctx)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}
	result := make([]string, 0)
	for _, item := range items {
		result = append(result, item.Id)
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}
