package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	yamlcmd "gitlab.com/gitlab-com/gl-infra/jsonnet-render/internal/cmd/yaml"
	"gopkg.in/yaml.v2"
)

var jpaths []string
var multiDir string
var priorityKeys []string
var header string
var filenamePrefix string

func init() {
	rootCmd.AddCommand(yamlCommand)
	yamlCommand.PersistentFlags().StringArrayVarP(&jpaths, "jpath", "J", nil, "Specify an additional library search dir")
	yamlCommand.PersistentFlags().StringArrayVarP(&priorityKeys, "priority-keys", "P", nil, "Order these keys first in YAML output")
	yamlCommand.PersistentFlags().StringVarP(&multiDir, "multi", "m", ".", "Write multiple files to the directory, list files on stdout")
	yamlCommand.PersistentFlags().StringVarP(&header, "header", "H", "", "Write header to each file")
	yamlCommand.PersistentFlags().StringVarP(&filenamePrefix, "prefix", "p", "", "Prefix to append to every emitted file")
}

func handleYAMLFile(k string, data string) error {
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return err
	}

	filePath := path.Join(multiDir, k)
	fileDir := path.Dir(filePath)
	fileBase := path.Base(filePath)
	filePathWithPrefix := path.Join(fileDir, filenamePrefix+fileBase)

	err = os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(filePathWithPrefix)
	if err != nil {
		return err
	}
	defer f.Close()

	if header != "" {
		_, err = f.WriteString(header)
		if err != nil {
			return err
		}
	}

	ordered := yamlcmd.ReorderKeys(m, priorityKeys)
	encoder := yaml.NewEncoder(f)
	encoder.Encode(ordered)

	fmt.Println(filePath)

	return nil
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
			JPaths: jpaths,
		})

		files, err := vm.EvaluateFileMulti(args[0])
		if err != nil {
			return err
		}

		for k, data := range files {
			err = handleYAMLFile(k, data)
			if err != nil {
				return err
			}

		}

		return nil
	},
}
