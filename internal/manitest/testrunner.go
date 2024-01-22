package manitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

var errSetupTestFailed = errors.New("unable to execute test")
var errTestFailed = errors.New("test failed")

type TestRunner struct {
	vm      *jsonnet.VM
	visitor TestVisitor
}

const runTestsSnippet = `
	local ts = import '%s';

	local testStart = std.native('testStart');
	local testCompleted = std.native('testCompleted');

	std.foldl(
		function(memo, k)
			local f = ts[k];
			memo {
				[k]: testStart(k, {}) + testCompleted(k, f())
			},
		std.objectFields(ts),
		{}
	)
`

func (c *TestRunner) RegisterNatives() {
	c.vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "testStart",
		Params: ast.Identifiers{"testName", "passthrough"},
		Func: func(s []interface{}) (interface{}, error) {
			testName, ok := s[0].(string)
			if !ok {
				return nil, fmt.Errorf("testStart requires a string: %w", errTestFailed)
			}

			err := c.visitor.TestCaseManifestationStarted("", testName)
			if err != nil {
				return nil, fmt.Errorf("StartTestCaseEvaluation visitor failed: %w", errTestFailed)
			}

			return s[1], nil
		},
	})

	c.vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "testCompleted",
		Params: ast.Identifiers{"testName", "testResult"},
		Func: func(s []interface{}) (interface{}, error) {
			testName, ok := s[0].(string)
			if !ok {
				return nil, fmt.Errorf("testStart requires a string: %w", errTestFailed)
			}

			err := c.visitor.TestCaseManifestationCompleted("", testName)
			if err != nil {
				return nil, fmt.Errorf("CompletedTestCaseEvaluation visitor failed: %w", errTestFailed)
			}

			return s[1], nil
		},
	})
}

func (c *TestRunner) RunTestFile(fileName string) error {
	cachedResult, err := c.visitor.CachedTestCaseResultLookup(fileName)
	if err != nil {
		fmt.Printf("warning: cache lookup failed: %v", err)
	}

	err = c.visitor.TestFileStarted(fileName)
	if err != nil {
		log.Printf("warning: %v", err)
	}

	allSuccessful := true

	defer func() {
		err2 := c.visitor.TestFileCompleted(fileName, allSuccessful)
		if err2 != nil {
			log.Printf("warning: %v", err)
		}
	}()

	// Only skip successfully cached results...
	if cachedResult != nil && cachedResult.Success {
		_ = c.visitor.TestCaseManifestationStarted(fileName, "")
		_ = c.visitor.TestCaseEvaluationCompleted(fileName, "", cachedResult)

		return nil
	}

	testManifest, err := c.vm.EvaluateAnonymousSnippet("testrunner.go", fmt.Sprintf(runTestsSnippet, fileName))
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
	err := c.visitor.TestCaseManifestationStarted(fileName, testcase)
	if err != nil {
		return false, fmt.Errorf("visitor failed: %w", err)
	}

	result := c.evaluateTestCaseType(fileName, testcase, t)

	err = c.visitor.TestCaseEvaluationCompleted(fileName, testcase, result)
	if err != nil {
		return false, fmt.Errorf("visitor failed: %w", err)
	}

	return result.Success, nil
}

func (c *TestRunner) evaluateTestCaseType(fileName string, testcase string, t *TestCase) *TestCaseResult {
	if t.ExpectJSON != nil {
		return c.evaluateTestCaseJSON(fileName, testcase, t)
	}

	if t.ExpectYAML != nil {
		return c.evaluateTestCaseYAML(fileName, testcase, t)
	}

	if t.ExpectPlainText != nil {
		return c.evaluateTestCasePlainText(fileName, testcase, t)
	}

	if t.Expect != nil {
		return c.evaluateTestCaseValue(fileName, testcase, t)
	}

	return testCaseResultForError(fmt.Errorf("malformed test expectation: %w", errSetupTestFailed))
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
