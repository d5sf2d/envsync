package envfile

import "fmt"

// PromoteOptions controls how values are promoted between environments.
type PromoteOptions struct {
	// OnlyKeys restricts promotion to these keys (nil = all keys).
	OnlyKeys []string
	// SkipKeys excludes these keys from promotion.
	SkipKeys []string
	// DryRun returns the result without writing to disk.
	DryRun bool
	// Overwrite replaces existing values in the destination.
	Overwrite bool
}

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	Overwritten []string
}

// Promote copies selected keys from src into dst according to opts.
// When DryRun is false the destination EnvFile entries are mutated in place;
// the caller is responsible for persisting the file.
func Promote(src, dst *EnvFile, opts PromoteOptions) (PromoteResult, error) {
	if src == nil {
		return PromoteResult{}, fmt.Errorf("promote: source EnvFile is nil")
	}
	if dst == nil {
		return PromoteResult{}, fmt.Errorf("promote: destination EnvFile is nil")
	}

	skipSet := make(map[string]bool, len(opts.SkipKeys))
	for _, k := range opts.SkipKeys {
		skipSet[k] = true
	}

	onlySet := make(map[string]bool, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		onlySet[k] = true
	}

	dstIndex := make(map[string]int, len(dst.Entries))
	for i, e := range dst.Entries {
		dstIndex[e.Key] = i
	}

	var result PromoteResult

	for _, entry := range src.Entries {
		if entry.Key == "" {
			continue
		}
		if skipSet[entry.Key] {
			result.Skipped = append(result.Skipped, entry.Key)
			continue
		}
		if len(onlySet) > 0 && !onlySet[entry.Key] {
			result.Skipped = append(result.Skipped, entry.Key)
			continue
		}

		if idx, exists := dstIndex[entry.Key]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, entry.Key)
				continue
			}
			if !opts.DryRun {
				dst.Entries[idx].Value = entry.Value
			}
			result.Overwritten = append(result.Overwritten, entry.Key)
		} else {
			if !opts.DryRun {
				dst.Entries = append(dst.Entries, entry)
				dstIndex[entry.Key] = len(dst.Entries) - 1
			}
			result.Promoted = append(result.Promoted, entry.Key)
		}
	}

	return result, nil
}
