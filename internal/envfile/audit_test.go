package envfile

import (
	"bytes"
	"strings"
	"testing"
)

func makeAuditDiffs() []DiffEntry {
	return []DiffEntry{
		{Key: "NEW_KEY", ValueA: "", ValueB: "hello", Status: StatusAdded},
		{Key: "OLD_KEY", ValueA: "bye", ValueB: "", Status: StatusRemoved},
		{Key: "CHANGED_KEY", ValueA: "v1", ValueB: "v2", Status: StatusChanged},
		{Key: "SAME_KEY", ValueA: "same", ValueB: "same", Status: StatusUnchanged},
	}
}

func TestNewAuditLog_Empty(t *testing.T) {
	log := NewAuditLog()
	if len(log.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(log.Entries))
	}
}

func TestAuditLog_Record(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditAdded, "API_KEY", "", "secret", "alice", ".env.prod")
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Key != "API_KEY" || e.Action != AuditAdded || e.NewValue != "secret" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestAuditLog_RecordDiff_SkipsUnchanged(t *testing.T) {
	log := NewAuditLog()
	log.RecordDiff(makeAuditDiffs(), "ci", ".env")
	// SAME_KEY (unchanged) should not be recorded
	if len(log.Entries) != 3 {
		t.Errorf("expected 3 entries (add/remove/change), got %d", len(log.Entries))
	}
}

func TestWriteAuditReport_NoEntries(t *testing.T) {
	log := NewAuditLog()
	var buf bytes.Buffer
	WriteAuditReport(&buf, log, false)
	if !strings.Contains(buf.String(), "No audit entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteAuditReport_MasksSecrets(t *testing.T) {
	log := NewAuditLog()
	log.Record(AuditChanged, "DB_PASS", "oldpass", "newpass", "dev", ".env")
	var buf bytes.Buffer
	WriteAuditReport(&buf, log, true)
	out := buf.String()
	if strings.Contains(out, "oldpass") || strings.Contains(out, "newpass") {
		t.Errorf("expected secrets to be masked, got: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected masked placeholder in output")
	}
}

func TestWriteAuditReport_ContainsActions(t *testing.T) {
	log := NewAuditLog()
	log.RecordDiff(makeAuditDiffs(), "bot", ".env.staging")
	var buf bytes.Buffer
	WriteAuditReport(&buf, log, false)
	out := buf.String()
	for _, action := range []string{"ADDED", "REMOVED", "CHANGED"} {
		if !strings.Contains(out, action) {
			t.Errorf("expected %s in report output", action)
		}
	}
}
