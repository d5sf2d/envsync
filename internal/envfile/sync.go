package envfile

import (
	"fmt"
	"os"
	"strings"
)

// SyncOptions controls how a sync operation behaves.
type SyncOptions struct {
	// DryRun prints what would change without writing.
	DryRun bool
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
	// AddMissing adds keys present in source but missing in destination.
	AddMissing bool
}

// SyncResult summarises what happened during a sync.
type SyncResult struct {
	Added   []string
	Updated []string
	Skipped []string
}

// Sync applies changes from src onto dst according to opts.
// If DryRun is false the destination file is rewritten.
func Sync(dst, src *EnvFile, destPath string, opts SyncOptions) (*SyncResult, error) {
	result := &SyncResult{}
	diffs := Diff(dst, src)

	updated := cloneEntries(dst.Entries)

	for _, d := range diffs {
		switch d.Type {
		case DiffAdded:
			if opts.AddMissing {
				updated = append(updated, Entry{Key: d.Key, Value: d.NewValue})
				result.Added = append(result.Added, d.Key)
			} else {
				result.Skipped = append(result.Skipped, d.Key)
			}
		case DiffChanged:
			if opts.Overwrite {
				for i, e := range updated {
					if e.Key == d.Key {
						updated[i].Value = d.NewValue
						break
					}
				}
				result.Updated = append(result.Updated, d.Key)
			} else {
				result.Skipped = append(result.Skipped, d.Key)
			}
		}
	}

	if !opts.DryRun {
		if err := writeEntries(destPath, updated); err != nil {
			return result, fmt.Errorf("sync: write %s: %w", destPath, err)
		}
	}

	return result, nil
}

func cloneEntries(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	return out
}

func writeEntries(path string, entries []Entry) error {
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
	}
	return os.WriteFile(path, []byte(sb.String()), 0o600)
}
