package manitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/google/go-jsonnet"
)

var errSetupTestFailed = errors.New("unable to execute test")
var errTestFailed = errors.New("test failed")

type TestRunner struct {
	vm      *jsonnet.VM
	visitor TestVisitor
}

func (c *TestRunner) RunTest(fileName string) error {
	cachedResult, err := c.visitor.CachedResult(fileName)
	if err != nil {
		fmt.Printf("warning: cache lookup failed: %v", err)
	}

	err = c.visitor.StartTestFile(fileName)
	if err != nil {
		log.Printf("warning: %v", err)
	}

	allSuccessful := true

	defer func() {
		err2 := c.visitor.TestFileComplete(fileName, allSuccessful)
		if err2 != nil {
			log.Printf("warning: %v", err)
		}
	}()

	// Only skip successfully cached results...
	if cachedResult != nil && cachedResult.Success {
		_ = c.visitor.StartTestCase(fileName, "")
		_ = c.visitor.TestCaseComplete(fileName, "", cachedResult)

		return nil
	}

	testManifest, err := c.vm.EvaluateFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to evaluate jsonnet: %w: %w", err, errSetupTestFailed)
	}

	testResults := TestCases{}

	err = json.Unmarshal([]byte(testManifest), &testResults)
	if err != nil {
		return fmt.Errorf("failed to evaluate jsonnet: %w: %w", err, errSetupTestFailed)
	}

	sortedKeys := getSortedKeys(testResults)

	for _, testcase := range sortedKeys {
		t := testResults[testcase]

		success, err := c.evaluateTestCase(fileName, testcase, t)
		if !success {
			allSuccessful = false
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Return the keys from the map, sorted.
func getSortedKeys(testResults TestCases) []string {
	sortedKeys := make([]string, len(testResults))
	i := 0

	for k := range testResults {
		sortedKeys[i] = k
		i = i + 1
	}

	slices.Sort(sortedKeys)

	return sortedKeys
}

func (c *TestRunner) evaluateTestCase(fileName string, testcase string, t *TestCase) (bool, error) {
	err := c.visitor.StartTestCase(fileName, testcase)
	if err != nil {
		return false, fmt.Errorf("visitor failed: %w", err)
	}

	var result *TestCaseResult
	if t.ExpectJSON != nil {
		result = c.evaluateTestCaseJSON(fileName, testcase, t)
	} else if t.ExpectYAML != nil {
		result = c.evaluateTestCaseYAML(fileName, testcase, t)
	} else {
		result = testCaseResultForError(fmt.Errorf("malformed test expectation: %w", errSetupTestFailed))
	}

	err = c.visitor.TestCaseComplete(fileName, testcase, result)
	if err != nil {
		return false, fmt.Errorf("visitor failed: %w", err)
	}

	return result.Success, nil
}

func testCaseResultForError(err error) *TestCaseResult {
	return &TestCaseResult{
		Success: false,
		Error:   err,
	}
}

func NewTestRunner(vm *jsonnet.VM, visitor TestVisitor) *TestRunner {
	return &TestRunner{vm, visitor}
}
