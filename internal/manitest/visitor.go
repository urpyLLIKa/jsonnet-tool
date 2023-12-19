package manitest

type TestVisitor interface {
	StartTestFile(fileName string) error
	TestFileComplete(fileName string, allSuccessful bool) error

	StartTestCase(fileName string, testcase string) error
	TestCaseComplete(fileName string, testcase string, result *TestCaseResult) error
	Delta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error

	CachedResult(fileName string) (*TestCaseResult, error)

	Complete() error
}
