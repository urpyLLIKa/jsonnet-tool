package manitest

import "os"

type ExitCodeVisitor struct {
	Total       int
	Success     int
	Failures    int
	HasFailures bool

	baseVisitor
}

var _ TestVisitor = &ExitCodeVisitor{}

func (av *ExitCodeVisitor) TestFileCompleted(fileName string, allSuccessful bool) error {
	if !allSuccessful {
		av.HasFailures = true
	}

	return nil
}

func (av *ExitCodeVisitor) AllTestsCompleted() error {
	if av.HasFailures {
		os.Exit(1)
	}

	return nil
}
