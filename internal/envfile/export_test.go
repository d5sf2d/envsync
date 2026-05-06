package envfile

import (
	"strings"
	"testing"
)

func makeExportEnv(pairs ...string) *EnvFile {
	ef := &EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestExport_Dotenv(t *testing.T) {
	ef := makeExportEnv("APP_ENV", "production", "PORT", "8080")
	var sb strings.Builder
	if err := Export(ef, &sb, ExportOptions{Format: FormatDotenv}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got: %s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got: %s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	ef := makeExportEnv("DB_PASS", "s3cr3t", "HOST", "localhost")
	var sb strings.Builder
	if err := Export(ef, &sb, ExportOptions{Format: FormatShell}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export DB_PASS='s3cr3t'") {
		t.Errorf("expected shell export for DB_PASS, got: %s", out)
	}
	if !strings.Contains(out, "export HOST='localhost'") {
		t.Errorf("expected shell export for HOST, got: %s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	ef := makeExportEnv("KEY", "value")
	var sb strings.Builder
	if err := Export(ef, &sb, ExportOptions{Format: FormatJSON}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, `"KEY"`) || !strings.Contains(out, `"value"`) {
		t.Errorf("expected JSON with KEY/value, got: %s", out)
	}
}

func TestExport_Sorted(t *testing.T) {
	ef := makeExportEnv("ZEBRA", "1", "ALPHA", "2", "MIDDLE", "3")
	var sb strings.Builder
	if err := Export(ef, &sb, ExportOptions{Format: FormatDotenv, Sorted: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected sorted output starting with ALPHA, got: %v", lines)
	}
}

func TestExport_MaskedValue(t *testing.T) {
	ef := makeExportEnv("SECRET_KEY", "topsecret")
	masker := NewMasker([]string{"SECRET_KEY"})
	var sb strings.Builder
	if err := Export(ef, &sb, ExportOptions{Format: FormatDotenv, Masker: masker}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if strings.Contains(out, "topsecret") {
		t.Errorf("expected secret to be masked, got: %s", out)
	}
}
