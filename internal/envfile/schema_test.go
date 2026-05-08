package envfile

import (
	"bytes"
	"strings"
	"testing"
)

func makeSchemaEnv(pairs ...string) *EnvFile {
	ef := &EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestValidateSchema_NoViolations(t *testing.T) {
	ef := makeSchemaEnv("APP_ENV", "production", "PORT", "8080")
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "APP_ENV", Required: true},
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	violations := ValidateSchema(ef, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidateSchema_RequiredMissing(t *testing.T) {
	ef := makeSchemaEnv("PORT", "8080")
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "APP_ENV", Required: true},
		},
	}
	violations := ValidateSchema(ef, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "missing") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	ef := makeSchemaEnv("PORT", "not-a-port")
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	violations := ValidateSchema(ef, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "does not match pattern") {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestValidateSchema_EmptyValueNotAllowed(t *testing.T) {
	ef := makeSchemaEnv("SECRET", "")
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "SECRET", Required: true, AllowEmpty: false},
		},
	}
	violations := ValidateSchema(ef, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_EmptyValueAllowed(t *testing.T) {
	ef := makeSchemaEnv("OPTIONAL", "")
	schema := &Schema{
		Fields: []SchemaField{
			{Key: "OPTIONAL", Required: false, AllowEmpty: true},
		},
	}
	violations := ValidateSchema(ef, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestWriteSchemaReport_Pass(t *testing.T) {
	var buf bytes.Buffer
	WriteSchemaReport(&buf, nil, false)
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected pass message, got: %s", buf.String())
	}
}

func TestWriteSchemaReport_Violations(t *testing.T) {
	var buf bytes.Buffer
	violations := []SchemaViolation{
		{Key: "APP_ENV", Message: "required key is missing"},
	}
	WriteSchemaReport(&buf, violations, false)
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "required key is missing") {
		t.Errorf("expected message in output, got: %s", out)
	}
}
