package envfile

import (
	"fmt"
	"io"
)

// WritePinReport writes a human-readable summary of Pin or CheckPins results.
func WritePinReport(w io.Writer, results []PinResult, dryRun bool) {
	if len(results) == 0 {
		fmt.Fprintln(w, "pin: no keys processed")
		return
	}

	if dryRun {
		fmt.Fprintln(w, colorize("yellow", "[dry-run] pin report:"))
	} else {
		fmt.Fprintln(w, "pin report:")
	}

	for _, r := range results {
		switch {
		case r.Expired:
			fmt.Fprintf(w, "  %s %-30s  %s\n",
				colorize("yellow", "EXPIRED"),
				r.Key,
				r.Reason,
			)
		case r.Skipped:
			fmt.Fprintf(w, "  %s  %-30s  %s\n",
				colorize("red", "SKIPPED"),
				r.Key,
				r.Reason,
			)
		case r.Pinned:
			label := "PINNED "
			if dryRun {
				label = "WOULD PIN"
			}
			fmt.Fprintf(w, "  %s  %-30s\n",
				colorize("green", label),
				r.Key,
			)
		default:
			fmt.Fprintf(w, "  %-10s %-30s\n", "OK", r.Key)
		}
	}
}
