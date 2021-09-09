package render

import (
	"fmt"

	yamlcmd "gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/cmd/yaml"
	"gopkg.in/yaml.v2"
)

// Render a string as a YAML file
func YAMLStringData(filenameKey string, data string, options Options) error {
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return err
	}

	f, filePath, err := openFileForRender(filenameKey, options)
	if err != nil {
		return err
	}
	defer f.Close()

	if options.Header != "" {
		_, err = f.WriteString(options.Header + "\n")
		if err != nil {
			return err
		}
	}

	ordered := yamlcmd.ReorderKeys(m, options.PriorityKeys)
	encoder := yaml.NewEncoder(f)
	encoder.Encode(ordered)

	fmt.Println(filePath)

	return nil
}
