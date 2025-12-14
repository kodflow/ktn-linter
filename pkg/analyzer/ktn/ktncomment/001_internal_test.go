// Internal tests for analyzer 001 in ktncomment package.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
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
	x := 1 // ` + strings.Repeat("a", 151) + `
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
		{
			name: "GenDecl doc comment",
			code: `package test
// This is a group comment
const (
	A = 1
	B = 2
)`,
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

// Test_containsURL tests the containsURL function
func Test_containsURL(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "text without URL",
			text: "This is a simple comment",
			want: false,
		},
		{
			name: "text with https URL",
			text: "See https://example.com/documentation for more info",
			want: true,
		},
		{
			name: "text with http URL",
			text: "Check http://localhost:8080/api",
			want: true,
		},
		{
			name: "text with file URL",
			text: "file:///etc/config is the config path",
			want: true,
		},
		{
			name: "text with long URL",
			text: "See https://github.com/org/repo/blob/main/path/to/file.go#L100-L200 for implementation",
			want: true,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsURL(tt.text)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("containsURL(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

// Test_checkMultiLineComment tests the checkMultiLineComment function.
//
// Params:
//   - t: testing context
func Test_checkMultiLineComment(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantErrs int
	}{
		{
			name:     "multi-line with short lines",
			text:     "/* short line 1\nshort line 2 */",
			wantErrs: 0,
		},
		{
			name:     "multi-line with long line",
			text:     "/* short line\n" + strings.Repeat("a", 155) + "\nshort line */",
			wantErrs: 1,
		},
		{
			name:     "multi-line with URL on long line",
			text:     "/* short line\nSee https://example.com/very/long/path/to/something */",
			wantErrs: 0,
		},
		{
			name:     "multi-line with empty lines",
			text:     "/* first\n\n\nlast */",
			wantErrs: 0,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			comment := &ast.Comment{
				Slash: token.NoPos,
				Text:  tt.text,
			}

			errorCount := 0
			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					errorCount++
				},
			}

			// Run the function
			checkMultiLineComment(pass, comment, tt.text, defaultMaxCommentLength)

			// Check error count matches expectation
			if errorCount != tt.wantErrs {
				t.Errorf("checkMultiLineComment() errors = %d, want %d", errorCount, tt.wantErrs)
			}
		})
	}
}

// Test_runComment001_ruleDisabled tests behavior when rule is disabled.
//
// Params:
//   - t: testing context
func Test_runComment001_ruleDisabled(t *testing.T) {
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
					"KTN-COMMENT-001": {Enabled: config.Bool(false)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			func main() {
			x := 1 // ` + strings.Repeat("a", 155) + `
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment001 failed: %v", err)
			}

			// Should report no errors when rule disabled
			if errorCount != 0 {
				t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
			}

		})
	}
}

// Test_runComment001_fileExcluded tests behavior when file is excluded.
//
// Params:
//   - t: testing context
func Test_runComment001_fileExcluded(t *testing.T) {
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
					"KTN-COMMENT-001": {
						Enabled: config.Bool(true),
						Exclude: []string{"*.go"},
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			func main() {
			x := 1 // ` + strings.Repeat("a", 155) + `
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment001 failed: %v", err)
			}

			// Should report no errors when file excluded
			if errorCount != 0 {
				t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
			}

		})
	}
}

// Test_runComment001_customThreshold tests behavior with custom threshold.
//
// Params:
//   - t: testing context
func Test_runComment001_customThreshold(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Configure custom threshold
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-001": {
						Enabled:   config.Bool(true),
						Threshold: config.Int(50),
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			func main() {
			x := 1 // ` + strings.Repeat("a", 51) + `
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment001 failed: %v", err)
			}

			// Should report error with custom threshold
			if errorCount != 1 {
				t.Errorf("expected 1 error with custom threshold, got %d", errorCount)
			}

		})
	}
}

// Test_isDocComment_edgeCases tests edge cases for isDocComment.
//
// Params:
//   - t: testing context
func Test_isDocComment_edgeCases(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "comment on same line as declaration",
			code: `package test
// Doc comment
const X = 1`,
			want: true,
		},
		{
			name: "comment not at line start",
			code: `package test
func main() {
		// Indented comment
		x := 1
		_ = x
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

// Test_runComment001_multiLineBlock tests multi-line block comments.
//
// Params:
//   - t: testing context
func Test_runComment001_multiLineBlock(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			func main() {
			x := 1 /* ` + strings.Repeat("a", 155) + ` */
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runComment001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment001 failed: %v", err)
			}

			// Should report error for long multi-line comment
			if errorCount != 1 {
				t.Errorf("expected 1 error for long multi-line comment, got %d", errorCount)
			}

		})
	}
}
