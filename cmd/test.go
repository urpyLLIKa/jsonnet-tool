package cmd

import (
	"fmt"

	"github.com/fatih/color"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-com/gl-infra/jsonnet-tool/internal/manitest"
)

var (
	testCommandJPaths []string
	writeFixtures     bool
	cacheResults      bool
	jsonnetExtVars    map[string]string
	emitAllTraces     bool
)

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

	testCommand.PersistentFlags().StringToStringVarP(
		&jsonnetExtVars, "ext-str", "V", map[string]string{},
		"Provide an external value as a string to jsonnet",
	)

	testCommand.PersistentFlags().BoolVarP(
		&emitAllTraces, "all-traces", "T", false,
		"Emit all traces. By default, only traces for failed tests will be emitted",
	)
}

var testCommand = NewTestCommand()

func newTestCommandRun(cmd *cobra.Command, args []string) error {
	traceVisitor := manitest.NewTraceVisitor(cmd.OutOrStdout(), cmd.ErrOrStderr())
	reporterVisitor := manitest.NewReporterVisitor(emitAllTraces, args, cmd.OutOrStdout(), cmd.ErrOrStderr())

	visitors := []manitest.TestVisitor{
		traceVisitor,
		reporterVisitor,
	}

	if writeFixtures {
		visitors = append(visitors, &manitest.WriterVisitor{})
	}

	vm := jsonnet.MakeVM()
	for k, v := range jsonnetExtVars {
		vm.ExtVar(k, v)
	}

	vm.SetTraceOut(traceVisitor)
	vm.ErrorFormatter.SetColorFormatter(color.New(color.FgRed).Fprintf)
	vm.Importer(&jsonnet.FileImporter{
		JPaths: testCommandJPaths,
	})

	var cacheManager *manitest.CacheManager
	if cacheResults {
		cacheManager = manitest.NewCacheManager(vm)

		err := cacheManager.LoadCachedResults()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "failed to load cached test results: %v\n", err)
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

	// Add required natives
	runner.RegisterNatives()

	// No errors as they are collected by the runner
	runTests(runner, args)

	if cacheManager != nil {
		err := cacheManager.SaveCachedResults()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "failed to save cached test results: %v\n", err)
		}
	}

	err := visitor.AllTestsCompleted()
	if err != nil {
		// AllTestsCompleted passes the error back to the caller, which may control the termination
		// of the program.
		return fmt.Errorf("visitor returned error: %w", err)
	}

	return nil
}

func silenceErrorsUsage(cmd *cobra.Command, args []string) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
}

func NewTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:              "test",
		Short:            "Run jsonnet tests",
		Args:             cobra.MinimumNArgs(1),
		PersistentPreRun: silenceErrorsUsage,
		RunE:             newTestCommandRun,
	}
}

// Given a test runner, run the tests.
func runTests(runner *manitest.TestRunner, args []string) {
	for _, a := range args {
		runner.RunTestFile(a)
	}
}
