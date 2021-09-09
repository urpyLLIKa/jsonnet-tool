package render

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Render a map as a YAML file
func YAMLMapData(filenameKey string, data map[string]interface{}, options Options) error {
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

	encoder := yaml.NewEncoder(f)
	encoder.Encode(data)

	fmt.Println(filePath)

	return nil
}
