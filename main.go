package main

import (
	"os"

	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
