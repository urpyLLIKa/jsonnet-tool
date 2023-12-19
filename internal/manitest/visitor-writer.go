package manitest

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type WriterVisitor struct {
}

var _ TestVisitor = &WriterVisitor{}

func (rv *WriterVisitor) StartTestFile(fileName string) error                        { return nil }
func (rv *WriterVisitor) TestFileComplete(fileName string, allSuccessful bool) error { return nil }

func (rv *WriterVisitor) StartTestCase(fileName string, testcase string) error {
	return nil
}

func (rv *WriterVisitor) TestCaseComplete(fileName string, testcase string, result *TestCaseResult) error {
	return nil
}

func (rv *WriterVisitor) Delta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	f, err := os.Create(fixturePath)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", fixturePath, err)
	}

	defer f.Close()

	if !strings.HasSuffix(canonicalActual, "\n") {
		// Always add a newline at the end of the file...
		canonicalActual = canonicalActual + "\n"
	}

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

func (rv *WriterVisitor) CachedResult(fileName string) (*TestCaseResult, error) {
	return nil, nil
}

func (rv *WriterVisitor) Complete() error { return nil }
