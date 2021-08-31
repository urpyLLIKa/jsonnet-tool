package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(renderCommand)
	renderCommand.PersistentFlags().StringArrayVarP(&jpaths, "jpath", "J", nil, "Specify an additional library search dir")
	renderCommand.PersistentFlags().StringVarP(&multiDir, "multi", "m", ".", "Write multiple files to the directory, list files on stdout")
	renderCommand.PersistentFlags().StringVarP(&filenamePrefix, "prefix", "p", "", "Prefix to append to every emitted file")
}

func handleRenderFile(k string, data interface{}) error {
	switch path.Ext(k) {
	case ".yml":
		mapData := data.(map[string]interface{})
		return handleYAMLData(k, mapData)
	case ".yaml":
		mapData := data.(map[string]interface{})
		return handleYAMLData(k, mapData)
	default:
		return handleDefaultData(k, data)
	}
}

func openFile(filename string) (*os.File, string, error) {
	filePath := path.Join(multiDir, filename)
	fileDir := path.Dir(filePath)
	fileBase := path.Base(filePath)
	filePathWithPrefix := path.Join(fileDir, filenamePrefix+fileBase)

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return nil, "", err
	}

	f, err := os.Create(filePathWithPrefix)
	return f, filePath, err
}

func handleYAMLData(k string, data map[string]interface{}) error {
	f, filePath, err := openFile(k)
	if err != nil {
		return err
	}
	defer f.Close()

	if header != "" {
		_, err = f.WriteString(header + "\n")
		if err != nil {
			return err
		}
	}

	encoder := yaml.NewEncoder(f)
	encoder.Encode(data)

	fmt.Println(filePath)

	return nil
}

func handleDefaultData(k string, data interface{}) error {
	f, filePath, err := openFile(k)
	if err != nil {
		return err
	}
	defer f.Close()

	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(marshalled))
	if err != nil {
		return err
	}

	fmt.Println(filePath)

	return nil
}

var renderCommand = &cobra.Command{
	Use:   "render",
	Short: "Render files from Jsonnet using sensible defaults",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vm := jsonnet.MakeVM()
		vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)

		vm.Importer(&jsonnet.FileImporter{
			JPaths: jpaths,
		})

		jsonData, err := vm.EvaluateFile(args[0])
		if err != nil {
			return err
		}

		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(jsonData), &m)
		if err != nil {
			return err
		}

		for k, data := range m {
			err = handleRenderFile(k, data)
			if err != nil {
				return err
			}

		}

		return nil
	},
}
