package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env file's entries.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   map[string]string `json:"entries"`
}

// TakeSnapshot reads the given EnvFile and returns a Snapshot.
func TakeSnapshot(ef *EnvFile, source string) *Snapshot {
	entries := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		entries[e.Key] = e.Value
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Entries:   entries,
	}
}

// SaveSnapshot writes a Snapshot to the given file path as JSON.
func SaveSnapshot(snap *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from the given JSON file path.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &snap, nil
}

// DiffSnapshot compares two snapshots and returns a list of DiffEntry values
// describing what changed between them.
func DiffSnapshot(base, head *Snapshot) []DiffEntry {
	var results []DiffEntry

	for k, hv := range head.Entries {
		if bv, ok := base.Entries[k]; !ok {
			results = append(results, DiffEntry{Key: k, Status: StatusAdded, HeadValue: hv})
		} else if bv != hv {
			results = append(results, DiffEntry{Key: k, Status: StatusChanged, BaseValue: bv, HeadValue: hv})
		}
	}

	for k, bv := range base.Entries {
		if _, ok := head.Entries[k]; !ok {
			results = append(results, DiffEntry{Key: k, Status: StatusRemoved, BaseValue: bv})
		}
	}

	return results
}
