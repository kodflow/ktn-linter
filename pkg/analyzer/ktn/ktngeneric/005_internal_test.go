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

func TestPredeclaredIdentifiers(t *testing.T) {
	tests := []struct {
		name       string
		identifier string
		expected   bool
	}{
		{
			name:       "string is predeclared",
			identifier: "string",
			expected:   true,
		},
		{
			name:       "int is predeclared",
			identifier: "int",
			expected:   true,
		},
		{
			name:       "error is predeclared",
			identifier: "error",
			expected:   true,
		},
		{
			name:       "bool is predeclared",
			identifier: "bool",
			expected:   true,
		},
		{
			name:       "len is predeclared",
			identifier: "len",
			expected:   true,
		},
		{
			name:       "make is predeclared",
			identifier: "make",
			expected:   true,
		},
		{
			name:       "nil is predeclared",
			identifier: "nil",
			expected:   true,
		},
		{
			name:       "T is not predeclared",
			identifier: "T",
			expected:   false,
		},
		{
			name:       "MyType is not predeclared",
			identifier: "MyType",
			expected:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Verification de l'identifiant predeclare
			if predeclaredIdentifiers[tt.identifier] != tt.expected {
				t.Errorf("predeclaredIdentifiers[%q] = %v, want %v", tt.identifier, predeclaredIdentifiers[tt.identifier], tt.expected)
			}
		})
	}
}

func TestCheckTypeParamList(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantCount int
	}{
		{
			name:      "no type parameters",
			code:      "package test\nfunc f() {}",
			wantCount: 0,
		},
		{
			name:      "valid type parameter T",
			code:      "package test\nfunc f[T any]() {}",
			wantCount: 0,
		},
		{
			name:      "predeclared type parameter string",
			code:      "package test\nfunc f[string any]() {}",
			wantCount: 1,
		},
		{
			name:      "multiple predeclared type parameters",
			code:      "package test\nfunc f[int, string any]() {}",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Verification d'erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Compteur de type parameters predeclares
			count := 0
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				// Verification du type de declaration
				if !ok {
					// Pas une fonction
					continue
				}
				// Verification des type parameters
				if funcDecl.Type.TypeParams != nil {
					for _, field := range funcDecl.Type.TypeParams.List {
						for _, name := range field.Names {
							// Verification si predeclare
							if predeclaredIdentifiers[name.Name] {
								count++
							}
						}
					}
				}
			}
			// Verification du nombre de type parameters predeclares
			if count != tt.wantCount {
				t.Errorf("got %d predeclared type params, want %d", count, tt.wantCount)
			}
		})
	}
}

func TestCheckFuncTypeParams(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
	}{
		{
			name: "nil type params",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					TypeParams: nil,
				},
			},
		},
		{
			name: "empty type params list",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "foo"},
				Type: &ast.FuncType{
					TypeParams: &ast.FieldList{List: nil},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			checkFuncTypeParams(nil, tt.funcDecl)
		})
	}
}

func TestCheckTypeSpecTypeParams(t *testing.T) {
	tests := []struct {
		name     string
		typeSpec *ast.TypeSpec
	}{
		{
			name: "nil type params",
			typeSpec: &ast.TypeSpec{
				Name:       &ast.Ident{Name: "Foo"},
				TypeParams: nil,
			},
		},
		{
			name: "empty type params list",
			typeSpec: &ast.TypeSpec{
				Name:       &ast.Ident{Name: "Foo"},
				TypeParams: &ast.FieldList{List: nil},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			checkTypeSpecTypeParams(nil, tt.typeSpec)
		})
	}
}

func TestCheckTypeParamListUnit(t *testing.T) {
	tests := []struct {
		name       string
		typeParams *ast.FieldList
	}{
		{
			name: "empty list",
			typeParams: &ast.FieldList{
				List: []*ast.Field{},
			},
		},
		{
			name: "non-predeclared names",
			typeParams: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "T"}, {Name: "U"}},
						Type:  &ast.Ident{Name: "any"},
					},
				},
			},
		},
		{
			name: "predeclared name with nil pass",
			typeParams: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "string"}},
						Type:  &ast.Ident{Name: "any"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic - pass is nil so no reporting
			checkTypeParamList(nil, tt.typeParams)
		})
	}
}

// Test_runGeneric005 tests the main runGeneric005 function.
func Test_runGeneric005(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "type param shadows predeclared",
			code: `package test
func foo[string any](s string) {}
`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Configure rule as enabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-005": {Enabled: config.Bool(true)},
				},
			})
			defer config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Verification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Create inspector
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
					// Expected to report
				},
			}

			// Execute analyzer
			result, err := runGeneric005(pass)
			// Verification erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Verification resultat nil
			if result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// Test_reportShadowing tests the reportShadowing function.
func Test_reportShadowing(t *testing.T) {
	tests := []struct {
		name  string
		ident *ast.Ident
	}{
		{
			name:  "predeclared identifier",
			ident: &ast.Ident{Name: "string"},
		},
		{
			name:  "another predeclared",
			ident: &ast.Ident{Name: "int"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic with nil pass
			reportShadowing(nil, tt.ident)
		})
	}
}

// TestReportShadowingEdgeCases tests edge cases for reportShadowing.
func TestReportShadowingEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		ident *ast.Ident
	}{
		{
			name:  "nil pass should not panic",
			ident: &ast.Ident{Name: "string"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic with nil pass
			reportShadowing(nil, tt.ident)
		})
	}
}

// Test_runGeneric005_disabled tests behavior when rule is disabled.
func Test_runGeneric005_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec regle desactivee
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-005": {Enabled: config.Bool(false)},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			// Creer un pass minimal
			result, err := runGeneric005(&analysis.Pass{})
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

// Test_runGeneric005_excludedFile tests behavior with excluded files.
func Test_runGeneric005_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Configuration avec fichier exclu
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-GENERIC-005": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config apres le test
			defer config.Reset()

			code := `package test
func foo[string any](s string) {}
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
			result, err := runGeneric005(pass)
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
