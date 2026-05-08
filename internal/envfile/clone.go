package envfile

import (
	"fmt"
	"slices"
)

// CloneOptions controls how an env file is cloned.
type CloneOptions struct {
	// OnlyKeys restricts cloning to these keys (empty = all keys).
	OnlyKeys []string
	// ExcludeKeys omits these keys from the clone.
	ExcludeKeys []string
	// Redact replaces values of sensitive keys with a placeholder.
	Redact bool
	// RedactPlaceholder is the string used when Redact is true.
	RedactPlaceholder string
	// DryRun skips writing the output file.
	DryRun bool
}

// CloneResult summarises the outcome of a Clone operation.
type CloneResult struct {
	Copied  []string
	Skipped []string
	Redacted []string
}

// Clone copies entries from src into a new EnvFile written to destPath,
// applying filtering and optional redaction according to opts.
func Clone(src *EnvFile, destPath string, opts CloneOptions) (*CloneResult, error) {
	if src == nil {
		return nil, fmt.Errorf("clone: source EnvFile must not be nil")
	}

	placeholder := opts.RedactPlaceholder
	if placeholder == "" {
		placeholder = "***"
	}

	result := &CloneResult{}
	var entries []Entry

	for _, e := range src.Entries {
		if len(opts.ExcludeKeys) > 0 && slices.Contains(opts.ExcludeKeys, e.Key) {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}
		if len(opts.OnlyKeys) > 0 && !slices.Contains(opts.OnlyKeys, e.Key) {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}

		cloned := Entry{Key: e.Key, Value: e.Value, Comment: e.Comment}
		if opts.Redact && isSensitiveKey(e.Key) {
			cloned.Value = placeholder
			result.Redacted = append(result.Redacted, e.Key)
		}
		entries = append(entries, cloned)
		result.Copied = append(result.Copied, e.Key)
	}

	if !opts.DryRun {
		dest := &EnvFile{Path: destPath, Entries: entries}
		if err := writeEntries(dest); err != nil {
			return nil, fmt.Errorf("clone: write %s: %w", destPath, err)
		}
	}

	return result, nil
}

// isSensitiveKey returns true for keys that look like secrets.
func isSensitiveKey(key string) bool {
	sensitive := []string{"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL", "AUTH", "API_KEY"}
	for _, s := range sensitive {
		if containsFold(key, s) {
			return true
		}
	}
	return false
}

func containsFold(s, sub string) bool {
	return len(s) >= len(sub) &&
		(s == sub ||
			len(s) > 0 && (stringContainsFold(s, sub)))
}

func stringContainsFold(s, sub string) bool {
	sU := toUpper(s)
	subU := toUpper(sub)
	for i := 0; i+len(subU) <= len(sU); i++ {
		if sU[i:i+len(subU)] == subU {
			return true
		}
	}
	return false
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
