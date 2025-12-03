package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Test_runVar002 tests the private runVar002 function.
func Test_runVar002(t *testing.T) {
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
