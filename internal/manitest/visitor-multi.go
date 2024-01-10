package manitest

import (
	"fmt"
)

type MultiVisitor struct {
	Visitors []TestVisitor
}

var _ TestVisitor = &MultiVisitor{}

func (mv *MultiVisitor) TestFileStarted(fileName string) error {
	for _, v := range mv.Visitors {
		err := v.TestFileStarted(fileName)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) TestFileCompleted(fileName string, allSuccessful bool) error {
	for _, v := range mv.Visitors {
		err := v.TestFileCompleted(fileName, allSuccessful)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) TestCaseManifestationStarted(fileName string, testcase string) error {
	for _, v := range mv.Visitors {
		err := v.TestCaseManifestationStarted(fileName, testcase)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) TestCaseManifestationCompleted(fileName string, testcase string) error {
	for _, v := range mv.Visitors {
		err := v.TestCaseManifestationCompleted(fileName, testcase)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) TestCaseEvaluationCompleted(fileName string, testcase string, result *TestCaseResult) error {
	for _, v := range mv.Visitors {
		err := v.TestCaseEvaluationCompleted(fileName, testcase, result)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) TestCaseEvaluationDelta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	for _, v := range mv.Visitors {
		err := v.TestCaseEvaluationDelta(fileName, testcase, fixturePath, canonicalActual, canonicalExpected)
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}

func (mv *MultiVisitor) CachedTestCaseResultLookup(fileName string) (*TestCaseResult, error) {
	for _, v := range mv.Visitors {
		result, err := v.CachedTestCaseResultLookup(fileName)
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

func (mv *MultiVisitor) AllTestsCompleted() error {
	for _, v := range mv.Visitors {
		err := v.AllTestsCompleted()
		if err != nil {
			return fmt.Errorf("visitor failed: %w", err)
		}
	}

	return nil
}
