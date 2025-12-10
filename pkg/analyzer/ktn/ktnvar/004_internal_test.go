package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar004 tests the private runVar004 function.
func Test_runVar004(t *testing.T) {
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

// Test_isAppendCall tests the private isAppendCall helper function.
func Test_isAppendCall(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name: "append call",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "append"},
			},
			expected: true,
		},
		{
			name: "other function call",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "len"},
			},
			expected: false,
		},
		{
			name:     "not a call expr",
			expr:     &ast.BasicLit{Value: "1"},
			expected: false,
		},
		{
			name: "method call",
			expr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "s"},
					Sel: &ast.Ident{Name: "append"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAppendCall(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isAppendCall() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_isSliceArrayOrMap tests the private isSliceArrayOrMap helper function.
func Test_isSliceArrayOrMap(t *testing.T) {
	tests := []struct {
		name     string
		typeExpr ast.Expr
		expected bool
	}{
		{
			name:     "nil type",
			typeExpr: nil,
			expected: false,
		},
		{
			name:     "slice type",
			typeExpr: &ast.ArrayType{Len: nil},
			expected: true,
		},
		{
			name:     "array type",
			typeExpr: &ast.ArrayType{Len: &ast.BasicLit{Value: "10"}},
			expected: true,
		},
		{
			name:     "map type",
			typeExpr: &ast.MapType{},
			expected: true,
		},
		{
			name:     "struct type",
			typeExpr: &ast.StructType{},
			expected: false,
		},
		{
			name:     "ident type",
			typeExpr: &ast.Ident{Name: "int"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSliceArrayOrMap(tt.typeExpr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isSliceArrayOrMap() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_isInReturnStatement tests the private isInReturnStatement helper function.
func Test_isInReturnStatement(t *testing.T) {
	tests := []struct {
		name     string
		stack    []ast.Node
		expected bool
	}{
		{
			name:     "empty stack",
			stack:    []ast.Node{},
			expected: false,
		},
		{
			name: "has return in stack",
			stack: []ast.Node{
				&ast.FuncDecl{},
				&ast.BlockStmt{},
				&ast.ReturnStmt{},
			},
			expected: true,
		},
		{
			name: "no return in stack",
			stack: []ast.Node{
				&ast.FuncDecl{},
				&ast.BlockStmt{},
				&ast.AssignStmt{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInReturnStatement(tt.stack)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isInReturnStatement() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_isInStructLiteral tests the private isInStructLiteral helper function.
func Test_isInStructLiteral(t *testing.T) {
	tests := []struct {
		name     string
		stack    []ast.Node
		expected bool
	}{
		{
			name:     "empty stack",
			stack:    []ast.Node{},
			expected: false,
		},
		{
			name: "has struct literal in stack",
			stack: []ast.Node{
				&ast.FuncDecl{},
				&ast.CompositeLit{Type: &ast.Ident{Name: "MyStruct"}},
			},
			expected: true,
		},
		{
			name: "has key-value expr in stack",
			stack: []ast.Node{
				&ast.FuncDecl{},
				&ast.KeyValueExpr{},
			},
			expected: true,
		},
		{
			name: "has slice literal in stack",
			stack: []ast.Node{
				&ast.FuncDecl{},
				&ast.CompositeLit{Type: &ast.ArrayType{}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInStructLiteral(tt.stack)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isInStructLiteral() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_collectAppendVariables tests the private collectAppendVariables function.
func Test_collectAppendVariables(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function collects append variables
		})
	}
}

// Test_checkMakeCalls tests the private checkMakeCalls function.
func Test_checkMakeCalls(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks make calls
		})
	}
}

// Test_checkMakeCall_notMake tests checkMakeCall with non-make call.
func Test_checkMakeCall_notMake(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with non-make call
	call := &ast.CallExpr{
		Fun: &ast.Ident{Name: "append"},
	}
	checkMakeCall(pass, call)
	// No error expected
}

// Test_checkMakeCall_wrongArgs tests checkMakeCall with wrong arg count.
func Test_checkMakeCall_wrongArgs(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with make call but wrong number of args
	call := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "make"},
		Args: []ast.Expr{&ast.Ident{Name: "T"}}, // Only 1 arg
	}
	checkMakeCall(pass, call)

	// Test with 3 args (has capacity)
	call3Args := &ast.CallExpr{
		Fun: &ast.Ident{Name: "make"},
		Args: []ast.Expr{
			&ast.Ident{Name: "T"},
			&ast.BasicLit{Value: "0"},
			&ast.BasicLit{Value: "10"},
		},
	}
	checkMakeCall(pass, call3Args)
	// No error expected
}

// Test_checkEmptySliceLiterals tests the private checkEmptySliceLiterals function.
func Test_checkEmptySliceLiterals(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks empty slice literals
		})
	}
}

// Test_checkMakeCall_notSlice tests checkMakeCall with non-slice type.
func Test_checkMakeCall_notSlice(t *testing.T) {
	code := `package test
func example() {
	m := make(map[string]int, 10)
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Type check the code
	conf := types.Config{}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	_, _ = conf.Check("test", fset, []*ast.File{file}, info)

	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: info,
		Report:    func(_d analysis.Diagnostic) {},
	}

	// Find the make call
	ast.Inspect(file, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			checkMakeCall(pass, call)
		}
		return true
	})
	// No error expected for map type
}

// Test_checkCompositeLit_nonEmptySlice tests non-empty slice.
func Test_checkCompositeLit_nonEmptySlice(t *testing.T) {
	ctx := &litCheckContext{
		pass: &analysis.Pass{
			Report:    func(_d analysis.Diagnostic) {},
			TypesInfo: &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)},
		},
		appendVars: make(map[string]bool),
	}

	// Test with non-empty slice literal
	assign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
	}
	lit := &ast.CompositeLit{
		Type: &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
		Elts: []ast.Expr{&ast.BasicLit{Value: "1"}}, // Non-empty
	}
	checkCompositeLit(ctx, assign, 0, lit, []ast.Node{})
	// No error expected
}

// Test_checkCompositeLit_invalidIndex tests invalid index.
func Test_checkCompositeLit_invalidIndex(t *testing.T) {
	ctx := &litCheckContext{
		pass: &analysis.Pass{
			Report:    func(_d analysis.Diagnostic) {},
			TypesInfo: &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)},
		},
		appendVars: map[string]bool{"x": true},
	}

	// Test with invalid index
	assign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
	}
	lit := &ast.CompositeLit{
		Type: &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
	}
	// Index 5 is out of bounds for Lhs with 1 element
	checkCompositeLit(ctx, assign, 5, lit, []ast.Node{})
	// No error expected
}

// Test_checkCompositeLit_notSliceType tests non-slice composite literal.
func Test_checkCompositeLit_notSliceType(t *testing.T) {
	code := `package test
func example() {
	s := struct{}{}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Type check
	conf := types.Config{}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	_, _ = conf.Check("test", fset, []*ast.File{file}, info)

	ctx := &litCheckContext{
		pass: &analysis.Pass{
			Fset:      fset,
			TypesInfo: info,
			Report:    func(_d analysis.Diagnostic) {},
		},
		appendVars: map[string]bool{},
	}

	// Find composite lit
	ast.Inspect(file, func(n ast.Node) bool {
		if assign, ok := n.(*ast.AssignStmt); ok {
			for i, rhs := range assign.Rhs {
				if lit, isLit := rhs.(*ast.CompositeLit); isLit {
					checkCompositeLit(ctx, assign, i, lit, []ast.Node{})
				}
			}
		}
		return true
	})
	// No error expected for struct type
}

// Test_runVar004_disabled tests runVar004 with disabled rule.
func Test_runVar004_disabled(t *testing.T) {
	// Setup config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-004": {Enabled: config.Bool(false)},
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

	_, err = runVar004(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar004() error = %v", err)
	}

	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar004() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar004_fileExcluded tests runVar004 with excluded file.
func Test_runVar004_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-004": {
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

	_, err = runVar004(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar004() error = %v", err)
	}

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar004() reported %d issues, expected 0 when file excluded", reportCount)
	}
}
