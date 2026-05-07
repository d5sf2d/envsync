package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

var interpolatePattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

// InterpolateResult holds the result of an interpolation operation.
type InterpolateResult struct {
	Key          string
	Original     string
	Interpolated string
	Unresolved   []string
}

// Interpolate expands variable references within env file values using the
// entries defined in the same file. References of the form ${VAR} or $VAR
// are supported. Returns a new EnvFile with expanded values and a slice of
// InterpolateResult describing each substitution made.
func Interpolate(ef *EnvFile) (*EnvFile, []InterpolateResult, error) {
	lookup := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		lookup[e.Key] = e.Value
	}

	results := make([]InterpolateResult, 0)
	cloned := cloneEntries(ef.Entries)

	for i, entry := range cloned {
		if !strings.ContainsAny(entry.Value, "$") {
			continue
		}

		unresolved := []string{}
		expanded := interpolatePattern.ReplaceAllStringFunc(entry.Value, func(match string) string {
			varName := strings.TrimPrefix(strings.TrimPrefix(strings.Trim(match, "${}"), "${"), "$")
			varName = strings.TrimSuffix(varName, "}")
			if val, ok := lookup[varName]; ok {
				return val
			}
			unresolved = append(unresolved, varName)
			return match
		})

		if len(unresolved) > 0 {
			return nil, nil, fmt.Errorf("interpolate: key %q references unresolved variables: %s",
				entry.Key, strings.Join(unresolved, ", "))
		}

		results = append(results, InterpolateResult{
			Key:          entry.Key,
			Original:     entry.Value,
			Interpolated: expanded,
			Unresolved:   unresolved,
		})
		cloned[i].Value = expanded
		lookup[entry.Key] = expanded
	}

	return &EnvFile{Path: ef.Path, Entries: cloned}, results, nil
}
