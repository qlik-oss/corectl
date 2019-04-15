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

// getPathFlagFromConfigFile returns a parameter from the config file.
// It modifies the parameter to actually be relative to the config file and not the working directory
func getPathFlagFromConfigFile(paramName string) string {
	pathInConfigFile := viper.GetString(paramName)
	if pathInConfigFile != "" {
		return internal.RelativeToProject(viper.ConfigFileUsed(), pathInConfigFile)
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
	globalFlags.StringP("app", "a", "", "App name, if no app is specified a session app is used instead")
	globalFlags.String("ttl", "30", "Qlik Associative Engine session time to live in seconds")
	globalFlags.Bool("no-data", false, "Open app without data")
	globalFlags.Bool("bash", false, "Bash flag used to adapt output to bash completion format")
	globalFlags.MarkHidden("bash")

	globalFlags.VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag(flag.Name, flag)
	})

	// not bound to viper
	globalFlags.StringVarP(&explicitConfigFile, "config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")
	globalFlags.StringToStringVar(&headersMap, "headers", nil, "Http headers to use when connecting to Qlik Associative Engine")

	// Set annotation to run bash completion function
	globalFlags.SetAnnotation("app", cobra.BashCompCustom, []string{"__corectl_get_apps"})
	globalFlags.SetAnnotation("engine", cobra.BashCompCustom, []string{"__corectl_get_local_engines"})

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
	localFlags.Bool("no-reload", false, "Do not run the reload script")
	localFlags.Bool("suppress", false, "Suppress confirmation dialogue")
	localFlags.String("catwalk-url", "https://catwalk.core.qlik.com", "Url to an instance of catwalk, if not provided the qlik one will be used")

	localFlags.VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag(flag.Name, flag)
	})

	// not bound to viper
	// Don't bind these to viper since paths are treated separately to support relative paths!
	localFlags.String("connections", "", "Path to a yml file containing the data connection definitions")
	localFlags.String("dimensions", "", "A list of generic dimension json paths")
	localFlags.String("measures", "", "A list of generic measures json paths")
	localFlags.String("objects", "", "A list of generic object json paths")
	localFlags.String("script", "", "Path to a qvs file containing the app data reload script")

	if runtime.GOOS != "windows" {
		// Set annotation to run bash completion function
		// Do not add bash completion annotations for paths and files as they are not compatible with windows. On windows
		// we instead rely on the default bash behavior
		localFlags.SetAnnotation("connections", cobra.BashCompFilenameExt, []string{"yml", "yaml"})
		localFlags.SetAnnotation("dimensions", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("measures", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("objects", cobra.BashCompFilenameExt, []string{"json"})
		localFlags.SetAnnotation("script", cobra.BashCompFilenameExt, []string{"qvs"})
	}

	// Add all local flags to the set of valid config properties.
	localFlags.VisitAll(func(flag *pflag.Flag) {
		internal.AddValidProp(flag.Name)
	})
}
