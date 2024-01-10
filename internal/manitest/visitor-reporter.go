package manitest

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/kr/text"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type ReporterVisitor struct {
	success int
	total   int
	fail    int
	cached  bool

	totalFiles        int
	totalFileFailures int
}

var _ TestVisitor = &ReporterVisitor{}

func (rv *ReporterVisitor) StartTestFile(fileName string) error {
	rv.success = 0
	rv.total = 0
	rv.fail = 0
	rv.cached = false

	fmt.Printf("▶️  Executing manitest %s\n", fileName)

	return nil
}

func (rv *ReporterVisitor) TestFileComplete(fileName string, allSuccessful bool) error {
	if rv.cached {
		return nil
	}

	rv.totalFiles = rv.totalFiles + 1
	if !allSuccessful {
		rv.totalFileFailures = rv.totalFileFailures + 1
	}

	if rv.fail > 0 {
		fmt.Printf("\r  %s %s %d test(s) completed with %d failure(s)\n\n", color.HiRedString("Failed"), fileName, rv.total, rv.fail)
	} else {
		fmt.Printf("\r  %s %s %d test(s) completed\n\n", color.BlueString("Completed"), fileName, rv.total)
	}

	return nil
}

func (rv *ReporterVisitor) StartTestCase(fileName string, testcase string) error {
	fmt.Printf("\r  ➡️  %s starting...", testcase)

	return nil
}

func (rv *ReporterVisitor) TestCaseComplete(fileName string, testcase string, result *TestCaseResult) error {
	if result.Cached {
		rv.cached = true

		if result.Success {
			fmt.Printf("\r  ✔️  %s (all tests) (cached)\n", fileName)
		} else {
			fmt.Printf("\r  ⨯  %s (all tests) (cached)\n", fileName)
		}

		return nil
	}

	rv.total = rv.total + 1

	if result.Success {
		rv.success = rv.success + 1

		fmt.Printf("\r  ✔️  %-20s %-40s\n", testcase, "")
	} else {
		rv.fail = rv.fail + 1

		fmt.Printf("\r  ❌  %-6s %-20s %-40s\n", testcase, color.HiRedString("failed"), color.YellowString(result.FixturePath))

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(result.Expected, result.Actual, false)

		fmt.Println(text.Indent(dmp.DiffPrettyText(diffs), "    "))
		fmt.Println()
	}

	return nil
}

func (rv *ReporterVisitor) Delta(fileName string, testcase string, fixturePath string, canonicalActual string, canonicalExpected string) error {
	return nil
}

func (rv *ReporterVisitor) CachedResult(fileName string) (*TestCaseResult, error) {
	return nil, nil
}

func (rv *ReporterVisitor) Complete() error {
	if rv.totalFileFailures > 0 {
		fmt.Printf("❌ %s: %d files tested, %s\n",
			color.HiRedString("Test Suite Failed"),
			rv.totalFiles,
			color.YellowString(fmt.Sprintf("%d files failed", rv.totalFileFailures)),
		)

		return nil
	}

	fmt.Printf("✅  Test Suite Passed: %d files tested\n",
		rv.totalFiles,
	)

	return nil
}
