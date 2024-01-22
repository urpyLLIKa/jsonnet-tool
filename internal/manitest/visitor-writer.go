package manitest

import (
	"bufio"
	"fmt"
	"os"
)

type WriterVisitor struct {
	baseVisitor
}

var _ TestVisitor = &WriterVisitor{}

func (rv *WriterVisitor) TestCaseEvaluationDelta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	// No fixture path? skip.
	if fixturePath == "" {
		return nil
	}

	f, err := os.Create(fixturePath)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", fixturePath, err)
	}

	defer f.Close()

	canonicalActual = normalizePlainTextString(canonicalActual)

	w := bufio.NewWriter(f)

	_, err = w.WriteString(canonicalActual)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fixturePath, err)
	}

	err = w.Flush()
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fixturePath, err)
	}

	return nil
}
