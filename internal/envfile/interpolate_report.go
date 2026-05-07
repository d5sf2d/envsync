package envfile

import (
	"fmt"
	"io"
)

// WriteInterpolateReport writes a human-readable summary of interpolation
// results to w. It lists each variable that was expanded, showing the
// original reference and the resolved value.
func WriteInterpolateReport(w io.Writer, results []InterpolateResult, color bool) {
	if len(results) == 0 {
		fmt.Fprintln(w, colorize("✔ No interpolations performed.", "green", color))
		return
	}

	fmt.Fprintf(w, colorize("Interpolated %d variable(s):\n", "cyan", color), len(results))

	for _, r := range results {
		key := colorize(r.Key, "yellow", color)
		orig := colorize(r.Original, "red", color)
		resolved := colorize(r.Interpolated, "green", color)
		fmt.Fprintf(w, "  %s: %s → %s\n", key, orig, resolved)
	}
}
