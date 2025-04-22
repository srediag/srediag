// Package diff provides functionality for generating and comparing diffs
package diff

import (
	"fmt"

	"github.com/pmezard/go-difflib/difflib"
)

// DiffGenerator encapsulates diff generation functionality
type DiffGenerator struct {
	FromFile string
	ToFile   string
	Context  int
}

// NewDiffGenerator creates a new instance of DiffGenerator
func NewDiffGenerator(fromFile, toFile string, context int) *DiffGenerator {
	return &DiffGenerator{
		FromFile: fromFile,
		ToFile:   toFile,
		Context:  context,
	}
}

// GenerateUnifiedDiff generates a unified diff between two string sequences
func (d *DiffGenerator) GenerateUnifiedDiff(a, b []string) (string, error) {
	if len(a) == 0 && len(b) == 0 {
		return fmt.Sprintf("--- %s\n+++ %s\n", d.FromFile, d.ToFile), nil
	}

	diff := difflib.UnifiedDiff{
		A:        a,
		B:        b,
		FromFile: d.FromFile,
		ToFile:   d.ToFile,
		Context:  d.Context,
	}
	return difflib.GetUnifiedDiffString(diff)
}

// GenerateContextDiff generates a context diff between two string sequences
func (d *DiffGenerator) GenerateContextDiff(a, b []string) (string, error) {
	if len(a) == 0 && len(b) == 0 {
		return fmt.Sprintf("*** %s\n--- %s\n", d.FromFile, d.ToFile), nil
	}

	diff := difflib.ContextDiff{
		A:        a,
		B:        b,
		FromFile: d.FromFile,
		ToFile:   d.ToFile,
		Context:  d.Context,
	}
	return difflib.GetContextDiffString(diff)
}
