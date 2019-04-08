package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var headersMap = make(map[string]string)
var explicitConfigFile = ""
var version = ""
var headers http.Header
var rootCtx = context.Background()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Hidden:                 true,
	Use:                    "corectl",
	Short:                  "",
	Long:                   `corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
	DisableAutoGenTag:      true,
	BashCompletionFunction: bashCompletionFunc,
	Annotations: map[string]string{
		"x-qlik-stability": "experimental",
	},

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		// if help, version or generate-docs command, no prerun is needed.
		if strings.Contains(ccmd.Use, "help") || ccmd.Use == "generate-docs" || ccmd.Use == "generate-spec" || ccmd.Use == "version" {
			return
		}
		internal.ReadConfigFile(explicitConfigFile)

		if len(headersMap) == 0 {
			headersMap = viper.GetStringMapString("headers")
		}
		headers = make(http.Header, 1)
		for key, value := range headersMap {
			headers.Set(key, value)
		}
	},

	Run: func(ccmd *cobra.Command, args []string) {
		ccmd.HelpFunc()(ccmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(mainVersion string) {
	version = mainVersion
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&explicitConfigFile, "config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolP("traffic", "t", false, "Log JSON websocket traffic to stdout")
	viper.BindPFlag("traffic", rootCmd.PersistentFlags().Lookup("traffic"))

	rootCmd.PersistentFlags().StringP("engine", "e", "localhost:9076", "URL to the Qlik Associative Engine")
	viper.BindPFlag("engine", rootCmd.PersistentFlags().Lookup("engine"))

	rootCmd.PersistentFlags().String("ttl", "30", "Qlik Associative Engine session time to live in seconds")
	viper.BindPFlag("ttl", rootCmd.PersistentFlags().Lookup("ttl"))

	rootCmd.PersistentFlags().Bool("no-data", false, "Open app without data")
	viper.BindPFlag("no-data", rootCmd.PersistentFlags().Lookup("no-data"))

	rootCmd.PersistentFlags().Bool("bash", false, "Bash flag used to adapt output to bash completion format")
	rootCmd.PersistentFlags().MarkHidden("bash")
	viper.BindPFlag("bash", rootCmd.PersistentFlags().Lookup("bash"))

	//not binding to viper since binding a map does not seem to work.
	rootCmd.PersistentFlags().StringToStringVar(&headersMap, "headers", nil, "Http headers to use when connecting to Qlik Associative Engine")

	rootCmd.PersistentFlags().StringP("app", "a", "", "App name, if no app is specified a session app is used instead.")
	viper.BindPFlag("app", rootCmd.PersistentFlags().Lookup("app"))
	// Set annotation to run bash completion function
	rootCmd.PersistentFlags().SetAnnotation("app", cobra.BashCompCustom, []string{"__corectl_get_apps"})

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("connections", "", "Path to a yml file containing the data connection definitions")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("dimensions", "", "A list of generic dimension json paths")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("measures", "", "A list of generic measures json paths")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("objects", "", "A list of generic object json paths")
	}

	for _, command := range []*cobra.Command{buildCmd, reloadCmd} {
		command.PersistentFlags().Bool("silent", false, "Do not log reload progress")
	}

	for _, command := range []*cobra.Command{reloadCmd, removeConnectionCmd, removeDimensionCmd, removeMeasureCmd, removeObjectCmd, setCmd} {
		command.PersistentFlags().Bool("no-save", false, "Do not save the app")
	}

	for _, command := range []*cobra.Command{buildCmd, setAllCmd} {
		// Don't bind these to viper since paths are treated separately to support relative paths!
		command.PersistentFlags().String("script", "", "path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded")
	}

	for _, command := range []*cobra.Command{getAppsCmd, getConnectionsCmd, getDimensionsCmd, getMeasuresCmd, getObjectsCmd} {
		command.PersistentFlags().Bool("json", false, "Prints the information in json format")
	}

	for _, command := range []*cobra.Command{removeCmd} {
		command.PersistentFlags().Bool("suppress", false, "Suppress all confirmation dialogues")
	}
	catwalkCmd.PersistentFlags().String("catwalk-url", "https://catwalk.core.qlik.com", "Url to an instance of catwalk, if not provided the qlik one will be used.")

	if runtime.GOOS != "windows" {
		// Do not add bash completion annotations for paths and files as they are not compatible with windows. On windows
		// we instead rely on the default bash behavior
		addFileRelatedBashAnnotations()
	}
}

// GetRelativeParameter returns a parameter from the config file.
// It modifies the parameter to actually be relative to the config file and not the working directory
func GetRelativeParameter(paramName string) string {
	pathInConfigFile := viper.GetString(paramName)
	if pathInConfigFile != "" {
		return internal.RelativeToProject(viper.ConfigFileUsed(), pathInConfigFile)
	}
	return ""
}

func getEntityProperties(ccmd *cobra.Command, args []string, entityType string) {
	if len(args) < 1 {
		fmt.Println("Expected an " + entityType + " id to specify what " + entityType + " to use as a parameter")
		ccmd.Usage()
		os.Exit(1)
	}
	state := internal.PrepareEngineState(rootCtx, headers, false)
	printer.PrintGenericEntityProperties(state, args[0], entityType)
}

func getEntityLayout(ccmd *cobra.Command, args []string, entityType string) {
	if len(args) < 1 {
		fmt.Println("Expected an " + entityType + " id to specify what " + entityType + " to use as a parameter")
		ccmd.Usage()
		os.Exit(1)
	}
	state := internal.PrepareEngineState(rootCtx, headers, false)
	printer.PrintGenericEntityLayout(state, args[0], entityType)
}

func getEntities(ccmd *cobra.Command, args []string, entityType string, printAsJSON bool) {
	state := internal.PrepareEngineState(rootCtx, headers, false)
	allInfos, err := state.Doc.GetAllInfos(rootCtx)
	if err != nil {
		internal.FatalError(err)
	}
	printer.PrintGenericEntities(allInfos, entityType, printAsJSON, viper.GetBool("bash"))
}

const bashCompletionFunc = `

	__custom_func()
	{
		case ${last_command} in
			corectl_get_dimension_properties | corectl_get_dimension_layout)
				__corectl_get_dimensions
				;;
			corectl_get_measure_properties | corectl_get_measure_layout)
				__corectl_get_measures
				;;
			corectl_get_object_data | corectl_get_object_properties | corectl_get_object_layout)
				__corectl_get_objects
				;;
			corectl_get_connection)
				__corectl_get_connections
				;;
      *)
				COMPREPLY+=( $( compgen -W "" -- "$cur" ) )
				;;
		esac
	}

  __extract_flags_to_forward()
	{
    local forward_flags
  	local result
	  forward_flags=( "--engine" "-e" "--app" "-a" "--config" "-c" "--headers" "--ttl" );
	  while [[ $# -gt 0 ]]; do
  	  for i in "${forward_flags[@]}"
			do
				case $1 in
				$i)
					# If there is a flag with spacing we need to check that an arg is passed
					if [[ $# -gt 1 ]]; then
						result+="$1=";
						shift;
						result+="$1 "
					fi
      	;;
      	$i=*)
        	result+="$1 "
      	;;
    	esac
			done
    	shift
  	done
    echo "$result";
	}

  __corectl_call_corectl() 
  {
    local flags=$(__extract_flags_to_forward ${words[@]})
		local corectl_out
		local errorcode
		corectl_out=$(corectl $1 $flags 2>/dev/null)
		errorcode=$?
		if [[ errorcode -eq 0 ]]; then
  		local IFS=$'\n'
  		COMPREPLY=( $(compgen -W "${corectl_out}" -- "$cur") )
		else
  		COMPREPLY=()
		fi;
  }

	__corectl_get_dimensions()
	{
		__corectl_call_corectl "get dimensions --bash"
	}

	__corectl_get_measures()
	{
		__corectl_call_corectl "get measures --bash"
	}

	__corectl_get_objects()
	{
		__corectl_call_corectl "get objects --bash"
	}

	__corectl_get_connections()
	{
		__corectl_call_corectl "get connections --bash"
	}

	__corectl_get_apps()
	{
		__corectl_call_corectl "get apps --bash"
	}
`
