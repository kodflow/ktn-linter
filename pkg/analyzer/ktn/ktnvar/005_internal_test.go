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

// TestMaxVarNameLength005 tests the maximum name length constant.
func TestMaxVarNameLength005(t *testing.T) {
	if maxVarNameLength005 != 30 {
		t.Errorf("maxVarNameLength005 = %d, want 30", maxVarNameLength005)
	}
}

// TestNameLengthLimits005 tests edge cases for name length.
func TestNameLengthLimits005(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		tooLong bool
	}{
		{
			name:    "exactly 30 chars",
			varName: "ThisIsExactlyThirtyCharacterss",
			tooLong: false,
		},
		{
			name:    "31 chars",
			varName: "ThisIsExactlyThirtyOneCharacter",
			tooLong: true,
		},
		{
			name:    "short name",
			varName: "ok",
			tooLong: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := len(tt.varName) > maxVarNameLength005
			if result != tt.tooLong {
				t.Errorf("len(%q) > %d = %v, want %v",
					tt.varName, maxVarNameLength005, result, tt.tooLong)
			}
		})
	}
}

// Test_runVar005_disabled tests rule disabled branch.
func Test_runVar005_disabled(t *testing.T) {
	config.Reset()
	cfg := config.Get()
	cfg.Rules = map[string]*config.RuleConfig{
		ruleCodeVar005: {Enabled: config.Bool(false)},
	}

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", `package test; var a = 1`, 0)
	insp := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) {},
	}

	result, err := runVar005(pass)
	if err != nil || result != nil {
		t.Errorf("runVar005() = (%v, %v), want (nil, nil)", result, err)
	}
	config.Reset()
}

// Test_checkVar005PackageLevel_nonVarDecl tests non-var declarations.
func Test_checkVar005PackageLevel_nonVarDecl(t *testing.T) {
	config.Reset()

	fset := token.NewFileSet()
	code := `package test
const x = 1
func foo() {}
`
	file, _ := parser.ParseFile(fset, "test.go", code, 0)
	insp := inspector.New([]*ast.File{file})
	cfg := config.Get()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	checkVar005PackageLevel(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports, got %d", reportCount)
	}
}

// Test_checkVar005LocalVars_nilBody tests function without body.
func Test_checkVar005LocalVars_nilBody(t *testing.T) {
	config.Reset()

	fset := token.NewFileSet()
	code := `package test
//go:linkname foo runtime.foo
func foo()
`
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	insp := inspector.New([]*ast.File{file})
	cfg := config.Get()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	checkVar005LocalVars(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for nil body, got %d", reportCount)
	}
}

// Test_checkVar005AssignStmt_notDefine tests non-DEFINE assignments.
func Test_checkVar005AssignStmt_notDefine(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
	}

	checkVar005AssignStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for = assignment, got %d", reportCount)
	}
}

// Test_checkVar005AssignStmt_nonIdent tests non-identifier LHS.
func Test_checkVar005AssignStmt_nonIdent(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.IndexExpr{
			X:     &ast.Ident{Name: "arr"},
			Index: &ast.BasicLit{Kind: token.INT, Value: "0"},
		}},
	}

	checkVar005AssignStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-ident LHS, got %d", reportCount)
	}
}

// Test_checkVar005RangeStmt_notDefine tests non-DEFINE range.
func Test_checkVar005RangeStmt_notDefine(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.RangeStmt{
		Tok: token.ASSIGN,
		Key: &ast.Ident{Name: "i"},
	}

	checkVar005RangeStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for = range, got %d", reportCount)
	}
}

// Test_checkVar005RangeStmt_nonIdentKey tests non-identifier key.
func Test_checkVar005RangeStmt_nonIdentKey(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.RangeStmt{
		Tok: token.DEFINE,
		Key: &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
	}

	checkVar005RangeStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-ident key, got %d", reportCount)
	}
}

// Test_checkVar005RangeStmt_nilValue tests nil value in range.
func Test_checkVar005RangeStmt_nilValue(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.RangeStmt{
		Tok:   token.DEFINE,
		Key:   &ast.Ident{Name: "idx"},
		Value: nil,
	}

	checkVar005RangeStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for nil value, got %d", reportCount)
	}
}

// Test_checkVar005RangeStmt_nonIdentValue tests non-identifier value.
func Test_checkVar005RangeStmt_nonIdentValue(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.RangeStmt{
		Tok:   token.DEFINE,
		Key:   &ast.Ident{Name: "idx"},
		Value: &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
	}

	checkVar005RangeStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-ident value, got %d", reportCount)
	}
}

// Test_checkVar005DeclStmt_nonGenDecl tests non-GenDecl.
func Test_checkVar005DeclStmt_nonGenDecl(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.DeclStmt{
		Decl: &ast.FuncDecl{Name: &ast.Ident{Name: "foo"}},
	}

	checkVar005DeclStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-GenDecl, got %d", reportCount)
	}
}

// Test_checkVar005DeclStmt_constDecl tests const declarations.
func Test_checkVar005DeclStmt_constDecl(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	stmt := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.CONST,
			Specs: []ast.Spec{
				&ast.ValueSpec{Names: []*ast.Ident{{Name: "x"}}},
			},
		},
	}

	checkVar005DeclStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for const decl, got %d", reportCount)
	}
}

// Test_checkVar005Name_blankIdent tests blank identifier.
func Test_checkVar005Name_blankIdent(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	ident := &ast.Ident{Name: "_"}
	checkVar005Name(pass, ident)

	if reportCount != 0 {
		t.Errorf("expected 0 reports for blank ident, got %d", reportCount)
	}
}

// Test_checkVar005Name_withinLimit tests names within limit.
func Test_checkVar005Name_withinLimit(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	ident := &ast.Ident{Name: "shortName"}
	checkVar005Name(pass, ident)

	if reportCount != 0 {
		t.Errorf("expected 0 reports for short name, got %d", reportCount)
	}
}

// Test_checkVar005Node_nilNode tests nil node handling.
func Test_checkVar005Node_nilNode(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	checkVar005Node(pass, nil)

	if reportCount != 0 {
		t.Errorf("expected 0 reports for nil node, got %d", reportCount)
	}
}

// Test_checkVar005Node_otherNodeTypes tests unhandled node types.
func Test_checkVar005Node_otherNodeTypes(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	checkVar005Node(pass, &ast.ReturnStmt{})
	checkVar005Node(pass, &ast.IfStmt{})

	if reportCount != 0 {
		t.Errorf("expected 0 reports for unhandled nodes, got %d", reportCount)
	}
}

// Test_checkVar005PackageLevel_excluded tests file exclusion.
func Test_checkVar005PackageLevel_excluded(t *testing.T) {
	config.Reset()

	cfg := config.Get()
	cfg.Rules = map[string]*config.RuleConfig{
		ruleCodeVar005: {Exclude: []string{"test.go"}},
	}

	fset := token.NewFileSet()
	code := `package test
var veryLongVariableNameThatExceedsLimitOf30Chars = 1
`
	file, _ := parser.ParseFile(fset, "test.go", code, 0)
	insp := inspector.New([]*ast.File{file})

	reportCount := 0
	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	checkVar005PackageLevel(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for excluded file, got %d", reportCount)
	}
	config.Reset()
}

// Test_checkVar005LocalVars_excluded tests file exclusion for local vars.
func Test_checkVar005LocalVars_excluded(t *testing.T) {
	config.Reset()

	cfg := config.Get()
	cfg.Rules = map[string]*config.RuleConfig{
		ruleCodeVar005: {Exclude: []string{"test.go"}},
	}

	fset := token.NewFileSet()
	code := `package test
func foo() {
	veryLongVariableNameThatExceedsLimitOf30Chars := 1
	_ = veryLongVariableNameThatExceedsLimitOf30Chars
}
`
	file, _ := parser.ParseFile(fset, "test.go", code, 0)
	insp := inspector.New([]*ast.File{file})

	reportCount := 0
	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	checkVar005LocalVars(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for excluded file, got %d", reportCount)
	}
	config.Reset()
}
