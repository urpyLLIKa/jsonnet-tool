package manitest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"slices"

	"github.com/google/go-jsonnet"
)

type CacheResult struct {
	Success bool   `json:"success"`
	Hash    string `json:"hash"`
}

type CacheResults map[string]*CacheResult

type CacheManager struct {
	vm           *jsonnet.VM
	cacheResults CacheResults
	hashCache    map[string]string
}

func (c *CacheManager) LoadCachedResults() error {
	jsonFile, err := os.Open(".jsonnet-tool-test-cache.json")
	if err != nil {
		return fmt.Errorf("failed to open cache results: %w", err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result CacheResults

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return fmt.Errorf("failed to load cache results: %w", err)
	}

	if result != nil {
		c.cacheResults = result
	}

	return nil
}

func (c *CacheManager) GetCachedResult(fileName string) (*bool, error) {
	result, ok := c.cacheResults[fileName]
	if !ok {
		return nil, nil
	}

	hash, err := c.getHash(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	if result.Hash == hash {
		return &result.Success, nil
	}

	return nil, nil
}

func (c *CacheManager) SaveCachedResults() error {
	b, err := json.Marshal(c.cacheResults)
	if err != nil {
		return fmt.Errorf("unable to marshall cache file %w", err)
	}

	err = os.WriteFile(".jsonnet-tool-test-cache.json", b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *CacheManager) RecordResult(fileName string, success bool) error {
	hash, err := c.getHash(fileName)
	if err != nil {
		return err
	}

	c.cacheResults[fileName] = &CacheResult{Success: success, Hash: hash}

	return nil
}

func (c *CacheManager) getHash(fileName string) (string, error) {
	hash, ok := c.hashCache[fileName]
	if ok {
		return hash, nil
	}

	hash, err := c.calculateHashSum(fileName)
	if err != nil {
		return "", err
	}

	c.hashCache[fileName] = hash

	return hash, nil
}

// calculateHashSum generates a unique hash based on the content of all files
// used in the test, including jsonnet, imports, test fixtures.
func (c *CacheManager) calculateHashSum(fileName string) (string, error) {
	deps, err := c.listAllDependencies(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate manitest: %w", err)
	}

	hash := sha256.New()
	for _, fileName := range deps {
		err = addFileForHashing(hash, fileName)
		if err != nil {
			return "", fmt.Errorf("failed to hash: %s: %w", fileName, err)
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// listAllDependencies function will inspect a test and return a stable set of all unique dependencies
// for that test file.
func (c *CacheManager) listAllDependencies(fileName string) ([]string, error) {
	results := map[string]struct{}{}

	results[fileName] = struct{}{}

	// Remove the "actual" value for much faster evaluation
	testManifest, err := c.vm.EvaluateAnonymousSnippet(fileName, "(import '"+fileName+"') { actual:: null }")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate manitest: %w", err)
	}

	testResults := TestCases{}

	err = json.Unmarshal([]byte(testManifest), &testResults)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate jsonnet: %w: %w", err, errSetupTestFailed)
	}

	dir := path.Dir(fileName)

	for _, v := range testResults {
		if v.ExpectJSON != nil {
			fixturePathJSON := path.Join(dir, *v.ExpectJSON)
			results[fixturePathJSON] = struct{}{}
		}

		if v.ExpectYAML != nil {
			fixturePathYAML := path.Join(dir, *v.ExpectYAML)

			results[fixturePathYAML] = struct{}{}
		}
	}

	deps, err := c.vm.FindDependencies("", []string{fileName})
	if err != nil {
		return nil, fmt.Errorf("failed to find dependencies: %s: %w", fileName, err)
	}

	for _, dep := range deps {
		results[dep] = struct{}{}
	}

	// Extract unique results into a slice
	i := 0
	uniqueResults := make([]string, len(results))

	for dep := range results {
		uniqueResults[i] = dep
		i = i + 1
	}

	// Sort the slice for consistency
	slices.Sort(uniqueResults)

	return uniqueResults, nil
}

func addFileForHashing(h hash.Hash, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", fileName, err)
	}

	defer file.Close()

	_, err = io.Copy(h, file)
	if err != nil {
		return fmt.Errorf("failed to generate hash: %w", err)
	}

	return nil
}

func NewCacheManager(vm *jsonnet.VM) *CacheManager {
	return &CacheManager{
		vm:           vm,
		cacheResults: CacheResults{},
		hashCache:    map[string]string{},
	}
}
