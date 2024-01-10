package manitest

import "os"

type ExitCodeVisitor struct {
	Total       int
	Success     int
	Failures    int
	HasFailures bool
}

var _ TestVisitor = &ExitCodeVisitor{}

func (av *ExitCodeVisitor) StartTestFile(fileName string) error                  { return nil }
func (av *ExitCodeVisitor) StartTestCase(fileName string, testcase string) error { return nil }

func (av *ExitCodeVisitor) TestFileComplete(fileName string, allSuccessful bool) error {
	if !allSuccessful {
		av.HasFailures = true
	}

	return nil
}

func (av *ExitCodeVisitor) TestCaseComplete(fileName string, testcase string, result *TestCaseResult) error {
	return nil
}

func (av *ExitCodeVisitor) Delta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	return nil
}

func (av *ExitCodeVisitor) CachedResult(fileName string) (*TestCaseResult, error) {
	return nil, nil
}

func (av *ExitCodeVisitor) Complete() error {
	if av.HasFailures {
		os.Exit(1)
	}

	return nil
}
