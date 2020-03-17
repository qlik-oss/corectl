package engine

import (
	"fmt"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var localFlags pflag.FlagSet
var initialized bool

// DefaultUnbuildFolder is the placeholder for unbuild folder location
var DefaultUnbuildFolder = "./<app name>-unbuild"

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
		dynconf.AddValidConfigFilePropertyName(flag.Name)
	})
}
