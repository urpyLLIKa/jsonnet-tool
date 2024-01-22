package manitest

import (
	"fmt"
	"reflect"
)

func (c *TestRunner) evaluateTestCaseValue(fileName string, testcase string, t *TestCase) *TestCaseResult {
	// If the "JSON" is a string, unmarshal it
	actual := t.Actual

	if !reflect.DeepEqual(actual, t.Expect) {
		canonicalActual, err := canonicalJSON(actual)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("failed to manifest actual JSON %w", err))
		}

		canonicalExpected, err := canonicalJSON(t.Expect)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("failed to manifest expected JSON %w", err))
		}

		err = c.visitor.TestCaseEvaluationDelta(fileName, testcase, "", canonicalActual, canonicalExpected)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("visitor failed: %w", err))
		}

		return &TestCaseResult{
			Success:     false,
			Error:       fmt.Errorf("values don't match: %w", errTestFailed),
			FixturePath: "",
			Actual:      canonicalActual,
			Expected:    canonicalExpected,
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
