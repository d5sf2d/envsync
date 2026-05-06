package envfile

import (
	"testing"
)

func makeEnvFile(entries map[string]string) *EnvFile {
	var e []Entry
	for k, v := range entries {
		e = append(e, Entry{Key: k, Value: v})
	}
	return &EnvFile{Entries: e}
}

func TestDiff_Added(t *testing.T) {
	src := makeEnvFile(map[string]string{"A": "1"})
	tgt := makeEnvFile(map[string]string{"A": "1", "B": "2"})
	result := Diff(src, tgt)
	if !result.HasChanges() {
		t.Error("expected changes")
	}
	found := false
	for _, e := range result.Entries {
		if e.Key == "B" && e.Kind == Added {
			found = true
		}
	}
	if !found {
		t.Error("expected B to be added")
	}
}

func TestDiff_Removed(t *testing.T) {
	src := makeEnvFile(map[string]string{"A": "1", "B": "2"})
	tgt := makeEnvFile(map[string]string{"A": "1"})
	result := Diff(src, tgt)
	for _, e := range result.Entries {
		if e.Key == "B" && e.Kind != Removed {
			t.Errorf("expected B to be removed, got %s", e.Kind)
		}
	}
}

func TestDiff_Changed(t *testing.T) {
	src := makeEnvFile(map[string]string{"A": "old"})
	tgt := makeEnvFile(map[string]string{"A": "new"})
	result := Diff(src, tgt)
	for _, e := range result.Entries {
		if e.Key == "A" {
			if e.Kind != Changed {
				t.Errorf("expected Changed, got %s", e.Kind)
			}
			if e.SourceValue != "old" || e.TargetValue != "new" {
				t.Errorf("unexpected values: src=%s tgt=%s", e.SourceValue, e.TargetValue)
			}
		}
	}
}

func TestDiff_Unchanged(t *testing.T) {
	src := makeEnvFile(map[string]string{"A": "1"})
	tgt := makeEnvFile(map[string]string{"A": "1"})
	result := Diff(src, tgt)
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}
