package envfile

import (
	"fmt"
	"time"
)

// AuditAction represents the type of change recorded in an audit log.
type AuditAction string

const (
	AuditAdded    AuditAction = "ADDED"
	AuditRemoved  AuditAction = "REMOVED"
	AuditChanged  AuditAction = "CHANGED"
	AuditExported AuditAction = "EXPORTED"
	AuditSynced   AuditAction = "SYNCED"
)

// AuditEntry records a single auditable event.
type AuditEntry struct {
	Timestamp time.Time
	Action    AuditAction
	Key       string
	OldValue  string
	NewValue  string
	Actor     string
	Source    string
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends a new entry to the audit log.
func (a *AuditLog) Record(action AuditAction, key, oldVal, newVal, actor, source string) {
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		OldValue:  oldVal,
		NewValue:  newVal,
		Actor:     actor,
		Source:    source,
	})
}

// RecordDiff records audit entries for all diffs between two EnvFiles.
func (a *AuditLog) RecordDiff(diffs []DiffEntry, actor, source string) {
	for _, d := range diffs {
		switch d.Status {
		case StatusAdded:
			a.Record(AuditAdded, d.Key, "", d.ValueB, actor, source)
		case StatusRemoved:
			a.Record(AuditRemoved, d.Key, d.ValueA, "", actor, source)
		case StatusChanged:
			a.Record(AuditChanged, d.Key, d.ValueA, d.ValueB, actor, source)
		}
	}
}

// Summary returns a brief count summary of the audit log.
func (a *AuditLog) Summary() string {
	return fmt.Sprintf("%d audit entries recorded", len(a.Entries))
}
