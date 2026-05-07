package envfile

import (
	"testing"
)

func makeRenameEnv(entries []Entry) *EnvFile {
	return &EnvFile{Entries: entries}
}

func TestRename_BasicRename(t *testing.T) {
	ef := makeRenameEnv([]Entry{
		{Key: "OLD_KEY", Value: "value1"},
		{Key: "KEEP", Value: "value2"},
	})
	results, err := Rename(ef, map[string]string{"OLD_KEY": "NEW_KEY"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Skipped {
		t.Fatalf("expected 1 non-skipped result, got %+v", results)
	}
	if ef.Entries[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %q", ef.Entries[0].Key)
	}
	if ef.Entries[1].Key != "KEEP" {
		t.Errorf("KEEP should be unchanged")
	}
}

func TestRename_MissingSourceSkips(t *testing.T) {
	ef := makeRenameEnv([]Entry{{Key: "A", Value: "1"}})
	results, err := Rename(ef, map[string]string{"MISSING": "B"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected skipped result for missing key")
	}
}

func TestRename_FailOnMissing(t *testing.T) {
	ef := makeRenameEnv([]Entry{{Key: "A", Value: "1"}})
	_, err := Rename(ef, map[string]string{"MISSING": "B"}, RenameOptions{FailOnMissing: true})
	if err == nil {
		t.Error("expected error for missing source key")
	}
}

func TestRename_CollisionSkips(t *testing.T) {
	ef := makeRenameEnv([]Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	})
	results, err := Rename(ef, map[string]string{"A": "B"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected skipped result on collision")
	}
}

func TestRename_FailOnCollision(t *testing.T) {
	ef := makeRenameEnv([]Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	})
	_, err := Rename(ef, map[string]string{"A": "B"}, RenameOptions{FailOnCollision: true})
	if err == nil {
		t.Error("expected error on collision")
	}
}

func TestRename_DryRunDoesNotMutate(t *testing.T) {
	ef := makeRenameEnv([]Entry{{Key: "OLD", Value: "val"}})
	results, err := Rename(ef, map[string]string{"OLD": "NEW"}, RenameOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Skipped {
		t.Error("dry-run result should not be marked skipped")
	}
	if ef.Entries[0].Key != "OLD" {
		t.Errorf("dry-run must not mutate: got %q", ef.Entries[0].Key)
	}
}
