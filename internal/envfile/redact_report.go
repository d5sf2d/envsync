package envfile

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteRedactReport writes a summary of which keys were redacted to w.
func WriteRedactReport(w io.Writer, original, redacted EnvFile, noColor bool) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	redactedCount := 0

	fmt.Fprintln(tw, "KEY\tSTATUS\tVALUE")
	fmt.Fprintln(tw, "---\t------\t-----")

	origMap := make(map[string]string, len(original))
	for _, e := range original {
		origMap[e.Key] = e.Value
	}

	for _, entry := range redacted {
		orig, exists := origMap[entry.Key]
		if !exists {
			continue
		}

		if entry.Value != orig {
			redactedCount++
			status := colorize("REDACTED", "yellow", noColor)
			fmt.Fprintf(tw, "%s\t%s\t%s\n", entry.Key, status, entry.Value)
		} else {
			status := colorize("ok", "green", noColor)
			fmt.Fprintf(tw, "%s\t%s\t%s\n", entry.Key, status, displayValue(entry.Value))
		}
	}

	tw.Flush()

	if redactedCount == 0 {
		fmt.Fprintln(w, "\nNo sensitive keys detected.")
	} else {
		fmt.Fprintf(w, "\n%d key(s) redacted.\n", redactedCount)
	}
}
