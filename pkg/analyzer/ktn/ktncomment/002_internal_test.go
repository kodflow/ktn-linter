// Internal tests for 002.go - checkFileComment function.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
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

			got := checkFileComment(file)
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
