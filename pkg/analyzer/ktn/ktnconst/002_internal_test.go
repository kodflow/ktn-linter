package ktnconst

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

// Test_runConst002 tests the private runConst002 function.
func Test_runConst002(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API
		})
	}
}

// Test_collectDeclarations tests the private collectDeclarations function.
func Test_collectDeclarations(t *testing.T) {
	tests := []struct {
		name      string
		wantConst int
		wantVar   int
		wantType  int
		wantFunc  int
		setupFile func() *ast.File
	}{
		{
			name:      "empty file",
			wantConst: 0,
			wantVar:   0,
			wantType:  0,
			wantFunc:  0,
			setupFile: func() *ast.File {
				// Return empty file
				return &ast.File{Decls: []ast.Decl{}}
			},
		},
		{
			name:      "file with const only",
			wantConst: 1,
			wantVar:   0,
			wantType:  0,
			wantFunc:  0,
			setupFile: func() *ast.File {
				// Return file with one const
				return &ast.File{
					Decls: []ast.Decl{
						&ast.GenDecl{Tok: token.CONST},
					},
				}
			},
		},
		{
			name:      "file with all declaration types",
			wantConst: 1,
			wantVar:   1,
			wantType:  1,
			wantFunc:  1,
			setupFile: func() *ast.File {
				// Return file with all types
				return &ast.File{
					Decls: []ast.Decl{
						&ast.GenDecl{Tok: token.CONST},
						&ast.GenDecl{Tok: token.VAR},
						&ast.GenDecl{Tok: token.TYPE},
						&ast.FuncDecl{
							Name: &ast.Ident{Name: "testFunc"},
							Type: &ast.FuncType{},
						},
					},
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			file := tt.setupFile()
			decls := collectDeclarations(file)

			// Verify const count
			if len(decls.constDecls) != tt.wantConst {
				t.Errorf("constDecls = %d, want %d", len(decls.constDecls), tt.wantConst)
			}
			// Verify var count
			if len(decls.varDecls) != tt.wantVar {
				t.Errorf("varDecls = %d, want %d", len(decls.varDecls), tt.wantVar)
			}
			// Verify type count
			if len(decls.typeDecls) != tt.wantType {
				t.Errorf("typeDecls = %d, want %d", len(decls.typeDecls), tt.wantType)
			}
			// Verify func count
			if len(decls.funcDecls) != tt.wantFunc {
				t.Errorf("funcDecls = %d, want %d", len(decls.funcDecls), tt.wantFunc)
			}
		})
	}
}

// Test_collectScatteredViolations tests the private collectScatteredViolations function.
func Test_collectScatteredViolations(t *testing.T) {
	tests := []struct {
		name           string
		constDecls     []token.Pos
		varDecls       []token.Pos
		wantViolations int
	}{
		{
			name:           "no const groups",
			constDecls:     []token.Pos{},
			varDecls:       []token.Pos{},
			wantViolations: 0,
		},
		{
			name:           "single const group",
			constDecls:     []token.Pos{token.Pos(100)},
			varDecls:       []token.Pos{},
			wantViolations: 0,
		},
		{
			name:           "two const groups no var (allowed)",
			constDecls:     []token.Pos{token.Pos(100), token.Pos(200)},
			varDecls:       []token.Pos{},
			wantViolations: 0,
		},
		{
			name:           "two const groups both before var (allowed)",
			constDecls:     []token.Pos{token.Pos(100), token.Pos(200)},
			varDecls:       []token.Pos{token.Pos(300)},
			wantViolations: 0,
		},
		{
			name:           "two const groups one after var (scattered)",
			constDecls:     []token.Pos{token.Pos(100), token.Pos(400)},
			varDecls:       []token.Pos{token.Pos(300)},
			wantViolations: 1,
		},
		{
			name:           "three const groups two after var (scattered)",
			constDecls:     []token.Pos{token.Pos(100), token.Pos(400), token.Pos(500)},
			varDecls:       []token.Pos{token.Pos(300)},
			wantViolations: 2,
		},
		{
			name:           "only NoPos values (defensive guard)",
			constDecls:     []token.Pos{token.NoPos, token.NoPos},
			varDecls:       []token.Pos{token.Pos(100)},
			wantViolations: 0,
		},
		{
			name:           "mixed NoPos and valid positions",
			constDecls:     []token.Pos{token.NoPos, token.Pos(100), token.Pos(400)},
			varDecls:       []token.Pos{token.Pos(300)},
			wantViolations: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create fileDeclarations with const and var positions
			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				varDecls:   tt.varDecls,
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			}

			// Create violations map
			violations := make(map[token.Pos]bool)
			collectScatteredViolations(decls, violations)

			// Verify violations count
			if len(violations) != tt.wantViolations {
				t.Errorf("collectScatteredViolations() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

// Test_minPos tests the private minPos function.
func Test_minPos(t *testing.T) {
	tests := []struct {
		name      string
		positions []token.Pos
		expected  token.Pos
	}{
		{
			name:      "empty slice",
			positions: []token.Pos{},
			expected:  token.NoPos,
		},
		{
			name:      "single element",
			positions: []token.Pos{token.Pos(100)},
			expected:  token.Pos(100),
		},
		{
			name:      "two elements ascending",
			positions: []token.Pos{token.Pos(100), token.Pos(200)},
			expected:  token.Pos(100),
		},
		{
			name:      "two elements descending",
			positions: []token.Pos{token.Pos(200), token.Pos(100)},
			expected:  token.Pos(100),
		},
		{
			name:      "multiple elements unsorted",
			positions: []token.Pos{token.Pos(300), token.Pos(100), token.Pos(200)},
			expected:  token.Pos(100),
		},
		{
			name:      "all same value",
			positions: []token.Pos{token.Pos(50), token.Pos(50), token.Pos(50)},
			expected:  token.Pos(50),
		},
		{
			name:      "slice with NoPos values",
			positions: []token.Pos{token.NoPos, token.Pos(100), token.NoPos},
			expected:  token.Pos(100),
		},
		{
			name:      "only NoPos values",
			positions: []token.Pos{token.NoPos, token.NoPos},
			expected:  token.NoPos,
		},
		{
			name:      "NoPos at start",
			positions: []token.Pos{token.NoPos, token.Pos(200), token.Pos(100)},
			expected:  token.Pos(100),
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := minPos(tt.positions)
			// Verify result
			if result != tt.expected {
				t.Errorf("minPos() = %d, want %d", result, tt.expected)
			}
		})
	}
}

// Test_findFirstNonConstPos tests the private findFirstNonConstPos function.
func Test_findFirstNonConstPos(t *testing.T) {
	tests := []struct {
		name      string
		varDecls  []token.Pos
		typeDecls []token.Pos
		funcDecls []token.Pos
		expected  token.Pos
	}{
		{
			name:      "no declarations",
			varDecls:  []token.Pos{},
			typeDecls: []token.Pos{},
			funcDecls: []token.Pos{},
			expected:  token.NoPos,
		},
		{
			name:      "only var",
			varDecls:  []token.Pos{token.Pos(100)},
			typeDecls: []token.Pos{},
			funcDecls: []token.Pos{},
			expected:  token.Pos(100),
		},
		{
			name:      "only type",
			varDecls:  []token.Pos{},
			typeDecls: []token.Pos{token.Pos(200)},
			funcDecls: []token.Pos{},
			expected:  token.Pos(200),
		},
		{
			name:      "only func",
			varDecls:  []token.Pos{},
			typeDecls: []token.Pos{},
			funcDecls: []token.Pos{token.Pos(300)},
			expected:  token.Pos(300),
		},
		{
			name:      "var before type before func",
			varDecls:  []token.Pos{token.Pos(100)},
			typeDecls: []token.Pos{token.Pos(200)},
			funcDecls: []token.Pos{token.Pos(300)},
			expected:  token.Pos(100),
		},
		{
			name:      "func first",
			varDecls:  []token.Pos{token.Pos(300)},
			typeDecls: []token.Pos{token.Pos(200)},
			funcDecls: []token.Pos{token.Pos(100)},
			expected:  token.Pos(100),
		},
		{
			name:      "type first",
			varDecls:  []token.Pos{token.Pos(300)},
			typeDecls: []token.Pos{token.Pos(100)},
			funcDecls: []token.Pos{token.Pos(200)},
			expected:  token.Pos(100),
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create fileDeclarations
			decls := &fileDeclarations{
				varDecls:  tt.varDecls,
				typeDecls: tt.typeDecls,
				funcDecls: tt.funcDecls,
			}

			// Call function
			result := findFirstNonConstPos(decls)

			// Verify result
			if result != tt.expected {
				t.Errorf("findFirstNonConstPos() = %d, want %d", result, tt.expected)
			}
		})
	}
}

// Test_collectVarOrderViolations tests the private collectVarOrderViolations function.
func Test_collectVarOrderViolations(t *testing.T) {
	tests := []struct {
		name           string
		constDecls     []token.Pos
		varDecls       []token.Pos
		wantViolations int
	}{
		{
			name:           "no var declarations",
			constDecls:     []token.Pos{token.Pos(100)},
			varDecls:       []token.Pos{},
			wantViolations: 0,
		},
		{
			name:           "const before var",
			constDecls:     []token.Pos{token.Pos(100)},
			varDecls:       []token.Pos{token.Pos(200)},
			wantViolations: 0,
		},
		{
			name:           "const after var",
			constDecls:     []token.Pos{token.Pos(200)},
			varDecls:       []token.Pos{token.Pos(100)},
			wantViolations: 1,
		},
		{
			name:           "multiple consts some after var",
			constDecls:     []token.Pos{token.Pos(50), token.Pos(200)},
			varDecls:       []token.Pos{token.Pos(100)},
			wantViolations: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				varDecls:   tt.varDecls,
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			}

			// Create violations map
			violations := make(map[token.Pos]bool)
			collectVarOrderViolations(decls, violations)

			// Verify violations count
			if len(violations) != tt.wantViolations {
				t.Errorf("collectVarOrderViolations() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

// Test_collectTypeOrderViolations tests the private collectTypeOrderViolations function.
func Test_collectTypeOrderViolations(t *testing.T) {
	tests := []struct {
		name           string
		constDecls     []token.Pos
		typeDecls      []token.Pos
		wantViolations int
	}{
		{
			name:           "no type declarations",
			constDecls:     []token.Pos{token.Pos(100)},
			typeDecls:      []token.Pos{},
			wantViolations: 0,
		},
		{
			name:           "const before type",
			constDecls:     []token.Pos{token.Pos(100)},
			typeDecls:      []token.Pos{token.Pos(200)},
			wantViolations: 0,
		},
		{
			name:           "const after type",
			constDecls:     []token.Pos{token.Pos(200)},
			typeDecls:      []token.Pos{token.Pos(100)},
			wantViolations: 1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				typeDecls:  tt.typeDecls,
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			}

			// Create violations map
			violations := make(map[token.Pos]bool)
			collectTypeOrderViolations(decls, violations)

			// Verify violations count
			if len(violations) != tt.wantViolations {
				t.Errorf("collectTypeOrderViolations() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

// Test_collectFuncOrderViolations tests the private collectFuncOrderViolations function.
func Test_collectFuncOrderViolations(t *testing.T) {
	tests := []struct {
		name           string
		constDecls     []token.Pos
		funcDecls      []token.Pos
		wantViolations int
	}{
		{
			name:           "no func declarations",
			constDecls:     []token.Pos{token.Pos(100)},
			funcDecls:      []token.Pos{},
			wantViolations: 0,
		},
		{
			name:           "const before func",
			constDecls:     []token.Pos{token.Pos(100)},
			funcDecls:      []token.Pos{token.Pos(200)},
			wantViolations: 0,
		},
		{
			name:           "const after func",
			constDecls:     []token.Pos{token.Pos(200)},
			funcDecls:      []token.Pos{token.Pos(100)},
			wantViolations: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				funcDecls:  tt.funcDecls,
			}

			// Create violations map
			violations := make(map[token.Pos]bool)
			collectFuncOrderViolations(decls, violations)

			// Verify violations count
			if len(violations) != tt.wantViolations {
				t.Errorf("collectFuncOrderViolations() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

// Test_checkConstOrder tests the private checkConstOrder function.
func Test_checkConstOrder(t *testing.T) {
	tests := []struct {
		name        string
		decls       *fileDeclarations
		wantReports int
	}{
		{
			name: "no const declarations",
			decls: &fileDeclarations{
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			},
			wantReports: 0,
		},
		{
			name: "correct order const before all",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(100)},
				varDecls:   []token.Pos{token.Pos(200)},
				typeDecls:  []token.Pos{token.Pos(300)},
				funcDecls:  []token.Pos{token.Pos(400)},
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			},
			wantReports: 0,
		},
		{
			name: "const after everything",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(500)},
				varDecls:   []token.Pos{token.Pos(100)},
				typeDecls:  []token.Pos{token.Pos(200)},
				funcDecls:  []token.Pos{token.Pos(300)},
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			},
			wantReports: 1, // Only one error per position after deduplication
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal mock pass
			reportCount := 0
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
			}

			checkConstOrder(pass, tt.decls)

			// Verify report count
			if reportCount != tt.wantReports {
				t.Errorf("checkConstOrder() reports = %d, want %d", reportCount, tt.wantReports)
			}
		})
	}
}

// Test_fileDeclarations tests the fileDeclarations type structure.
func Test_fileDeclarations(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"type structure validation"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create declarations
			decls := &fileDeclarations{
				constDecls: []token.Pos{},
				varDecls:   []token.Pos{},
				typeDecls:  []token.Pos{},
				funcDecls:  []token.Pos{},
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			}

			// Verify fields are initialized
			if decls.constDecls == nil {
				t.Error("fileDeclarations.constDecls should not be nil")
			}
			// Verify fields are initialized
			if decls.varDecls == nil {
				t.Error("fileDeclarations.varDecls should not be nil")
			}
			// Verify fields are initialized
			if decls.typeDecls == nil {
				t.Error("fileDeclarations.typeDecls should not be nil")
			}
			// Verify fields are initialized
			if decls.funcDecls == nil {
				t.Error("fileDeclarations.funcDecls should not be nil")
			}
			// Verify map fields are initialized
			if decls.typeNames == nil {
				t.Error("fileDeclarations.typeNames should not be nil")
			}
			// Verify map fields are initialized
			if decls.constTypes == nil {
				t.Error("fileDeclarations.constTypes should not be nil")
			}
		})
	}
}

// Test_runConst002_disabled tests that the rule is skipped when disabled.
func Test_runConst002_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup: disable the rule
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-CONST-002": {Enabled: config.Bool(false)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			// Create minimal pass - should not report anything
			reportCount := 0
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Files: []*ast.File{
					{
						Decls: []ast.Decl{
							&ast.GenDecl{Tok: token.VAR},
							&ast.GenDecl{Tok: token.CONST},
						},
					},
				},
				Report: func(_ analysis.Diagnostic) {
					reportCount++
					t.Error("Unexpected error reported when rule is disabled")
				},
			}

			// Run the analyzer - should not report anything
			_, err := runConst002(pass)
			if err != nil {
				t.Errorf("runConst002() error = %v", err)
			}
			// Verify no reports
			if reportCount != 0 {
				t.Errorf("Expected 0 reports, got %d", reportCount)
			}

		})
	}
}

// Test_extractConstTypeName tests the private extractConstTypeName function.
func Test_extractConstTypeName(t *testing.T) {
	tests := []struct {
		name     string
		genDecl  *ast.GenDecl
		expected string
	}{
		{
			name: "const with explicit type",
			genDecl: &ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{{Name: "StatusOK"}},
						Type:  &ast.Ident{Name: "Status"},
						Values: []ast.Expr{
							&ast.Ident{Name: "iota"},
						},
					},
				},
			},
			expected: "Status",
		},
		{
			name: "const without explicit type",
			genDecl: &ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{{Name: "maxSize"}},
						Values: []ast.Expr{
							&ast.BasicLit{Kind: token.INT, Value: "100"},
						},
					},
				},
			},
			expected: "",
		},
		{
			name: "const with non-ident type",
			genDecl: &ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{{Name: "ptr"}},
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: "int"},
						},
					},
				},
			},
			expected: "",
		},
		{
			name:     "empty specs",
			genDecl:  &ast.GenDecl{Tok: token.CONST, Specs: []ast.Spec{}},
			expected: "",
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := extractConstTypeName(tt.genDecl)
			// Verify result
			if result != tt.expected {
				t.Errorf("extractConstTypeName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_collectTypeNames tests the private collectTypeNames function.
func Test_collectTypeNames(t *testing.T) {
	tests := []struct {
		name          string
		genDecl       *ast.GenDecl
		expectedNames []string
	}{
		{
			name: "single type declaration",
			genDecl: &ast.GenDecl{
				Tok:    token.TYPE,
				TokPos: token.Pos(100),
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "Status"},
						Type: &ast.Ident{Name: "int"},
					},
				},
			},
			expectedNames: []string{"Status"},
		},
		{
			name: "multiple type declarations",
			genDecl: &ast.GenDecl{
				Tok:    token.TYPE,
				TokPos: token.Pos(100),
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "Status"},
						Type: &ast.Ident{Name: "int"},
					},
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "Priority"},
						Type: &ast.Ident{Name: "int"},
					},
				},
			},
			expectedNames: []string{"Status", "Priority"},
		},
		{
			name: "empty specs",
			genDecl: &ast.GenDecl{
				Tok:   token.TYPE,
				Specs: []ast.Spec{},
			},
			expectedNames: []string{},
		},
		{
			name: "non-type spec ignored",
			genDecl: &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{{Name: "x"}},
					},
				},
			},
			expectedNames: []string{},
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			typeNames := make(map[string]token.Pos)
			collectTypeNames(tt.genDecl, typeNames)

			// Verify count
			if len(typeNames) != len(tt.expectedNames) {
				t.Errorf("collectTypeNames() collected %d names, want %d", len(typeNames), len(tt.expectedNames))
			}

			// Verify each expected name exists
			for _, name := range tt.expectedNames {
				// Check if name exists in map
				if _, exists := typeNames[name]; !exists {
					t.Errorf("collectTypeNames() missing expected name %q", name)
				}
			}
		})
	}
}

// Test_collectTypeOrderViolations_iotaPattern tests the iota pattern exception.
func Test_collectTypeOrderViolations_iotaPattern(t *testing.T) {
	tests := []struct {
		name           string
		decls          *fileDeclarations
		wantViolations int
	}{
		{
			name: "const after type using same type (valid iota pattern)",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(200)},
				typeDecls:  []token.Pos{token.Pos(100)},
				typeNames:  map[string]token.Pos{"Status": token.Pos(100)},
				constTypes: map[token.Pos]string{token.Pos(200): "Status"},
			},
			wantViolations: 0,
		},
		{
			name: "const after type not using that type",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(200)},
				typeDecls:  []token.Pos{token.Pos(100)},
				typeNames:  map[string]token.Pos{"Status": token.Pos(100)},
				constTypes: map[token.Pos]string{},
			},
			wantViolations: 1,
		},
		{
			name: "const after type using different type",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(200)},
				typeDecls:  []token.Pos{token.Pos(100)},
				typeNames:  map[string]token.Pos{"Status": token.Pos(100)},
				constTypes: map[token.Pos]string{token.Pos(200): "OtherType"},
			},
			wantViolations: 1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create violations map
			violations := make(map[token.Pos]bool)
			collectTypeOrderViolations(tt.decls, violations)

			// Verify violations count
			if len(violations) != tt.wantViolations {
				t.Errorf("collectTypeOrderViolations() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

// Test_runConst002_excludedFile tests that excluded files are skipped.
func Test_runConst002_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup: exclude test.go
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-CONST-002": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			// Create fset and file with bad order (var before const)
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			// Create pass with file that would trigger error if not excluded
			reportCount := 0
			pass := &analysis.Pass{
				Fset: fset,
				Files: []*ast.File{
					{
						Package: file.Pos(1),
						Decls: []ast.Decl{
							&ast.GenDecl{Tok: token.VAR, TokPos: file.Pos(10)},
							&ast.GenDecl{Tok: token.CONST, TokPos: file.Pos(50)},
						},
					},
				},
				Report: func(_ analysis.Diagnostic) {
					reportCount++
					t.Error("Unexpected error reported for excluded file")
				},
			}

			// Run the analyzer
			_, err := runConst002(pass)
			if err != nil {
				t.Errorf("runConst002() error = %v", err)
			}
			// Verify no reports for excluded file
			if reportCount != 0 {
				t.Errorf("Expected 0 reports for excluded file, got %d", reportCount)
			}

		})
	}
}
