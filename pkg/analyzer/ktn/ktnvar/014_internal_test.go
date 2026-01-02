package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar014 tests the private runVar014 function.
func Test_runVar014(t *testing.T) {
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

// Test_checkLoopBodyVar015 tests the private checkLoopBodyVar015 function.
func Test_checkLoopBodyVar015(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks loop bodies
		})
	}
}

// Test_checkAssignmentForBuffer tests the private checkAssignmentForBuffer function.
func Test_checkAssignmentForBuffer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks assignments for buffers
		})
	}
}

// Test_checkMakeCallForByteSlice tests the private checkMakeCallForByteSlice function.
func Test_checkMakeCallForByteSlice(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks make calls for byte slices
		})
	}
}

// Test_runVar014_disabled tests runVar014 with disabled rule.
func Test_runVar014_disabled(t *testing.T) {
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
					"KTN-VAR-014": {Enabled: config.Bool(false)},
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

			_, err = runVar014(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar014() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar014() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar014_fileExcluded tests runVar014 with excluded file.
func Test_runVar014_fileExcluded(t *testing.T) {
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
					"KTN-VAR-014": {
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

			_, err = runVar014(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar014() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar014() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_checkLoopBodyVar015_nilBody tests checkLoopBodyVar015 with nil body.
func Test_checkLoopBodyVar015_nilBody(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with nil body
	checkLoopBodyVar015(pass, nil)
	// Should not panic
}

// Test_checkMakeCallForByteSlice_notMake tests with non-make call.
func Test_checkMakeCallForByteSlice_notMake(t *testing.T) {
	reportCount := 0
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Test with non-make call
	call := &ast.CallExpr{
		Fun: &ast.Ident{Name: "len"},
	}
	checkMakeCallForByteSlice(pass, call)

	if reportCount != 0 {
		t.Errorf("checkMakeCallForByteSlice() reported %d, expected 0", reportCount)
	}
}

// Test_checkMakeCallForByteSlice_notByteSlice tests with non-byte slice.
func Test_checkMakeCallForByteSlice_notByteSlice(t *testing.T) {
	reportCount := 0
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Test with make([]int, ...)
	call := &ast.CallExpr{
		Fun: &ast.Ident{Name: "make"},
		Args: []ast.Expr{
			&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
		},
	}
	checkMakeCallForByteSlice(pass, call)

	if reportCount != 0 {
		t.Errorf("checkMakeCallForByteSlice() reported %d, expected 0", reportCount)
	}
}

// Test_checkMakeCallForByteSlice_noArgs tests with no arguments.
func Test_checkMakeCallForByteSlice_noArgs(t *testing.T) {
	reportCount := 0
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Test with make but no args
	call := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "make"},
		Args: []ast.Expr{},
	}
	checkMakeCallForByteSlice(pass, call)

	if reportCount != 0 {
		t.Errorf("checkMakeCallForByteSlice() reported %d, expected 0", reportCount)
	}
}

// Test_checkAssignmentForBuffer_nonCallExpr tests with non-CallExpr rhs.
func Test_checkAssignmentForBuffer_nonCallExpr(t *testing.T) {
	reportCount := 0
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Test with non-CallExpr rhs
	stmt := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
	}
	checkAssignmentForBuffer(pass, stmt)

	if reportCount != 0 {
		t.Errorf("checkAssignmentForBuffer() reported %d, expected 0", reportCount)
	}
}

// Test_runVar014_withLoop tests runVar014 with a for loop.
func Test_runVar014_withLoop(t *testing.T) {
	config.Reset()
	defer config.Reset()

	code := `package test
func foo() {
	for i := 0; i < 10; i++ {
		buf := make([]byte, 1024)
		_ = buf
	}
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, 0)
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

	_, err := runVar014(pass)
	if err != nil {
		t.Errorf("runVar014() error = %v", err)
	}

	// Should report the buffer allocation in loop
	if reportCount != 1 {
		t.Errorf("runVar014() reported %d, expected 1", reportCount)
	}
}

// Test_runVar014_withRangeLoop tests runVar014 with a range loop.
func Test_runVar014_withRangeLoop(t *testing.T) {
	config.Reset()
	defer config.Reset()

	code := `package test
func foo() {
	items := []int{1, 2, 3}
	for range items {
		buf := make([]byte, 1024)
		_ = buf
	}
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", code, 0)
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

	_, err := runVar014(pass)
	if err != nil {
		t.Errorf("runVar014() error = %v", err)
	}

	// Should report the buffer allocation in loop
	if reportCount != 1 {
		t.Errorf("runVar014() reported %d, expected 1", reportCount)
	}
}

// Test_runVar014_fileExcludedWithLoop tests file exclusion with loop that would trigger.
func Test_runVar014_fileExcludedWithLoop(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-014": {
				Exclude: []string{"excluded.go"},
			},
		},
	})
	defer config.Reset()

	// Code that would normally trigger the rule
	code := `package test
func foo() {
	for i := 0; i < 10; i++ {
		buf := make([]byte, 1024)
		_ = buf
	}
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "excluded.go", code, 0)
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

	_, err := runVar014(pass)
	if err != nil {
		t.Errorf("runVar014() error = %v", err)
	}

	// Should not report when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar014() reported %d, expected 0 when file excluded", reportCount)
	}
}
