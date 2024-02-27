package cmd

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/exitcode"

	"github.com/stretchr/testify/assert"
)

var testTestFixtures = []struct {
	name        string
	args        []string
	deleteCache bool
	exitCode    int
	wantOutput  string
}{
	{
		name:       "test1",
		args:       []string{"../examples/tests/test1.manitest.jsonnet"},
		exitCode:   0,
		wantOutput: "✅ Test suite completed: 1 file tested, 1 file passed\n",
	},
	{
		name:       "test2",
		args:       []string{"../examples/tests/test2.manitest.jsonnet"},
		exitCode:   0,
		wantOutput: "✅ Test suite completed: 1 file tested, 1 file passed\n",
	},
	{
		name:       "test3",
		args:       []string{"../examples/tests/test3.fail.manitest.jsonnet"},
		exitCode:   1,
		wantOutput: "❌ Test suite completed: 1 file tested, 0 files passed, 1 file failed\n",
	},
	{
		name:       "test4",
		args:       []string{"../examples/tests/test4.fail.manitest.jsonnet"},
		exitCode:   1,
		wantOutput: "❌ Test suite completed: 1 file tested, 0 files passed, 1 file failed",
	},
	{
		name:       "test5",
		args:       []string{"../examples/tests/test5.invalid.manitest.jsonnet"},
		exitCode:   3,
		wantOutput: "💥 Test suite completed: 1 file tested, 0 files passed, 1 file invalid\n",
	},
	{
		name:       "multiple_success",
		args:       []string{"../examples/tests/test1.manitest.jsonnet", "../examples/tests/test2.manitest.jsonnet"},
		exitCode:   0,
		wantOutput: "✅ Test suite completed: 2 files tested, 2 files passed\n",
	},
	{
		name:       "partial_failure",
		args:       []string{"../examples/tests/test1.manitest.jsonnet", "../examples/tests/test4.fail.manitest.jsonnet"},
		exitCode:   1,
		wantOutput: "❌ Test suite completed: 2 files tested, 1 file passed, 1 file failed\n",
	},
	{
		name:       "partial_failure_last_pass",
		args:       []string{"../examples/tests/test4.fail.manitest.jsonnet", "../examples/tests/test1.manitest.jsonnet"},
		exitCode:   1,
		wantOutput: "❌ Test suite completed: 2 files tested, 1 file passed, 1 file failed\n",
	},
	{
		name: "successful_failure_invalid",
		args: []string{
			"../examples/tests/test1.manitest.jsonnet",
			"../examples/tests/test4.fail.manitest.jsonnet",
			"../examples/tests/test5.invalid.manitest.jsonnet",
		},
		exitCode:   3,
		wantOutput: "💥 Test suite completed: 3 files tested, 1 file passed, 1 file failed, 1 file invalid\n",
	},
	{
		name:        "cache_success",
		deleteCache: true,
		args:        []string{"--cache", "../examples/tests/test2.manitest.jsonnet"},
		exitCode:    0,
		wantOutput:  "✅ Test suite completed: 1 file tested, 1 file passed\n",
	},
	{
		name:        "cache_successful_failure_invalid",
		deleteCache: true,
		args: []string{
			"--cache",
			"../examples/tests/test1.manitest.jsonnet",
			"../examples/tests/test4.fail.manitest.jsonnet",
			"../examples/tests/test5.invalid.manitest.jsonnet",
		},
		exitCode:   3,
		wantOutput: "💥 Test suite completed: 3 files tested, 1 file passed, 1 file failed, 1 file invalid",
	},
}

// TestTestCommand runs table-driven tests for the testCommand Cobra command.
func TestTestCommand(t *testing.T) {
	t.Parallel()

	for _, tt := range testTestFixtures {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.deleteCache {
				_ = os.Remove(".jsonnet-tool-test-cache.json")
			}

			output, err := executeTestCommand(tt.args)

			if tt.exitCode == 0 {
				require.NoError(t, err)
			} else {
				var errWithExitCode *exitcode.Error
				if errors.As(err, &errWithExitCode) {
					assert.EqualValues(t, tt.exitCode, errWithExitCode.ExitCode)
				} else {
					assert.NoError(t, err, "unexpected error response did not include exit code")
				}
			}

			assert.Contains(t, output, tt.wantOutput)
		})
	}
}

// executeTestCommand executes the testCommand with given arguments and flags, and returns the output.
func executeTestCommand(args []string) (string, error) {
	testCommand := NewTestCommand()

	buf := new(bytes.Buffer)
	testCommand.SetOut(buf)
	testCommand.SetErr(buf)
	testCommand.SetArgs(args)

	err := testCommand.Execute()

	return buf.String(), err
}
