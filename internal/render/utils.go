package render

import (
	"fmt"
	"os"
	"path"
)

func openFileForRender(filename string, options Options) (*os.File, string, error) {
	filePath := path.Join(options.MultiDir, filename)
	fileDir := path.Dir(filePath)
	fileBase := path.Base(filePath)
	filePathWithPrefix := path.Join(fileDir, options.FilenamePrefix+fileBase)

	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return nil, "", fmt.Errorf("unable to MkdirAll for %s: %w", fileDir, err)
	}

	f, err := os.Create(filePathWithPrefix)
	return f, filePathWithPrefix, err
}
