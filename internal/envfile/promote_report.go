package envfile

import (
	"fmt"
	"io"
)

// WritePromoteReport writes a human-readable summary of a PromoteResult to w.
func WritePromoteReport(w io.Writer, result PromoteResult, masked bool) {
	if len(result.Promoted) == 0 && len(result.Overwritten) == 0 && len(result.Skipped) == 0 {
		fmt.Fprintln(w, "promote: nothing to do")
		return
	}

	for _, k := range result.Promoted {
		line := fmt.Sprintf("  + %-30s  (promoted)", k)
		fmt.Fprintln(w, colorize("green", line))
	}

	for _, k := range result.Overwritten {
		val := "***"
		if !masked {
			val = "(overwritten)"
		}
		line := fmt.Sprintf("  ~ %-30s  %s", k, val)
		fmt.Fprintln(w, colorize("yellow", line))
	}

	for _, k := range result.Skipped {
		line := fmt.Sprintf("  - %-30s  (skipped)", k)
		fmt.Fprintln(w, colorize("cyan", line))
	}

	fmt.Fprintf(w, "\npromote: %d promoted, %d overwritten, %d skipped\n",
		len(result.Promoted), len(result.Overwritten), len(result.Skipped))
}
