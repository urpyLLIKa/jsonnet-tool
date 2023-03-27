package cmd

import (
	"crypto/sha256"
	"fmt"
	"sort"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/checksum"
)

var checksumCommandJPaths []string

func init() {
	rootCmd.AddCommand(checksumCommand)
	checksumCommand.PersistentFlags().StringArrayVarP(
		&checksumCommandJPaths, "jpath", "J", nil,
		"Specify an additional library search dir",
	)
}

var checksumCommand = &cobra.Command{
	Use:   "checksum",
	Short: "checksum jsonnet files for caching purposes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		importer := &jsonnet.FileImporter{
			JPaths: checksumCommandJPaths,
		}
		seen := make(map[string][sha256.Size]byte)

		err := checksum.Parse(args[0], importer, seen)
		if err != nil {
			return fmt.Errorf("checksum.Parse: %w", err)
		}

		var lines []string
		for filename, sum := range seen {
			lines = append(lines, fmt.Sprintf("%x  %s", sum, filename))
		}

		sort.Strings(lines)
		for _, line := range lines {
			fmt.Println(line)
		}

		return nil
	},
}
