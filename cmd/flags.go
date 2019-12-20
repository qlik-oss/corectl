package cmd

import (
	"fmt"
	"runtime"

	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var localFlags pflag.FlagSet
var initialized bool

// DefaultUnbuildFolder is the placeholder for unbuild folder location
var DefaultUnbuildFolder = "./<app name>-unbuild"

// getPathFlagFromConfigFile returns a parameter from the config file.
// It modifies the parameter to actually be relative to the config file and not the working directory
func getPathFlagFromConfigFile(paramName string) string {
	pathInConfigFile := viper.GetString(paramName)
	if pathInConfigFile != "" {
		return internal.RelativeToProject(pathInConfigFile)
	}
	return ""
}

func withLocalFlags(ccmd *cobra.Command, localFlagNames ...string) *cobra.Command {
	if !initialized {
		initLocalFlags()
		initialized = true
	}
	for _, flagName := range localFlagNames {
		flag := localFlags.Lookup(flagName)
		if flag != nil {
			ccmd.PersistentFlags().AddFlag(flag)
		} else {
			fmt.Println("Unknown flag:", flagName)
			panic("")
		}
	}
	return ccmd
}

func initGlobalFlags(globalFlags *pflag.FlagSet) {
	// bound to viper
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

	globalFlags.VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag(flag.Name, flag)
	})

	// Not bound to viper. Certificates are handled similarly to config as the path might be relative.
	globalFlags.StringVarP(&explicitConfigFile, "config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")
	globalFlags.StringToStringVar(&headersMap, "headers", nil, "Http headers to use when connecting to Qlik Associative Engine")
	globalFlags.StringVar(&explicitCertificatePath, "certificates", "", "path/to/folder containing client.pem, client_key.pem and root.pem certificates")

	// Set annotation to run bash completion function
	globalFlags.SetAnnotation("app", cobra.BashCompCustom, []string{"__corectl_get_apps"})
	globalFlags.SetAnnotation("engine", cobra.BashCompCustom, []string{"__corectl_get_local_engines"})
	globalFlags.SetAnnotation("context", cobra.BashCompCustom, []string{"__corectl_get_contexts"})

	if runtime.GOOS != "windows" {
		// Do not add bash completion annotations for paths and files as they are not compatible with windows. On windows
		// we instead rely on the default bash behavior
		globalFlags.SetAnnotation("config", cobra.BashCompFilenameExt, []string{"yaml", "yml"})
	}

	// Add all global flags to the set of valid config properties.
	globalFlags.VisitAll(func(flag *pflag.Flag) {
		internal.AddValidProp(flag.Name)
	})
}

func initLocalFlags() {
	//bound to viper
	localFlags.Bool("no-save", false, "Do not save the app")
	localFlags.Bool("silent", false, "Do not log reload output")
	localFlags.Int("limit", 0, "Limit the number of rows to load")
	localFlags.Bool("no-reload", false, "Do not run the reload script")
	localFlags.Bool("suppress", false, "Suppress confirmation dialogue")
	localFlags.String("catwalk-url", "https://catwalk.core.qlik.com", "Url to an instance of catwalk, if not provided the qlik one will be used")
	localFlags.Bool("minimum", false, "Only print properties required by engine")
	localFlags.Bool("full", false, "Using 'GetFullPropertyTree' to retrieve properties for children as well")
	localFlags.String("comment", "", "Comment for the context")
	localFlags.BoolP("quiet", "q", false, "Only print IDs. Useful for scripting")
	localFlags.String("user", "", "Username to be used when logging in to Qlik Sense Enterprise")
	localFlags.String("password", "", "Password to be used when logging in to Qlik Sense Enterprise (use with caution)")

	localFlags.VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag(flag.Name, flag)
	})

	// not bound to viper
	// Don't bind these to viper since paths are treated separately to support relative paths!
	localFlags.String("connections", "", "Path to a yml file containing the data connection definitions")
	localFlags.String("dimensions", "", "A list of generic dimension json paths")
	localFlags.String("variables", "", "A list of generic variable json paths")
	localFlags.String("bookmarks", "", "A list of generic bookmark json paths")
	localFlags.String("measures", "", "A list of generic measures json paths")
	localFlags.String("objects", "", "A list of generic object json paths")
	localFlags.String("script", "", "Path to a qvs file containing the app data reload script")
	localFlags.String("app-properties", "", "Path to a json file containing the app properties")
	localFlags.String("dir", DefaultUnbuildFolder, "Path to a the folder where the unbuilt app is exported")

	if runtime.GOOS != "windows" {
		// Set annotation to run bash completion function
		// Do not add bash completion annotations for paths and files as they are not compatible with windows. On windows
		// we instead rely on the default bash behavior
		localFlags.SetAnnotation("connections", cobra.BashCompFilenameExt, []string{"yml", "yaml"})
		localFlags.SetAnnotation("dimensions", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("variables", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("bookmarks", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("measures", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("objects", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("script", cobra.BashCompFilenameExt, []string{"qvs"})
	}

	// Add all local flags to the set of valid config properties.
	localFlags.VisitAll(func(flag *pflag.Flag) {
		internal.AddValidProp(flag.Name)
	})
}
