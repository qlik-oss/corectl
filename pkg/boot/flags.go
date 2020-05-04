package boot

import (
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"runtime"
)

func InjectGlobalFlags(command *cobra.Command, hideEngineSpecificFlags bool) {
	globalFlags := command.PersistentFlags()
	globalFlags.BoolP("verbose", "v", false, "Log extra information")
	globalFlags.BoolP("traffic", "t", false, "Log JSON websocket traffic to stdout")
	globalFlags.StringP("server", "s", "", "URL to a Qlik Product, a local engine, cluster or sense-enterprise")
	globalFlags.String("ttl", "0", "Qlik Associative Engine session time to live in seconds")
	globalFlags.Bool("json", false, "Returns output in JSON format if possible, disables verbose and traffic output")
	globalFlags.Bool("no-data", false, "Open app without data")
	globalFlags.Bool("bash", false, "Bash flag used to adapt output to bash completion format")
	globalFlags.MarkHidden("bash")
	globalFlags.String("context", "", "Name of the context used when connecting to Qlik Associative Engine")
	globalFlags.Bool("insecure", false, "Enabling insecure will make it possible to connect using self signed certificates")
	globalFlags.StringP("config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")
	globalFlags.StringToString("headers", nil, "Http headers to use when connecting to Qlik Associative Engine")
	globalFlags.String("certificates", "", "path/to/folder containing client.pem, client_key.pem and root.pem certificates")

	// Set annotation to run bash completion function
	if !hideEngineSpecificFlags {
		globalFlags.SetAnnotation("server", cobra.BashCompCustom, []string{"__corectl_get_local_engines"})
		globalFlags.SetAnnotation("context", cobra.BashCompCustom, []string{"__corectl_get_contexts"})
	}

	if runtime.GOOS != "windows" {
		// Do not add bash completion annotations for paths and files as they are not compatible with windows. On windows
		// we instead rely on the default bash behavior
		globalFlags.SetAnnotation("config", cobra.BashCompFilenameExt, []string{"yaml", "yml"})
	}

	if hideEngineSpecificFlags {
		globalFlags.MarkHidden("no-data")
		globalFlags.MarkHidden("ttl")
		globalFlags.MarkHidden("traffic")
		globalFlags.MarkHidden("certificates")
	}

	// Add all global flags to the set of valid config properties.
	globalFlags.VisitAll(func(flag *pflag.Flag) {
		dynconf.AddValidConfigFilePropertyName(flag.Name)
	})
}

func InjectAppWebSocketFlags(command *cobra.Command, hideEngineSpecificAppCompletion bool) {
	globalFlags := command.PersistentFlags()
	globalFlags.StringP("app", "a", "", "Name or identifier of the app")
	if !hideEngineSpecificAppCompletion {
		globalFlags.SetAnnotation("app", cobra.BashCompCustom, []string{"__corectl_get_apps"})
	}
	// Add all global flags to the set of valid config properties.
	globalFlags.VisitAll(func(flag *pflag.Flag) {
		dynconf.AddValidConfigFilePropertyName(flag.Name)
	})
}
