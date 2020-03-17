package cmd

import (
	"context"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/commands/engine"
	"github.com/qlik-oss/corectl/pkg/commands/standard"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCtx = context.Background()

func CreateRootCommand(version, branch, commit string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Hidden:                 true,
		Use:                    "corectl",
		Short:                  "",
		Long:                   `corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information`,
		DisableAutoGenTag:      true,
		BashCompletionFunction: bashCompletionFunc,

		Annotations: map[string]string{
			"x-qlik-stability": "stable",
		},

		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	//App Building Commands
	appBuildCommands := []*cobra.Command{
		engine.CreateBuildCommand(),
		engine.CreateReloadCommand(),
		engine.CreateUnbuildCommand(),
	}
	Annotate("command_category", "build", appBuildCommands...)
	rootCmd.AddCommand(appBuildCommands...)

	// Common commands
	commonCommands := []*cobra.Command{
		engine.CreategetAssociationsCommand(),
		engine.CreateCatwalkCommand(),
		engine.CreateEvalCommand(),
		engine.CreateGetFieldsCommand(),
		engine.CreateGetValuesCommand(),
		engine.CreateGetMetaCommand(),
		engine.CreateGetKeysCommand(),
		engine.CreategetTablesCommand(),
	}
	//Annotate("command_category", "common", getAssociationsCmd, catwalkCmd, evalCmd,
	//		getFieldsCmd, getValuesCmd, getMetaCmd, getKeysCmd, getTablesCmd)
	rootCmd.AddCommand(commonCommands...)

	// Subcommands
	subCommands := []*cobra.Command{
		engine.CreateAppCommand(),
		engine.CreateBookmarkCommand(),
		engine.CreateConnectionCommand(),
		engine.CreateDimensionCommand(),
		engine.CreateMeasureCommand(),
		engine.CreateObjectCommand(),
		engine.CreateScriptCommand(),
		engine.CreateAlternateStateCommand(),
		engine.CreateVariableCommand(),
	}
	Annotate("command_category", "sub", subCommands...)
	rootCmd.AddCommand(subCommands...)

	// Other
	otherCommands := []*cobra.Command{
		standard.CreateCompletionCommand("corectl"),
		standard.CreateContextCommand(),
		standard.CreateStatusCommand(),
		standard.CreateVersionCommand(version, branch, commit),
	}
	Annotate("command_category", "other", otherCommands...)
	rootCmd.AddCommand(otherCommands...)

	// Hidden administrative commands
	rootCmd.AddCommand(standard.CreateGenerateDocsCommand())
	rootCmd.AddCommand(standard.CreateGenerateSpecCommand(version))

	boot.InjectGlobalFlags(rootCmd, false)

	//initGlobalFlags(rootCmd.PersistentFlags())
	patchRootCommandUsageTemplate(rootCmd)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, branch, commit string) {
	rootCmd := CreateRootCommand(version, branch, commit)
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints an error message so we just want to exit
		os.Exit(1)
	}
}

func patchRootCommandUsageTemplate(rootCmd *cobra.Command) {
	var originalUsageSnippet = `Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var rootSnippetMainSection = `App Building Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (eq (index .Annotations "command_category") "build"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

App Analysis Commands:{{range .Commands}}{{if (and .IsAvailableCommand (eq (index .Annotations "command_category") ""))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Advanced Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (eq (index .Annotations "command_category") "sub"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Other Commands:{{range .Commands}}{{if (and (or .IsAvailableCommand (eq .Name "help")) (or (eq (index .Annotations "command_category") "other") (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}`

	var newUsageSnippet = `{{if (eq .Name "corectl")}}` + rootSnippetMainSection + `{{else}}` + originalUsageSnippet + "{{end}}"

	var patchedUsageTemplate = strings.Replace(rootCmd.UsageTemplate(), originalUsageSnippet, newUsageSnippet, 1)
	rootCmd.SetUsageTemplate(patchedUsageTemplate)
}

func Annotate(key, value string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		if cmd.Annotations == nil {
			cmd.Annotations = map[string]string{}
		}
		cmd.Annotations[key] = value
	}
}

const bashCompletionFunc = `

__corectl_custom_func()
{
	case ${last_command} in
		corectl_dimension_properties | corectl_dimension_layout | corectl_dimension_rm)
			__corectl_get_dimensions
			;;
		corectl_measure_properties | corectl_measure_layout | corectl_measure_rm)
			__corectl_get_measures
			;;
		corectl_bookmark_properties | corectl_bookmark_layout | corectl_bookmark_rm)
			__corectl_get_bookmarks
			;;
		corectl_variable_properties | corectl_variable_layout | corectl_variable_rm)
			__corectl_get_variables
			;;
		corectl_object_data | corectl_object_properties | corectl_object_layout | corectl_object_rm)
			__corectl_get_objects
			;;
		corectl_connection_get | corectl_connection_rm)
			__corectl_get_connections
			;;
		corectl_state_rm)
			__corectl_state_ls
			;;
		corectl_app_rm)
			__corectl_get_apps
			;;
		corectl_context_rm | corectl_context_set | corectl_context_get | corectl_context_use)
			__corectl_get_contexts
			;;
		corectl_app_import)
			__corectl_handle_filename_extension_flag "qvf"
			;;
		corectl_script_set)
			__corectl_handle_filename_extension_flag "qvs"
			;;
		corectl_dimension_set | corectl_measure_set | corectl_bookmark_set | corectl_variable_set | corectl_object_set)
			__corectl_handle_filename_extension_flag "json"
			;;
		corectl_connection_set)
			__corectl_handle_filename_extension_flag "yaml|yml"
			;;
		corectl_values)
			__corectl_fields
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
	          result+="$1"
	        fi
	        ;;
	      $i=*)
	        result+="$1"
	        ;;
	    esac
			# Since host:port gets treated as 3 words by cobra we have to puzzle it back to an url again
			# Also, the 'words' contain a lot of trailing whitespaces hence the sed trim
	    if [[ $# -gt 2 ]]; then
	      if [ "$1" = ":" ]; then
	        shift;
	        result=$(echo $result | sed 's/[ \t]*$//')
	        result+=":$1"
	      fi
	    fi
	  result+=" "
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
	__corectl_call_corectl "dimension ls --bash"
}

__corectl_get_measures()
{
	__corectl_call_corectl "measure ls --bash"
}

__corectl_get_bookmarks()
{
	__corectl_call_corectl "bookmark ls --bash"
}

__corectl_get_variables()
{
	__corectl_call_corectl "variable ls --bash"
}

__corectl_get_objects()
{
	__corectl_call_corectl "object ls --bash"
}

__corectl_get_connections()
{
	__corectl_call_corectl "connection ls --bash"
}

__corectl_state_ls()
{
	__corectl_call_corectl "state ls --bash"
}

__corectl_get_apps()
{
	__corectl_call_corectl "app ls --bash"
}

__corectl_get_local_engines()
{
	local docker_out
	local errorcode
	docker_out=$(docker ps 2>/dev/null | grep /engine:|sed \ 's/.*0.0.0.0:/localhost:/g'|sed 's/->.*//g')
	errorcode=$?
	if [[ errorcode -eq 0 ]]; then
		local IFS=$'\n'
		COMPREPLY=( $(compgen -W "${docker_out}" -- "$cur") )
	else
		COMPREPLY=()
	fi;
}

__corectl_get_contexts()
{
	__corectl_call_corectl "context ls --bash"
}

__corectl_fields()
{
	__corectl_call_corectl "fields --bash"
}

`
