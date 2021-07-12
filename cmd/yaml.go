package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var jpaths []string
var multiDir string

func init() {
	rootCmd.AddCommand(yamlCommand)
	yamlCommand.PersistentFlags().StringArrayVarP(&jpaths, "jpath", "J", nil, "Specify an additional library search dir")
	yamlCommand.PersistentFlags().StringVarP(&multiDir, "multi", "m", ".", "Write multiple files to the directory, list files on stdout")
}

func handleYAMLFile(k string, data string) error {
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return err
	}

	filePath := path.Join(multiDir, "autogenerated-"+k)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString("# HEADER\n")
	if err != nil {
		return err
	}

	encoder := yaml.NewEncoder(f)
	encoder.Encode(m)

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
