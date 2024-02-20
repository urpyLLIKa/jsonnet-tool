package manitest

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type TraceVisitor struct {
	buf              *bytes.Buffer
	previousTestCase string
	traces           map[string]string
	stderr           io.Writer
	stdout           io.Writer

	baseVisitor
}

var _ TestVisitor = &TraceVisitor{}
var _ io.Writer = &TraceVisitor{}

func (c *TraceVisitor) TestCaseManifestationStarted(fileName string, testcase string) error {
	if c.buf != nil {
		c.traces[c.previousTestCase] = c.buf.String()
	}

	c.buf = &bytes.Buffer{}
	c.previousTestCase = testcase

	return nil
}

func (c *TraceVisitor) TestCaseEvaluationCompleted(fileName string, testcase string, result *TestCaseResult) error {
	result.Trace = c.traces[c.previousTestCase]

	return nil
}

func (c *TraceVisitor) Write(p []byte) (int, error) {
	var b io.Writer
	if c.buf == nil {
		b = os.Stderr
	} else {
		b = c.buf
	}

	n, err := b.Write(p)
	if err != nil {
		return n, fmt.Errorf("write failed: %w", err)
	}

	return n, nil
}

func NewTraceVisitor(stdout io.Writer, stderr io.Writer) *TraceVisitor {
	return &TraceVisitor{
		stderr:           stderr,
		stdout:           stdout,
		buf:              nil,
		previousTestCase: "",
		traces:           map[string]string{},
	}
}
