package envfile

import (
	"fmt"
	"sort"
	"time"
)

// PinEntry represents a pinned key-value pair with metadata.
type PinEntry struct {
	Key       string
	Value     string
	PinnedAt  time.Time
	ExpiresAt *time.Time
	Comment   string
}

// PinResult holds the outcome of a pin operation for a single key.
type PinResult struct {
	Key     string
	Pinned  bool
	Expired bool
	Skipped bool
	Reason  string
}

// PinOptions controls the behaviour of Pin.
type PinOptions struct {
	// Keys to pin; if empty, all keys are pinned.
	Keys []string
	// TTL optionally sets an expiry duration from now.
	TTL *time.Duration
	// DryRun reports what would happen without modifying the EnvFile.
	DryRun bool
	// FailOnMissing returns an error when a requested key is absent.
	FailOnMissing bool
}

// Pin locks the values of specified keys in ef, recording them as PinEntries.
// It returns the list of results and the canonical pin set.
func Pin(ef *EnvFile, opts PinOptions) ([]PinResult, []PinEntry, error) {
	now := time.Now().UTC()

	targets := buildTargetSet(opts.Keys, ef)

	var results []PinResult
	var pins []PinEntry

	for _, key := range sortedKeys(targets) {
		entry, found := findEntry(ef, key)
		if !found {
			if opts.FailOnMissing {
				return nil, nil, fmt.Errorf("pin: key %q not found", key)
			}
			results = append(results, PinResult{Key: key, Skipped: true, Reason: "key not found"})
			continue
		}

		pin := PinEntry{
			Key:      key,
			Value:    entry.Value,
			PinnedAt: now,
			Comment:  entry.Comment,
		}
		if opts.TTL != nil {
			exp := now.Add(*opts.TTL)
			pin.ExpiresAt = &exp
		}

		pins = append(pins, pin)
		results = append(results, PinResult{Key: key, Pinned: !opts.DryRun})
	}

	return results, pins, nil
}

// CheckPins validates that pinned values in ef still match the recorded pins.
func CheckPins(ef *EnvFile, pins []PinEntry) []PinResult {
	now := time.Now().UTC()
	var results []PinResult
	for _, pin := range pins {
		if pin.ExpiresAt != nil && now.After(*pin.ExpiresAt) {
			results = append(results, PinResult{Key: pin.Key, Expired: true, Reason: "pin expired"})
			continue
		}
		entry, found := findEntry(ef, pin.Key)
		if !found {
			results = append(results, PinResult{Key: pin.Key, Skipped: true, Reason: "key not found"})
			continue
		}
		if entry.Value != pin.Value {
			results = append(results, PinResult{Key: pin.Key, Skipped: true, Reason: "value drift detected"})
			continue
		}
		results = append(results, PinResult{Key: pin.Key, Pinned: true})
	}
	return results
}

func buildTargetSet(keys []string, ef *EnvFile) map[string]struct{} {
	targets := make(map[string]struct{})
	if len(keys) == 0 {
		for _, e := range ef.Entries {
			targets[e.Key] = struct{}{}
		}
	} else {
		for _, k := range keys {
			targets[k] = struct{}{}
		}
	}
	return targets
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func findEntry(ef *EnvFile, key string) (Entry, bool) {
	for _, e := range ef.Entries {
		if e.Key == key {
			return e, true
		}
	}
	return Entry{}, false
}
