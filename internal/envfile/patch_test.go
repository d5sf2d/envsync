package envfile

import (
	"testing"
)

func makePatchEnv(t *testing.T, pairs ...string) *EnvFile {
	t.Helper()
	ef := &EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Set(pairs[i], pairs[i+1])
	}
	return ef
}

func TestPatch_SetNewKey(t *testing.T) {
	ef := makePatchEnv(t)
	results, err := Patch(ef, []PatchInstruction{{Op: PatchSet, Key: "FOO", Value: "bar"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Errorf("expected Applied=true")
	}
	val, ok := ef.Lookup("FOO")
	if !ok || val != "bar" {
		t.Errorf("expected FOO=bar, got %q ok=%v", val, ok)
	}
}

func TestPatch_DeleteExistingKey(t *testing.T) {
	ef := makePatchEnv(t, "TO_DELETE", "value")
	results, err := Patch(ef, []PatchInstruction{{Op: PatchDelete, Key: "TO_DELETE"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Errorf("expected Applied=true")
	}
	if _, ok := ef.Lookup("TO_DELETE"); ok {
		t.Errorf("expected key to be deleted")
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	ef := makePatchEnv(t)
	results, err := Patch(ef, []PatchInstruction{{Op: PatchDelete, Key: "GHOST"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Errorf("expected Applied=false for missing key")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	ef := makePatchEnv(t, "OLD_NAME", "hello")
	results, err := Patch(ef, []PatchInstruction{{Op: PatchRename, Key: "OLD_NAME", NewKey: "NEW_NAME"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Errorf("expected Applied=true")
	}
	if _, ok := ef.Lookup("OLD_NAME"); ok {
		t.Errorf("old key should not exist after rename")
	}
	val, ok := ef.Lookup("NEW_NAME")
	if !ok || val != "hello" {
		t.Errorf("expected NEW_NAME=hello, got %q ok=%v", val, ok)
	}
}

func TestPatch_DryRun(t *testing.T) {
	ef := makePatchEnv(t, "KEEP", "original")
	_, err := Patch(ef, []PatchInstruction{{Op: PatchSet, Key: "KEEP", Value: "changed"}}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, _ := ef.Lookup("KEEP")
	if val != "original" {
		t.Errorf("dry run should not mutate: got %q", val)
	}
}

func TestPatch_NilEnvFile(t *testing.T) {
	_, err := Patch(nil, []PatchInstruction{{Op: PatchSet, Key: "X", Value: "1"}}, false)
	if err == nil {
		t.Error("expected error for nil EnvFile")
	}
}

func TestPatch_UnknownOp(t *testing.T) {
	ef := makePatchEnv(t)
	_, err := Patch(ef, []PatchInstruction{{Op: "explode", Key: "X"}}, false)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}
