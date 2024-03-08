package manitest

import (
	"fmt"
	"io"
	"os"
	"path"
	"reflect"

	yamlfmt "github.com/google/yamlfmt/formatters/basic"
	yaml "gopkg.in/yaml.v2"
)

func (c *TestRunner) evaluateTestCaseYAML(fileName string, testcase string, t *TestCase) *TestCaseResult {
	dir := path.Dir(fileName)
	fixturePath := path.Join(dir, *t.ExpectYAML)

	// If the "YAML" is a string, unmarshal it
	actual := t.Actual

	actualString, ok := actual.(string)
	if ok {
		err := yaml.Unmarshal([]byte(actualString), &actual)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("unable to unmarshal YAML from string: %w", err))
		}
	}

	fixtureFile, err := os.Open(fixturePath)
	if err != nil {
		canonicalActual, _ := canonicalYAML(actual)
		_ = c.visitor.TestCaseEvaluationDelta(fileName, testcase, fixturePath, canonicalActual, "")

		return testCaseResultForError(fmt.Errorf("unable to open fixture %s: %w", fixturePath, err))
	}

	defer fixtureFile.Close()

	var expectedYAMLValue interface{}

	b, err := io.ReadAll(fixtureFile)
	if err != nil {
		return testCaseResultForError(fmt.Errorf("unable to read fixture %s: %w", fixturePath, err))
	}

	err = yaml.Unmarshal(b, &expectedYAMLValue)
	if err != nil {
		return testCaseResultForError(fmt.Errorf("unable to parse fixture %s: %w", fixturePath, err))
	}

	if !reflect.DeepEqual(actual, expectedYAMLValue) {
		canonicalActual, err := canonicalYAML(actual)
		if err != nil {
			return testCaseResultForError(fmt.Errorf("failed to manifest actual YAML %s: %w", fixturePath, err))
		}

		canonicalExpected, err := canonicalYAML(expectedYAMLValue)
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

func canonicalYAML(input interface{}) (string, error) {
	bytes, err := yaml.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	factory := yamlfmt.BasicFormatterFactory{}

	// Configuration from https://github.com/google/yamlfmt/blob/main/docs/config-file.md
	formatter, err := factory.NewFormatter(
		map[string]interface{}{
			"pad_line_comments":         2,
			"retain_line_breaks_single": true,
			"line_ending":               "lf",
			"indent":                    2,
			"indentless_arrays":         false,
			"max_line_length":           80,
		})
	if err != nil {
		return "", fmt.Errorf("failed to create formatter for YAML: %w", err)
	}

	formatted, err := formatter.Format(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to format YAML: %w", err)
	}

	return string(formatted), nil
}
