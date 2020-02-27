package standard

import (
	"github.com/qlik-oss/corectl/cmd"
	"github.com/spf13/cobra"
)

func ContextCommand() *cobra.Command {
	return cmd.ContextCommand()
}
func CompletionCommand() *cobra.Command {
	return cmd.CompletionCommand()
}
func StatusCommand() *cobra.Command {
	return cmd.StatusCommand()
}

func GenerateDocsCommand() *cobra.Command {
	return cmd.GenerateDocsCommand()
}
func GenerateSpecCommand() *cobra.Command {
	return cmd.GenerateSpecCommand()
}
