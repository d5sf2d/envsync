package envfile

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteRotateReport_AllRotated(t *testing.T) {
	results := []RotateResult{
		{Key: "DB_PASS", OldValue: "***", NewValue: "***", Rotated: true},
		{Key: "API_KEY", OldValue: "***", NewValue: "***", Rotated: true},
	}
	var buf bytes.Buffer
	WriteRotateReport(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "DB_PASS") {
		t.Error("expected DB_PASS in output")
	}
	if !strings.Contains(out, "2 rotated") {
		t.Errorf("expected '2 rotated', got: %s", out)
	}
}

func TestWriteRotateReport_WithSkipped(t *testing.T) {
	results := []RotateResult{
		{Key: "FOUND", Rotated: true},
		{Key: "MISSING", Skipped: true, Reason: "key not found"},
	}
	var buf bytes.Buffer
	WriteRotateReport(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "key not found") {
		t.Error("expected skip reason in output")
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected '1 skipped', got: %s", out)
	}
}

func TestWriteRotateReport_DryRun(t *testing.T) {
	var buf bytes.Buffer
	WriteRotateReport(&buf, []RotateResult{}, true)
	out := buf.String()
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run notice, got: %s", out)
	}
}

func TestWriteRotateReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteRotateReport(&buf, []RotateResult{}, false)
	out := buf.String()
	if !strings.Contains(out, "0 rotated") {
		t.Errorf("expected '0 rotated', got: %s", out)
	}
}
