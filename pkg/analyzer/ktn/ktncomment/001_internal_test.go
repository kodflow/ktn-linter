// Internal tests for analyzer 001 in ktncomment package.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runComment001 tests the private runComment001 function
func Test_runComment001(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErrs int
	}{
		{
			name: "short inline comment is OK",
			code: `package test
func main() {
	x := 1 // short comment
}`,
			wantErrs: 0,
		},
		{
			name: "long inline comment",
			code: `package test
func main() {
	x := 1 // ` + strings.Repeat("a", 81) + `
}`,
			wantErrs: 1,
		},
		{
			name: "doc comment should be ignored",
			code: `package test
// ` + strings.Repeat("a", 100) + `
func main() {}`,
			wantErrs: 0,
		},
		{
			name: "comment at line start",
			code: `package test
// ` + strings.Repeat("a", 100) + `
var x = 1`,
			wantErrs: 0,
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

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer first
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment001(pass)
			// Check for execution errors
			if err != nil {
				t.Fatalf("runComment001 failed: %v", err)
			}

			// Check error count matches expectation
			if errorCount != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, errorCount)
			}
		})
	}
}

// Test_isDocComment tests the isDocComment function
func Test_isDocComment(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "function doc comment",
			code: `package test
// Documentation for function
func main() {}`,
			want: true,
		},
		{
			name: "inline comment",
			code: `package test
func main() {
	x := 1 // inline
}`,
			want: false,
		},
		{
			name: "comment at line start",
			code: `package test
// Comment at start
var x = 1`,
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

			pass := &analysis.Pass{
				Fset: fset,
			}

			// Get first comment
			if len(file.Comments) == 0 {
				t.Fatal("no comments found in code")
			}
			comment := file.Comments[0].List[0]

			got := isDocComment(pass, file, comment)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("isDocComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_getCommentLine tests the getCommentLine function
func Test_getCommentLine(t *testing.T) {
	code := `package test
// Comment on line 2
func main() {}
// Comment on line 4
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	tests := []struct {
		name        string
		commentIdx  int
		wantLine    int
	}{
		{"first comment", 0, 2},
		{"second comment", 1, 4},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check comment exists
			if tt.commentIdx >= len(file.Comments) {
				t.Fatalf("comment index %d out of range", tt.commentIdx)
			}

			comment := file.Comments[tt.commentIdx].List[0]
			got := getCommentLine(fset, comment)
			// Check result matches expectation
			if got != tt.wantLine {
				t.Errorf("getCommentLine() = %d, want %d", got, tt.wantLine)
			}
		})
	}
}

// Test_isCommentAtLineStart tests the isCommentAtLineStart function
func Test_isCommentAtLineStart(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "comment at column 1",
			code: `package test
// start
`,
			want: true,
		},
		{
			name: "inline comment",
			code: `package test
func main() { x := 1 // inline
}`,
			want: false,
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

			// Check comment exists
			if len(file.Comments) == 0 {
				t.Fatal("no comments found")
			}

			pass := &analysis.Pass{Fset: fset}
			comment := file.Comments[0].List[0]
			got := isCommentAtLineStart(pass, comment)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("isCommentAtLineStart() = %v, want %v", got, tt.want)
			}
		})
	}
}
