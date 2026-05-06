package envfile

// DiffKind describes the type of difference between two env files.
type DiffKind string

const (
	Added    DiffKind = "added"    // key exists in target but not in source
	Removed  DiffKind = "removed"  // key exists in source but not in target
	Changed  DiffKind = "changed"  // key exists in both but values differ
	Unchanged DiffKind = "unchanged"
)

// DiffEntry represents a single difference between two env files.
type DiffEntry struct {
	Key         string
	Kind        DiffKind
	SourceValue string
	TargetValue string
}

// DiffResult holds the complete diff between a source and target env file.
type DiffResult struct {
	Entries []DiffEntry
}

// HasChanges returns true if there are any added, removed, or changed entries.
func (d *DiffResult) HasChanges() bool {
	for _, e := range d.Entries {
		if e.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Diff computes the difference between source and target EnvFiles.
func Diff(source, target *EnvFile) *DiffResult {
	srcMap := source.ToMap()
	tgtMap := target.ToMap()

	seen := make(map[string]bool)
	var entries []DiffEntry

	for k, sv := range srcMap {
		seen[k] = true
		if tv, ok := tgtMap[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Kind: Removed, SourceValue: sv})
		} else if sv != tv {
			entries = append(entries, DiffEntry{Key: k, Kind: Changed, SourceValue: sv, TargetValue: tv})
		} else {
			entries = append(entries, DiffEntry{Key: k, Kind: Unchanged, SourceValue: sv, TargetValue: tv})
		}
	}

	for k, tv := range tgtMap {
		if !seen[k] {
			entries = append(entries, DiffEntry{Key: k, Kind: Added, TargetValue: tv})
		}
	}

	return &DiffResult{Entries: entries}
}
