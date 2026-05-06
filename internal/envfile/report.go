package envfile

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// ReportOptions controls how the diff report is rendered.
type ReportOptions struct {
	Color  bool
	Masked bool
}

// WriteDiffReport writes a human-readable diff report to w.
func WriteDiffReport(w io.Writer, diffs []DiffEntry, opts ReportOptions) {
	if len(diffs) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	m := NewMasker(nil)

	for _, d := range diffs {
		var line string
		switch d.Type {
		case DiffAdded:
			v := displayValue(d.NewValue, d.Key, opts.Masked, m)
			line = fmt.Sprintf("+ %s=%s", d.Key, v)
			line = colorize(line, colorGreen, opts.Color)
		case DiffRemoved:
			v := displayValue(d.OldValue, d.Key, opts.Masked, m)
			line = fmt.Sprintf("- %s=%s", d.Key, v)
			line = colorize(line, colorRed, opts.Color)
		case DiffChanged:
			old := displayValue(d.OldValue, d.Key, opts.Masked, m)
			new := displayValue(d.NewValue, d.Key, opts.Masked, m)
			line = fmt.Sprintf("~ %s: %s → %s", d.Key, old, new)
			line = colorize(line, colorYellow, opts.Color)
		case DiffUnchanged:
			v := displayValue(d.OldValue, d.Key, opts.Masked, m)
			line = fmt.Sprintf("  %s=%s", d.Key, v)
			line = colorize(line, colorCyan, opts.Color)
		}
		fmt.Fprintln(w, line)
	}
}

// WriteSyncReport writes a summary of a completed sync operation.
func WriteSyncReport(w io.Writer, res *SyncResult) {
	fmt.Fprintf(w, "Sync complete: +%d added, ~%d updated, =%d skipped\n",
		len(res.Added), len(res.Updated), len(res.Skipped))
	if len(res.Added) > 0 {
		fmt.Fprintf(w, "  Added:   %s\n", strings.Join(res.Added, ", "))
	}
	if len(res.Updated) > 0 {
		fmt.Fprintf(w, "  Updated: %s\n", strings.Join(res.Updated, ", "))
	}
	if len(res.Skipped) > 0 {
		fmt.Fprintf(w, "  Skipped: %s\n", strings.Join(res.Skipped, ", "))
	}
}

func displayValue(val, key string, masked bool, m *Masker) string {
	if masked {
		return m.Mask(key, val)
	}
	return val
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}
