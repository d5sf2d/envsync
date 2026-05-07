package envfile

import (
	"fmt"
	"io"
)

// WriteTemplateReport writes a human-readable summary of a TemplateResult to w.
func WriteTemplateReport(w io.Writer, res TemplateResult) {
	if len(res.Missing) == 0 && len(res.Unused) == 0 {
		fmt.Fprintln(w, colorize("green", "✔ Template rendered successfully with no issues."))
		return
	}

	if len(res.Missing) > 0 {
		fmt.Fprintln(w, colorize("red", fmt.Sprintf("✖ %d missing placeholder(s):", len(res.Missing))))
		for _, k := range res.Missing {
			fmt.Fprintf(w, "    %s ${%s}\n", colorize("red", "-"), k)
		}
	}

	if len(res.Unused) > 0 {
		fmt.Fprintln(w, colorize("yellow", fmt.Sprintf("⚠ %d unused env key(s):", len(res.Unused))))
		for _, k := range res.Unused {
			fmt.Fprintf(w, "    %s %s\n", colorize("yellow", "~"), k)
		}
	}
}
