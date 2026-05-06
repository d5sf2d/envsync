package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeSnapshotEnv(t *testing.T, content string) *EnvFile {
	t.Helper()
	f := writeTempEnv(t, content)
	ef, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return ef
}

func TestTakeSnapshot_CapturesEntries(t *testing.T) {
	ef := makeSnapshotEnv(t, "FOO=bar\nBAZ=qux\n")
	snap := TakeSnapshot(ef, "test.env")

	if snap.Source != "test.env" {
		t.Errorf("expected source %q, got %q", "test.env", snap.Source)
	}
	if snap.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", snap.Entries["FOO"])
	}
	if snap.Entries["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", snap.Entries["BAZ"])
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	ef := makeSnapshotEnv(t, "KEY=value\nSECRET=s3cr3t\n")
	snap := TakeSnapshot(ef, "prod.env")

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := SaveSnapshot(snap, tmp); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(tmp)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if loaded.Source != snap.Source {
		t.Errorf("source mismatch: want %q got %q", snap.Source, loaded.Source)
	}
	if loaded.Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Entries["KEY"])
	}
	if loaded.Entries["SECRET"] != "s3cr3t" {
		t.Errorf("expected SECRET=s3cr3t, got %q", loaded.Entries["SECRET"])
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	base := &Snapshot{
		Entries: map[string]string{"FOO": "old", "KEEP": "same", "GONE": "bye"},
	}
	head := &Snapshot{
		Entries: map[string]string{"FOO": "new", "KEEP": "same", "ADDED": "hi"},
	}

	diffs := DiffSnapshot(base, head)

	statuses := map[string]DiffStatus{}
	for _, d := range diffs {
		statuses[d.Key] = d.Status
	}

	if statuses["FOO"] != StatusChanged {
		t.Errorf("expected FOO changed, got %v", statuses["FOO"])
	}
	if statuses["ADDED"] != StatusAdded {
		t.Errorf("expected ADDED added, got %v", statuses["ADDED"])
	}
	if statuses["GONE"] != StatusRemoved {
		t.Errorf("expected GONE removed, got %v", statuses["GONE"])
	}
	if _, ok := statuses["KEEP"]; ok {
		t.Error("expected KEEP to be absent from diffs")
	}
}

func TestSaveSnapshot_BadPath(t *testing.T) {
	snap := &Snapshot{Entries: map[string]string{}}
	err := SaveSnapshot(snap, filepath.Join(t.TempDir(), "nosuchdir", "snap.json"))
	if err == nil {
		t.Error("expected error for bad path, got nil")
	}
	_ = os.Remove("snap.json")
}
