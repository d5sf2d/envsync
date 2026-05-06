package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeSyncEnv(entries []Entry) *EnvFile {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return &EnvFile{Entries: entries, index: m}
}

func TestSync_AddMissing(t *testing.T) {
	dst := makeSyncEnv([]Entry{{Key: "A", Value: "1"}})
	src := makeSyncEnv([]Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}})

	tmp := filepath.Join(t.TempDir(), ".env")
	res, err := Sync(dst, src, tmp, SyncOptions{AddMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Errorf("expected B to be added, got %v", res.Added)
	}
	parsed, _ := Parse(tmp)
	if parsed.index["B"] != "2" {
		t.Errorf("expected B=2 in written file")
	}
}

func TestSync_Overwrite(t *testing.T) {
	dst := makeSyncEnv([]Entry{{Key: "A", Value: "old"}})
	src := makeSyncEnv([]Entry{{Key: "A", Value: "new"}})

	tmp := filepath.Join(t.TempDir(), ".env")
	res, err := Sync(dst, src, tmp, SyncOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "A" {
		t.Errorf("expected A to be updated, got %v", res.Updated)
	}
	parsed, _ := Parse(tmp)
	if parsed.index["A"] != "new" {
		t.Errorf("expected A=new in written file")
	}
}

func TestSync_DryRun(t *testing.T) {
	dst := makeSyncEnv([]Entry{{Key: "A", Value: "1"}})
	src := makeSyncEnv([]Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}})

	tmp := filepath.Join(t.TempDir(), ".env")
	_, err := Sync(dst, src, tmp, SyncOptions{AddMissing: true, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(tmp); !os.IsNotExist(statErr) {
		t.Errorf("dry run should not create file")
	}
}

func TestSync_SkipsWithoutFlags(t *testing.T) {
	dst := makeSyncEnv([]Entry{{Key: "A", Value: "old"}})
	src := makeSyncEnv([]Entry{{Key: "A", Value: "new"}, {Key: "B", Value: "2"}})

	tmp := filepath.Join(t.TempDir(), ".env")
	res, err := Sync(dst, src, tmp, SyncOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(res.Skipped))
	}
}
