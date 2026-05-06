package envfile

import (
	"strings"
	"testing"
)

func TestWriteDiffReport_NoDiffs(t *testing.T) {
	var sb strings.Builder
	WriteDiffReport(&sb, []DiffEntry{}, ReportOptions{})
	if !strings.Contains(sb.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %s", sb.String())
	}
}

func TestWriteDiffReport_Added(t *testing.T) {
	var sb strings.Builder
	diffs := []DiffEntry{{Type: DiffAdded, Key: "NEW_KEY", NewValue: "val"}}
	WriteDiffReport(&sb, diffs, ReportOptions{})
	out := sb.String()
	if !strings.Contains(out, "+ NEW_KEY=val") {
		t.Errorf("expected added line, got: %s", out)
	}
}

func TestWriteDiffReport_Removed(t *testing.T) {
	var sb strings.Builder
	diffs := []DiffEntry{{Type: DiffRemoved, Key: "OLD_KEY", OldValue: "old"}}
	WriteDiffReport(&sb, diffs, ReportOptions{})
	out := sb.String()
	if !strings.Contains(out, "- OLD_KEY=old") {
		t.Errorf("expected removed line, got: %s", out)
	}
}

func TestWriteDiffReport_Changed(t *testing.T) {
	var sb strings.Builder
	diffs := []DiffEntry{{Type: DiffChanged, Key: "K", OldValue: "a", NewValue: "b"}}
	WriteDiffReport(&sb, diffs, ReportOptions{})
	out := sb.String()
	if !strings.Contains(out, "~ K:") {
		t.Errorf("expected changed line, got: %s", out)
	}
}

func TestWriteDiffReport_Masked(t *testing.T) {
	var sb strings.Builder
	diffs := []DiffEntry{{Type: DiffAdded, Key: "SECRET_TOKEN", NewValue: "supersecret"}}
	WriteDiffReport(&sb, diffs, ReportOptions{Masked: true})
	out := sb.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("secret value should be masked, got: %s", out)
	}
}

func TestWriteSyncReport(t *testing.T) {
	var sb strings.Builder
	res := &SyncResult{
		Added:   []string{"X"},
		Updated: []string{"Y"},
		Skipped: []string{"Z"},
	}
	WriteSyncReport(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "+1 added") || !strings.Contains(out, "~1 updated") {
		t.Errorf("unexpected sync report: %s", out)
	}
}
