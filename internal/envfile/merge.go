package envfile

import "fmt"

// MergeStrategy controls how conflicting keys are handled during a merge.
type MergeStrategy int

const (
	// MergeStrategyKeepBase keeps the base file's value on conflict.
	MergeStrategyKeepBase MergeStrategy = iota
	// MergeStrategyPreferOverride uses the override file's value on conflict.
	MergeStrategyPreferOverride
	// MergeStrategyError returns an error on any conflict.
	MergeStrategyError
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Entries  []Entry
	Conflicts []MergeConflict
}

// MergeConflict describes a key that existed in both files with different values.
type MergeConflict struct {
	Key       string
	BaseValue string
	OverValue string
}

// Merge combines base and override EnvFiles according to the given strategy.
// Keys present only in base or only in override are always included.
// Conflicts are resolved per strategy.
func Merge(base, override *EnvFile, strategy MergeStrategy) (*MergeResult, error) {
	result := &MergeResult{}

	baseMap := make(map[string]string, len(base.Entries))
	for _, e := range base.Entries {
		baseMap[e.Key] = e.Value
	}

	overrideMap := make(map[string]string, len(override.Entries))
	for _, e := range override.Entries {
		overrideMap[e.Key] = e.Value
	}

	seen := make(map[string]bool)

	// Process base entries first to preserve order.
	for _, e := range base.Entries {
		seen[e.Key] = true
		if overVal, exists := overrideMap[e.Key]; exists && overVal != e.Value {
			conflict := MergeConflict{Key: e.Key, BaseValue: e.Value, OverValue: overVal}
			result.Conflicts = append(result.Conflicts, conflict)

			switch strategy {
			case MergeStrategyError:
				return nil, fmt.Errorf("merge conflict on key %q: base=%q override=%q", e.Key, e.Value, overVal)
			case MergeStrategyPreferOverride:
				result.Entries = append(result.Entries, Entry{Key: e.Key, Value: overVal})
			default: // MergeStrategyKeepBase
				result.Entries = append(result.Entries, e)
			}
		} else {
			result.Entries = append(result.Entries, e)
		}
	}

	// Append keys only present in override.
	for _, e := range override.Entries {
		if !seen[e.Key] {
			result.Entries = append(result.Entries, e)
		}
	}

	return result, nil
}
