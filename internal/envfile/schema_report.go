package envfile

import (
	"fmt"
	"io"
)

// WriteSchemaReport writes a human-readable schema validation report to w.
func WriteSchemaReport(w io.Writer, violations []SchemaViolation, masked bool) {
	if len(violations) == 0 {
		fmt.Fprintln(w, colorize("green", "✔ schema validation passed — no violations found"))
		return
	}

	fmt.Fprintf(w, colorize("red", "✘ schema validation failed — %d violation(s):\n"), len(violations))
	for _, v := range violations {
		key := colorize("yellow", v.Key)
		msg := v.Message
		if masked {
			msg = maskSchemaMessage(msg)
		}
		fmt.Fprintf(w, "  • %s: %s\n", key, msg)
	}
}

// maskSchemaMessage replaces quoted values in violation messages with "***".
func maskSchemaMessage(msg string) string {
	// Replace content inside double quotes with masked placeholder.
	var result []byte
	inQuote := false
	for i := 0; i < len(msg); i++ {
		ch := msg[i]
		if ch == '"' {
			if inQuote {
				result = append(result, []byte(`"***"`)...)
				inQuote = false
				// skip until closing quote already consumed
			} else {
				inQuote = true
			}
			continue
		}
		if !inQuote {
			result = append(result, ch)
		}
	}
	return string(result)
}
