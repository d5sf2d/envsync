package envfile

import (
	"testing"
)

func makeInterpolateEnv(entries []Entry) *EnvFile {
	return &EnvFile{Path: "test.env", Entries: entries}
}

func TestInterpolate_NoReferences(t *testing.T) {
	ef := makeInterpolateEnv([]Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
	})
	out, results, err := Interpolate(ef)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
	if out.Entries[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", out.Entries[0].Value)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	ef := makeInterpolateEnv([]Entry{
		{Key: "HOST", Value: "db.example.com"},
		{Key: "DSN", Value: "postgres://${HOST}/mydb"},
	})
	out, results, err := Interpolate(ef)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Interpolated != "postgres://db.example.com/mydb" {
		t.Errorf("unexpected interpolation: %q", results[0].Interpolated)
	}
	if out.Entries[1].Value != "postgres://db.example.com/mydb" {
		t.Errorf("entry not updated: %q", out.Entries[1].Value)
	}
}

func TestInterpolate_BareStyle(t *testing.T) {
	ef := makeInterpolateEnv([]Entry{
		{Key: "USER", Value: "admin"},
		{Key: "GREETING", Value: "Hello $USER"},
	})
	out, _, err := Interpolate(ef)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[1].Value != "Hello admin" {
		t.Errorf("expected 'Hello admin', got %q", out.Entries[1].Value)
	}
}

func TestInterpolate_ChainedReferences(t *testing.T) {
	ef := makeInterpolateEnv([]Entry{
		{Key: "SCHEME", Value: "https"},
		{Key: "HOST", Value: "${SCHEME}://api.example.com"},
		{Key: "URL", Value: "${HOST}/v1"},
	})
	out, _, err := Interpolate(ef)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[2].Value != "https://api.example.com/v1" {
		t.Errorf("chained interpolation failed: %q", out.Entries[2].Value)
	}
}

func TestInterpolate_UnresolvedReturnsError(t *testing.T) {
	ef := makeInterpolateEnv([]Entry{
		{Key: "DSN", Value: "postgres://${MISSING_HOST}/db"},
	})
	_, _, err := Interpolate(ef)
	if err == nil {
		t.Fatal("expected error for unresolved variable, got nil")
	}
}

func TestInterpolate_OriginalPreserved(t *testing.T) {
	ef := makeInterpolateEnv([]Entry{
		{Key: "BASE", Value: "example.com"},
		{Key: "FULL", Value: "https://${BASE}"},
	})
	_, results, err := Interpolate(ef)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Original != "https://${BASE}" {
		t.Errorf("original not preserved: %q", results[0].Original)
	}
}
