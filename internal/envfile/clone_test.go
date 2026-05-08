package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeCloneEnv(t *testing.T, lines string) *EnvFile {
	t.Helper()
	f := writeTempEnv(t, lines)
	env, err := Parse(f)
	if err != nil {
		t.Fatalf("makeCloneEnv: parse: %v", err)
	}
	return env
}

func TestClone_CopiesAllKeys(t *testing.T) {
	src := makeCloneEnv(t, "APP=hello\nDB_HOST=localhost\n")
	dest := filepath.Join(t.TempDir(), "out.env")

	res, err := Clone(src, dest, CloneOptions{})
	if err != nil {
		t.Fatalf("Clone error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("dest file not created: %v", err)
	}
}

func TestClone_OnlyKeys(t *testing.T) {
	src := makeCloneEnv(t, "APP=hello\nDB_HOST=localhost\nSECRET_KEY=abc\n")
	dest := filepath.Join(t.TempDir(), "out.env")

	res, err := Clone(src, dest, CloneOptions{OnlyKeys: []string{"APP"}})
	if err != nil {
		t.Fatalf("Clone error: %v", err)
	}
	if len(res.Copied) != 1 || res.Copied[0] != "APP" {
		t.Errorf("expected only APP copied, got %v", res.Copied)
	}
	if len(res.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(res.Skipped))
	}
}

func TestClone_ExcludeKeys(t *testing.T) {
	src := makeCloneEnv(t, "APP=hello\nDB_HOST=localhost\n")
	dest := filepath.Join(t.TempDir(), "out.env")

	res, err := Clone(src, dest, CloneOptions{ExcludeKeys: []string{"DB_HOST"}})
	if err != nil {
		t.Fatalf("Clone error: %v", err)
	}
	if len(res.Copied) != 1 || res.Copied[0] != "APP" {
		t.Errorf("expected APP only, got %v", res.Copied)
	}
}

func TestClone_RedactSensitive(t *testing.T) {
	src := makeCloneEnv(t, "APP=hello\nSECRET_KEY=supersecret\n")
	dest := filepath.Join(t.TempDir(), "out.env")

	res, err := Clone(src, dest, CloneOptions{Redact: true})
	if err != nil {
		t.Fatalf("Clone error: %v", err)
	}
	if len(res.Redacted) != 1 || res.Redacted[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY redacted, got %v", res.Redacted)
	}

	out, err := Parse(dest)
	if err != nil {
		t.Fatalf("parse dest: %v", err)
	}
	for _, e := range out.Entries {
		if e.Key == "SECRET_KEY" && e.Value != "***" {
			t.Errorf("expected redacted value, got %q", e.Value)
		}
	}
}

func TestClone_DryRunDoesNotWrite(t *testing.T) {
	src := makeCloneEnv(t, "APP=hello\n")
	dest := filepath.Join(t.TempDir(), "out.env")

	_, err := Clone(src, dest, CloneOptions{DryRun: true})
	if err != nil {
		t.Fatalf("Clone dry-run error: %v", err)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Error("expected dest file NOT to be created in dry-run")
	}
}

func TestClone_NilSource(t *testing.T) {
	_, err := Clone(nil, "/tmp/x.env", CloneOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}
