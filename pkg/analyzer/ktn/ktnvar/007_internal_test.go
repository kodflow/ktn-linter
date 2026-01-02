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

// Test_hasInitWithoutType tests the hasInitWithoutType helper function.
func Test_hasInitWithoutType(t *testing.T) {
	tests := []struct {
		name     string
		spec     *ast.ValueSpec
		expected bool
	}{
		{
			name: "has init without type",
			spec: &ast.ValueSpec{
				Names:  []*ast.Ident{{Name: "x"}},
				Type:   nil,
				Values: []ast.Expr{&ast.BasicLit{Value: "1"}},
			},
			expected: true,
		},
		{
			name: "has init with type",
			spec: &ast.ValueSpec{
				Names:  []*ast.Ident{{Name: "x"}},
				Type:   &ast.Ident{Name: "int"},
				Values: []ast.Expr{&ast.BasicLit{Value: "1"}},
			},
			expected: true,
		},
		{
			name: "no init",
			spec: &ast.ValueSpec{
				Names:  []*ast.Ident{{Name: "x"}},
				Type:   &ast.Ident{Name: "int"},
				Values: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasInitWithoutType(tt.spec)
			// Check result
			if result != tt.expected {
				t.Errorf("hasInitWithoutType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test_runVar007_disabled tests runVar007 with disabled rule.
func Test_runVar007_disabled(t *testing.T) {
	// Setup config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-007": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	code := `package test
func example() { var x = 42; _ = x }
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

	result, err := runVar007(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar007() error = %v", err)
	}
	// Should return nil
	if result != nil {
		t.Errorf("runVar007() result = %v, expected nil", result)
	}
	// Should not report when disabled
	if reportCount != 0 {
		t.Errorf("runVar007() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar007_nilInspector tests runVar007 with nil inspector.
func Test_runVar007_nilInspector(t *testing.T) {
	// Reset config
	config.Reset()

	fset := token.NewFileSet()
	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: nil,
		},
		Report: func(_d analysis.Diagnostic) {},
	}

	result, err := runVar007(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar007() error = %v", err)
	}
	// Should return nil
	if result != nil {
		t.Errorf("runVar007() result = %v, expected nil", result)
	}
}

// Test_runVar007_invalidInspector tests runVar007 with wrong type.
func Test_runVar007_invalidInspector(t *testing.T) {
	// Reset config
	config.Reset()

	fset := token.NewFileSet()
	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: "not an inspector",
		},
		Report: func(_d analysis.Diagnostic) {},
	}

	result, err := runVar007(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar007() error = %v", err)
	}
	// Should return nil
	if result != nil {
		t.Errorf("runVar007() result = %v, expected nil", result)
	}
}

// Test_runVar007_nilFset tests runVar007 with nil Fset.
func Test_runVar007_nilFset(t *testing.T) {
	// Reset config
	config.Reset()

	code := `package test
func example() { var x = 42; _ = x }
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	pass := &analysis.Pass{
		Fset: nil, // nil Fset
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {},
	}

	result, err := runVar007(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar007() error = %v", err)
	}
	// Should return nil
	if result != nil {
		t.Errorf("runVar007() result = %v, expected nil", result)
	}
}

// Test_runVar007_wrongInspectorType tests runVar007 with wrong inspector type.
func Test_runVar007_wrongInspectorType(t *testing.T) {
	config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: "not an inspector"},
		Report:   func(_ analysis.Diagnostic) {},
	}

	result, err := runVar007(pass)
	if err != nil {
		t.Errorf("runVar007() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar007() result = %v, want nil", result)
	}
}

// Test_runVar007_fileExcluded tests runVar007 with excluded file.
func Test_runVar007_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-007": {
				Exclude: []string{"test.go"},
			},
		},
	})
	defer config.Reset()

	code := `package test
func example() { var x = 42; _ = x }
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

	_, err = runVar007(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar007() error = %v", err)
	}
	// Should not report when file excluded
	if reportCount != 0 {
		t.Errorf("runVar007() reported %d issues, expected 0 when file excluded", reportCount)
	}
}

// Test_runVar007_nilFuncBody tests runVar007 with nil function body.
func Test_runVar007_nilFuncBody(t *testing.T) {
	// Reset config
	config.Reset()

	// External function declaration has nil body
	code := `package test
func external() int
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

	_, err = runVar007(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar007() error = %v", err)
	}
	// Should not report for nil body
	if reportCount != 0 {
		t.Errorf("runVar007() reported %d issues, expected 0 for nil body", reportCount)
	}
}

// Test_checkControlFlowStmt_returnFalse tests the false return path.
func Test_checkControlFlowStmt_returnFalse(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with non-control-flow statement
	exprStmt := &ast.ExprStmt{
		X: &ast.Ident{Name: "x"},
	}
	result := checkControlFlowStmt(pass, exprStmt)
	// Should return false
	if result {
		t.Errorf("checkControlFlowStmt() = true, expected false for non-control-flow statement")
	}
}

// Test_checkControlFlowStmt_allCases tests all control flow statement types.
func Test_checkControlFlowStmt_allCases(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test IfStmt
	ifResult := checkControlFlowStmt(pass, &ast.IfStmt{Body: &ast.BlockStmt{}})
	// Should return true
	if !ifResult {
		t.Errorf("checkControlFlowStmt(IfStmt) = false, expected true")
	}

	// Test ForStmt
	forResult := checkControlFlowStmt(pass, &ast.ForStmt{Body: &ast.BlockStmt{}})
	// Should return true
	if !forResult {
		t.Errorf("checkControlFlowStmt(ForStmt) = false, expected true")
	}

	// Test RangeStmt
	rangeResult := checkControlFlowStmt(pass, &ast.RangeStmt{Body: &ast.BlockStmt{}})
	// Should return true
	if !rangeResult {
		t.Errorf("checkControlFlowStmt(RangeStmt) = false, expected true")
	}

	// Test SwitchStmt
	switchResult := checkControlFlowStmt(pass, &ast.SwitchStmt{Body: &ast.BlockStmt{}})
	// Should return true
	if !switchResult {
		t.Errorf("checkControlFlowStmt(SwitchStmt) = false, expected true")
	}

	// Test TypeSwitchStmt
	typeSwitchResult := checkControlFlowStmt(pass, &ast.TypeSwitchStmt{Body: &ast.BlockStmt{}})
	// Should return true
	if !typeSwitchResult {
		t.Errorf("checkControlFlowStmt(TypeSwitchStmt) = false, expected true")
	}
}

// Test_checkCaseClause_withBody tests checkCaseClause with statements.
func Test_checkCaseClause_withBody(t *testing.T) {
	// Reset config
	config.Reset()

	fset := token.NewFileSet()
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Create case clause with var declaration
	varDecl := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  []*ast.Ident{{Name: "x"}},
					Values: []ast.Expr{&ast.BasicLit{Value: "1"}},
				},
			},
		},
	}
	caseClause := &ast.CaseClause{
		Body: []ast.Stmt{varDecl},
	}
	checkCaseClause(pass, caseClause)
	// Should report the var declaration
	if reportCount != 1 {
		t.Errorf("checkCaseClause() reported %d issues, expected 1", reportCount)
	}
}

// Test_checkCommClause_withBody tests checkCommClause with statements.
func Test_checkCommClause_withBody(t *testing.T) {
	// Reset config
	config.Reset()

	fset := token.NewFileSet()
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Create comm clause with var declaration
	varDecl := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  []*ast.Ident{{Name: "y"}},
					Values: []ast.Expr{&ast.BasicLit{Value: "2"}},
				},
			},
		},
	}
	commClause := &ast.CommClause{
		Body: []ast.Stmt{varDecl},
	}
	checkCommClause(pass, commClause)
	// Should report the var declaration
	if reportCount != 1 {
		t.Errorf("checkCommClause() reported %d issues, expected 1", reportCount)
	}
}

// Test_reportVarErrors_missingMessage tests the fallback message path.
func Test_reportVarErrors_missingMessage(t *testing.T) {
	// Reset config
	config.Reset()

	// Store original message and remove it
	originalMsg, hasOriginal := messages.Get("KTN-VAR-007")
	messages.Unregister("KTN-VAR-007")
	defer func() {
		// Restore original message
		if hasOriginal {
			messages.Register(originalMsg)
		}
	}()

	fset := token.NewFileSet()
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	spec := &ast.ValueSpec{
		Names:  []*ast.Ident{{Name: "x"}},
		Values: []ast.Expr{&ast.BasicLit{Value: "1"}},
	}
	reportVarErrors(pass, spec)
	// Should report with fallback message
	if reportCount != 1 {
		t.Errorf("reportVarErrors() reported %d issues, expected 1 with fallback", reportCount)
	}
}

// Test_checkStatement_emptyStmt tests checkStatement with EmptyStmt.
func Test_checkStatement_emptyStmt(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Call checkStatement with empty statement
	checkStatement(pass, &ast.EmptyStmt{})
	// No error expected
}

// Test_checkStatement_badDecl tests non-GenDecl in DeclStmt.
func Test_checkStatement_badDecl(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Create DeclStmt with FuncDecl (not GenDecl)
	badDecl := &ast.DeclStmt{
		Decl: &ast.FuncDecl{Name: &ast.Ident{Name: "test"}},
	}
	checkStatement(pass, badDecl)
	// No error expected - should skip non-GenDecl
}

// Test_checkStatement_constDecl tests const declaration.
func Test_checkStatement_constDecl(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Create const declaration
	constDecl := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.CONST,
		},
	}
	checkStatement(pass, constDecl)
	// No error expected - should skip const declarations
}

// Test_checkNestedBlocks_allTypes tests all statement types.
func Test_checkNestedBlocks_allTypes(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test switch statement
	checkNestedBlocks(pass, &ast.SwitchStmt{Body: &ast.BlockStmt{}})

	// Test type switch statement
	checkNestedBlocks(pass, &ast.TypeSwitchStmt{Body: &ast.BlockStmt{}})

	// Test select statement
	checkNestedBlocks(pass, &ast.SelectStmt{Body: &ast.BlockStmt{}})

	// Test for statement
	checkNestedBlocks(pass, &ast.ForStmt{Body: &ast.BlockStmt{}})

	// Test range statement
	checkNestedBlocks(pass, &ast.RangeStmt{Body: &ast.BlockStmt{}})

	// Test block statement
	checkNestedBlocks(pass, &ast.BlockStmt{List: []ast.Stmt{}})

	// Test case clause
	checkNestedBlocks(pass, &ast.CaseClause{Body: []ast.Stmt{}})

	// Test comm clause
	checkNestedBlocks(pass, &ast.CommClause{Body: []ast.Stmt{}})
}

// Test_checkIfStmt_allBranches tests checkIfStmt with all branches.
func Test_checkIfStmt_allBranches(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with nil body
	checkIfStmt(pass, &ast.IfStmt{Body: nil, Else: nil})

	// Test with body but no else
	checkIfStmt(pass, &ast.IfStmt{Body: &ast.BlockStmt{}, Else: nil})

	// Test with both body and else
	checkIfStmt(pass, &ast.IfStmt{Body: &ast.BlockStmt{}, Else: &ast.BlockStmt{}})
}

// Test_checkBlockIfNotNil_nilBlock tests checkBlockIfNotNil with nil.
func Test_checkBlockIfNotNil_nilBlock(t *testing.T) {
	pass := &analysis.Pass{
		Report: func(_d analysis.Diagnostic) {},
	}

	// Test with nil block
	checkBlockIfNotNil(pass, nil)
	// No error expected
}
