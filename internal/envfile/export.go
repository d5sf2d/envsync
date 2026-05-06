package envfile

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported env data.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatShell  ExportFormat = "shell"
)

// ExportOptions controls how the export is rendered.
type ExportOptions struct {
	Format  ExportFormat
	Masker  *Masker
	Sorted  bool
}

// Export writes the entries of an EnvFile to w in the requested format.
func Export(ef *EnvFile, w io.Writer, opts ExportOptions) error {
	entries := ef.Entries
	if opts.Sorted {
		entries = sortedEntries(entries)
	}

	switch opts.Format {
	case FormatJSON:
		return exportJSON(entries, w, opts.Masker)
	case FormatShell:
		return exportShell(entries, w, opts.Masker)
	default:
		return exportDotenv(entries, w, opts.Masker)
	}
}

func exportDotenv(entries []Entry, w io.Writer, m *Masker) error {
	for _, e := range entries {
		val := resolveValue(e.Value, e.Key, m)
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func exportShell(entries []Entry, w io.Writer, m *Masker) error {
	for _, e := range entries {
		val := resolveValue(e.Value, e.Key, m)
		quoted := strings.ReplaceAll(val, "'", "'\\''")
		if _, err := fmt.Fprintf(w, "export %s='%s'\n", e.Key, quoted); err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(entries []Entry, w io.Writer, m *Masker) error {
	kv := make(map[string]string, len(entries))
	for _, e := range entries {
		kv[e.Key] = resolveValue(e.Value, e.Key, m)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(kv)
}

func resolveValue(val, key string, m *Masker) string {
	if m != nil {
		return m.Mask(key, val)
	}
	return val
}

func sortedEntries(entries []Entry) []Entry {
	clone := make([]Entry, len(entries))
	copy(clone, entries)
	sort.Slice(clone, func(i, j int) bool {
		return clone[i].Key < clone[j].Key
	})
	return clone
}
