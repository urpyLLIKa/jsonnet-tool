package render

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	yamlcmd "gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/cmd/yaml"
)

// YAMLStringData will render a string as a YAML file.
func YAMLStringData(filenameKey string, data string, options Options) error {
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w: %w", err, errRenderFailure)
	}

	f, filePath, err := openFileForRender(filenameKey, options)
	if err != nil {
		return fmt.Errorf("open file failed: %w: %w", err, errRenderFailure)
	}

	defer f.Close()

	if options.Header != "" {
		_, err = f.WriteString(options.Header + "\n")
		if err != nil {
			return fmt.Errorf("write failed: %w: %w", err, errRenderFailure)
		}
	}

	ordered := yamlcmd.ReorderKeys(m, options.PriorityKeys)
	encoder := yaml.NewEncoder(f)

	err = encoder.Encode(ordered)
	if err != nil {
		return fmt.Errorf("encode failed: %w: %w", err, errRenderFailure)
	}

	fmt.Println(filePath)

	return nil
}
