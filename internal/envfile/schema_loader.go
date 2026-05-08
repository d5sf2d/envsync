package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadSchema reads a simple schema definition file.
// Each non-blank, non-comment line has the format:
//
//	KEY [required] [nonempty] [pattern=REGEX]
func LoadSchema(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("schema: open %q: %w", path, err)
	}
	defer f.Close()

	var schema Schema
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		field := SchemaField{Key: parts[0]}

		for _, token := range parts[1:] {
			switch {
			case token == "required":
				field.Required = true
			case token == "nonempty":
				field.AllowEmpty = false
			case token == "allowempty":
				field.AllowEmpty = true
			case strings.HasPrefix(token, "pattern="):
				field.Pattern = strings.TrimPrefix(token, "pattern=")
			default:
				return nil, fmt.Errorf("schema: line %d: unknown token %q", lineNum, token)
			}
		}

		schema.Fields = append(schema.Fields, field)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("schema: scan %q: %w", path, err)
	}

	return &schema, nil
}
