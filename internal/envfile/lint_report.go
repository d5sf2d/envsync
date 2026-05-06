package envfile

import (
	"fmt"
	"io"
)

// WriteLintReport writes a human-readable lint report to w.
// Returns true if any errors were found.
func WriteLintReport(w io.Writer, filename string, result *LintResult) bool {
	if len(result.Issues) == 0 {
		fmt.Fprintf(w, "%s\n", colorize("green", fmt.Sprintf("✔ %s: no lint issues found", filename)))
		return false
	}

	warnCount := 0
	errorCount := 0

	fmt.Fprintf(w, "Lint report for %s:\n", colorize("cyan", filename))
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, issue := range result.Issues {
		var line string
		switch issue.Severity {
		case LintError:
			line = colorize("red", issue.String())
			errorCount++
		case LintWarn:
			line = colorize("yellow", issue.String())
			warnCount++
		}
		fmt.Fprintln(w, line)
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintf(w, "Summary: %s, %s\n",
		colorize("red", fmt.Sprintf("%d error(s)", errorCount)),
		colorize("yellow", fmt.Sprintf("%d warning(s)", warnCount)),
	)

	return errorCount > 0
}
