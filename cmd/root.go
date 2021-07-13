package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jsonnet-tool",
	Short: "A tool for rendering jsonnet",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
