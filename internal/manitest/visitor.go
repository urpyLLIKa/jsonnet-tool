package manitest

type TestVisitor interface {
	// --- Test File Events ---

	// TestFileStarted event happens when a new manitest test file starts getting processed.
	TestFileStarted(fileName string) error

	// TestFileCompleted event happens when a manitest test file is complete.
	TestFileCompleted(fileName string, allSuccessful bool) error

	// --- Test Case Events ---

	// TestCaseManifestationStarted event happens when a test case begins manifesting.
	TestCaseManifestationStarted(fileName string, testcase string) error

	// TestCaseManifestationCompleted event happens when a test case completes manifesting.
	TestCaseManifestationCompleted(fileName string, testcase string) error

	// TestCaseEvaluationCompleted event happens when the manifested output from a test case is evaluated against it's fixture.
	TestCaseEvaluationCompleted(fileName string, testcase string, result *TestCaseResult) error

	// TestCaseEvaluationDelta event happens when the manifested output doesn't not match the expected fixture.
	TestCaseEvaluationDelta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error

	// --- Misc Events ---

	// CachedTestCaseResultLookup happens when the runner is looking for a cached result.
	CachedTestCaseResultLookup(fileName string) (*TestCaseResult, error)

	// TestSuiteCompleted happens when all test files have completed.
	AllTestsCompleted() error
}
