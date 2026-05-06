package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two env files.
type CompareResult struct {
	MatchCount    int
	MismatchCount int
	MissingInA    []string
	MissingInB    []string
	Mismatches    map[string][2]string // key -> [valueA, valueB]
	Similarity    float64              // 0.0 - 1.0
}

// Compare performs a detailed comparison between two EnvFile maps,
// returning a CompareResult with similarity scoring.
func Compare(a, b *EnvFile) CompareResult {
	result := CompareResult{
		Mismatches: make(map[string][2]string),
	}

	aMap := make(map[string]string, len(a.Entries))
	for _, e := range a.Entries {
		aMap[e.Key] = e.Value
	}

	bMap := make(map[string]string, len(b.Entries))
	for _, e := range b.Entries {
		bMap[e.Key] = e.Value
	}

	allKeys := make(map[string]struct{})
	for k := range aMap {
		allKeys[k] = struct{}{}
	}
	for k := range bMap {
		allKeys[k] = struct{}{}
	}

	for k := range allKeys {
		va, inA := aMap[k]
		vb, inB := bMap[k]

		switch {
		case inA && inB:
			if va == vb {
				result.MatchCount++
			} else {
				result.MismatchCount++
				result.Mismatches[k] = [2]string{va, vb}
			}
		case inA && !inB:
			result.MissingInB = append(result.MissingInB, k)
		case !inA && inB:
			result.MissingInA = append(result.MissingInA, k)
		}
	}

	sort.Strings(result.MissingInA)
	sort.Strings(result.MissingInB)

	total := len(allKeys)
	if total > 0 {
		result.Similarity = float64(result.MatchCount) / float64(total)
	} else {
		result.Similarity = 1.0
	}

	return result
}

// Summary returns a human-readable one-line summary of the comparison.
func (r CompareResult) Summary() string {
	parts := []string{
		fmt.Sprintf("%d matching", r.MatchCount),
		fmt.Sprintf("%d mismatched", r.MismatchCount),
		fmt.Sprintf("%d missing in A", len(r.MissingInA)),
		fmt.Sprintf("%d missing in B", len(r.MissingInB)),
		fmt.Sprintf("similarity=%.0f%%", r.Similarity*100),
	}
	return strings.Join(parts, ", ")
}
