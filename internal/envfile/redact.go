package envfile

import (
	"regexp"
	"strings"
)

// RedactOptions controls how sensitive values are redacted.
type RedactOptions struct {
	// Patterns are additional key patterns (regex) to treat as sensitive.
	Patterns []string
	// Placeholder replaces redacted values. Defaults to "[REDACTED]".
	Placeholder string
}

var defaultSensitivePatterns = []string{
	`(?i)secret`,
	`(?i)password`,
	`(?i)passwd`,
	`(?i)token`,
	`(?i)api[_]?key`,
	`(?i)private[_]?key`,
	`(?i)auth`,
	`(?i)credential`,
}

// Redactor replaces sensitive values in an EnvFile.
type Redactor struct {
	compiled    []*regexp.Regexp
	placeholder string
}

// NewRedactor creates a Redactor from the given options.
func NewRedactor(opts RedactOptions) (*Redactor, error) {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}

	patterns := append(defaultSensitivePatterns, opts.Patterns...)
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}

	return &Redactor{compiled: compiled, placeholder: placeholder}, nil
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	for _, re := range r.compiled {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// Redact returns a copy of the EnvFile with sensitive values replaced.
func (r *Redactor) Redact(ef EnvFile) EnvFile {
	redacted := make(EnvFile, len(ef))
	for i, entry := range ef {
		if r.IsSensitive(strings.TrimSpace(entry.Key)) {
			redacted[i] = Entry{
				Key:     entry.Key,
				Value:   r.placeholder,
				Comment: entry.Comment,
			}
		} else {
			redacted[i] = entry
		}
	}
	return redacted
}
