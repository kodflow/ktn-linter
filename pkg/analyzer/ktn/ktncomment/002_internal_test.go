// Internal tests for 002.go - checkFileComment function.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

// Test_checkFileComment tests the checkFileComment function.
//
// Params:
//   - t: testing context
func Test_checkFileComment(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   bool
	}{
		{
			name: "file with valid package comment",
			source: `// Package example provides utilities.
package example`,
			want: true,
		},
		{
			name: "file without package comment",
			source: `package example`,
			want: false,
		},
		{
			name: "file with empty comment",
			source: `//
package example`,
			want: false,
		},
		{
			name: "file with short comment",
			source: `// ab
package example`,
			want: false,
		},
		{
			name: "file with minimum valid comment",
			source: `// abc
package example`,
			want: true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.source, parser.ParseComments)
			// Check parse error
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			got := checkFileComment(file, defaultMinPackageCommentLength)
			// Check result
			if got != tt.want {
				t.Errorf("checkFileComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_runComment002 tests the runComment002 function indirectly via checkFileComment.
// The actual analyzer is tested via analysistest in 002_external_test.go.
//
// Params:
//   - t: testing context
func Test_runComment002(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment002 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer002 is properly configured
			if Analyzer002 == nil {
				t.Error("Analyzer002 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer002.Name != "ktncomment002" {
				t.Errorf("Analyzer002.Name = %q, want %q", Analyzer002.Name, "ktncomment002")
			}
		})
	}
}

// helperParseFile parses source code and returns AST file.
//
// Params:
//   - t: testing context
//   - source: source code to parse
//
// Returns:
//   - *ast.File: parsed file
func helperParseFile(t *testing.T, source string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", source, parser.ParseComments)
	// Check parse error
	if err != nil {
		t.Fatalf("failed to parse source: %v", err)
	}
	// Return parsed file
	return file
}

// Test_runComment002_ruleDisabled tests behavior when rule is disabled.
//
// Params:
//   - t: testing context
func Test_runComment002_ruleDisabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-002": {Enabled: config.Bool(false)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment002(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment002 failed: %v", err)
			}

			// Should report no errors when rule disabled
			if errorCount != 0 {
				t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
			}

		})
	}
}

// Test_runComment002_fileExcluded tests behavior when file is excluded.
//
// Params:
//   - t: testing context
func Test_runComment002_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-002": {
						Enabled: config.Bool(true),
						Exclude: []string{"*.go"},
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment002(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment002 failed: %v", err)
			}

			// Should report no errors when file excluded
			if errorCount != 0 {
				t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
			}

		})
	}
}
