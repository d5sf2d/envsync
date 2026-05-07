package envfile

import (
	"fmt"
	"io"
	"strings"
)

// WriteAuditReport writes a human-readable audit log to w.
func WriteAuditReport(w io.Writer, log *AuditLog, maskSecrets bool) {
	if len(log.Entries) == 0 {
		fmt.Fprintln(w, "No audit entries recorded.")
		return
	}

	fmt.Fprintf(w, "Audit Log (%d entries)\n", len(log.Entries))
	fmt.Fprintln(w, strings.Repeat("-", 60))

	for _, e := range log.Entries {
		ts := e.Timestamp.Format("2006-01-02T15:04:05Z")
		actor := e.Actor
		if actor == "" {
			actor = "unknown"
		}
		source := e.Source
		if source == "" {
			source = "unknown"
		}

		oldVal := e.OldValue
		newVal := e.NewValue
		if maskSecrets {
			if oldVal != "" {
				oldVal = "***"
			}
			if newVal != "" {
				newVal = "***"
			}
		}

		switch e.Action {
		case AuditAdded:
			fmt.Fprintf(w, "[%s] %s | %-8s | %s = %q (actor: %s, src: %s)\n",
				ts, colorize("green", "+"), e.Action, e.Key, newVal, actor, source)
		case AuditRemoved:
			fmt.Fprintf(w, "[%s] %s | %-8s | %s = %q (actor: %s, src: %s)\n",
				ts, colorize("red", "-"), e.Action, e.Key, oldVal, actor, source)
		case AuditChanged:
			fmt.Fprintf(w, "[%s] %s | %-8s | %s: %q -> %q (actor: %s, src: %s)\n",
				ts, colorize("yellow", "~"), e.Action, e.Key, oldVal, newVal, actor, source)
		default:
			fmt.Fprintf(w, "[%s] * | %-8s | %s (actor: %s, src: %s)\n",
				ts, e.Action, e.Key, actor, source)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 60))
	fmt.Fprintln(w, log.Summary())
}
