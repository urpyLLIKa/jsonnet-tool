package main

import (
	"errors"
	"os"

	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/exitcode"

	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		var errWithExitCode *exitcode.Error
		if errors.As(err, &errWithExitCode) {
			os.Exit(errWithExitCode.ExitCode)
		}

		os.Exit(2)
	}
}
