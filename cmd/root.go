package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func silenceErrorsUsage(cmd *cobra.Command, args []string) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
}

var rootCmd = &cobra.Command{
	Use:              "jsonnet-tool",
	Short:            "A tool for rendering jsonnet",
	PersistentPreRun: silenceErrorsUsage,
}

// Execute executes the root command.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("execution failed %w: %w", err, errCommandFailed)
	}

	return nil
}
