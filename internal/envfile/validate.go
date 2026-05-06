package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a single validation issue found in an env file.
type ValidationError struct {
	Line    int
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: %s: %s", e.Line, e.Key, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) add(line int, key, msg string) {
	r.Errors = append(r.Errors, ValidationError{Line: line, Key: key, Message: msg})
}

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks an EnvFile for common issues:
// - duplicate keys
// - invalid key names
// - empty values for required keys (if requiredKeys is non-nil)
func Validate(ef *EnvFile, requiredKeys []string) *ValidationResult {
	result := &ValidationResult{}
	seen := make(map[string]int) // key -> first line number

	for i, entry := range ef.Entries {
		lineNum := i + 1

		if !validKeyRe.MatchString(entry.Key) {
			result.add(lineNum, entry.Key, "invalid key name (must match [A-Za-z_][A-Za-z0-9_]*)")
		}

		if first, dup := seen[entry.Key]; dup {
			result.add(lineNum, entry.Key, fmt.Sprintf("duplicate key (first defined at line %d)", first))
		} else {
			seen[entry.Key] = lineNum
		}
	}

	for _, req := range requiredKeys {
		val, ok := ef.Get(req)
		if !ok {
			result.add(0, req, "required key is missing")
		} else if strings.TrimSpace(val) == "" {
			result.add(0, req, "required key has an empty value")
		}
	}

	return result
}
