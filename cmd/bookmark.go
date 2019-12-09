package cmd

import (
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setBookmarksCmd = withLocalFlags(&cobra.Command{
	Use:     "set <glob-pattern-path-to-bookmark-files.json>",
	Args:    cobra.ExactArgs(1),
	Short:   "Set or update the bookmarks in the current app",
	Long:    "Set or update the bookmarks in the current app",
	Example: "corectl bookmark set ./my-bookmarks-glob-path.json",
	Hidden:  true,

	Run: func(ccmd *cobra.Command, args []string) {
		commandLineBookmarks := args[0]
		if commandLineBookmarks == "" {
			log.Fatalln("no bookmarks specified")
		}
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, true, false)
		internal.SetBookmarks(rootCtx, state.Doc, commandLineBookmarks)
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var removeBookmarkCmd = withLocalFlags(&cobra.Command{
	Use:     "rm <bookmark-id>...",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove one or many bookmarks in the current app",
	Long:    "Remove one or many bookmarks in the current app",
	Example: "corectl dimension rm ID-1",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		for _, entity := range args {
			destroyed, err := state.Doc.DestroyBookmark(rootCtx, entity)
			if err != nil {
				log.Fatalf("could not remove generic bookmark '%s': %s\n", entity, err)
			} else if !destroyed {
				log.Fatalf("could not remove generic bookmark '%s'\n", entity)
			}
		}
		if !viper.GetBool("no-save") {
			internal.Save(rootCtx, state.Doc)
		}
	},
}, "no-save")

var listBookmarksCmd = &cobra.Command{
	Use:     "ls",
	Args:    cobra.ExactArgs(0),
	Short:   "Print a list of all generic bookmarks in the current app",
	Long:    "Print a list of all generic bookmarks in the current app",
	Example: "corectl bookmark ls",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		items := internal.ListBookmarks(state.Ctx, state.Doc)
		printer.PrintNamedItemsList(items, viper.GetBool("bash"), false)
	},
}

var getBookmarkPropertiesCmd = withLocalFlags(&cobra.Command{
	Use:     "properties <bookmark-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Print the properties of the generic bookmark",
	Long:    "Print the properties of the generic bookmark",
	Example: "corectl bookmark properties BOOKMARK-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityProperties(state, args[0], "bookmark", viper.GetBool("minimum"))
	},
}, "minimum")

var getBookmarkLayoutCmd = &cobra.Command{
	Use:     "layout <bookmark-id>",
	Args:    cobra.ExactArgs(1),
	Short:   "Evaluate the layout of an generic bookmark",
	Long:    "Evaluate the layout of an generic bookmark",
	Example: "corectl bBookmark layout BOOKMARK-ID",

	Run: func(ccmd *cobra.Command, args []string) {
		state := internal.PrepareEngineState(rootCtx, headers, tlsClientConfig, false, false)
		printer.PrintGenericEntityLayout(state, args[0], "bookmark")
	},
}

var bookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "Explore and manage bookmarks",
	Long:  "Explore and manage bookmarks",
	Annotations: map[string]string{
		"command_category": "sub",
		"x-qlik-stability": "experimental",
	},
}

func init() {
	bookmarkCmd.AddCommand(setBookmarksCmd, removeBookmarkCmd, listBookmarksCmd, getBookmarkPropertiesCmd, getBookmarkLayoutCmd)
}
