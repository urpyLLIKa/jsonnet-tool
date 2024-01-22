package manitest

import (
	"fmt"
	"io"
	"os"
	"path"
)

func (c *TestRunner) evaluateTestCasePlainText(fileName string, testcase string, t *TestCase) *TestCaseResult {
	dir := path.Dir(fileName)
	fixturePath := path.Join(dir, *t.ExpectPlainText)

	actualString, ok := t.Actual.(string)
	if !ok {
		return testCaseResultForError(fmt.Errorf("actual value must be a string: %w", errTestFailed))
	}

	actualString = normalizePlainTextString(actualString)

	fixtureFile, err := os.Open(fixturePath)
	if err != nil {
		_ = c.visitor.TestCaseEvaluationDelta(fileName, testcase, fixturePath, actualString, "")

		return testCaseResultForError(fmt.Errorf("unable to open fixture %s: %w", fixturePath, err))
	}

	defer fixtureFile.Close()

	b, err := io.ReadAll(fixtureFile)
	if err != nil {
		return testCaseResultForError(fmt.Errorf("unable to read fixture %s: %w", fixturePath, err))
	}

	expectedString := normalizePlainTextString(string(b))

	if actualString != expectedString {
		err = c.visitor.TestCaseEvaluationDelta(fileName, testcase, fixturePath, actualString, expectedString)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("visitor failed: %w", err))
		}

		return &TestCaseResult{
			Success:     false,
			Error:       fmt.Errorf("values don't match: %w", errTestFailed),
			FixturePath: fixturePath,
			Actual:      actualString,
			Expected:    expectedString,
		}
	}

	return &TestCaseResult{
		Success:     true,
		Error:       nil,
		FixturePath: "",
		Actual:      "",
		Expected:    "",
	}
}
