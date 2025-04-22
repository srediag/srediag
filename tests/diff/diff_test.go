package diff_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/srediag/srediag/internal/diff"
)

func TestDiffGenerator_Creation(t *testing.T) {
	generator := diff.NewDiffGenerator("arquivo1.txt", "arquivo2.txt", 3)
	assert.Equal(t, "arquivo1.txt", generator.FromFile)
	assert.Equal(t, "arquivo2.txt", generator.ToFile)
	assert.Equal(t, 3, generator.Context)
}

func TestDiffGenerator_UnifiedDiff(t *testing.T) {
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
		{
			name: "arquivos vazios",
			a:    []string{},
			b:    []string{},
			contains: []string{
				"---",
				"+++",
			},
		},
	}

	generator := diff.NewDiffGenerator("original.txt", "modificado.txt", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.GenerateUnifiedDiff(tt.a, tt.b)
			assert.NoError(t, err)

			for _, s := range tt.contains {
				assert.True(t, strings.Contains(result, s), "diff deve conter '%s'", s)
			}

			for _, s := range tt.notContains {
				assert.False(t, strings.Contains(result, s), "diff não deve conter '%s'", s)
			}
		})
	}
}

func TestDiffGenerator_ContextDiff(t *testing.T) {
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
		{
			name: "arquivos vazios",
			a:    []string{},
			b:    []string{},
			contains: []string{
				"***",
				"---",
			},
		},
	}

	generator := diff.NewDiffGenerator("original.txt", "modificado.txt", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.GenerateContextDiff(tt.a, tt.b)
			assert.NoError(t, err)

			for _, s := range tt.contains {
				assert.True(t, strings.Contains(result, s), "diff deve conter '%s'", s)
			}
		})
	}
}
