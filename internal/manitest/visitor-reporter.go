package manitest

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/kr/text"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"

	"github.com/alessio/shellescape"
)

type ReporterVisitor struct {
	// EmitAllTraces controls whether all traces should be emitted, or only those for failed tests.
	EmitAllTraces bool

	// All the positional args passed to the test command
	Args []string

	passes      int
	total       int
	fail        int
	invalid     int
	cached      bool
	fileInvalid bool

	manifestStart time.Time

	totalFiles        int
	totalFilesPassed  int
	totalFilesCached  int
	totalFileFailures int
	totalFileInvalid  int

	stdout io.Writer
	stderr io.Writer

	baseVisitor
}

var _ TestVisitor = &ReporterVisitor{}

func (rv *ReporterVisitor) TestFileStarted(fileName string) error {
	rv.passes = 0
	rv.total = 0
	rv.fail = 0
	rv.invalid = 0
	rv.cached = false
	rv.fileInvalid = false

	_, _ = fmt.Fprintf(rv.stdout, "▶️  Executing test file %s\n", fileName)

	return nil
}

func (rv *ReporterVisitor) TestFileInvalid(name string, err error) error {
	_, _ = fmt.Fprintf(rv.stdout, "💥  Invalid test %s\n%v\n", name, err)
	rv.fileInvalid = true

	return nil
}

func (rv *ReporterVisitor) TestFileCompleted(fileName string, _ bool) error {
	rv.totalFiles = rv.totalFiles + 1

	reportElements := []string{
		color.HiWhiteString("%d %s tested", rv.total, plural(rv.totalFiles, "test", "files")),
		color.HiBlueString("%d %s passed", rv.passes, plural(rv.totalFilesPassed, "test", "tests")),
	}

	if rv.fail > 0 {
		reportElements = append(reportElements,
			color.HiYellowString("%d %s failed", rv.fail, plural(rv.fail, "test", "tests")),
		)
	}

	if rv.invalid > 0 {
		reportElements = append(reportElements,
			color.HiRedString("%d %s invalid", rv.invalid, plural(rv.invalid, "test", "tests")),
		)
	}

	if rv.cached {
		reportElements = append(reportElements,
			color.HiBlueString("cached"),
		)
	}

	_, _ = fmt.Printf("\r  %s\n\n", strings.Join(reportElements, ", "))

	if rv.fail > 0 || rv.invalid > 0 || rv.fileInvalid {
		rerunCommand := rv.getRerunCommand(fileName)
		_, _ = fmt.Fprintln(rv.stdout, text.Indent(rerunCommand, "      "))
	}

	if rv.invalid > 0 || rv.fileInvalid {
		rv.totalFileInvalid = rv.totalFileInvalid + 1
	} else if rv.fail > 0 {
		rv.totalFileFailures = rv.totalFileFailures + 1
	} else {
		rv.totalFilesPassed = rv.totalFilesPassed + 1
	}

	if rv.cached {
		rv.totalFilesCached = rv.totalFilesCached + 1
		return nil
	}

	return nil
}

func (rv *ReporterVisitor) TestCaseManifestationStarted(fileName string, testcase string) error {
	_, _ = fmt.Fprintf(rv.stdout, "  ➡️  %s manifesting...", testcase)

	rv.manifestStart = time.Now()

	return nil
}

func (rv *ReporterVisitor) TestCaseManifestationCompleted(fileName string, testcase string) error {
	duration := time.Since(rv.manifestStart)

	_, _ = fmt.Fprintf(rv.stdout, "\r  ➡️  %s manifestation completed in %dms\n", testcase, int64(duration/time.Millisecond))

	return nil
}

func (rv *ReporterVisitor) TestCaseEvaluationCompleted(fileName string, testcase string, result *TestCaseResult) error {
	if result.Cached {
		rv.cached = true

		if result.Success {
			_, _ = fmt.Fprintf(rv.stdout, "\r  ✔️  %s (all tests) (cached)\n", fileName)
		} else {
			_, _ = fmt.Fprintf(rv.stdout, "\r  ⨯  %s (all tests) (cached)\n", fileName)
		}

		return nil
	}

	rv.total = rv.total + 1

	if result.Success {
		rv.passes = rv.passes + 1

		_, _ = fmt.Fprintf(rv.stdout, "\r  ✔️  %-20s %-40s\n", testcase, "")
	} else {
		rv.fail = rv.fail + 1

		_, _ = fmt.Fprintf(rv.stdout, "\r  ❌  %-6s %-20s %-40s\n", testcase, color.HiRedString("failed"), color.YellowString(result.FixturePath))
		_, _ = fmt.Fprintf(rv.stdout, "\r      %s\n", result.Error)

		expected := normalizePlainTextString(result.Expected)
		actual := normalizePlainTextString(result.Actual)

		if expected != actual {
			prettyDiff := rv.generatePrettyDiff(result, expected, actual)

			_, _ = fmt.Fprintf(rv.stdout, "%s\n\n", text.Indent(prettyDiff, "      "))
		}
	}

	// If there's a trace and either we're showing all traces, or this test failed
	// show the trace.
	if result.Trace != "" && (!result.Success || rv.EmitAllTraces) {
		_, _ = fmt.Fprintln(rv.stdout, text.Indent(result.Trace, "      "))
	}

	return nil
}

func (rv *ReporterVisitor) TestCaseInvalid(name string, testcase string, err error) error {
	rv.invalid = rv.invalid + 1
	_, _ = fmt.Fprintf(rv.stdout, "💥  Invalid test case %s\n%v\n", testcase, err)

	return nil
}

// generatePrettyDiff will generate a diff for display.
func (rv *ReporterVisitor) generatePrettyDiff(result *TestCaseResult, expected, actual string) string {
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
	finalReportElements := []string{
		color.HiWhiteString("%d %s tested", rv.totalFiles, plural(rv.totalFiles, "file", "files")),
		color.HiBlueString("%d %s passed", rv.totalFilesPassed, plural(rv.totalFilesPassed, "file", "files")),
	}
	icon := "✅"

	if rv.totalFileFailures > 0 {
		icon = "❌"

		finalReportElements = append(finalReportElements,
			color.HiYellowString("%d %s failed", rv.totalFileFailures, plural(rv.totalFileFailures, "file", "files")),
		)
	}

	if rv.totalFileInvalid > 0 {
		icon = "💥"

		finalReportElements = append(finalReportElements,
			color.HiRedString("%d %s invalid", rv.totalFileInvalid, plural(rv.totalFileInvalid, "file", "files")),
		)
	}

	if rv.totalFilesCached > 0 {
		finalReportElements = append(finalReportElements,
			color.CyanString("%d %s cached", rv.totalFilesCached, plural(rv.totalFilesCached, "file", "files")),
		)
	}

	_, _ = fmt.Fprintf(rv.stdout, "--------------------------------------------------------\n%s Test suite completed: %s\n\n",
		icon,
		strings.Join(finalReportElements, ", "),
	)

	return nil
}

func plural(count int, sing string, plur string) string {
	if count == 1 {
		return sing
	}

	return plur
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

func NewReporterVisitor(emitAllTraces bool, args []string, stdout io.Writer, stderr io.Writer) *ReporterVisitor {
	return &ReporterVisitor{
		EmitAllTraces: emitAllTraces,
		Args:          args,
		stdout:        stdout,
		stderr:        stderr,
	}
}
