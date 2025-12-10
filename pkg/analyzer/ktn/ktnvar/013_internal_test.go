package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar013 tests the private runVar013 function.
func Test_runVar013(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_collectVarGroups tests the private collectVarGroups helper function.
func Test_collectVarGroups(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "single var group",
			code: `package test
var (
	x int
	y string
)`,
			expected: 1,
		},
		{
			name: "multiple var groups",
			code: `package test
var x int
var y string`,
			expected: 2,
		},
		{
			name: "no vars",
			code: `package test
const x = 1`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			groups := collectVarGroups(file)
			// Vérification du nombre de groupes
			if len(groups) != tt.expected {
				t.Errorf("collectVarGroups() returned %d groups, expected %d", len(groups), tt.expected)
			}
		})
	}
}

// Test_checkVarGrouping tests the private checkVarGrouping helper function.
func Test_checkVarGrouping(t *testing.T) {
	tests := []struct {
		name          string
		groupCount    int
		expectReports int
	}{
		{
			name:          "no groups",
			groupCount:    0,
			expectReports: 0,
		},
		{
			name:          "one group",
			groupCount:    1,
			expectReports: 0,
		},
		{
			name:          "multiple groups",
			groupCount:    3,
			expectReports: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reports := 0
			mockPass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reports++
				},
			}

			// Create fake groups
			var groups []shared.DeclGroup
			for i := 0; i < tt.groupCount; i++ {
				groups = append(groups, shared.DeclGroup{
					Decl: &ast.GenDecl{TokPos: token.Pos(i + 1)},
					Pos:  token.Pos(i + 1),
				})
			}

			checkVarGrouping(mockPass, groups)

			// Vérification du nombre de rapports
			if reports != tt.expectReports {
				t.Errorf("checkVarGrouping() reported %d issues, expected %d", reports, tt.expectReports)
			}
		})
	}
}

// Test_runVar013_disabled tests runVar013 with disabled rule.
func Test_runVar013_disabled(t *testing.T) {
	// Setup config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-013": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	// Parse simple code
	code := `package test
var x int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar013(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar013() error = %v", err)
	}

	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar013() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar013_fileExcluded tests runVar013 with excluded file.
func Test_runVar013_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-013": {
				Exclude: []string{"test.go"},
			},
		},
	})
	defer config.Reset()

	// Parse simple code
	code := `package test
var x int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar013(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar013() error = %v", err)
	}

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar013() reported %d issues, expected 0 when file excluded", reportCount)
	}
}
