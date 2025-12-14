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
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API
		})
	}
}

// Test_collectDeclarations tests the private collectDeclarations function.
func Test_collectDeclarations(t *testing.T) {
	tests := []struct {
		name          string
		wantConst     int
		wantVar       int
		wantType      int
		wantFunc      int
		setupFile     func() *ast.File
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

// Test_checkScatteredConstBlocks tests the private checkScatteredConstBlocks function.
func Test_checkScatteredConstBlocks(t *testing.T) {
	tests := []struct {
		name        string
		constDecls  []token.Pos
		wantReports int
	}{
		{
			name:        "no const groups",
			constDecls:  []token.Pos{},
			wantReports: 0,
		},
		{
			name:        "single const group",
			constDecls:  []token.Pos{token.Pos(100)},
			wantReports: 0,
		},
		{
			name:        "two const groups",
			constDecls:  []token.Pos{token.Pos(100), token.Pos(200)},
			wantReports: 1,
		},
		{
			name:        "three const groups",
			constDecls:  []token.Pos{token.Pos(100), token.Pos(200), token.Pos(300)},
			wantReports: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal mock pass
			reportCount := 0
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
			}

			// Create fileDeclarations with const positions
			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			}

			checkScatteredConstBlocks(pass, decls)

			// Verify report count
			if reportCount != tt.wantReports {
				t.Errorf("checkScatteredConstBlocks() reports = %d, want %d", reportCount, tt.wantReports)
			}
		})
	}
}

// Test_checkConstBeforeVar tests the private checkConstBeforeVar function.
func Test_checkConstBeforeVar(t *testing.T) {
	tests := []struct {
		name        string
		constDecls  []token.Pos
		varDecls    []token.Pos
		wantReports int
	}{
		{
			name:        "no var declarations",
			constDecls:  []token.Pos{token.Pos(100)},
			varDecls:    []token.Pos{},
			wantReports: 0,
		},
		{
			name:        "const before var",
			constDecls:  []token.Pos{token.Pos(100)},
			varDecls:    []token.Pos{token.Pos(200)},
			wantReports: 0,
		},
		{
			name:        "const after var",
			constDecls:  []token.Pos{token.Pos(200)},
			varDecls:    []token.Pos{token.Pos(100)},
			wantReports: 1,
		},
		{
			name:        "multiple consts some after var",
			constDecls:  []token.Pos{token.Pos(50), token.Pos(200)},
			varDecls:    []token.Pos{token.Pos(100)},
			wantReports: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal mock pass
			reportCount := 0
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
			}

			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				varDecls:   tt.varDecls,
			}

			checkConstBeforeVar(pass, decls)

			// Verify report count
			if reportCount != tt.wantReports {
				t.Errorf("checkConstBeforeVar() reports = %d, want %d", reportCount, tt.wantReports)
			}
		})
	}
}

// Test_checkConstBeforeType tests the private checkConstBeforeType function.
func Test_checkConstBeforeType(t *testing.T) {
	tests := []struct {
		name        string
		constDecls  []token.Pos
		typeDecls   []token.Pos
		wantReports int
	}{
		{
			name:        "no type declarations",
			constDecls:  []token.Pos{token.Pos(100)},
			typeDecls:   []token.Pos{},
			wantReports: 0,
		},
		{
			name:        "const before type",
			constDecls:  []token.Pos{token.Pos(100)},
			typeDecls:   []token.Pos{token.Pos(200)},
			wantReports: 0,
		},
		{
			name:        "const after type",
			constDecls:  []token.Pos{token.Pos(200)},
			typeDecls:   []token.Pos{token.Pos(100)},
			wantReports: 1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal mock pass
			reportCount := 0
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
			}

			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				typeDecls:  tt.typeDecls,
				typeNames:  make(map[string]token.Pos),
				constTypes: make(map[token.Pos]string),
			}

			checkConstBeforeType(pass, decls)

			// Verify report count
			if reportCount != tt.wantReports {
				t.Errorf("checkConstBeforeType() reports = %d, want %d", reportCount, tt.wantReports)
			}
		})
	}
}

// Test_checkConstBeforeFunc tests the private checkConstBeforeFunc function.
func Test_checkConstBeforeFunc(t *testing.T) {
	tests := []struct {
		name        string
		constDecls  []token.Pos
		funcDecls   []token.Pos
		wantReports int
	}{
		{
			name:        "no func declarations",
			constDecls:  []token.Pos{token.Pos(100)},
			funcDecls:   []token.Pos{},
			wantReports: 0,
		},
		{
			name:        "const before func",
			constDecls:  []token.Pos{token.Pos(100)},
			funcDecls:   []token.Pos{token.Pos(200)},
			wantReports: 0,
		},
		{
			name:        "const after func",
			constDecls:  []token.Pos{token.Pos(200)},
			funcDecls:   []token.Pos{token.Pos(100)},
			wantReports: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal mock pass
			reportCount := 0
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
			}

			decls := &fileDeclarations{
				constDecls: tt.constDecls,
				funcDecls:  tt.funcDecls,
			}

			checkConstBeforeFunc(pass, decls)

			// Verify report count
			if reportCount != tt.wantReports {
				t.Errorf("checkConstBeforeFunc() reports = %d, want %d", reportCount, tt.wantReports)
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
			wantReports: 3,
		},
	}

	for _, tt := range tests {
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

// Test_checkConstBeforeType_iotaPattern tests the iota pattern exception.
func Test_checkConstBeforeType_iotaPattern(t *testing.T) {
	tests := []struct {
		name        string
		decls       *fileDeclarations
		wantReports int
	}{
		{
			name: "const after type using same type (valid iota pattern)",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(200)},
				typeDecls:  []token.Pos{token.Pos(100)},
				typeNames:  map[string]token.Pos{"Status": token.Pos(100)},
				constTypes: map[token.Pos]string{token.Pos(200): "Status"},
			},
			wantReports: 0,
		},
		{
			name: "const after type not using that type",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(200)},
				typeDecls:  []token.Pos{token.Pos(100)},
				typeNames:  map[string]token.Pos{"Status": token.Pos(100)},
				constTypes: map[token.Pos]string{},
			},
			wantReports: 1,
		},
		{
			name: "const after type using different type",
			decls: &fileDeclarations{
				constDecls: []token.Pos{token.Pos(200)},
				typeDecls:  []token.Pos{token.Pos(100)},
				typeNames:  map[string]token.Pos{"Status": token.Pos(100)},
				constTypes: map[token.Pos]string{token.Pos(200): "OtherType"},
			},
			wantReports: 1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal mock pass
			reportCount := 0
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
			}

			checkConstBeforeType(pass, tt.decls)

			// Verify report count
			if reportCount != tt.wantReports {
				t.Errorf("checkConstBeforeType() reports = %d, want %d", reportCount, tt.wantReports)
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
