package envfile

import (
	"testing"
)

func makeRotateEnv(t *testing.T, pairs ...string) *EnvFile {
	t.Helper()
	lines := ""
	for i := 0; i+1 < len(pairs); i += 2 {
		lines += pairs[i] + "=" + pairs[i+1] + "\n"
	}
	return writeTempAndParse(t, lines)
}

func writeTempAndParse(t *testing.T, content string) *EnvFile {
	t.Helper()
	path := writeTempEnv(t, content)
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return ef
}

func TestRotate_ReplacesExistingKeys(t *testing.T) {
	ef := makeRotateEnv(t, "DB_PASS", "old123", "API_KEY", "oldkey")
	opts := RotateOptions{
		Replacements: map[string]string{"DB_PASS": "new456", "API_KEY": "newkey"},
	}
	results, err := Rotate(ef, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Rotated {
			t.Errorf("key %s should be rotated", r.Key)
		}
	}
	entry, _ := ef.Lookup("DB_PASS")
	if entry.Value != "new456" {
		t.Errorf("expected new456, got %s", entry.Value)
	}
}

func TestRotate_SkipsMissingKeys(t *testing.T) {
	ef := makeRotateEnv(t, "EXISTING", "val")
	opts := RotateOptions{
		Replacements: map[string]string{"MISSING": "newval"},
	}
	results, err := Rotate(ef, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected skipped result for missing key")
	}
}

func TestRotate_DryRunDoesNotMutate(t *testing.T) {
	ef := makeRotateEnv(t, "SECRET", "original")
	opts := RotateOptions{
		Replacements: map[string]string{"SECRET": "rotated"},
		DryRun:       true,
	}
	_, err := Rotate(ef, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry, _ := ef.Lookup("SECRET")
	if entry.Value != "original" {
		t.Errorf("dry run should not mutate; got %s", entry.Value)
	}
}

func TestRotate_MaskValues(t *testing.T) {
	ef := makeRotateEnv(t, "TOKEN", "supersecret")
	opts := RotateOptions{
		Replacements: map[string]string{"TOKEN": "newtoken"},
		MaskValues:   true,
	}
	results, err := Rotate(ef, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].OldValue != "***" || results[0].NewValue != "***" {
		t.Errorf("expected masked values, got old=%s new=%s", results[0].OldValue, results[0].NewValue)
	}
}

func TestRotate_NilEnvFile(t *testing.T) {
	_, err := Rotate(nil, RotateOptions{Replacements: map[string]string{"K": "v"}})
	if err == nil {
		t.Error("expected error for nil EnvFile")
	}
}
