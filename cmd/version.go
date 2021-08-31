package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = ""
var commit = ""
var date = ""

func init() {
	rootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Version: %v\n", version)
		fmt.Printf("Commit: %v\n", commit)
		fmt.Printf("Date: %v\n", date)

		return nil
	},
}
