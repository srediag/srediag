// Package diff provides functionality for generating and comparing diffs
package diff

import (
	"fmt"

	"github.com/pmezard/go-difflib/difflib"
)

// DiffGenerator encapsula a funcionalidade de geração de diffs
type DiffGenerator struct {
	FromFile string
	ToFile   string
	Context  int
}

// NewDiffGenerator cria uma nova instância de DiffGenerator
func NewDiffGenerator(fromFile, toFile string, context int) *DiffGenerator {
	return &DiffGenerator{
		FromFile: fromFile,
		ToFile:   toFile,
		Context:  context,
	}
}

// GenerateUnifiedDiff gera um diff unificado entre duas sequências de strings
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

// GenerateContextDiff gera um diff de contexto entre duas sequências de strings
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
