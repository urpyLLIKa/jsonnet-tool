package manitest

type baseVisitor struct{}

var _ TestVisitor = &baseVisitor{}

func (c *baseVisitor) TestFileStarted(fileName string) error { return nil }

func (c *baseVisitor) TestFileCompleted(fileName string, allSuccessful bool) error { return nil }

func (c *baseVisitor) TestCaseManifestationStarted(fileName string, testcase string) error {
	return nil
}

func (c *baseVisitor) TestCaseManifestationCompleted(fileName string, testcase string) error {
	return nil
}

func (c *baseVisitor) TestCaseEvaluationCompleted(fileName string, testcase string, result *TestCaseResult) error {
	return nil
}

func (c *baseVisitor) TestCaseEvaluationDelta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	return nil
}

func (c *baseVisitor) CachedTestCaseResultLookup(fileName string) (*TestCaseResult, error) {
	return nil, nil
}

func (c *baseVisitor) AllTestsCompleted() error { return nil }
