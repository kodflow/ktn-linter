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

// TestIdiomaticOneChar004 tests the idiomaticOneChar004 map.
func TestIdiomaticOneChar004(t *testing.T) {
	expectedOneChar := []string{"i", "j", "k", "n", "b", "c", "f", "m", "r", "s", "t", "w", "_"}

	for _, v := range expectedOneChar {
		if !idiomaticOneChar004[v] {
			t.Errorf("idiomaticOneChar004 should contain %q", v)
		}
	}
}

// TestIdiomaticShort004 tests the idiomaticShort004 map.
func TestIdiomaticShort004(t *testing.T) {
	expectedIdioms := []string{"ok"}

	for _, v := range expectedIdioms {
		if !idiomaticShort004[v] {
			t.Errorf("idiomaticShort004 should contain %q", v)
		}
	}
}

// TestMinVarNameLength004 tests the minimum name length constant.
func TestMinVarNameLength004(t *testing.T) {
	if minVarNameLength004 != 2 {
		t.Errorf("minVarNameLength004 = %d, want 2", minVarNameLength004)
	}
}

// Test_runVar004_disabled tests rule disabled branch.
func Test_runVar004_disabled(t *testing.T) {
	config.Reset()
	// Disable rule
	cfg := config.Get()
	cfg.Rules = map[string]*config.RuleConfig{
		ruleCodeVar004: {Enabled: config.Bool(false)},
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

	result, err := runVar004(pass)
	if err != nil || result != nil {
		t.Errorf("runVar004() = (%v, %v), want (nil, nil)", result, err)
	}
	config.Reset()
}

// Test_checkVar004PackageLevel_nonVarDecl tests non-var declarations.
func Test_checkVar004PackageLevel_nonVarDecl(t *testing.T) {
	config.Reset()

	fset := token.NewFileSet()
	// Code with const and func declarations (not var)
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

	checkVar004PackageLevel(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports, got %d", reportCount)
	}
}

// Test_checkVar004LocalVars_nilBody tests function without body.
func Test_checkVar004LocalVars_nilBody(t *testing.T) {
	config.Reset()

	fset := token.NewFileSet()
	// External function with no body
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

	checkVar004LocalVars(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for nil body, got %d", reportCount)
	}
}

// Test_checkVar004AssignStmt_notDefine tests non-DEFINE assignments.
func Test_checkVar004AssignStmt_notDefine(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// Regular assignment (=), not short declaration (:=)
	stmt := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
	}

	checkVar004AssignStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for = assignment, got %d", reportCount)
	}
}

// Test_checkVar004AssignStmt_nonIdent tests non-identifier LHS.
func Test_checkVar004AssignStmt_nonIdent(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// Short declaration with non-identifier LHS (e.g., arr[0])
	stmt := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.IndexExpr{
			X:     &ast.Ident{Name: "arr"},
			Index: &ast.BasicLit{Kind: token.INT, Value: "0"},
		}},
	}

	checkVar004AssignStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-ident LHS, got %d", reportCount)
	}
}

// Test_checkVar004DeclStmt_nonGenDecl tests non-GenDecl.
func Test_checkVar004DeclStmt_nonGenDecl(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// Create a DeclStmt with a FuncDecl (not GenDecl)
	stmt := &ast.DeclStmt{
		Decl: &ast.FuncDecl{Name: &ast.Ident{Name: "foo"}},
	}

	checkVar004DeclStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-GenDecl, got %d", reportCount)
	}
}

// Test_checkVar004DeclStmt_constDecl tests const declarations.
func Test_checkVar004DeclStmt_constDecl(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// GenDecl with const (not var)
	stmt := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.CONST,
			Specs: []ast.Spec{
				&ast.ValueSpec{Names: []*ast.Ident{{Name: "x"}}},
			},
		},
	}

	checkVar004DeclStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for const decl, got %d", reportCount)
	}
}

// Test_checkVar004Name_blankIdent tests blank identifier.
func Test_checkVar004Name_blankIdent(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	ident := &ast.Ident{Name: "_"}
	checkVar004Name(pass, ident, false)
	checkVar004Name(pass, ident, true)

	if reportCount != 0 {
		t.Errorf("expected 0 reports for blank ident, got %d", reportCount)
	}
}

// Test_checkVar004Name_longEnough tests names that are long enough.
func Test_checkVar004Name_longEnough(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	ident := &ast.Ident{Name: "ok"}
	checkVar004Name(pass, ident, false)
	checkVar004Name(pass, ident, true)

	if reportCount != 0 {
		t.Errorf("expected 0 reports for 'ok', got %d", reportCount)
	}
}

// Test_checkVar004Name_idiomaticNames tests idiomatic 1-char names.
func Test_checkVar004Name_idiomaticNames(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// Test all idiomatic 1-char names at function level
	for name := range idiomaticOneChar004 {
		ident := &ast.Ident{Name: name}
		checkVar004Name(pass, ident, false)
	}

	if reportCount != 0 {
		t.Errorf("expected 0 reports for idiomatic names, got %d", reportCount)
	}
}

// Test_checkVar004Node_nilNode tests nil node handling.
func Test_checkVar004Node_nilNode(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	checkVar004Node(pass, nil)

	if reportCount != 0 {
		t.Errorf("expected 0 reports for nil node, got %d", reportCount)
	}
}

// Test_checkVar004Node_otherNodeTypes tests unhandled node types.
func Test_checkVar004Node_otherNodeTypes(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// Test with a node type that isn't handled
	checkVar004Node(pass, &ast.ReturnStmt{})
	checkVar004Node(pass, &ast.IfStmt{})

	if reportCount != 0 {
		t.Errorf("expected 0 reports for unhandled nodes, got %d", reportCount)
	}
}

// Test_checkVar004PackageLevel_excluded tests file exclusion.
func Test_checkVar004PackageLevel_excluded(t *testing.T) {
	config.Reset()

	cfg := config.Get()
	cfg.Rules = map[string]*config.RuleConfig{
		ruleCodeVar004: {Exclude: []string{"test.go"}},
	}

	fset := token.NewFileSet()
	code := `package test
var a = 1
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

	checkVar004PackageLevel(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for excluded file, got %d", reportCount)
	}
	config.Reset()
}

// Test_checkVar004LocalVars_excluded tests file exclusion for local vars.
func Test_checkVar004LocalVars_excluded(t *testing.T) {
	config.Reset()

	cfg := config.Get()
	cfg.Rules = map[string]*config.RuleConfig{
		ruleCodeVar004: {Exclude: []string{"test.go"}},
	}

	fset := token.NewFileSet()
	code := `package test
func foo() {
	x := 1
	_ = x
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

	checkVar004LocalVars(pass, insp, cfg)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for excluded file, got %d", reportCount)
	}
	config.Reset()
}

// Test_checkVar004DeclStmt_nonValueSpec tests non-ValueSpec in specs.
func Test_checkVar004DeclStmt_nonValueSpec(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// GenDecl with var but non-ValueSpec (should not happen but defensive)
	stmt := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ImportSpec{Name: &ast.Ident{Name: "foo"}},
			},
		},
	}

	checkVar004DeclStmt(pass, stmt)
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-ValueSpec, got %d", reportCount)
	}
}
