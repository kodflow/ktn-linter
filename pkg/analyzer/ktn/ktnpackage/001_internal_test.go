// Internal tests for analyzer 001 in ktnpackage package.
package ktnpackage

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_runPackage001 tests the private runPackage001 function
func Test_runPackage001(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		filename string
		wantErrs int
	}{
		{
			name: "missing package comment",
			code: `package test

func main() {}`,
			filename: "main.go",
			wantErrs: 1,
		},
		{
			name: "has package comment",
			code: `// Package test is a test package.
package test

func main() {}`,
			filename: "main.go",
			wantErrs: 0,
		},
		{
			name: "test file should be skipped",
			code: `package test

func TestSomething() {}`,
			filename: "main_internal_test.go",
			wantErrs: 0,
		},
		{
			name: "short comment is invalid",
			code: `// x
package test`,
			filename: "main.go",
			wantErrs: 1,
		},
		{
			name: "valid multiline comment",
			code: `/*
Package test does something.
*/
package test`,
			filename: "main.go",
			wantErrs: 0,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tt.filename, tt.code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runPackage001(pass)
			// Check for execution errors
			if err != nil {
				t.Fatalf("runPackage001 failed: %v", err)
			}

			// Check error count matches expectation
			if errorCount != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, errorCount)
			}
		})
	}
}

// Test_checkFileComment tests the checkFileComment function
func Test_checkFileComment(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "valid comment",
			code: `// Package test is valid.
package test`,
			want: true,
		},
		{
			name: "no comment",
			code: `package test`,
			want: false,
		},
		{
			name: "empty comment",
			code: `//
package test`,
			want: false,
		},
		{
			name: "too short comment",
			code: `// ab
package test`,
			want: false,
		},
		{
			name: "valid block comment",
			code: `/* Package test */
package test`,
			want: true,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			got := checkFileComment(file)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("checkFileComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
