package envfile

import (
	"fmt"
	"io"
)

// WriteCloneReport writes a human-readable summary of a CloneResult to w.
func WriteCloneReport(w io.Writer, result *CloneResult, destPath string, dryRun bool) {
	if dryRun {
		fmt.Fprintln(w, colorize("yellow", "[dry-run] Clone preview — no file written"))
	} else {
		fmt.Fprintf(w, "Clone → %s\n", destPath)
	}

	fmt.Fprintf(w, "  Copied:   %d key(s)\n", len(result.Copied))
	for _, k := range result.Copied {
		redacted := false
		for _, r := range result.Redacted {
			if r == k {
				redacted = true
				break
			}
		}
		if redacted {
			fmt.Fprintf(w, "    %s %s %s\n",
				colorize("cyan", k),
				colorize("yellow", "(redacted)"),
				"")
		} else {
			fmt.Fprintf(w, "    %s\n", colorize("cyan", k))
		}
	}

	if len(result.Skipped) > 0 {
		fmt.Fprintf(w, "  Skipped:  %d key(s)\n", len(result.Skipped))
		for _, k := range result.Skipped {
			fmt.Fprintf(w, "    %s\n", colorize("yellow", k))
		}
	}

	if len(result.Redacted) > 0 {
		fmt.Fprintf(w, "  Redacted: %d key(s)\n", len(result.Redacted))
	}
}
