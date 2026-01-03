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

// Test_runVar005 tests the runVar005 function.
func Test_runVar005(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		ruleEnabled bool
		expectCount int
	}{
		{
			name:        "enabled with violation",
			code:        `package test; var veryLongVariableNameThatExceedsLimitOf30CharactersTotal = 1`,
			ruleEnabled: true,
			expectCount: 1,
		},
		{
			name:        "enabled without violation",
			code:        `package test; var shortName = 1`,
			ruleEnabled: true,
			expectCount: 0,
		},
		{
			name:        "disabled",
			code:        `package test; var veryLongVariableNameThatExceedsLimitOf30CharactersTotal = 1`,
			ruleEnabled: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.ruleEnabled {
				config.Reset()
			} else {
				config.Set(&config.Config{
					Rules: map[string]*config.RuleConfig{
						ruleCodeVar005: {Enabled: config.Bool(false)},
					},
				})
			}
			defer config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérifier l'erreur de parsing
			if err != nil || file == nil {
				t.Fatalf("failed to parse test code: %v", err)
			}
			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
				Report:   func(_ analysis.Diagnostic) { reportCount++ },
			}

			_, _ = runVar005(pass)

			if reportCount != tt.expectCount {
				t.Errorf("runVar005() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005PackageLevel tests the checkVar005PackageLevel function.
func Test_checkVar005PackageLevel(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectCount int
	}{
		{
			name:        "long name violation",
			code:        `package test; var veryLongVariableNameThatExceedsLimitOf30CharactersTotal = 1`,
			expectCount: 1,
		},
		{
			name:        "valid name",
			code:        `package test; var shortName = 1`,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérifier l'erreur de parsing
			if err != nil || file == nil {
				t.Fatalf("failed to parse test code: %v", err)
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

			checkVar005PackageLevel(pass, insp, cfg)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005PackageLevel() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005LocalVars tests the checkVar005LocalVars function.
func Test_checkVar005LocalVars(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectCount int
	}{
		{
			name:        "long name violation",
			code:        `package test; func f() { veryLongVariableNameThatExceedsLimitOf30 := 1; _ = veryLongVariableNameThatExceedsLimitOf30 }`,
			expectCount: 1,
		},
		{
			name:        "valid name",
			code:        `package test; func f() { shortName := 1; _ = shortName }`,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérifier l'erreur de parsing
			if err != nil || file == nil {
				t.Fatalf("failed to parse test code: %v", err)
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

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005LocalVars() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005Node tests the checkVar005Node function.
func Test_checkVar005Node(t *testing.T) {
	tests := []struct {
		name        string
		node        ast.Node
		expectCount int
	}{
		{
			name:        "nil node",
			node:        nil,
			expectCount: 0,
		},
		{
			name:        "unhandled node",
			node:        &ast.ReturnStmt{},
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			reportCount := 0
			pass := &analysis.Pass{
				Fset:   token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) { reportCount++ },
			}

			checkVar005Node(pass, tt.node)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005Node() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005AssignStmt tests the checkVar005AssignStmt function.
func Test_checkVar005AssignStmt(t *testing.T) {
	tests := []struct {
		name        string
		stmt        *ast.AssignStmt
		expectCount int
	}{
		{
			name:        "regular assign",
			stmt:        &ast.AssignStmt{Tok: token.ASSIGN, Lhs: []ast.Expr{&ast.Ident{Name: "x"}}},
			expectCount: 0,
		},
		{
			name:        "define short name",
			stmt:        &ast.AssignStmt{Tok: token.DEFINE, Lhs: []ast.Expr{&ast.Ident{Name: "x"}}},
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			reportCount := 0
			pass := &analysis.Pass{
				Fset:   token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) { reportCount++ },
			}

			checkVar005AssignStmt(pass, tt.stmt)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005AssignStmt() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005RangeStmt tests the checkVar005RangeStmt function.
func Test_checkVar005RangeStmt(t *testing.T) {
	tests := []struct {
		name        string
		stmt        *ast.RangeStmt
		expectCount int
	}{
		{
			name:        "regular assign",
			stmt:        &ast.RangeStmt{Tok: token.ASSIGN, Key: &ast.Ident{Name: "i"}},
			expectCount: 0,
		},
		{
			name:        "define short name",
			stmt:        &ast.RangeStmt{Tok: token.DEFINE, Key: &ast.Ident{Name: "i"}},
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			reportCount := 0
			pass := &analysis.Pass{
				Fset:   token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) { reportCount++ },
			}

			checkVar005RangeStmt(pass, tt.stmt)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005RangeStmt() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005DeclStmt tests the checkVar005DeclStmt function.
func Test_checkVar005DeclStmt(t *testing.T) {
	tests := []struct {
		name        string
		stmt        *ast.DeclStmt
		expectCount int
	}{
		{
			name:        "func decl",
			stmt:        &ast.DeclStmt{Decl: &ast.FuncDecl{Name: &ast.Ident{Name: "f"}}},
			expectCount: 0,
		},
		{
			name:        "const decl",
			stmt:        &ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.CONST}},
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			reportCount := 0
			pass := &analysis.Pass{
				Fset:   token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) { reportCount++ },
			}

			checkVar005DeclStmt(pass, tt.stmt)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005DeclStmt() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005Spec tests the checkVar005Spec function.
func Test_checkVar005Spec(t *testing.T) {
	tests := []struct {
		name        string
		spec        *ast.ValueSpec
		expectCount int
	}{
		{
			name:        "short name",
			spec:        &ast.ValueSpec{Names: []*ast.Ident{{Name: "ok"}}},
			expectCount: 0,
		},
		{
			name:        "long name",
			spec:        &ast.ValueSpec{Names: []*ast.Ident{{Name: "veryLongVariableNameExceeds30Chars"}}},
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			reportCount := 0
			pass := &analysis.Pass{
				Fset:   token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) { reportCount++ },
			}

			checkVar005Spec(pass, tt.spec)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005Spec() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar005Name tests the checkVar005Name function.
func Test_checkVar005Name(t *testing.T) {
	tests := []struct {
		name        string
		ident       *ast.Ident
		expectCount int
	}{
		{
			name:        "blank identifier",
			ident:       &ast.Ident{Name: "_"},
			expectCount: 0,
		},
		{
			name:        "short name",
			ident:       &ast.Ident{Name: "ok"},
			expectCount: 0,
		},
		{
			name:        "long name",
			ident:       &ast.Ident{Name: "veryLongVariableNameExceeds30Chars"},
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()
			reportCount := 0
			pass := &analysis.Pass{
				Fset:   token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) { reportCount++ },
			}

			checkVar005Name(pass, tt.ident)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar005Name() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

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
		tt := tt
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
	file, parseErr := parser.ParseFile(fset, "test.go", `package test; var a = 1`, 0)
	// Vérifier l'erreur de parsing
	if parseErr != nil || file == nil {
		t.Fatalf("failed to parse test code: %v", parseErr)
	}
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
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Vérifier l'erreur de parsing
	if err != nil || file == nil {
		t.Fatalf("failed to parse test code: %v", err)
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
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Vérifier l'erreur de parsing
	if err != nil || file == nil {
		t.Fatalf("failed to parse test code: %v", err)
	}
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
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Vérifier l'erreur de parsing
	if err != nil || file == nil {
		t.Fatalf("failed to parse test code: %v", err)
	}
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
