package envfile

import (
	"fmt"
	"strings"
)

// PatchOp represents a single patch operation kind.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
	PatchRename PatchOp = "rename"
)

// PatchInstruction describes one change to apply to an EnvFile.
type PatchInstruction struct {
	Op      PatchOp
	Key     string
	Value   string // used by set
	NewKey  string // used by rename
}

// PatchResult records the outcome of a single instruction.
type PatchResult struct {
	Instruction PatchInstruction
	Applied     bool
	Note        string
}

// Patch applies a slice of PatchInstructions to the given EnvFile.
// It returns a list of PatchResults describing what happened.
// When dryRun is true the EnvFile is not mutated.
func Patch(ef *EnvFile, instructions []PatchInstruction, dryRun bool) ([]PatchResult, error) {
	if ef == nil {
		return nil, fmt.Errorf("patch: nil EnvFile")
	}

	results := make([]PatchResult, 0, len(instructions))

	for _, inst := range instructions {
		result := PatchResult{Instruction: inst}

		switch inst.Op {
		case PatchSet:
			if strings.TrimSpace(inst.Key) == "" {
				result.Note = "skipped: empty key"
			} else {
				if !dryRun {
					ef.Set(inst.Key, inst.Value)
				}
				result.Applied = true
				result.Note = fmt.Sprintf("set %s", inst.Key)
			}

		case PatchDelete:
			if _, exists := ef.Lookup(inst.Key); !exists {
				result.Note = fmt.Sprintf("skipped: key %q not found", inst.Key)
			} else {
				if !dryRun {
					ef.Delete(inst.Key)
				}
				result.Applied = true
				result.Note = fmt.Sprintf("deleted %s", inst.Key)
			}

		case PatchRename:
			val, exists := ef.Lookup(inst.Key)
			if !exists {
				result.Note = fmt.Sprintf("skipped: key %q not found", inst.Key)
			} else if strings.TrimSpace(inst.NewKey) == "" {
				result.Note = "skipped: empty new key"
			} else {
				if !dryRun {
					ef.Delete(inst.Key)
					ef.Set(inst.NewKey, val)
				}
				result.Applied = true
				result.Note = fmt.Sprintf("renamed %s -> %s", inst.Key, inst.NewKey)
			}

		default:
			return nil, fmt.Errorf("patch: unknown op %q", inst.Op)
		}

		results = append(results, result)
	}

	return results, nil
}
