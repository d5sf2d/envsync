package envfile

import (
	"strings"
)

// TrimResult holds the outcome of trimming a single entry.
type TrimResult struct {
	Key      string
	Original string
	Trimmed  string
	Changed  bool
}

// TrimOptions controls which kinds of whitespace trimming are applied.
type TrimOptions struct {
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// TrimKeys removes leading/trailing whitespace from keys.
	TrimKeys bool
	// NormalizeEmpty replaces whitespace-only values with empty string.
	NormalizeEmpty bool
	// DryRun reports changes without mutating the EnvFile.
	DryRun bool
}

// Trim applies whitespace trimming to entries in an EnvFile according to opts.
// It returns a slice of TrimResult describing every entry that was examined.
func Trim(ef *EnvFile, opts TrimOptions) ([]TrimResult, error) {
	if ef == nil {
		return nil, nil
	}

	results := make([]TrimResult, 0, len(ef.Entries))

	for i := range ef.Entries {
		entry := &ef.Entries[i]
		origKey := entry.Key
		origVal := entry.Value

		newKey := origKey
		newVal := origVal

		if opts.TrimKeys {
			newKey = strings.TrimSpace(origKey)
		}

		if opts.TrimValues {
			newVal = strings.TrimSpace(origVal)
		}

		if opts.NormalizeEmpty && strings.TrimSpace(newVal) == "" {
			newVal = ""
		}

		changed := newKey != origKey || newVal != origVal

		results = append(results, TrimResult{
			Key:      origKey,
			Original: origVal,
			Trimmed:  newVal,
			Changed:  changed,
		})

		if changed && !opts.DryRun {
			entry.Key = newKey
			entry.Value = newVal
		}
	}

	return results, nil
}
