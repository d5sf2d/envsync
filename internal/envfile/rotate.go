package envfile

import (
	"fmt"
	"time"
)

// RotateResult describes the outcome of rotating a single key.
type RotateResult struct {
	Key      string
	OldValue string
	NewValue string
	Rotated  bool
	Skipped  bool
	Reason   string
}

// RotateOptions controls how key rotation behaves.
type RotateOptions struct {
	// Replacements maps key names to their new values.
	Replacements map[string]string
	// DryRun previews changes without writing.
	DryRun bool
	// MaskValues hides old/new values in results.
	MaskValues bool
}

// Rotate applies new values to specified keys in an EnvFile.
// Keys not present in Replacements are left untouched.
func Rotate(ef *EnvFile, opts RotateOptions) ([]RotateResult, error) {
	if ef == nil {
		return nil, fmt.Errorf("rotate: nil EnvFile")
	}
	if opts.Replacements == nil {
		return nil, fmt.Errorf("rotate: no replacements provided")
	}

	results := make([]RotateResult, 0, len(opts.Replacements))

	for key, newVal := range opts.Replacements {
		res := RotateResult{Key: key, NewValue: newVal}

		entry, exists := ef.Lookup(key)
		if !exists {
			res.Skipped = true
			res.Reason = "key not found"
			results = append(results, res)
			continue
		}

		res.OldValue = entry.Value
		if opts.MaskValues {
			res.OldValue = "***"
			res.NewValue = "***"
		}

		if !opts.DryRun {
			ef.Set(key, newVal)
		}

		res.Rotated = true
		results = append(results, res)
	}

	if !opts.DryRun && len(results) > 0 {
		ef.UpdatedAt = time.Now()
	}

	return results, nil
}
