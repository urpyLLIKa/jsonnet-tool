package manitest

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/kr/text"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	shellescape "gopkg.in/alessio/shellescape.v1"
)

type ReporterVisitor struct {
	// EmitAllTraces controls whether all traces should be emitted, or only those for failed tests.
	EmitAllTraces bool

	// All the positional args passed to the test command
	Args []string

	success int
	total   int
	fail    int
	cached  bool

	manifestStart time.Time

	totalFiles        int
	totalFilesCached  int
	totalFileFailures int

	baseVisitor
}

var _ TestVisitor = &ReporterVisitor{}

func (rv *ReporterVisitor) TestFileStarted(fileName string) error {
	rv.success = 0
	rv.total = 0
	rv.fail = 0
	rv.cached = false

	fmt.Printf("▶️  Executing test file %s\n", fileName)

	return nil
}

func (rv *ReporterVisitor) TestFileCompleted(fileName string, allSuccessful bool) error {
	rv.totalFiles = rv.totalFiles + 1
	if !allSuccessful {
		rv.totalFileFailures = rv.totalFileFailures + 1
	}

	if rv.cached {
		rv.totalFilesCached = rv.totalFilesCached + 1
		return nil
	}

	if rv.fail > 0 {
		fmt.Printf("\r  %s %s %d test(s) completed with %d failure(s)\n\n", color.HiRedString("Failed"), fileName, rv.total, rv.fail)

		rerunCommand := rv.getRerunCommand(fileName)
		fmt.Println(text.Indent(rerunCommand, "      "))
	}

	return nil
}

func (rv *ReporterVisitor) TestCaseManifestationStarted(fileName string, testcase string) error {
	fmt.Printf("  ➡️  %s manifesting...", testcase)

	rv.manifestStart = time.Now()

	return nil
}

func (rv *ReporterVisitor) TestCaseManifestationCompleted(fileName string, testcase string) error {
	duration := time.Since(rv.manifestStart)

	fmt.Printf("\r  ➡️  %s manifestation completed in %dms\n", testcase, int64(duration/time.Millisecond))

	return nil
}

func (rv *ReporterVisitor) TestCaseEvaluationCompleted(fileName string, testcase string, result *TestCaseResult) error {
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
		fmt.Printf("\r      %s\n", result.Error)

		expected := normalizePlainTextString(result.Expected)
		actual := normalizePlainTextString(result.Actual)

		if expected != actual {
			prettyDiff := generatePrettyDiff(result, expected, actual)

			fmt.Println(text.Indent(prettyDiff, "      "))
			fmt.Println()
		}
	}

	// If there's a trace and either we're showing all traces, or this test failed
	// show the trace.
	if result.Trace != "" && (!result.Success || rv.EmitAllTraces) {
		fmt.Println(text.Indent(result.Trace, "      "))
	}

	return nil
}

// generatePrettyDiff will generate a diff for display.
func generatePrettyDiff(result *TestCaseResult, expected, actual string) string {
	edits := myers.ComputeEdits(span.URIFromPath(result.FixturePath), expected, actual)
	diff := fmt.Sprint(gotextdiff.ToUnified(result.FixturePath, result.FixturePath, expected, edits))

	out := ""

	// Remove the first two lines of the unified diff
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(diff))

	for scanner.Scan() {
		count = count + 1

		if count > 3 {
			out = out + "\n"
		}

		if count > 2 {
			line := scanner.Text()
			if strings.HasPrefix(line, "-") {
				out = out + color.RedString(scanner.Text())
			} else if strings.HasPrefix(line, "+") {
				out = out + color.YellowString(scanner.Text())
			} else {
				out = out + scanner.Text()
			}
		}
	}

	return out
}

func (rv *ReporterVisitor) AllTestsCompleted() error {
	if rv.totalFileFailures > 0 {
		fmt.Printf("\n\n  ❌ %s: %d file(s) tested, %d cached, %s\n\n",
			color.HiRedString("Test run completed. Some tests failed."),
			rv.totalFiles,
			rv.totalFilesCached,
			color.YellowString(fmt.Sprintf("%d file(s) failed", rv.totalFileFailures)),
		)

		return nil
	}

	fmt.Printf("  ✅ Testing run completed. All tests passed: %d file(s) tested, %d file(s) cached\n\n",
		rv.totalFiles,
		rv.totalFilesCached,
	)

	return nil
}

func (rv *ReporterVisitor) getRerunCommand(fileName string) string {
	cmd := []string{os.Args[0]}

	// This is very hacky: recreate the full command line by excluding the
	// positional args.
	set := map[string]struct{}{}
	for _, v := range rv.Args {
		set[v] = struct{}{}
	}

	for i, v := range os.Args {
		_, ok := set[v]
		if i > 0 && !ok {
			cmd = append(cmd, v)
		}
	}

	cmd = append(cmd, shellescape.Quote(fileName))

	return "To rerun this test on it's own, use the following command:\n" +
		color.CyanString("%s\n\n", strings.Join(cmd, " ")) +
		color.MagentaString("%s", "NOTE: adding the --write-fixtures will auto-update the fixtures based on the actual values.")
}

func normalizePlainTextString(s string) string {
	return strings.TrimSpace(s) + "\n"
}
