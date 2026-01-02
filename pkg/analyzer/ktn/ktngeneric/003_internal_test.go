package ktngeneric

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestCheckDeprecatedConstraintsImport(t *testing.T) {
	tests := []struct {
		name       string
		importSpec *ast.ImportSpec
	}{
		{
			name: "nil path",
			importSpec: &ast.ImportSpec{
				Path: nil,
			},
		},
		{
			name: "non-deprecated import",
			importSpec: &ast.ImportSpec{
				Path: &ast.BasicLit{Value: `"fmt"`},
			},
		},
		{
			name: "cmp import (ok)",
			importSpec: &ast.ImportSpec{
				Path: &ast.BasicLit{Value: `"cmp"`},
			},
		},
		{
			name: "deprecated constraints import with nil pass",
			importSpec: &ast.ImportSpec{
				Path: &ast.BasicLit{Value: `"golang.org/x/exp/constraints"`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic - pass is nil so no reporting
			checkDeprecatedConstraintsImport(nil, tt.importSpec)
		})
	}
}

// Test_runGeneric003_disabled tests behavior when rule is disabled.
func Test_runGeneric003_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec regle desactivee
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-003": {Enabled: config.Bool(false)},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			// Creer un pass minimal
			result, err := runGeneric003(&analysis.Pass{})
			// Verification de l'erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Verification du resultat nil
			if result != nil {
				t.Errorf("Expected nil result when rule disabled, got %v", result)
			}
		})
	}
}

// Test_runGeneric003_excludedFile tests behavior with excluded files.
func Test_runGeneric003_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec fichier exclu
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-003": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			code := `package test
import "golang.org/x/exp/constraints"
var _ = constraints.Ordered
`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Verification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Creer un inspector
			files := []*ast.File{file}
			inspectResult, _ := inspect.Analyzer.Run(&analysis.Pass{
				Fset:  fset,
				Files: files,
			})

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Expected no diagnostics for excluded file, got: %s", d.Message)
				},
			}

			// Executer l'analyseur
			result, err := runGeneric003(pass)
			// Verification de l'erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Verification du resultat nil
			if result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}
