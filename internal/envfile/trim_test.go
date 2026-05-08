package envfile

import (
	"testing"
)

func makeTrimEnv(entries []Entry) *EnvFile {
	return &EnvFile{Entries: append([]Entry(nil), entries...)}
}

func TestTrim_NoChanges(t *testing.T) {
	ef := makeTrimEnv([]Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	})
	results, err := Trim(ef, TrimOptions{TrimValues: true, TrimKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no change for key %q", r.Key)
		}
	}
}

func TestTrim_TrimsValues(t *testing.T) {
	ef := makeTrimEnv([]Entry{
		{Key: "HOST", Value: "  localhost  "},
		{Key: "PORT", Value: "8080"},
	})
	_, err := Trim(ef, TrimOptions{TrimValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Entries[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", ef.Entries[0].Value)
	}
	if ef.Entries[1].Value != "8080" {
		t.Errorf("expected '8080', got %q", ef.Entries[1].Value)
	}
}

func TestTrim_TrimsKeys(t *testing.T) {
	ef := makeTrimEnv([]Entry{
		{Key: "  HOST  ", Value: "localhost"},
	})
	_, err := Trim(ef, TrimOptions{TrimKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Entries[0].Key != "HOST" {
		t.Errorf("expected key 'HOST', got %q", ef.Entries[0].Key)
	}
}

func TestTrim_NormalizeEmpty(t *testing.T) {
	ef := makeTrimEnv([]Entry{
		{Key: "EMPTY", Value: "   "},
		{Key: "HOST", Value: "localhost"},
	})
	results, err := Trim(ef, TrimOptions{TrimValues: true, NormalizeEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Entries[0].Value != "" {
		t.Errorf("expected empty string, got %q", ef.Entries[0].Value)
	}
	if !results[0].Changed {
		t.Error("expected EMPTY to be marked changed")
	}
}

func TestTrim_DryRunDoesNotMutate(t *testing.T) {
	ef := makeTrimEnv([]Entry{
		{Key: "HOST", Value: "  localhost  "},
	})
	results, err := Trim(ef, TrimOptions{TrimValues: true, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Changed {
		t.Error("expected Changed=true in dry run result")
	}
	if ef.Entries[0].Value != "  localhost  " {
		t.Errorf("dry run mutated value: got %q", ef.Entries[0].Value)
	}
}

func TestTrim_NilEnvFile(t *testing.T) {
	results, err := Trim(nil, TrimOptions{TrimValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results != nil {
		t.Error("expected nil results for nil EnvFile")
	}
}
