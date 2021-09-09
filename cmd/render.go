package cmd

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/render"
)

var renderCommandJPaths []string
var renderCommandRenderOptions render.Options

func init() {
	rootCmd.AddCommand(renderCommand)
	renderCommand.PersistentFlags().StringArrayVarP(&renderCommandJPaths, "jpath", "J", nil, "Specify an additional library search dir")
	renderCommand.PersistentFlags().StringVarP(&renderCommandRenderOptions.MultiDir, "multi", "m", ".", "Write multiple files to the directory, list files on stdout")
	renderCommand.PersistentFlags().StringVarP(&renderCommandRenderOptions.Header, "header", "H", "", "Write header to each file")
	renderCommand.PersistentFlags().StringVarP(&renderCommandRenderOptions.FilenamePrefix, "prefix", "p", "", "Prefix to append to every emitted file")
}

func handleYAMLFileType(k string, data interface{}) error {
	switch v := data.(type) {
	case string:
		return render.YAMLStringData(k, v, renderCommandRenderOptions)
	case map[string]interface{}:
		return render.YAMLMapData(k, v, renderCommandRenderOptions)
	default:
		return fmt.Errorf("unexpected type in map for key `%v`: %T", k, v)
	}
}

func handleRenderFile(k string, data interface{}) error {
	switch path.Ext(k) {
	case ".yml":
		return handleYAMLFileType(k, data)
	case ".yaml":
		return handleYAMLFileType(k, data)
	default:
		return render.JSONData(k, data, renderCommandRenderOptions)
	}
}

var renderCommand = &cobra.Command{
	Use:   "render",
	Short: "Render files from Jsonnet using sensible defaults",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vm := jsonnet.MakeVM()
		vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)

		vm.Importer(&jsonnet.FileImporter{
			JPaths: yamlCommandJPaths,
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
