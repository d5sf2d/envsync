package envfile

import (
	"bytes"
	"testing"
	"time"
)

func makePinEnv(entries []Entry) *EnvFile {
	return &EnvFile{Entries: entries}
}

func TestPin_AllKeys(t *testing.T) {
	ef := makePinEnv([]Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
	})
	results, pins, err := Pin(ef, PinOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pins) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(pins))
	}
	for _, r := range results {
		if !r.Pinned {
			t.Errorf("expected key %q to be pinned", r.Key)
		}
	}
}

func TestPin_SelectedKeys(t *testing.T) {
	ef := makePinEnv([]Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
		{Key: "DEBUG", Value: "true"},
	})
	results, pins, err := Pin(ef, PinOptions{Keys: []string{"HOST", "PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pins) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(pins))
	}
	_ = results
}

func TestPin_MissingKeySkips(t *testing.T) {
	ef := makePinEnv([]Entry{{Key: "HOST", Value: "localhost"}})
	results, _, err := Pin(ef, PinOptions{Keys: []string{"MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected skipped result for missing key")
	}
}

func TestPin_FailOnMissing(t *testing.T) {
	ef := makePinEnv([]Entry{{Key: "HOST", Value: "localhost"}})
	_, _, err := Pin(ef, PinOptions{Keys: []string{"GHOST"}, FailOnMissing: true})
	if err == nil {
		t.Error("expected error for missing key with FailOnMissing")
	}
}

func TestPin_DryRun(t *testing.T) {
	ef := makePinEnv([]Entry{{Key: "API_KEY", Value: "secret"}})
	results, pins, err := Pin(ef, PinOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin recorded even in dry-run")
	}
	if results[0].Pinned {
		t.Error("dry-run should not mark result as Pinned=true")
	}
}

func TestPin_WithTTL(t *testing.T) {
	ef := makePinEnv([]Entry{{Key: "TOKEN", Value: "abc"}})
	ttl := 24 * time.Hour
	_, pins, err := Pin(ef, PinOptions{TTL: &ttl})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pins[0].ExpiresAt == nil {
		t.Error("expected ExpiresAt to be set")
	}
}

func TestCheckPins_Matching(t *testing.T) {
	ef := makePinEnv([]Entry{{Key: "HOST", Value: "localhost"}})
	now := time.Now().UTC()
	pins := []PinEntry{{Key: "HOST", Value: "localhost", PinnedAt: now}}
	results := CheckPins(ef, pins)
	if len(results) != 1 || !results[0].Pinned {
		t.Errorf("expected matching pin to be Pinned=true")
	}
}

func TestCheckPins_ValueDrift(t *testing.T) {
	ef := makePinEnv([]Entry{{Key: "HOST", Value: "changed"}})
	now := time.Now().UTC()
	pins := []PinEntry{{Key: "HOST", Value: "original", PinnedAt: now}}
	results := CheckPins(ef, pins)
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected value drift to produce Skipped result")
	}
}

func TestCheckPins_Expired(t *testing.T) {
	ef := makePinEnv([]Entry{{Key: "TOKEN", Value: "abc"}})
	past := time.Now().UTC().Add(-time.Hour)
	pins := []PinEntry{{Key: "TOKEN", Value: "abc", PinnedAt: past, ExpiresAt: &past}}
	results := CheckPins(ef, pins)
	if len(results) != 1 || !results[0].Expired {
		t.Errorf("expected expired pin to produce Expired result")
	}
}

func TestWritePinReport_Smoke(t *testing.T) {
	results := []PinResult{
		{Key: "HOST", Pinned: true},
		{Key: "GHOST", Skipped: true, Reason: "key not found"},
		{Key: "TOKEN", Expired: true, Reason: "pin expired"},
	}
	var buf bytes.Buffer
	WritePinReport(&buf, results, false)
	if buf.Len() == 0 {
		t.Error("expected non-empty report")
	}
}
