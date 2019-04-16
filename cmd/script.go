package cmd

import (
	"fmt"

	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setScriptCmd = withLocalFlags(&cobra.Command{
	Use:     "set <path-to-script-file.qvs>",
	Short:   "Set the script in the current app",
	Long:    "Set the script in the current app",
	Example: "corectl script set ./my-script-file.qvs",

	Args: cobra.ExactArgs(1),
	Run: func(ccmd *cobra.Command, args []string) {

		state := internal.PrepareEngineState(rootCtx, headers, true)
		scriptFile := args[0]
		if scriptFile != "" {
			internal.SetScript(rootCtx, state.Doc, scriptFile)
		} else {
			internal.FatalError("Error: No loadscript (.qvs) file specified.")
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var getScriptCmd = &cobra.Command{
	Use:     "get",
	Short:   "Print the reload script",
	Long:    "Print the reload script currently set in the app",
	Example: `corectl script get`,

	Args: cobra.ExactArgs(0),
	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		script, err := state.Doc.GetScript(rootCtx)
		if err != nil {
			internal.FatalError(err)
		}
		fmt.Println(script)
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
