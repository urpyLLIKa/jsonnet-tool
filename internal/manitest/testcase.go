package manitest

type TestCase struct {
	Actual          interface{} `json:"actual"`
	ExpectJSON      *string     `json:"expectJSON"`
	ExpectYAML      *string     `json:"expectYAML"`
	ExpectPlainText *string     `json:"expectPlainText"`
	Expect          interface{} `json:"expect"`
}

type TestCases map[string]*TestCase

type TestCaseResult struct {
	Success     bool
	Cached      bool
	Error       error
	FixturePath string
	Actual      string
	Expected    string
	Trace       string
}
