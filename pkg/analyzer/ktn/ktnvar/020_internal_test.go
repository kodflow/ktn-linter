package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// TestFormatSelectorType tests the formatSelectorType function.
func TestFormatSelectorType(t *testing.T) {
	// Test with valid selector expression
	tests := []struct {
		name     string
		sel      *ast.SelectorExpr
		expected string
	}{
		{
			name: "with ident X",
			sel: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "time"},
				Sel: &ast.Ident{Name: "Time"},
			},
			expected: "time.Time",
		},
		{
			name: "with non-ident X",
			sel: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "foo"},
					Sel: &ast.Ident{Name: "bar"},
				},
				Sel: &ast.Ident{Name: "Baz"},
			},
			expected: "Baz",
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := formatSelectorType(tt.sel)
			// Check result
			if result != tt.expected {
				t.Errorf("formatSelectorType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestFormatMapType tests the formatMapType function.
func TestFormatMapType(t *testing.T) {
	// Test with valid map type
	tests := []struct {
		name     string
		m        *ast.MapType
		expected string
	}{
		{
			name: "simple map",
			m: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "int"},
			},
			expected: "map[string]int",
		},
		{
			name: "nested map",
			m: &ast.MapType{
				Key: &ast.Ident{Name: "string"},
				Value: &ast.MapType{
					Key:   &ast.Ident{Name: "int"},
					Value: &ast.Ident{Name: "bool"},
				},
			},
			expected: "map[string]map[int]bool",
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := formatMapType(tt.m)
			// Check result
			if result != tt.expected {
				t.Errorf("formatMapType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsSliceType tests the isSliceType function.
func TestIsSliceType(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "array type (slice)",
			expr:     &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
			expected: true,
		},
		{
			name:     "ident (not slice)",
			expr:     &ast.Ident{Name: "int"},
			expected: false,
		},
		{
			name:     "map type (not slice)",
			expr:     &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isSliceType(tt.expr)
			// Check result
			if result != tt.expected {
				t.Errorf("isSliceType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsMakeCall tests the isMakeCall function.
func TestIsMakeCall(t *testing.T) {
	tests := []struct {
		name     string
		call     *ast.CallExpr
		expected bool
	}{
		{
			name: "make call",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
			},
			expected: true,
		},
		{
			name: "append call",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "append"},
			},
			expected: false,
		},
		{
			name: "selector expr (not ident)",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "pkg"},
					Sel: &ast.Ident{Name: "Func"},
				},
			},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isMakeCall(tt.call)
			// Check result
			if result != tt.expected {
				t.Errorf("isMakeCall() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsZeroLiteral tests the isZeroLiteral function.
func TestIsZeroLiteral(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "zero literal",
			expr:     &ast.BasicLit{Value: "0"},
			expected: true,
		},
		{
			name:     "non-zero literal",
			expr:     &ast.BasicLit{Value: "1"},
			expected: false,
		},
		{
			name:     "ident (not literal)",
			expr:     &ast.Ident{Name: "x"},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isZeroLiteral(tt.expr)
			// Check result
			if result != tt.expected {
				t.Errorf("isZeroLiteral() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestFormatSliceType tests the formatSliceType function.
func TestFormatSliceType(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple slice",
			expr:     &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
			expected: "[]int",
		},
		{
			name:     "pointer slice",
			expr:     &ast.ArrayType{Elt: &ast.StarExpr{X: &ast.Ident{Name: "string"}}},
			expected: "[]*string",
		},
		{
			name:     "not array type",
			expr:     &ast.Ident{Name: "int"},
			expected: "[]T",
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := formatSliceType(tt.expr)
			// Check result
			if result != tt.expected {
				t.Errorf("formatSliceType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestFormatElementType tests the formatElementType function.
func TestFormatElementType(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "ident",
			expr:     &ast.Ident{Name: "int"},
			expected: "int",
		},
		{
			name: "selector",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "time"},
				Sel: &ast.Ident{Name: "Time"},
			},
			expected: "time.Time",
		},
		{
			name:     "pointer",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "string"}},
			expected: "*string",
		},
		{
			name:     "nested slice",
			expr:     &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
			expected: "[]int",
		},
		{
			name: "map",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "int"},
			},
			expected: "map[string]int",
		},
		{
			name:     "interface",
			expr:     &ast.InterfaceType{},
			expected: "interface{}",
		},
		{
			name:     "unknown type",
			expr:     &ast.ChanType{Value: &ast.Ident{Name: "int"}},
			expected: "T",
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := formatElementType(tt.expr)
			// Check result
			if result != tt.expected {
				t.Errorf("formatElementType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCheckMakeSliceZero tests the checkMakeSliceZero function.
func TestCheckMakeSliceZero(t *testing.T) {
	tests := []struct {
		name        string
		call        *ast.CallExpr
		expectError bool
	}{
		{
			name: "not make call",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "append"},
			},
			expectError: false,
		},
		{
			name: "make with one arg",
			call: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "make"},
				Args: []ast.Expr{&ast.ArrayType{Elt: &ast.Ident{Name: "int"}}},
			},
			expectError: false,
		},
		{
			name: "make with non-slice type",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.Ident{Name: "map"},
					&ast.BasicLit{Value: "0"},
				},
			},
			expectError: false,
		},
		{
			name: "make with non-zero length",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.BasicLit{Value: "10"},
				},
			},
			expectError: false,
		},
		{
			name: "make with capacity",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.BasicLit{Value: "0"},
					&ast.BasicLit{Value: "10"},
				},
			},
			expectError: false,
		},
		{
			name: "make with zero no capacity - should report",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.BasicLit{Value: "0"},
				},
			},
			expectError: true,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reported := false
			pass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reported = true
				},
			}
			// Call function
			checkMakeSliceZero(pass, tt.call)
			// Check result
			if reported != tt.expectError {
				t.Errorf("checkMakeSliceZero() reported = %v, want %v", reported, tt.expectError)
			}
		})
	}
}

// TestCheckEmptySliceLiteral tests the checkEmptySliceLiteral function.
func TestCheckEmptySliceLiteral(t *testing.T) {
	tests := []struct {
		name        string
		lit         *ast.CompositeLit
		expectError bool
	}{
		{
			name: "not slice type",
			lit: &ast.CompositeLit{
				Type: &ast.Ident{Name: "MyStruct"},
			},
			expectError: false,
		},
		{
			name: "non-empty slice",
			lit: &ast.CompositeLit{
				Type: &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
				Elts: []ast.Expr{&ast.BasicLit{Value: "1"}},
			},
			expectError: false,
		},
		{
			name: "empty slice - should report",
			lit: &ast.CompositeLit{
				Type: &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
				Elts: nil,
			},
			expectError: true,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reported := false
			pass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reported = true
				},
			}
			// Call function
			checkEmptySliceLiteral(pass, tt.lit)
			// Check result
			if reported != tt.expectError {
				t.Errorf("checkEmptySliceLiteral() reported = %v, want %v", reported, tt.expectError)
			}
		})
	}
}

// Test_runVar020_disabled tests runVar020 with disabled rule.
func Test_runVar020_disabled(t *testing.T) {
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
					"KTN-VAR-020": {Enabled: config.Bool(false)},
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

			_, err = runVar020(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar020() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar020() reported %d issues, expected 0 when disabled", reportCount)
			}
		})
	}
}

// Test_runVar020_fileExcluded tests runVar020 with excluded file.
func Test_runVar020_fileExcluded(t *testing.T) {
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
					"KTN-VAR-020": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse code with empty slice
			code := `package test
var s = []int{}
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

			_, err = runVar020(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar020() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar020() reported %d issues, expected 0 when file excluded", reportCount)
			}
		})
	}
}

// Test_runVar020_nilFset tests runVar020 with nil Fset.
func Test_runVar020_nilFset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"nil fset returns early"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Reset config to default
			config.Reset()

			// Parse simple code
			code := `package test
var s = []int{}
`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			// Pass with nil Fset
			pass := &analysis.Pass{
				Fset: nil, // Nil Fset to test early return
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			result, err := runVar020(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar020() error = %v", err)
			}

			// Should return nil
			if result != nil {
				t.Errorf("runVar020() result = %v, expected nil", result)
			}

			// Should not report anything when Fset is nil
			if reportCount != 0 {
				t.Errorf("runVar020() reported %d issues, expected 0 when Fset is nil", reportCount)
			}
		})
	}
}

// Test_checkEmptySliceLiteral_missingMessage tests with missing message.
func Test_checkEmptySliceLiteral_missingMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"fallback message when not registered"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Unregister the message temporarily
			messages.Unregister(ruleCodeVar020)
			defer func() {
				// Re-register the message
				messages.Register(messages.Message{
					Code:    ruleCodeVar020,
					Short:   "préférer nil slice à %s{}",
					Verbose: "préférer nil slice à %s{}",
				})
			}()

			reported := false
			pass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reported = true
				},
			}

			// Create empty slice literal
			lit := &ast.CompositeLit{
				Type: &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
				Elts: nil,
			}

			// Call function
			checkEmptySliceLiteral(pass, lit)

			// Should report using fallback message
			if !reported {
				t.Error("checkEmptySliceLiteral() did not report with fallback message")
			}
		})
	}
}

// Test_checkMakeSliceZero_missingMessage tests with missing message.
func Test_checkMakeSliceZero_missingMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"fallback message when not registered"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Unregister the message temporarily
			messages.Unregister(ruleCodeVar020)
			defer func() {
				// Re-register the message
				messages.Register(messages.Message{
					Code:    ruleCodeVar020,
					Short:   "préférer nil slice à %s{}",
					Verbose: "préférer nil slice à %s{}",
				})
			}()

			reported := false
			pass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reported = true
				},
			}

			// Create make([]int, 0) call
			call := &ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.BasicLit{Value: "0"},
				},
			}

			// Call function
			checkMakeSliceZero(pass, call)

			// Should report using fallback message
			if !reported {
				t.Error("checkMakeSliceZero() did not report with fallback message")
			}
		})
	}
}

// Test_runVar020_fullExecution tests runVar020 with actual code triggering reports.
func Test_runVar020_fullExecution(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedCount int
	}{
		{
			name: "empty slice literal triggers report",
			code: `package test
var s = []int{}
`,
			expectedCount: 1,
		},
		{
			name: "make slice with zero no capacity triggers report",
			code: `package test
var s = make([]int, 0)
`,
			expectedCount: 1,
		},
		{
			name: "both patterns trigger reports",
			code: `package test
var s1 = []int{}
var s2 = make([]string, 0)
`,
			expectedCount: 2,
		},
		{
			name: "valid code no report",
			code: `package test
var s1 = []int{1, 2, 3}
var s2 = make([]int, 0, 10)
var s3 []int
`,
			expectedCount: 0,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Reset config to default
			config.Reset()

			// Parse code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
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

			_, err = runVar020(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar020() error = %v", err)
			}

			// Check report count
			if reportCount != tt.expectedCount {
				t.Errorf("runVar020() reported %d issues, expected %d", reportCount, tt.expectedCount)
			}
		})
	}
}
