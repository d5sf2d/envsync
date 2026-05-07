package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTemplateEnv(pairs ...string) *EnvFile {
	ef := &EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestRenderTemplate_AllResolved(t *testing.T) {
	env := makeTemplateEnv("HOST", "localhost", "PORT", "5432")
	res, err := RenderTemplate("connect to ${HOST}:${PORT}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "connect to localhost:5432" {
		t.Errorf("got %q", res.Rendered)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing, got %v", res.Missing)
	}
}

func TestRenderTemplate_MissingKey(t *testing.T) {
	env := makeTemplateEnv("HOST", "localhost")
	res, err := RenderTemplate("${HOST}:${PORT}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "PORT" {
		t.Errorf("expected PORT missing, got %v", res.Missing)
	}
	if res.Rendered != "localhost:${PORT}" {
		t.Errorf("placeholder not preserved, got %q", res.Rendered)
	}
}

func TestRenderTemplate_UnusedKeys(t *testing.T) {
	env := makeTemplateEnv("HOST", "localhost", "SECRET", "abc")
	res, err := RenderTemplate("${HOST}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unused) != 1 || res.Unused[0] != "SECRET" {
		t.Errorf("expected SECRET unused, got %v", res.Unused)
	}
}

func TestRenderTemplate_NilEnv(t *testing.T) {
	_, err := RenderTemplate("${HOST}", nil)
	if err == nil {
		t.Error("expected error for nil env")
	}
}

func TestRenderTemplateFile_Basic(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.txt")
	if err := os.WriteFile(p, []byte("db=${DB_URL}"), 0644); err != nil {
		t.Fatal(err)
	}
	env := makeTemplateEnv("DB_URL", "postgres://localhost/mydb")
	res, err := RenderTemplateFile(p, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "db=postgres://localhost/mydb" {
		t.Errorf("got %q", res.Rendered)
	}
}

func TestRenderTemplateFile_MissingFile(t *testing.T) {
	_, err := RenderTemplateFile("/nonexistent/tmpl.txt", makeTemplateEnv())
	if err == nil {
		t.Error("expected error for missing file")
	}
}
