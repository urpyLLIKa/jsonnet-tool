package manitest

import (
	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/exitcode"
)

type ExitCodeVisitor struct {
	hasFailures bool
	hasInvalid  bool

	baseVisitor
}

var _ TestVisitor = &ExitCodeVisitor{}

func (e *ExitCodeVisitor) TestFileInvalid(name string, err error) error {
	e.hasInvalid = true
	return nil
}

func (e *ExitCodeVisitor) TestCaseInvalid(name string, testcase string, err error) error {
	e.hasInvalid = true
	return nil
}

func (e *ExitCodeVisitor) TestFileCompleted(fileName string, allSuccessful bool) error {
	if !allSuccessful {
		e.hasFailures = true
	}

	return nil
}

func (e *ExitCodeVisitor) AllTestsCompleted() error {
	if e.hasInvalid {
		return exitcode.Invalid()
	}

	if e.hasFailures {
		return exitcode.Failed()
	}

	return nil
}
