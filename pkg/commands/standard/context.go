package standard

import (
	"fmt"

	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/commands/engine"
	"github.com/qlik-oss/corectl/pkg/commands/login"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
)

func CreateContextCommand(binaryName string) *cobra.Command {
	oneArgCompletion := contextCompletion(1)
	manyArgCompletion := contextCompletion(-1)

	var createContextCmd = engine.WithLocalFlags(&cobra.Command{
		Use:   "create <context name>",
		Args:  cobra.ExactArgs(1),
		Short: "Create a context with the specified configuration",
		Long: `Create a context with the specified configuration

This command creates a context by using the supplied flags.
The information stored will be server url, headers and certificates
(if present) along with comment and the context-name.`,

		Example: fmt.Sprintf(`%[1]s context create local-engine
%[1]s context create rd-sense --server localhost:9076 --comment "R&D Qlik Sense deployment"`, binaryName),

		Run: func(ccmd *cobra.Command, args []string) {

			//Check the validity of the certificates folder
			cfg := dynconf.ReadSettingsWithoutContext(ccmd)

			headers := cfg.GetStringMap("headers")
			if len(headers) == 0 {
				headers = nil
			}

			newSettings := map[string]interface{}{
				"headers": headers,
				"comment": cfg.GetString("comment"),
			}

			certPath := cfg.GetAbsolutePath("certificates")
			if certPath != "" {
				newSettings["certificates"] = certPath
				cfg.GetTLSConfigFromPath("certificates")
			}

			if !cfg.IsUsingDefaultValue("server") {
				newSettings["server"] = cfg.GetString("server")
			}

			dynconf.CreateContext(args[0], newSettings)
		},
	}, "comment", "api-key")

	var updateContextCmd = engine.WithLocalFlags(&cobra.Command{
		Use:   "update <context name>",
		Args:  cobra.ExactArgs(1),
		Short: "Update a context with the specified configuration",
		Long:  "Update a context with the specified configuration",

		Example: fmt.Sprintf(`%[1]s context update local-engine
%[1]s context update rd-sense --server localhost:9076 --comment "R&D Qlik Sense deployment"`, binaryName),
		ValidArgsFunction: oneArgCompletion,

		Run: func(ccmd *cobra.Command, args []string) {

			// Check the validity of the certificates folder.
			cfg := dynconf.ReadSettingsWithoutContext(ccmd)

			headers := cfg.GetHeaders()
			if len(headers) == 0 {
				headers = nil
			}

			// Convert from http.Header to map[string]string
			headerMap := map[string]string{}
			for key := range headers {
				headerMap[key] = headers.Get(key)
			}

			newSettings := map[string]interface{}{
				"headers": headerMap,
				"comment": cfg.GetString("comment"),
			}

			certPath := cfg.GetAbsolutePath("certificates")
			if certPath != "" {
				newSettings["certificates"] = certPath
				cfg.GetTLSConfigFromPath("certificates")
			}

			if !cfg.IsUsingDefaultValue("server") {
				newSettings["server"] = cfg.GetString("server")
			}

			dynconf.UpdateContext(args[0], newSettings)
		},
	}, "comment", "api-key")

	var removeContextCmd = &cobra.Command{
		Use:   "rm <context name>...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Remove one or more contexts",
		Long:  "Remove one or more contexts",
		Example: fmt.Sprintf(`%[1]s context rm local-engine
%[1]s context rm ctx1 ctx2`, binaryName),
		ValidArgsFunction: manyArgCompletion,

		Run: func(ccmd *cobra.Command, args []string) {
			var removedCurrent bool
			for _, arg := range args {
				_, wasCurrent := dynconf.RemoveContext(arg)
				if wasCurrent {
					removedCurrent = true
				}
			}
			if removedCurrent {
				printer.PrintCurrentContext("")
			}
		},
	}

	var getContextCmd = &cobra.Command{
		Use:   "get [context name]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Get context, current context by default",
		Long:  "Get context, current context by default",
		Example: fmt.Sprintf(`%[1]s context get
%[1]s context get local-engine`, binaryName),
		ValidArgsFunction: oneArgCompletion,

		Run: func(ccmd *cobra.Command, args []string) {
			handler := dynconf.NewContextHandler()
			var name string

			if len(args) == 1 {
				name = args[0]
			}
			printer.PrintContext(name, handler)
		},
	}

	var listContextsCmd = &cobra.Command{
		Use:     "ls",
		Args:    cobra.ExactArgs(0),
		Short:   "List all contexts",
		Long:    "List all contexts",
		Example: binaryName + " context ls",

		Run: func(ccmd *cobra.Command, args []string) {
			handler := dynconf.NewContextHandler()
			comm := boot.NewCommunicator(ccmd)
			printer.PrintContexts(handler, comm.PrintMode())
		},
	}

	var useContextCmd = &cobra.Command{
		Use:               "use <context-name>",
		Args:              cobra.ExactArgs(1),
		Short:             "Specify what context to use",
		Long:              "Specify what context to use",
		Example:           binaryName + " context use local-engine",
		ValidArgsFunction: oneArgCompletion,

		Run: func(ccmd *cobra.Command, args []string) {
			name := dynconf.UseContext(args[0])
			printer.PrintCurrentContext(name)
		},
	}

	var clearContextCmd = &cobra.Command{
		Use:     "clear",
		Args:    cobra.ExactArgs(0),
		Short:   "Set the current context to none",
		Long:    "Set the current context to none",
		Example: binaryName + " context clear",

		Run: func(ccmd *cobra.Command, args []string) {
			previous := dynconf.ClearContext()
			if previous != "" {
				printer.PrintCurrentContext("")
			}
		},
	}

	var loginContextCmd = engine.WithLocalFlags(&cobra.Command{
		Use:   "login <context-name>",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Login and set cookie for the named context",
		Long: `Login and set cookie for the named context

This is only applicable when connecting to 'Qlik Sense Enterprise for Windows' through its proxy using HTTPS.
If no 'context-name' is used as argument the 'current-context' defined in the config will be used instead.`,
		Example: fmt.Sprintf(`%[1]s context login
%[1]s context login context-name`, binaryName),
		ValidArgsFunction: oneArgCompletion,

		Run: func(ccmd *cobra.Command, args []string) {
			contextName := ""

			if len(args) > 0 {
				contextName = args[0]
			}

			comm := boot.NewCommunicator(ccmd)

			dynconf.LoginContext(comm.TlsConfig(), contextName, comm.GetString("user"), comm.GetString("password"))
		},
	}, "user", "password")

	var contextCmd = &cobra.Command{
		Use:   "context",
		Short: "Create, update and use contexts",
		Long: `Create, update and use contexts

Contexts store connection information such as server url, certificates and headers,
similar to a config. The main difference between contexts and configs is that they
can be used globally. Use the context subcommands to configure contexts which
facilitate app development in environments where certificates and headers are needed.

The current context is the one that is being used. You can use "context get" to
display the contents of the current context and switch context with "context set"
or unset the current context with "context unset".

Note that contexts have the lowest precedence. This means that e.g. an --server flag
(or a server field in a config) will override the server url in the current context.

Contexts are stored locally in your ~/` + dynconf.ContextDir + `/contexts.yml file.`,
		Annotations: map[string]string{
			"command_category": "other",
			"x-qlik-stability": "experimental",
		},
	}

	contextCmd.AddCommand(createContextCmd, updateContextCmd, removeContextCmd,
		listContextsCmd, useContextCmd, getContextCmd,
		clearContextCmd, loginContextCmd, login.CreateInitCommand())

	return contextCmd
}

type Completion func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)

func contextCompletion(n int) Completion {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if n > 0 && n <= len(args) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		handler := dynconf.NewContextHandler()
		for _, arg := range args {
			if _, ok := handler.Contexts[arg]; ok {
				delete(handler.Contexts, arg)
			}
		}
		contexts := make([]string, len(handler.Contexts))
		i := 0
		for context := range handler.Contexts {
			contexts[i] = context
			i++
		}
		return contexts, cobra.ShellCompDirectiveNoFileComp
	}
}
