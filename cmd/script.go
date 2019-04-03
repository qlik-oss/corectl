package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var setScriptCmd = withCommonLocalFlags(&cobra.Command{
	Use:     "set <path-to-script-file.yml>",
	Short:   "Sets the script in the current app",
	Long:    "Sets the script in the current app",
	Example: "corectl set script ./my-script-file",

	Run: func(ccmd *cobra.Command, args []string) {

		state := internal.PrepareEngineState(rootCtx, headers, true)
		scriptFile := ""
		if len(args) > 0 {
			scriptFile = args[0]
		}
		if scriptFile == "" {
			scriptFile = GetRelativeParameter("script")
		}
		if scriptFile != "" {
			internal.SetScript(rootCtx, state.Doc, scriptFile)
		} else {
			fmt.Println("Expected the path to a file containing the qlik script")
			ccmd.Usage()
			os.Exit(1)
		}
		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var getScriptCmd = &cobra.Command{
	Use:   "get",
	Short: "Print the reload script",
	Long:  "Fetches the script currently set in the app and prints it in plain text.",
	Example: `corectl get script
corectl get script --app=my-app.qvf`,

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
