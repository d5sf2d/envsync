package envfile

import (
	"testing"
)

func makeCompareEnv(entries []Entry) *EnvFile {
	return &EnvFile{Entries: entries}
}

func TestCompare_IdenticalFiles(t *testing.T) {
	a := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}})
	b := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}})

	r := Compare(a, b)

	if r.MatchCount != 2 {
		t.Errorf("expected 2 matches, got %d", r.MatchCount)
	}
	if r.MismatchCount != 0 {
		t.Errorf("expected 0 mismatches, got %d", r.MismatchCount)
	}
	if r.Similarity != 1.0 {
		t.Errorf("expected similarity 1.0, got %f", r.Similarity)
	}
}

func TestCompare_MissingInB(t *testing.T) {
	a := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}, {Key: "ONLY_A", Value: "x"}})
	b := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}})

	r := Compare(a, b)

	if len(r.MissingInB) != 1 || r.MissingInB[0] != "ONLY_A" {
		t.Errorf("expected ONLY_A missing in B, got %v", r.MissingInB)
	}
	if len(r.MissingInA) != 0 {
		t.Errorf("expected no keys missing in A, got %v", r.MissingInA)
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}})
	b := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}, {Key: "ONLY_B", Value: "y"}})

	r := Compare(a, b)

	if len(r.MissingInA) != 1 || r.MissingInA[0] != "ONLY_B" {
		t.Errorf("expected ONLY_B missing in A, got %v", r.MissingInA)
	}
}

func TestCompare_Mismatches(t *testing.T) {
	a := makeCompareEnv([]Entry{{Key: "FOO", Value: "old"}})
	b := makeCompareEnv([]Entry{{Key: "FOO", Value: "new"}})

	r := Compare(a, b)

	if r.MismatchCount != 1 {
		t.Errorf("expected 1 mismatch, got %d", r.MismatchCount)
	}
	vals, ok := r.Mismatches["FOO"]
	if !ok {
		t.Fatal("expected FOO in mismatches")
	}
	if vals[0] != "old" || vals[1] != "new" {
		t.Errorf("unexpected mismatch values: %v", vals)
	}
	if r.Similarity != 0.0 {
		t.Errorf("expected similarity 0.0, got %f", r.Similarity)
	}
}

func TestCompare_EmptyFiles(t *testing.T) {
	a := makeCompareEnv([]Entry{})
	b := makeCompareEnv([]Entry{})

	r := Compare(a, b)

	if r.Similarity != 1.0 {
		t.Errorf("expected similarity 1.0 for two empty files, got %f", r.Similarity)
	}
}

func TestCompare_Summary(t *testing.T) {
	a := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}})
	b := makeCompareEnv([]Entry{{Key: "FOO", Value: "bar"}})

	r := Compare(a, b)
	s := r.Summary()

	if s == "" {
		t.Error("expected non-empty summary")
	}
}
