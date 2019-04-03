package cmd

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var commonFlags pflag.FlagSet
var initialized bool

// GetRelativeParameter returns a parameter from the config file.
// It modifies the parameter to actually be relative to the config file and not the working directory
func GetRelativeParameter(paramName string) string {
	pathInConfigFile := viper.GetString(paramName)
	if pathInConfigFile != "" {
		return internal.RelativeToProject(viper.ConfigFileUsed(), pathInConfigFile)
	}
	return ""
}

func withCommonLocalFlags(ccmd *cobra.Command, localFlagNames ...string) *cobra.Command {
	if !initialized {
		initCommonLocalFlags()
		initialized = true
	}
	for _, flagName := range localFlagNames {
		flag := commonFlags.Lookup(flagName)
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
	//not bound to viper
	rootCmd.PersistentFlags().StringVarP(&explicitConfigFile, "config", "c", "", "path/to/config.yml where parameters can be set instead of on the command line")
	rootCmd.PersistentFlags().StringToStringVar(&headersMap, "headers", nil, "Http headers to use when connecting to Qlik Associative Engine")

	//bound to viper
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Logs extra information")
	rootCmd.PersistentFlags().BoolP("traffic", "t", false, "Log JSON websocket traffic to stdout")
	rootCmd.PersistentFlags().StringP("engine", "e", "localhost:9076", "URL to the Qlik Associative Engine")
	rootCmd.PersistentFlags().StringP("app", "a", "", "App name, if no app is specified a session app is used instead.")
	rootCmd.PersistentFlags().String("ttl", "30", "Qlik Associative Engine session time to live in seconds")
	rootCmd.PersistentFlags().Bool("no-data", false, "Open app without data")
	rootCmd.PersistentFlags().Bool("bash", false, "Bash flag used to adapt output to bash completion format")
	rootCmd.PersistentFlags().MarkHidden("bash")

	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name != "config" && //not binding to viper since this configures viper itself
			flag.Name != "headers" { //not binding to viper since binding a map does not seem to work.
			viper.BindPFlag(flag.Name, flag)
		}
	})
}

func initCommonLocalFlags() {
	commonFlags.Bool("no-save", false, "Do not save the app")
	commonFlags.Bool("silent", false, "Do not log reload output")
	commonFlags.Bool("no-reload", false, "Do not run the reload script")
	commonFlags.Bool("suppress", false, "Suppress confirmation dialogue")
	commonFlags.String("catwalk-url", "https://catwalk.core.qlik.com", "Url to an instance of catwalk, if not provided the qlik one will be used.")

	commonFlags.VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag(flag.Name, flag)
	})
}
