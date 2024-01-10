package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/manitest"
)

var testCommandJPaths []string
var writeFixtures bool
var cacheResults bool
var wholeDir bool

func init() {
	rootCmd.AddCommand(testCommand)

	testCommand.PersistentFlags().StringArrayVarP(
		&testCommandJPaths, "jpath", "J", nil,
		"Specify an additional library search dir",
	)

	testCommand.PersistentFlags().BoolVarP(
		&writeFixtures, "write-fixtures", "w", false,
		"Automatically write actual values to fixtures",
	)

	testCommand.PersistentFlags().BoolVarP(
		&cacheResults, "cache", "", false,
		"Cache tests for unchanged files to improve test speed",
	)

	testCommand.PersistentFlags().BoolVarP(
		&wholeDir, "dir", "d", false,
		"Run all manitest declarations in directory",
	)
}

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Run jsonnet tests",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		reporterVisitor := &manitest.ReporterVisitor{}
		visitors := []manitest.TestVisitor{
			reporterVisitor,
		}

		if writeFixtures {
			visitors = append(visitors, &manitest.WriterVisitor{})
		}

		vm := jsonnet.MakeVM()
		vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)
		vm.Importer(&jsonnet.FileImporter{
			JPaths: testCommandJPaths,
		})

		var cacheManager *manitest.CacheManager
		if cacheResults {
			cacheManager = manitest.NewCacheManager(vm)
			err := cacheManager.LoadCachedResults()
			if err != nil {
				log.Printf("failed to load cached test results: %v\n", err)
			}

			cacheVisitor := manitest.NewCacheVisitor(vm, cacheManager)
			visitors = append(visitors, cacheVisitor)
		}

		// ExitCodeVisitor should always go last
		// so that it doesn't exit before other visitors
		// have run
		exitCodeVisitor := &manitest.ExitCodeVisitor{}
		visitors = append(visitors, exitCodeVisitor)

		visitor := &manitest.MultiVisitor{Visitors: visitors}
		runner := manitest.NewTestRunner(vm, visitor)

		err := runTests(runner, args)
		if err != nil {
			return fmt.Errorf("test failed: %w", err)
		}

		if cacheManager != nil {
			err = cacheManager.SaveCachedResults()
			if err != nil {
				log.Printf("failed to save cached test results: %v\n", err)
			}
		}

		err = visitor.Complete()
		if err != nil {
			log.Printf("visitor failed: %v\n", err)
		}

		return nil
	},
}

// Given a test runner, run the tests.
func runTests(runner *manitest.TestRunner, args []string) error {
	for _, a := range args {
		if wholeDir {
			err := runDirTests(runner, a)
			if err != nil {
				return fmt.Errorf("failed to run tests in %s: %w", a, err)
			}
		} else {
			err := runner.RunTest(a)

			if err != nil {
				return fmt.Errorf("failed to run tests in %s: %w", a, err)
			}
		}
	}

	return nil
}

func runDirTests(runner *manitest.TestRunner, arg string) error {
	files, err := os.ReadDir(arg)
	if err != nil {
		return fmt.Errorf("failed to list files in %s: %w", arg, err)
	}

	for _, file := range files {
		if file.Type().IsRegular() && strings.HasSuffix(file.Name(), ".manitest.jsonnet") {
			fileName := path.Join(arg, file.Name())
			err := runner.RunTest(fileName)

			if err != nil {
				return fmt.Errorf("failed to test %s: %w", fileName, err)
			}
		}
	}

	return nil
}
