package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//Remove commands
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove entities (connections, dimensions, measures, objects) in the app or the app itself",
	Long:  "Remove one or mores generic entities (connections, dimensions, measures, objects) in the app",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
		viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
		viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
		viper.BindPFlag("headers", ccmd.PersistentFlags().Lookup("headers"))
		viper.BindPFlag("suppress", ccmd.PersistentFlags().Lookup("suppress"))
	},
}

var removeAppCmd = &cobra.Command{
	Use:     "app <app-id>",
	Short:   "removes the specified app.",
	Long:    `removes the specified app.`,
	Example: "corectl remove app APP-ID",
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		removeCmd.PersistentPreRun(removeCmd, args)
	},
	Run: func(ccmd *cobra.Command, args []string) {
		app := viper.GetString("app")

		if len(args) != 1 && app == "" {
			fmt.Println("Expected an identifier of the app to delete.")
			ccmd.Usage()
			os.Exit(1)
		}

		if len(args) == 1 {
			app = args[0]
		}

		confirmed := askForConfirmation(fmt.Sprintf("Do you really want to delete the app: %s?", app))

		if confirmed {
			internal.DeleteApp(rootCtx, viper.GetString("engine"), app, viper.GetString("ttl"), headers)
		}
	},
}

var removeConnectionCmd = &cobra.Command{
	Use:     "connection <connection-id>",
	Short:   "removes the specified connection.",
	Long:    `removes the specified connection.`,
	Example: "corectl remove connection CONNECTION-ID",
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		removeCmd.PersistentPreRun(removeCmd, args)
		viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
	},
	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Expected an identifier of the connection to delete.")
			ccmd.Usage()
			os.Exit(1)
		}

		state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)
		err := state.Doc.DeleteConnection(rootCtx, args[0])
		if err != nil {
			internal.FatalError("Failed to remove connection", args[0])
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}

	},
}

var removeDimensionsCmd = &cobra.Command{
	Use:   "dimensions <dimension-id>...",
	Short: "Removes the specified generic dimensions in the current app",
	Long:  "Removes the specified generic dimensions in the current app. Example: corectl remove dimension ID-1 ID-2",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		removeCmd.PersistentPreRun(removeCmd, args)
		viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected atleast one dimension-id specify what dimension to remove from the app")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyDimension(rootCtx, entity)
			if err != nil {
				internal.FatalError("Failed to remove generic dimension ", entity+" with error: "+err.Error())
			} else if !destroyed {
				internal.FatalError("Failed to remove generic dimension ", entity)
			}
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

var removeMeasuresCmd = &cobra.Command{
	Use:   "measures <measure-id>...",
	Short: "Removes the specified generic measures in the current app",
	Long:  "Removes the specified generic measures in the current app. Example: corectl remove measures ID-1 ID-2",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		removeCmd.PersistentPreRun(removeCmd, args)
		viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected atleast one measure-id specify what measure to remove from the app")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyMeasure(rootCtx, entity)
			if err != nil {
				internal.FatalError("Failed to remove generic measure ", entity+" with error: "+err.Error())
			} else if !destroyed {
				internal.FatalError("Failed to remove generic measure ", entity)
			}
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

var removeObjectsCmd = &cobra.Command{
	Use:   "objects <object-id>...",
	Short: "Removes the specified generic objects in the current app",
	Long:  "Removes the specified generic objects in the current app. Example: corectl remove objects ID-1 ID-2",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		removeCmd.PersistentPreRun(removeCmd, args)
		viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Expected atleast one object-id specify what object to remove from the app")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyObject(rootCtx, entity)
			if err != nil {
				internal.FatalError("Failed to remove generic object ", entity+" with error: "+err.Error())
			} else if !destroyed {
				internal.FatalError("Failed to remove generic object ", entity)
			}
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.AddCommand(removeAppCmd)
	removeCmd.AddCommand(removeConnectionCmd)
	removeCmd.AddCommand(removeDimensionsCmd)
	removeCmd.AddCommand(removeMeasuresCmd)
	removeCmd.AddCommand(removeObjectsCmd)
}

func askForConfirmation(s string) bool {
	if viper.GetString("suppress") == "true" {
		return true
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
