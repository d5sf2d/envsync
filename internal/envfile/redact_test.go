package envfile

import (
	"testing"
)

func makeRedactEnv(pairs ...string) EnvFile {
	var ef EnvFile
	for i := 0; i+1 < len(pairs); i += 2 {
		ef = append(ef, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestNewRedactor_DefaultPatterns(t *testing.T) {
	r, err := NewRedactor(RedactOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil redactor")
	}
}

func TestNewRedactor_InvalidPattern(t *testing.T) {
	_, err := NewRedactor(RedactOptions{Patterns: []string{"[invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestIsSensitive_MatchesDefaults(t *testing.T) {
	r, _ := NewRedactor(RedactOptions{})
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"SECRET_TOKEN", true},
		{"AUTH_HEADER", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"DEBUG", false},
	}
	for _, tc := range cases {
		got := r.IsSensitive(tc.key)
		if got != tc.expected {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}

func TestRedact_ReplacesValues(t *testing.T) {
	r, _ := NewRedactor(RedactOptions{Placeholder: "***"})
	ef := makeRedactEnv(
		"APP_NAME", "myapp",
		"DB_PASSWORD", "supersecret",
		"PORT", "8080",
		"API_KEY", "abc123",
	)

	result := r.Redact(ef)

	if result[0].Value != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", result[0].Value)
	}
	if result[1].Value != "***" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", result[1].Value)
	}
	if result[2].Value != "8080" {
		t.Errorf("expected PORT unchanged, got %q", result[2].Value)
	}
	if result[3].Value != "***" {
		t.Errorf("expected API_KEY redacted, got %q", result[3].Value)
	}
}

func TestRedact_DefaultPlaceholder(t *testing.T) {
	r, _ := NewRedactor(RedactOptions{})
	ef := makeRedactEnv("SECRET", "topsecret")
	result := r.Redact(ef)
	if result[0].Value != "[REDACTED]" {
		t.Errorf("expected default placeholder, got %q", result[0].Value)
	}
}

func TestRedact_CustomPattern(t *testing.T) {
	r, _ := NewRedactor(RedactOptions{Patterns: []string{`(?i)internal`}})
	ef := makeRedactEnv("INTERNAL_URL", "http://internal", "PUBLIC_URL", "http://public")
	result := r.Redact(ef)
	if result[0].Value != "[REDACTED]" {
		t.Errorf("expected INTERNAL_URL redacted, got %q", result[0].Value)
	}
	if result[1].Value != "http://public" {
		t.Errorf("expected PUBLIC_URL unchanged, got %q", result[1].Value)
	}
}
