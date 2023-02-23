package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var errRenderFailure = errors.New("render failed")

func writeStringData(f *os.File, data string) error {
	_, err := f.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("write failed: %w: %w", err, errRenderFailure)
	}

	return nil
}

func writeJSONData(f *os.File, data interface{}) error {
	marshalled, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal failed: %w: %w", err, errRenderFailure)
	}

	_, err = f.Write(marshalled)
	if err != nil {
		return fmt.Errorf("write failed: %w: %w", err, errRenderFailure)
	}

	return nil
}

// JSONData renders data, either as a string or in JSON format.
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
			return fmt.Errorf("failed to write string data: %w: %w", err, errRenderFailure)
		}

	default:
		err = writeJSONData(f, v)
		if err != nil {
			return fmt.Errorf("failed to write JSON data: %w: %w", err, errRenderFailure)
		}
	}

	fmt.Println(filePath)

	return nil
}
