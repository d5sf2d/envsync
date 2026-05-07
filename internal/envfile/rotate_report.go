package envfile

import (
	"fmt"
	"io"
)

// WriteRotateReport writes a human-readable summary of rotation results.
func WriteRotateReport(w io.Writer, results []RotateResult, dryRun bool) {
	if dryRun {
		fmt.Fprintln(w, colorize("yellow", "[dry-run] No changes written."))
	}

	rotated := 0
	skipped := 0

	for _, r := range results {
		switch {
		case r.Rotated:
			rotated++
			fmt.Fprintf(w, "  %s %s\n",
				colorize("green", "~"),
				r.Key,
			)
		case r.Skipped:
			skipped++
			fmt.Fprintf(w, "  %s %s (%s)\n",
				colorize("yellow", "?"),
				r.Key,
				r.Reason,
			)
		}
	}

	fmt.Fprintf(w, "\nRotation complete: %d rotated, %d skipped.\n", rotated, skipped)
}
