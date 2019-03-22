package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/pkg/browser"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-latest"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Reloads and saves the app after updating connections, dimensions, measures, objects and the script",
	Example: "corectl build --connections ./myconnections.yml --script ./myscript.qvs",
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
		viper.BindPFlag("silent", ccmd.PersistentFlags().Lookup("silent"))
	},
	Run: func(ccmd *cobra.Command, args []string) {
		build(ccmd, args)
	},
}

var catwalkCmd = &cobra.Command{
	Use:   "catwalk",
	Short: "Opens the specified app in catwalk",
	Long:  `Opens the specified app in catwalk. If no app is specified the catwalk hub will be opened.`,
	Example: `corectl catwalk --app my-app.qvf
corectl catwalk --app my-app.qvf --catwalk-url http://localhost:8080`,
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
		viper.BindPFlag("catwalk-url", ccmd.PersistentFlags().Lookup("catwalk-url"))
	},
	Run: func(ccmd *cobra.Command, args []string) {
		catwalkURL := viper.GetString("catwalk-url") + "?engine_url=" + internal.TidyUpEngineURL(viper.GetString("engine")) + "/apps/" + viper.GetString("app")
		if !strings.HasPrefix(catwalkURL, "www") && !strings.HasPrefix(catwalkURL, "https://") && !strings.HasPrefix(catwalkURL, "http://") {
			fmt.Println("Please provide a valid URL starting with 'https://', 'http://' or 'www'")
			os.Exit(1)
		}
		err := browser.OpenURL(catwalkURL)
		if err != nil {
			fmt.Println("Could not open URL", err)
			os.Exit(1)
		}
	},
}

var evalCmd = &cobra.Command{
	Use:   "eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...]",
	Short: "Evaluates a list of measures and dimensions",
	Long:  `Evaluates a list of measures and dimensions. To evaluate a measure for a specific dimension use the <measure> by <dimension> notation. If dimensions are omitted then the eval will be evaluated over all dimensions.`,
	Example: `corectl eval "Count(a)" // returns the number of values in field "a"
corectl eval "1+1" // returns the calculated value for 1+1
corectl eval "Avg(Sales)" by "Region" // returns the average of measure "Sales" for dimension "Region"
corectl eval by "Region" // Returns the values for dimension "Region"`,

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Expected at least one dimension or measure")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, headers, false)
		internal.Eval(rootCtx, state.Doc, args)
	},
}

var reloadCmd = &cobra.Command{
	Use:     "reload",
	Short:   "Reloads the app.",
	Long:    "Reloads the app.",
	Example: "corectl reload",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)

		viper.BindPFlag("silent", ccmd.PersistentFlags().Lookup("silent"))
		viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, false)
		silent := viper.GetBool("silent")

		internal.Reload(rootCtx, state.Doc, state.Global, silent, true)

		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version of corectl",
	Example: "corectl version",

	Run: func(_ *cobra.Command, args []string) {

		if version != "development build" {
			checkLatestVersion()
		}

		fmt.Printf("corectl version: %s\n", version)
	},
}

// completionCmd generates auto completion commands
var completionCmd = &cobra.Command{
	Use:       "completion <shell>",
	ValidArgs: []string{"zsh", "bash"},
	Args:      cobra.MinimumNArgs(1),
	Short:     "Generates auto completion scripts",
	Long: `Generates a shell completion script for the specified shell (bash or zsh). The shell script must be evaluated to provide
interactive completion. This can be done by sourcing it in your ~/.bashrc or ~/.zshrc file. 
Note that jq and bash-completion are required and needs to be installed on your system.`,
	Example: `   Add the following to your ~/.bashrc or ~/.zshrc file

   . <(corectl completion zsh)

   or

   . <(corectl completion bash)`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case args[0] == "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case args[0] == "zsh":
			genZshCompletion()
		default:
			fmt.Printf("%s is not a supported shell", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(catwalkCmd)
	rootCmd.AddCommand(evalCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(completionCmd)
}

func build(ccmd *cobra.Command, args []string) {
	ctx := rootCtx
	state := internal.PrepareEngineState(ctx, headers, true)

	separateConnectionsFile := ccmd.Flag("connections").Value.String()
	if separateConnectionsFile == "" {
		separateConnectionsFile = GetRelativeParameter("connections")
	}
	internal.SetupConnections(ctx, state.Doc, separateConnectionsFile, viper.ConfigFileUsed())
	internal.SetupEntities(ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("dimensions").Value.String(), "dimension")
	internal.SetupEntities(ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("measures").Value.String(), "measure")
	internal.SetupEntities(ctx, state.Doc, viper.ConfigFileUsed(), ccmd.Flag("objects").Value.String(), "object")
	scriptFile := ccmd.Flag("script").Value.String()
	if scriptFile == "" {
		scriptFile = GetRelativeParameter("script")
	}
	if scriptFile != "" {
		internal.SetScript(ctx, state.Doc, scriptFile)
	}

	silent := viper.GetBool("silent")

	internal.Reload(ctx, state.Doc, state.Global, silent, true)

	if state.AppID != "" {
		internal.Save(ctx, state.Doc, state.AppID)
	}
}

// Function for checking current version against latest released version on github
func checkLatestVersion() {
	githubTag := &latest.GithubTag{
		Owner:      "qlik-oss",
		Repository: "corectl",
	}

	res, err := latest.Check(githubTag, version)

	if err == nil && res.Outdated {

		// Find absolute path of executable
		executable, _ := os.Executable()

		// Format a download string depending on OS
		var dwnl string
		if runtime.GOOS == "windows" {
			dwnl = fmt.Sprintf(`curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/v%s/corectl-windows-x86_64.zip" > corectl.zip && unzip ./corectl.zip -d "%s" && rm ./corectl.zip`, res.Current, path.Dir(executable))
		} else {
			dwnl = fmt.Sprintf(`curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/v%s/corectl-%s-x86_64.tar.gz" | tar xz -C /tmp && mv /tmp/corectl %s`, res.Current, runtime.GOOS, path.Dir(executable))
		}

		fmt.Println("-------------------------------------------------")
		fmt.Printf("There is a new version available! Please upgrade for the latest features and bug fixes. You are on %s, latest version is %s. \n", version, res.Current)
		fmt.Printf("To download the latest version you can use this command: \n")
		fmt.Printf(`'%s'`, dwnl)
		fmt.Println("\n-------------------------------------------------")
	}
}

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
