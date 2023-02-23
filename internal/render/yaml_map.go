package render

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// YAMLMapData will render a map as a YAML file.
func YAMLMapData(filenameKey string, data map[string]interface{}, options Options) error {
	f, filePath, err := openFileForRender(filenameKey, options)
	if err != nil {
		return err
	}
	defer f.Close()

	if options.Header != "" {
		_, err = f.WriteString(options.Header + "\n")
		if err != nil {
			return fmt.Errorf("write failed: %s: %w", err, errRenderFailure)
		}
	}

	encoder := yaml.NewEncoder(f)

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("encode failure: %s: %w", err, errRenderFailure)
	}

	fmt.Println(filePath)

	return nil
}
