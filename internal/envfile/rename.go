package envfile

import (
	"fmt"
)

// RenameResult describes the outcome of a single rename operation.
type RenameResult struct {
	From    string
	To      string
	Skipped bool
	Reason  string
}

// RenameOptions controls the behaviour of Rename.
type RenameOptions struct {
	// DryRun reports what would change without mutating the EnvFile.
	DryRun bool
	// FailOnMissing returns an error when a source key does not exist.
	FailOnMissing bool
	// FailOnCollision returns an error when the destination key already exists.
	FailOnCollision bool
}

// Rename renames one or more keys inside ef according to the pairs map
// (map[oldKey]newKey). It preserves the original insertion order of entries,
// replacing the key name in-place.
func Rename(ef *EnvFile, pairs map[string]string, opts RenameOptions) ([]RenameResult, error) {
	var results []RenameResult

	// Build a quick lookup of existing keys.
	existing := make(map[string]bool, len(ef.Entries))
	for _, e := range ef.Entries {
		existing[e.Key] = true
	}

	for from, to := range pairs {
		if !existing[from] {
			if opts.FailOnMissing {
				return results, fmt.Errorf("rename: source key %q not found", from)
			}
			results = append(results, RenameResult{From: from, To: to, Skipped: true, Reason: "source key not found"})
			continue
		}
		if existing[to] && from != to {
			if opts.FailOnCollision {
				return results, fmt.Errorf("rename: destination key %q already exists", to)
			}
			results = append(results, RenameResult{From: from, To: to, Skipped: true, Reason: "destination key already exists"})
			continue
		}

		if !opts.DryRun {
			for i, e := range ef.Entries {
				if e.Key == from {
					ef.Entries[i].Key = to
					delete(existing, from)
					existing[to] = true
					break
				}
			}
		}
		results = append(results, RenameResult{From: from, To: to})
	}

	return results, nil
}
