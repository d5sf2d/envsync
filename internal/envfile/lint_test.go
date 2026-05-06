package envfile

import (
	"testing"
)

func makeLintEnv(entries []Entry) *EnvFile {
	return &EnvFile{Entries: entries}
}

func TestLint_CleanFile(t *testing.T) {
	ef := makeLintEnv([]Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "PORT", Value: "8080"},
	})
	result := Lint(ef)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	ef := makeLintEnv([]Entry{
		{Key: "SECRET", Value: ""},
	})
	result := Lint(ef)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != LintWarn {
		t.Errorf("expected WARN, got %s", result.Issues[0].Severity)
	}
	if result.Issues[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", result.Issues[0].Key)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	ef := makeLintEnv([]Entry{
		{Key: "app_name", Value: "test"},
	})
	result := Lint(ef)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != LintWarn {
		t.Errorf("expected WARN, got %s", result.Issues[0].Severity)
	}
}

func TestLint_WhitespaceValue(t *testing.T) {
	ef := makeLintEnv([]Entry{
		{Key: "HOST", Value: " localhost "},
	})
	result := Lint(ef)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != LintError {
		t.Errorf("expected ERROR, got %s", result.Issues[0].Severity)
	}
	if !result.HasErrors() {
		t.Error("expected HasErrors() to return true")
	}
}

func TestLint_SkipsComments(t *testing.T) {
	ef := makeLintEnv([]Entry{
		{IsComment: true, Raw: "# this is a comment"},
		{Key: "VALID", Value: "yes"},
	})
	result := Lint(ef)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(result.Issues))
	}
}

func TestLint_HasErrors_False(t *testing.T) {
	ef := makeLintEnv([]Entry{
		{Key: "lower_key", Value: "val"},
	})
	result := Lint(ef)
	if result.HasErrors() {
		t.Error("expected HasErrors() to return false for WARN-only issues")
	}
}
