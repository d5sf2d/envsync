package envfile

import (
	"testing"
)

func makeValidateEnv(entries []Entry) *EnvFile {
	return &EnvFile{Entries: entries}
}

func TestValidate_ValidFile(t *testing.T) {
	ef := makeValidateEnv([]Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "PORT", Value: "8080"},
	})
	result := Validate(ef, nil)
	if result.HasErrors() {
		t.Fatalf("expected no errors, got: %v", result.Errors)
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	ef := makeValidateEnv([]Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_HOST", Value: "remotehost"},
	})
	result := Validate(ef, nil)
	if !result.HasErrors() {
		t.Fatal("expected duplicate key error")
	}
	if result.Errors[0].Key != "DB_HOST" {
		t.Errorf("expected error for DB_HOST, got %s", result.Errors[0].Key)
	}
}

func TestValidate_InvalidKeyName(t *testing.T) {
	ef := makeValidateEnv([]Entry{
		{Key: "123INVALID", Value: "value"},
	})
	result := Validate(ef, nil)
	if !result.HasErrors() {
		t.Fatal("expected invalid key name error")
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	ef := makeValidateEnv([]Entry{
		{Key: "APP_NAME", Value: "envsync"},
	})
	result := Validate(ef, []string{"SECRET_KEY"})
	if !result.HasErrors() {
		t.Fatal("expected missing required key error")
	}
	if result.Errors[0].Key != "SECRET_KEY" {
		t.Errorf("expected error for SECRET_KEY, got %s", result.Errors[0].Key)
	}
}

func TestValidate_RequiredKeyEmpty(t *testing.T) {
	ef := makeValidateEnv([]Entry{
		{Key: "SECRET_KEY", Value: "   "},
	})
	result := Validate(ef, []string{"SECRET_KEY"})
	if !result.HasErrors() {
		t.Fatal("expected empty required key error")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	ef := makeValidateEnv([]Entry{
		{Key: "GOOD_KEY", Value: "val"},
		{Key: "bad-key", Value: "val"},
		{Key: "GOOD_KEY", Value: "dup"},
	})
	result := Validate(ef, []string{"MISSING"})
	if len(result.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}
