package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Blank   bool
}

// EnvFile holds all entries parsed from an .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
}

// Parse reads and parses an .env file from the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		entry := parseLine(line)
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return &EnvFile{Path: path, Entries: entries}, nil
}

// ToMap converts the parsed entries into a key-value map (skips blanks/comments).
func (e *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(e.Entries))
	for _, entry := range e.Entries {
		if entry.Key != "" {
			m[entry.Key] = entry.Value
		}
	}
	return m
}

func parseLine(line string) Entry {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return Entry{Blank: true}
	}
	if strings.HasPrefix(trimmed, "#") {
		return Entry{Comment: trimmed}
	}
	parts := strings.SplitN(trimmed, "=", 2)
	if len(parts) != 2 {
		return Entry{Comment: line}
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	value = stripQuotes(value)
	return Entry{Key: key, Value: value}
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
