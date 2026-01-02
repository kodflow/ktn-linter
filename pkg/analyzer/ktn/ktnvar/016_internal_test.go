package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar016 tests the private runVar016 function.
func Test_runVar016(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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

// Test_collectVarGroups_nonGenDecl tests with non-GenDecl nodes.
func Test_collectVarGroups_nonGenDecl(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create file with function declaration (not GenDecl)
			file := &ast.File{
				Name: &ast.Ident{Name: "test"},
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Name: &ast.Ident{Name: "testFunc"},
						Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}},
						Body: &ast.BlockStmt{List: []ast.Stmt{}},
					},
				},
			}

			groups := collectVarGroups(file)

			// Should return empty slice
			if len(groups) != 0 {
				t.Errorf("collectVarGroups() returned %d groups, expected 0", len(groups))
			}
		})
	}
}

// Test_runVar016_disabled tests runVar016 with disabled rule.
func Test_runVar016_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-016": {Enabled: config.Bool(false)},
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

			_, err = runVar016(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar016() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar016() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar016_fileExcluded tests runVar016 with excluded file.
func Test_runVar016_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-016": {
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

			_, err = runVar016(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar016() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar016() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_runVar016_withMultipleVars tests runVar016 with multiple var declarations.
func Test_runVar016_withMultipleVars(t *testing.T) {
	config.Reset()

	code := `package test
var x int = 1
var y int = 2
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar016(pass)
	if err != nil {
		t.Errorf("runVar016() error = %v", err)
	}

	// Should report the second var declaration
	if reportCount != 1 {
		t.Errorf("runVar016() reported %d, expected 1", reportCount)
	}
}

// Test_checkVarGrouping_withVerbose tests checkVarGrouping with verbose mode.
func Test_checkVarGrouping_withVerbose(t *testing.T) {
	config.Set(&config.Config{
		Verbose: true,
	})
	defer config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Create 3 var groups
	groups := []shared.DeclGroup{
		{Decl: &ast.GenDecl{TokPos: token.Pos(1)}, Pos: token.Pos(1)},
		{Decl: &ast.GenDecl{TokPos: token.Pos(2)}, Pos: token.Pos(2)},
		{Decl: &ast.GenDecl{TokPos: token.Pos(3)}, Pos: token.Pos(3)},
	}

	checkVarGrouping(pass, groups)

	// Should report 2 issues (groups 2 and 3)
	if reportCount != 2 {
		t.Errorf("checkVarGrouping() with verbose reported %d, expected 2", reportCount)
	}
}

// Test_runVar016_withVerbose tests runVar016 with verbose mode.
func Test_runVar016_withVerbose(t *testing.T) {
	config.Set(&config.Config{
		Verbose: true,
	})
	defer config.Reset()

	code := `package test
var x int = 1
var y int = 2
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar016(pass)
	if err != nil {
		t.Errorf("runVar016() error = %v", err)
	}

	// Should report the second var declaration with verbose
	if reportCount != 1 {
		t.Errorf("runVar016() with verbose reported %d, expected 1", reportCount)
	}
}

// Test_runVar016_fileExcludedWithVars tests file exclusion with vars that would trigger.
func Test_runVar016_fileExcludedWithVars(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-016": {
				Exclude: []string{"excluded.go"},
			},
		},
	})
	defer config.Reset()

	// Code that would normally trigger the rule
	code := `package test
var x int = 1
var y int = 2
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar016(pass)
	if err != nil {
		t.Errorf("runVar016() error = %v", err)
	}

	// Should not report when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar016() reported %d, expected 0 when file excluded", reportCount)
	}
}

// Test_checkVarGrouping_missingMessage tests fallback when message is not found.
func Test_checkVarGrouping_missingMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation with missing message"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Get the current message to restore later
			originalMsg, hadMsg := messages.Get(ruleCodeVar016)

			// Remove the message to trigger fallback
			messages.Unregister(ruleCodeVar016)
			defer func() {
				// Restore the message
				if hadMsg {
					messages.Register(originalMsg)
				}
			}()

			reportCount := 0
			var reportedMsg string
			mockPass := &analysis.Pass{
				Report: func(d analysis.Diagnostic) {
					reportCount++
					reportedMsg = d.Message
				},
			}

			// Create 2 var groups to trigger reporting
			groups := []shared.DeclGroup{
				{Decl: &ast.GenDecl{TokPos: token.Pos(1)}, Pos: token.Pos(1)},
				{Decl: &ast.GenDecl{TokPos: token.Pos(2)}, Pos: token.Pos(2)},
			}

			checkVarGrouping(mockPass, groups)

			// Should report 1 issue with fallback message
			if reportCount != 1 {
				t.Errorf("checkVarGrouping() reported %d issues, expected 1", reportCount)
			}

			// Should use fallback message format
			expectedPrefix := "KTN-VAR-016: regrouper les variables"
			if !strings.HasPrefix(reportedMsg, expectedPrefix) {
				t.Errorf("checkVarGrouping() message = %q, expected prefix %q", reportedMsg, expectedPrefix)
			}
		})
	}
}
