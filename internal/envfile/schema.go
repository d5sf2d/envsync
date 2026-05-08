package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField defines the expected shape of a single env key.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // optional regex the value must match
	AllowEmpty bool
}

// Schema holds a collection of field definitions.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation describes a single schema violation.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// ValidateSchema checks an EnvFile against a Schema and returns any violations.
func ValidateSchema(ef *EnvFile, schema *Schema) []SchemaViolation {
	var violations []SchemaViolation

	index := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		index[e.Key] = e.Value
	}

	for _, field := range schema.Fields {
		val, exists := index[field.Key]

		if field.Required && !exists {
			violations = append(violations, SchemaViolation{
				Key:     field.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if !field.AllowEmpty && strings.TrimSpace(val) == "" {
			violations = append(violations, SchemaViolation{
				Key:     field.Key,
				Message: "value must not be empty",
			})
			continue
		}

		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("invalid schema pattern: %v", err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				})
			}
		}
	}

	return violations
}
