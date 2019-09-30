package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/spf13/cobra"
)

// completionCmd generates auto completion commands
var completionCmd = &cobra.Command{
	Use:       "completion <shell>",
	ValidArgs: []string{"zsh", "bash", "ps"},
	Args:      cobra.ExactValidArgs(1),
	Short:     "Generate auto completion scripts",
	Long: `Generate a shell completion script for the specified shell (bash or zsh). The shell script must be evaluated to provide
interactive completion. This can be done by sourcing it in your ~/.bashrc or ~/.zshrc file.
Note that bash-completion is required and needs to be installed on your system.`,
	Example: `   Add the following to your ~/.bashrc or ~/.zshrc file

   . <(corectl completion zsh)

   or

   . <(corectl completion bash)`,
	Annotations: map[string]string{
		"command_category": "other",
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//Prevent root command
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case args[0] == "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case args[0] == "zsh":
			genZshCompletion()
		case args[0] == "ps":
			rootCmd.GenPowerShellCompletion(os.Stdout)
		default:
			log.Fatalf("'%s' is not a supported shell\n", args[0])
		}
	},
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

// Code for generating zsh bash completion script
// Inspired by https://github.com/kubernetes/kubernetes/blob/e2c1f435516085ef17f222fb7f89cd3ba13aa944/pkg/kubectl/cmd/completion/completion.go
const zshHead = `#compdef corectl
__corectl_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
 	source "$@"
}
 __corectl_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
 		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__corectl_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
 __corectl_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
 	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
 __corectl_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
 __corectl_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
 __corectl_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__corectl_filedir() {
	local RET OLD_IFS w qw
 	__corectl_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
 	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
 	IFS="," __corectl_debug "RET=${RET[@]} len=${#RET[@]}"
 	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__corectl_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__corectl_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
		printf %q "$1"
	fi
}

autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__corectl_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__corectl_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__corectl_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__corectl_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__corectl_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__corectl_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__corectl_type/g" \
	<<'BASH_COMPLETION_EOF'
`

const zshTail = `
BASH_COMPLETION_EOF
}
__corectl_bash_source <(__corectl_convert_bash_to_zsh)
_complete corectl 2>/dev/null
`

// Function for generating zsh completion for corectl
func genZshCompletion() {
	fmt.Fprint(os.Stdout, zshHead)
	buf := new(bytes.Buffer)
	rootCmd.GenBashCompletion(buf)
	fmt.Fprint(os.Stdout, buf.String())
	fmt.Fprint(os.Stdout, zshTail)
}
