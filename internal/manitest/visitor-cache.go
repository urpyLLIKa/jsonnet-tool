package manitest

import (
	"fmt"

	"github.com/google/go-jsonnet"
)

type CacheVisitor struct {
	vm           *jsonnet.VM
	cacheManager *CacheManager
}

var _ TestVisitor = &CacheVisitor{}

func (cv *CacheVisitor) StartTestFile(fileName string) error {
	return nil
}

func (cv *CacheVisitor) TestFileComplete(fileName string, allSuccessful bool) error {
	err := cv.cacheManager.RecordResult(fileName, allSuccessful)
	if err != nil {
		return fmt.Errorf("failed to record cache result: %w", err)
	}

	return nil
}

func (cv *CacheVisitor) StartTestCase(fileName string, testcase string) error {
	return nil
}

func (cv *CacheVisitor) TestCaseComplete(fileName string, testcase string, result *TestCaseResult) error {
	return nil
}

func (cv *CacheVisitor) Delta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	return nil
}

func (cv *CacheVisitor) CachedResult(fileName string) (*TestCaseResult, error) {
	result, err := cv.cacheManager.GetCachedResult(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache result: %w", err)
	}

	if result != nil {
		cachedResult := &TestCaseResult{
			Success:     *result,
			Cached:      true,
			Error:       nil,
			FixturePath: "",
			Actual:      "",
			Expected:    "",
		}

		return cachedResult, nil
	}

	return nil, nil
}

func (cv *CacheVisitor) Complete() error { return nil }

func NewCacheVisitor(vm *jsonnet.VM, cacheManager *CacheManager) *CacheVisitor {
	return &CacheVisitor{
		vm:           vm,
		cacheManager: cacheManager,
	}
}
