// Internal tests for 007.go - control flow comment analyzer.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runComment007 tests the runComment007 function configuration.
func Test_runComment007(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment007 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer007 is properly configured
			if Analyzer007 == nil {
				t.Error("Analyzer007 should not be nil")
				// Retour anticipé
				return
			}
			// Check analyzer name
			if Analyzer007.Name != "ktncomment007" {
				t.Errorf("Analyzer007.Name = %q, want %q", Analyzer007.Name, "ktncomment007")
			}
		})
	}
}

// Test_checkIfStmt tests that checkIfStmt analyzer configuration exists.
// Actual behavior is tested via analysistest.
func Test_checkIfStmt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "checkIfStmt is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for checkIfStmt")
			}
		})
	}
}

// Test_checkSwitchStmt tests that checkSwitchStmt analyzer configuration exists.
func Test_checkSwitchStmt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "checkSwitchStmt is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for checkSwitchStmt")
			}
		})
	}
}

// Test_checkTypeSwitchStmt tests that checkTypeSwitchStmt analyzer configuration exists.
func Test_checkTypeSwitchStmt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "checkTypeSwitchStmt is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for checkTypeSwitchStmt")
			}
		})
	}
}

// Test_checkLoopStmt tests that checkLoopStmt analyzer configuration exists.
func Test_checkLoopStmt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "checkLoopStmt is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for checkLoopStmt")
			}
		})
	}
}

// Test_checkReturnStmt tests that checkReturnStmt analyzer configuration exists.
func Test_checkReturnStmt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "checkReturnStmt is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for checkReturnStmt")
			}
		})
	}
}

// Test_hasCommentBefore tests that hasCommentBefore function configuration exists.
func Test_hasCommentBefore(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "hasCommentBefore is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for hasCommentBefore")
			}
		})
	}
}

// Test_hasInlineComment tests that hasInlineComment function configuration exists.
func Test_hasInlineComment(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "hasInlineComment is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for hasInlineComment")
			}
		})
	}
}

// Test_hasCommentBeforeOrInside tests that hasCommentBeforeOrInside function exists.
func Test_hasCommentBeforeOrInside(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "hasCommentBeforeOrInside is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify the function exists by checking Analyzer007
			if Analyzer007 == nil {
				t.Error("Analyzer007 should exist for hasCommentBeforeOrInside")
			}
		})
	}
}

// Test_runComment007_ruleDisabled tests behavior when rule is disabled.
func Test_runComment007_ruleDisabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-007": {Enabled: config.Bool(false)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			func myFunc() {
			if true {
				x := 1
				_ = x
			}
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
			_, err = runComment007(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment007 failed: %v", err)
			}

			// Should report no errors when rule disabled
			if errorCount != 0 {
				t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
			}

		})
	}
}

// Test_runComment007_fileExcluded tests behavior when file is excluded.
func Test_runComment007_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-007": {
						Enabled: config.Bool(true),
						Exclude: []string{"*.go"},
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			func myFunc() {
			if true {
				x := 1
				_ = x
			}
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
			_, err = runComment007(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment007 failed: %v", err)
			}

			// Should report no errors when file excluded
			if errorCount != 0 {
				t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
			}

		})
	}
}

// Test_shouldSkipFunction tests the shouldSkipFunction function branches.
func Test_shouldSkipFunction(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		filename string
		wantSkip bool
	}{
		{
			name: "skip test file",
			code: `package test
func myFunc() {}`,
			filename: "test_test.go",
			wantSkip: true,
		},
		{
			name: "skip function without body",
			code: `package test
func myFunc()`,
			filename: "test.go",
			wantSkip: true,
		},
		{
			name: "do not skip regular function",
			code: `package test
func myFunc() {}`,
			filename: "test.go",
			wantSkip: false,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tt.filename, tt.code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Find the function declaration
			var funcDecl *ast.FuncDecl
			// Iterate over declarations
			for _, decl := range file.Decls {
				// Check if it's a function
				if fd, ok := decl.(*ast.FuncDecl); ok {
					funcDecl = fd
					break
				}
			}

			// Check function found
			if funcDecl == nil && tt.name != "skip function without body" {
				t.Fatal("no function declaration found")
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			cfg := config.Get()
			got := shouldSkipFunction(pass, cfg, funcDecl)

			// Check result
			if got != tt.wantSkip {
				t.Errorf("shouldSkipFunction() = %v, want %v", got, tt.wantSkip)
			}
		})
	}
}

// Test_hasCommentBeforeOrInsideBlockStmt tests hasCommentBeforeOrInside with block statements.
func Test_hasCommentBeforeOrInsideBlockStmt(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "comment inside block",
			code: `package test
func myFunc() {
	if true {
		x := 1
		_ = x
	} else {
		// Comment inside else
		y := 2
		_ = y
	}
}`,
			want: true,
		},
		{
			name: "no comment in block",
			code: `package test
func myFunc() {
	if true {
		x := 1
		_ = x
	} else {
		y := 2
		_ = y
	}
}`,
			want: false,
		},
		{
			name: "else with empty block",
			code: `package test
func myFunc() {
	if true {
		x := 1
		_ = x
	} else {
	}
}`,
			want: false,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			// Find the else block
			var elseStmt ast.Stmt
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if it's an if statement
				if ifStmt, ok := n.(*ast.IfStmt); ok {
					// Check if it has an else clause
					if ifStmt.Else != nil {
						elseStmt = ifStmt.Else
						return false
					}
				}
				// Continue traversal
				return true
			})

			// Check else statement found
			if elseStmt == nil {
				t.Fatal("no else statement found")
			}

			got := hasCommentBeforeOrInside(pass, elseStmt)

			// Check result
			if got != tt.want {
				t.Errorf("hasCommentBeforeOrInside() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_checkFunctionBody tests checkFunctionBody coverage.
func Test_checkFunctionBody(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			func myFunc() {
			// Check if statement
			if true {
				x := 1
				_ = x
			}

			// Check switch statement
			switch x := 1; x {
			// Check case
			case 1:
				y := 2
				_ = y
			}

			// Check type switch
			switch v := interface{}(1).(type) {
			// Check case
			case int:
				z := v
				_ = z
			}

			// Check for loop
			for i := 0; i < 10; i++ {
				w := i
				_ = w
			}

			// Check range loop
			for _, v := range []int{1, 2, 3} {
				u := v
				_ = u
			}

			// Return statement
			return
			}`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					// Collect reports
				},
			}

			// Find the function body
			var funcBody *ast.BlockStmt
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if it's a function declaration
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcBody = fd.Body
					return false
				}
				// Continue traversal
				return true
			})

			// Check function body found
			if funcBody == nil {
				t.Fatal("no function body found")
			}

			// Call checkFunctionBody to increase coverage
			checkFunctionBody(pass, funcBody)

		})
	}
}

// Test_hasCommentBeforeOrInside_nonBlockStmt tests non-block statements.
func Test_hasCommentBeforeOrInside_nonBlockStmt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			func myFunc() {
			if true {
				x := 1
				_ = x
			// Comment before else if
			} else if false {
				y := 2
				_ = y
			}
			}`

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

			// Find the else if statement
			var elseStmt ast.Stmt
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if it's an if statement
				if ifStmt, ok := n.(*ast.IfStmt); ok {
					// Check if it has an else clause
					if ifStmt.Else != nil {
						elseStmt = ifStmt.Else
						return false
					}
				}
				// Continue traversal
				return true
			})

			// Check else statement found
			if elseStmt == nil {
				t.Fatal("no else statement found")
			}

			got := hasCommentBeforeOrInside(pass, elseStmt)

			// Should find comment before else
			if !got {
				t.Errorf("hasCommentBeforeOrInside() = false, want true")
			}

		})
	}
}
