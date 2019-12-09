package cmd

import (
	"fmt"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setScriptCmd = withLocalFlags(&cobra.Command{
	Use:     "set <path-to-script-file.qvs>",
	Args:    cobra.ExactArgs(1),
	Short:   "Set the script in the current app",
	Long:    "Set the script in the current app",
	Example: "corectl script set ./my-script-file.qvs",

	Run: func(ccmd *cobra.Command, args []string) {

		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, true, false)
		scriptFile := args[0]
		if scriptFile != "" {
			internal.SetScript(rootCtx, state.Doc, scriptFile)
		} else {
			log.Fatalln("no loadscript (.qvs) file specified.")
		}
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var getScriptCmd = &cobra.Command{
	Use:     "get",
	Args:    cobra.ExactArgs(0),
	Short:   "Print the reload script",
	Long:    "Print the reload script currently set in the app",
	Example: "corectl script get",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		script, err := state.Doc.GetScript(rootCtx)
		if err != nil {
			log.Fatalf("could not retrieve script: %s\n", err)
		}
		if len(script) == 0 { // This happens if the script is set to an empty file
			fmt.Println("The loadscript is empty")
		} else {
			fmt.Println(script)
		}
	},
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Explore and manage the script",
	Long:  "Explore and manage the script",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

func init() {
	scriptCmd.AddCommand(setScriptCmd, getScriptCmd)
}
