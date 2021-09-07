package render

import (
	"encoding/json"
	"fmt"
	"os"
)

func writeStringData(f *os.File, data string) error {
	_, err := f.Write([]byte(data))
	return err
}

func writeJSONData(f *os.File, data interface{}) error {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(marshalled))
	return err
}

// JSONData renders data, either as a string or in JSON format
func JSONData(filenameKey string, data interface{}, options Options) error {
	f, filePath, err := openFileForRender(filenameKey, options)
	if err != nil {
		return err
	}
	defer f.Close()

	switch v := data.(type) {
	case string:
		err = writeStringData(f, v)
		if err != nil {
			return err
		}

	default:
		err = writeJSONData(f, v)
		if err != nil {
			return err
		}
	}

	fmt.Println(filePath)

	return nil
}
