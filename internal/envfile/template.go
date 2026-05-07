package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateResult holds the result of rendering a template against an EnvFile.
type TemplateResult struct {
	Rendered  string
	Missing   []string
	Unused    []string
}

var placeholderRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// RenderTemplate replaces ${KEY} placeholders in the template string with
// values from the provided EnvFile. It reports missing keys and unused entries.
func RenderTemplate(tmpl string, env *EnvFile) (TemplateResult, error) {
	if env == nil {
		return TemplateResult{}, fmt.Errorf("env must not be nil")
	}

	used := make(map[string]bool)
	missing := []string{}

	rendered := placeholderRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := placeholderRe.FindStringSubmatch(match)[1]
		for _, e := range env.Entries {
			if e.Key == key {
				used[key] = true
				return e.Value
			}
		}
		missing = append(missing, key)
		return match
	})

	unused := []string{}
	for _, e := range env.Entries {
		if !used[e.Key] {
			unused = append(unused, e.Key)
		}
	}

	return TemplateResult{
		Rendered: rendered,
		Missing:  missing,
		Unused:   unused,
	}, nil
}

// RenderTemplateFile reads a template from disk and renders it using RenderTemplate.
func RenderTemplateFile(path string, env *EnvFile) (TemplateResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TemplateResult{}, fmt.Errorf("read template %q: %w", path, err)
	}
	return RenderTemplate(strings.TrimRight(string(data), "\n"), env)
}
