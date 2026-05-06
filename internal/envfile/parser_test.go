package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if m["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", m["APP_ENV"])
	}
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", m["DB_HOST"])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"\nTOKEN='bearer-token'\n`)
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if m["SECRET"] != "my secret value" {
		t.Errorf("expected 'my secret value', got %q", m["SECRET"])
	}
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# This is a comment\n\nAPP=test\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(ef.Entries))
	}
	m := ef.ToMap()
	if m["APP"] != "test" {
		t.Errorf("expected test, got %s", m["APP"])
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseLine_NoEquals(t *testing.T) {
	entry := parseLine("INVALID_LINE")
	if entry.Key != "" {
		t.Errorf("expected empty key, got %q", entry.Key)
	}
}
