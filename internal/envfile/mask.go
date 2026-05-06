package envfile

import "strings"

// DefaultSecretPatterns holds common key substrings that indicate secret values.
var DefaultSecretPatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL",
	"AUTH", "PWD", "CERT", "KEY",
}

const maskedValue = "***"

// Masker controls which keys are treated as secrets.
type Masker struct {
	Patterns []string
}

// NewMasker creates a Masker with the default secret patterns.
func NewMasker() *Masker {
	return &Masker{Patterns: DefaultSecretPatterns}
}

// IsSecret returns true if the key matches any secret pattern.
func (m *Masker) IsSecret(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range m.Patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked placeholder if the key is a secret, otherwise the original value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSecret(key) {
		return maskedValue
	}
	return value
}

// MaskDiff returns a copy of the DiffResult with secret values replaced by the mask.
func (m *Masker) MaskDiff(d *DiffResult) *DiffResult {
	masked := make([]DiffEntry, len(d.Entries))
	for i, e := range d.Entries {
		me := e
		if m.IsSecret(e.Key) {
			if me.SourceValue != "" {
				me.SourceValue = maskedValue
			}
			if me.TargetValue != "" {
				me.TargetValue = maskedValue
			}
		}
		masked[i] = me
	}
	return &DiffResult{Entries: masked}
}
