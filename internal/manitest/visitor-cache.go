package manitest

import (
	"fmt"
	"log"

	"github.com/google/go-jsonnet"
)

type CacheVisitor struct {
	vm           *jsonnet.VM
	cacheManager *CacheManager

	baseVisitor
}

var _ TestVisitor = &CacheVisitor{}

func (cv *CacheVisitor) TestFileCompleted(fileName string, allSuccessful bool) error {
	err := cv.cacheManager.RecordResult(fileName, allSuccessful)
	if err != nil {
		log.Printf("failed to record cache result: %v", err)
	}

	return nil
}

func (cv *CacheVisitor) CachedTestCaseResultLookup(fileName string) (*TestCaseResult, error) {
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

func NewCacheVisitor(vm *jsonnet.VM, cacheManager *CacheManager) *CacheVisitor {
	return &CacheVisitor{
		vm:           vm,
		cacheManager: cacheManager,
	}
}
