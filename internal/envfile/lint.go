package envfile

import (
	"fmt"
	"strings"
)

// LintSeverity represents the severity level of a lint issue.
type LintSeverity string

const (
	LintWarn  LintSeverity = "WARN"
	LintError LintSeverity = "ERROR"
)

// LintIssue describes a single linting problem found in an env file.
type LintIssue struct {
	Line     int
	Key      string
	Message  string
	Severity LintSeverity
}

func (i LintIssue) String() string {
	if i.Line > 0 {
		return fmt.Sprintf("%s [line %d] %s: %s", i.Severity, i.Line, i.Key, i.Message)
	}
	return fmt.Sprintf("%s %s: %s", i.Severity, i.Key, i.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == LintError {
			return true
		}
	}
	return false
}

// Lint checks an EnvFile for style and quality issues beyond basic validation.
func Lint(ef *EnvFile) *LintResult {
	result := &LintResult{}

	for i, entry := range ef.Entries {
		lineNum := i + 1

		// Warn on empty values
		if entry.Value == "" && !entry.IsComment && entry.Key != "" {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  "value is empty",
				Severity: LintWarn,
			})
		}

		// Warn on keys that are not uppercase
		if entry.Key != "" && entry.Key != strings.ToUpper(entry.Key) {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  "key is not uppercase",
				Severity: LintWarn,
			})
		}

		// Error on values with unquoted leading/trailing whitespace
		if entry.Value != strings.TrimSpace(entry.Value) {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  "value has unquoted leading or trailing whitespace",
				Severity: LintError,
			})
		}
	}

	return result
}
