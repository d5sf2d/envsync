package envfile

import (
	"fmt"
	"io"
)

// WritePatchReport writes a human-readable summary of patch results to w.
func WritePatchReport(w io.Writer, results []PatchResult, dryRun bool) {
	if dryRun {
		fmt.Fprintln(w, colorize("yellow", "[dry-run] No changes were written."))
	}

	if len(results) == 0 {
		fmt.Fprintln(w, "No patch instructions.")
		return
	}

	applied := 0
	skipped := 0

	for _, r := range results {
		if r.Applied {
			applied++
			label := appliedLabel(r.Instruction.Op)
			fmt.Fprintf(w, "  %s %s\n", colorize("green", label), r.Note)
		} else {
			skipped++
			fmt.Fprintf(w, "  %s %s\n", colorize("yellow", "SKIP"), r.Note)
		}
	}

	fmt.Fprintf(w, "\nPatch summary: %d applied, %d skipped.\n", applied, skipped)
}

func appliedLabel(op PatchOp) string {
	switch op {
	case PatchSet:
		return "SET   "
	case PatchDelete:
		return "DELETE"
	case PatchRename:
		return "RENAME"
	default:
		return "OP    "
	}
}
