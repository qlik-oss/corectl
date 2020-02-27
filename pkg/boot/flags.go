package boot

import (
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"runtime"
)

func InjectGlobalFlags(command *cobra.Command) {
	globalFlags := command.PersistentFlags()
	globalFlags.BoolP("verbose", "v", false, "Log extra information")
	globalFlags.BoolP("traffic", "t", false, "Log JSON websocket traffic to stdout")
	globalFlags.StringP("engine", "e", "localhost:9076", "URL to the Qlik Associative Engine")
	globalFlags.StringP("app", "a", "", "Name or identifier of the app")
	globalFlags.String("ttl", "0", "Qlik Associative Engine session time to live in seconds")
	globalFlags.Bool("json", false, "Returns output in JSON format if possible, disables verbose and traffic output")
	globalFlags.Bool("no-data", false, "Open app without data")
	globalFlags.Bool("bash", false, "Bash flag used to adapt output to bash completion format")
	globalFlags.MarkHidden("bash")
	globalFlags.String("context", "", "Name of the context used when connecting to Qlik Associative Engine")
	globalFlags.Bool("insecure", false, "Enabling insecure will make it possible to connect using self signed certificates")

	// Not bound to viper. Certificates are handled similarly to config as the path might be relative.
	globalFlags.StringP("config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")
	globalFlags.StringToString("headers", nil, "Http headers to use when connecting to Qlik Associative Engine")
	globalFlags.String("certificates", "", "path/to/folder containing client.pem, client_key.pem and root.pem certificates")

	// Set annotation to run bash completion function
	globalFlags.SetAnnotation("app", cobra.BashCompCustom, []string{"__corectl_get_apps"})
	globalFlags.SetAnnotation("engine", cobra.BashCompCustom, []string{"__corectl_get_local_engines"})
	globalFlags.SetAnnotation("context", cobra.BashCompCustom, []string{"__corectl_get_contexts"})

	//globalFlags.StringP("sense", "", "", "URL to Qlik Sense")
	//globalFlags.StringP("key", "", "", "API Key to Qlik Sense")

	//TODO sort out how to handle this
	globalFlags.StringP("catwalk-web-integration-id", "", "", "Web integration id for catwalk")
	globalFlags.Lookup("catwalk-web-integration-id").Hidden = true

	if runtime.GOOS != "windows" {
		// Do not add bash completion annotations for paths and files as they are not compatible with windows. On windows
		// we instead rely on the default bash behavior
		globalFlags.SetAnnotation("config", cobra.BashCompFilenameExt, []string{"yaml", "yml"})
	}

	// Add all global flags to the set of valid config properties.
	globalFlags.VisitAll(func(flag *pflag.Flag) {
		dynconf.AddValidConfigFilePropertyName(flag.Name)
	})
}
