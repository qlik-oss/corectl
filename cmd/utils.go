package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/browser"
	"github.com/qlik-oss/corectl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Reloads and saves the app after updating connections, dimensions, measures, objects and the script",
	Long: `Builds the app. Example: corectl build --connections ./myconnections.yml --script ./myscript.qvs
		
`,
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
		viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
		viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
		viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
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
		viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
		viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
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
		viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
		viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
		viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Expected at least one dimension or measure")
			ccmd.Usage()
			os.Exit(1)
		}
		state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
		internal.Eval(rootCtx, state.Doc, args)
	},
}

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reloads the app.",
	Long:  "Reloads the app. Example: corectl reload",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(rootCmd, args)
		viper.BindPFlag("app", ccmd.PersistentFlags().Lookup("app"))
		viper.BindPFlag("engine", ccmd.PersistentFlags().Lookup("engine"))
		viper.BindPFlag("ttl", ccmd.PersistentFlags().Lookup("ttl"))
		viper.BindPFlag("silent", ccmd.PersistentFlags().Lookup("silent"))
		viper.BindPFlag("no-save", ccmd.PersistentFlags().Lookup("no-save"))
	},

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, false)
		silent := viper.GetBool("silent")

		internal.Reload(rootCtx, state.Doc, state.Global, silent, true)

		if state.AppID != "" && !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc, state.AppID)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of corectl",

	Run: func(_ *cobra.Command, args []string) {
		fmt.Printf("corectl version %s\n", version)
	},
}

var generateDocsCommand = &cobra.Command{
	Use:    "generate-docs",
	Short:  "Generate markdown docs based on cobra commands",
	Long:   "Generate markdown docs based on cobra commands",
	Hidden: true,

	Run: func(ccmd *cobra.Command, args []string) {
		fmt.Println("Generating documentation")
		doc.GenMarkdownTree(rootCmd, "./docs")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(catwalkCmd)
	rootCmd.AddCommand(evalCmd)
	rootCmd.AddCommand(generateDocsCommand)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(reloadCmd)
}

func build(ccmd *cobra.Command, args []string) {
	ctx := rootCtx
	state := internal.PrepareEngineState(ctx, viper.GetString("engine"), viper.GetString("app"), viper.GetString("ttl"), headers, true)

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
