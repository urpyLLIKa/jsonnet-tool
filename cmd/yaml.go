package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/render"
)

var errCommandFailed = errors.New("command failed")

var yamlCommandJPaths []string
var yamlCommandRenderOptions render.Options

func init() {
	rootCmd.AddCommand(yamlCommand)
	yamlCommand.PersistentFlags().StringArrayVarP(
		&yamlCommandJPaths, "jpath", "J", nil,
		"Specify an additional library search dir",
	)
	yamlCommand.PersistentFlags().StringArrayVarP(
		&yamlCommandRenderOptions.PriorityKeys, "priority-keys", "P", nil,
		"Order these keys first in YAML output",
	)
	yamlCommand.PersistentFlags().StringVarP(
		&yamlCommandRenderOptions.MultiDir, "multi", "m", ".",
		"Write multiple files to the directory, list files on stdout",
	)
	yamlCommand.PersistentFlags().StringVarP(
		&yamlCommandRenderOptions.Header, "header", "H", "",
		"Write header to each file",
	)
	yamlCommand.PersistentFlags().StringVarP(
		&yamlCommandRenderOptions.FilenamePrefix, "prefix", "p", "",
		"Prefix to append to every emitted file",
	)
}

var yamlCommand = &cobra.Command{
	Use:   "yaml",
	Short: "Generate YAML from Jsonnet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vm := jsonnet.MakeVM()
		vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)
		vm.StringOutput = true

		vm.Importer(&jsonnet.FileImporter{
			JPaths: yamlCommandJPaths,
		})

		files, err := vm.EvaluateFileMulti(args[0])
		if err != nil {
			return fmt.Errorf("failed to evaluate jsonnet: %w: %w", err, errCommandFailed)
		}

		for k, data := range files {
			err = render.YAMLStringData(k, data, yamlCommandRenderOptions)
			if err != nil {
				return fmt.Errorf("failed to write data: %w: %w", err, errCommandFailed)
			}

		}

		return nil
	},
}
