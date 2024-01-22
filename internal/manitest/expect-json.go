package manitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strings"
)

func (c *TestRunner) evaluateTestCaseJSON(fileName string, testcase string, t *TestCase) *TestCaseResult {
	dir := path.Dir(fileName)
	fixturePath := path.Join(dir, *t.ExpectJSON)

	var err error

	// If the "JSON" is a string, unmarshal it
	actual := t.Actual

	actualString, ok := actual.(string)
	if ok {
		actual, err = decodeReader(strings.NewReader(actualString))
		if err != nil {
			return testCaseResultForError(fmt.Errorf("unable to unmarshal JSON from string: %w", err))
		}
	}

	fixtureFile, err := os.Open(fixturePath)
	if err != nil {
		canonicalActual, _ := canonicalJSON(actual)
		_ = c.visitor.TestCaseEvaluationDelta(fileName, testcase, fixturePath, canonicalActual, "")

		return testCaseResultForError(fmt.Errorf("unable to open fixture %s: %w", fixturePath, err))
	}

	defer fixtureFile.Close()

	expectedJSONValue, err := decodeReader(fixtureFile)
	if err != nil {
		// Can't read the expected value...
		canonicalActual, _ := canonicalJSON(actual)
		_ = c.visitor.TestCaseEvaluationDelta(fileName, testcase, fixturePath, canonicalActual, "")

		return testCaseResultForError(fmt.Errorf("unable to parse fixture %s: %w", fixturePath, err))
	}

	if !reflect.DeepEqual(actual, expectedJSONValue) {
		canonicalActual, err := canonicalJSON(actual)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("failed to manifest actual JSON %s: %w", fixturePath, err))
		}

		canonicalExpected, err := canonicalJSON(expectedJSONValue)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("failed to manifest expected JSON %s: %w", fixturePath, err))
		}

		err = c.visitor.TestCaseEvaluationDelta(fileName, testcase, fixturePath, canonicalActual, canonicalExpected)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("visitor failed: %w", err))
		}

		return &TestCaseResult{
			Success:     false,
			Error:       fmt.Errorf("values don't match: %w", errTestFailed),
			FixturePath: fixturePath,
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

func decodeReader(reader io.Reader) (interface{}, error) {
	var v interface{}
	decoder := json.NewDecoder(reader)

	for {
		err := decoder.Decode(&v)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("unable to parse decode YAML: %w", err)
		}
	}

	return v, nil
}

func canonicalJSON(input interface{}) (string, error) {
	bytes, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(bytes), nil
}
