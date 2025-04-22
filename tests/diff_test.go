package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/srediag/srediag/internal/diff"
)

func TestDiffGeneratorCreation(t *testing.T) {
	diffGen := diff.NewDiffGenerator("teste1.txt", "teste2.txt", 3)
	assert.Equal(t, "teste1.txt", diffGen.FromFile)
	assert.Equal(t, "teste2.txt", diffGen.ToFile)
	assert.Equal(t, 3, diffGen.Context)
}

func TestDiffGeneratorUnifiedDiff(t *testing.T) {
	tests := []struct {
		name        string
		a           []string
		b           []string
		contains    []string
		notContains []string
	}{
		{
			name: "modificação simples",
			a:    []string{"linha1", "linha2", "linha3"},
			b:    []string{"linha1", "linha2 modificada", "linha3"},
			contains: []string{
				"-linha2",
				"+linha2 modificada",
			},
		},
		{
			name: "adição de linha",
			a:    []string{"linha1", "linha2"},
			b:    []string{"linha1", "linha2", "linha3"},
			contains: []string{
				"+linha3",
			},
			notContains: []string{
				"-linha3",
			},
		},
		{
			name: "remoção de linha",
			a:    []string{"linha1", "linha2", "linha3"},
			b:    []string{"linha1", "linha3"},
			contains: []string{
				"-linha2",
			},
		},
	}

	diffGen := diff.NewDiffGenerator("arquivo_original.txt", "arquivo_novo.txt", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := diffGen.GenerateUnifiedDiff(tt.a, tt.b)
			assert.NoError(t, err)

			for _, s := range tt.contains {
				assert.True(t, strings.Contains(diff, s), "diff deve conter '%s'", s)
			}

			for _, s := range tt.notContains {
				assert.False(t, strings.Contains(diff, s), "diff não deve conter '%s'", s)
			}
		})
	}
}

func TestDiffGeneratorContextDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		contains []string
	}{
		{
			name: "modificação simples",
			a:    []string{"linha1", "linha2", "linha3"},
			b:    []string{"linha1", "linha2 modificada", "linha3"},
			contains: []string{
				"! linha2 modificada",
			},
		},
		{
			name: "adição de linha",
			a:    []string{"linha1", "linha2"},
			b:    []string{"linha1", "linha2", "linha3"},
			contains: []string{
				"+ linha3",
			},
		},
		{
			name: "remoção de linha",
			a:    []string{"linha1", "linha2", "linha3"},
			b:    []string{"linha1", "linha3"},
			contains: []string{
				"- linha2",
			},
		},
	}

	diffGen := diff.NewDiffGenerator("arquivo_original.txt", "arquivo_novo.txt", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := diffGen.GenerateContextDiff(tt.a, tt.b)
			assert.NoError(t, err)

			for _, s := range tt.contains {
				assert.True(t, strings.Contains(diff, s), "diff deve conter '%s'", s)
			}
		})
	}
}

func TestDiffGeneratorWithEmptyInputs(t *testing.T) {
	diffGen := diff.NewDiffGenerator("vazio1.txt", "vazio2.txt", 3)

	// Teste com slices vazios
	unifiedDiff, err := diffGen.GenerateUnifiedDiff([]string{}, []string{})
	assert.NoError(t, err)
	assert.NotEmpty(t, unifiedDiff)

	contextDiff, err := diffGen.GenerateContextDiff([]string{}, []string{})
	assert.NoError(t, err)
	assert.NotEmpty(t, contextDiff)
}
