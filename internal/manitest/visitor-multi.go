package manitest

import (
	"fmt"
)

type MultiVisitor struct {
	Visitors []TestVisitor
}

var _ TestVisitor = &MultiVisitor{}

func (mv *MultiVisitor) StartTestFile(fileName string) error {
	for _, v := range mv.Visitors {
		err := v.StartTestFile(fileName)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) TestFileComplete(fileName string, allSuccessful bool) error {
	for _, v := range mv.Visitors {
		err := v.TestFileComplete(fileName, allSuccessful)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) StartTestCase(fileName string, testcase string) error {
	for _, v := range mv.Visitors {
		err := v.StartTestCase(fileName, testcase)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) TestCaseComplete(fileName string, testcase string, result *TestCaseResult) error {
	for _, v := range mv.Visitors {
		err := v.TestCaseComplete(fileName, testcase, result)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) Delta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	for _, v := range mv.Visitors {
		err := v.Delta(fileName, testcase, fixturePath, canonicalActual, canonicalExpected)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) CachedResult(fileName string) (*TestCaseResult, error) {
	for _, v := range mv.Visitors {
		result, err := v.CachedResult(fileName)
		if err != nil {
			return nil, fmt.Errorf("visitor failed: %w", err)
		}

		// Return the first non-nil result
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func (mv *MultiVisitor) Complete() error {
	for _, v := range mv.Visitors {
		err := v.Complete()
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}
