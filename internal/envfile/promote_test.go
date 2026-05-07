package envfile

import (
	"bytes"
	"strings"
	"testing"
)

func makePromoteEnv(pairs ...string) *EnvFile {
	ef := &EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestPromote_NewKeys(t *testing.T) {
	src := makePromoteEnv("NEW_KEY", "hello", "ANOTHER", "world")
	dst := makePromoteEnv("EXISTING", "keep")

	result, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(result.Promoted))
	}
	if len(dst.Entries) != 3 {
		t.Errorf("expected 3 entries in dst, got %d", len(dst.Entries))
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := makePromoteEnv("KEY", "new_value")
	dst := makePromoteEnv("KEY", "old_value")

	result, err := Promote(src, dst, PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if dst.Entries[0].Value != "old_value" {
		t.Errorf("value should not be overwritten")
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := makePromoteEnv("KEY", "new_value")
	dst := makePromoteEnv("KEY", "old_value")

	result, err := Promote(src, dst, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(result.Overwritten))
	}
	if dst.Entries[0].Value != "new_value" {
		t.Errorf("expected value to be overwritten")
	}
}

func TestPromote_OnlyKeys(t *testing.T) {
	src := makePromoteEnv("A", "1", "B", "2", "C", "3")
	dst := makePromoteEnv()

	result, err := Promote(src, dst, PromoteOptions{OnlyKeys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(result.Promoted))
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
}

func TestPromote_DryRun(t *testing.T) {
	src := makePromoteEnv("X", "val")
	dst := makePromoteEnv()

	_, err := Promote(src, dst, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst.Entries) != 0 {
		t.Errorf("dry run should not modify dst")
	}
}

func TestWritePromoteReport(t *testing.T) {
	result := PromoteResult{
		Promoted:    []string{"NEW"},
		Overwritten: []string{"CHANGED"},
		Skipped:     []string{"OLD"},
	}
	var buf bytes.Buffer
	WritePromoteReport(&buf, result, false)
	out := buf.String()
	if !strings.Contains(out, "promoted") {
		t.Errorf("expected 'promoted' in output")
	}
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' in output")
	}
}
