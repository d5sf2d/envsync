package envfile

import (
	"testing"
)

func makeMergeEnv(entries []Entry) *EnvFile {
	return &EnvFile{Entries: entries}
}

func TestMerge_NoConflicts(t *testing.T) {
	base := makeMergeEnv([]Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}})
	over := makeMergeEnv([]Entry{{Key: "C", Value: "3"}})

	res, err := Merge(base, over, MergeStrategyKeepBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Entries))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_KeepBase(t *testing.T) {
	base := makeMergeEnv([]Entry{{Key: "HOST", Value: "localhost"}})
	over := makeMergeEnv([]Entry{{Key: "HOST", Value: "prod.example.com"}})

	res, err := Merge(base, over, MergeStrategyKeepBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Entries[0].Value != "localhost" {
		t.Errorf("expected base value 'localhost', got %q", res.Entries[0].Value)
	}
}

func TestMerge_PreferOverride(t *testing.T) {
	base := makeMergeEnv([]Entry{{Key: "HOST", Value: "localhost"}})
	over := makeMergeEnv([]Entry{{Key: "HOST", Value: "prod.example.com"}})

	res, err := Merge(base, over, MergeStrategyPreferOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "prod.example.com" {
		t.Errorf("expected override value 'prod.example.com', got %q", res.Entries[0].Value)
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	base := makeMergeEnv([]Entry{{Key: "SECRET", Value: "abc"}})
	over := makeMergeEnv([]Entry{{Key: "SECRET", Value: "xyz"}})

	_, err := Merge(base, over, MergeStrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_OrderPreserved(t *testing.T) {
	base := makeMergeEnv([]Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}})
	over := makeMergeEnv([]Entry{{Key: "C", Value: "3"}, {Key: "D", Value: "4"}})

	res, err := Merge(base, over, MergeStrategyKeepBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := []string{"A", "B", "C", "D"}
	for i, e := range res.Entries {
		if e.Key != keys[i] {
			t.Errorf("position %d: expected key %q, got %q", i, keys[i], e.Key)
		}
	}
}
