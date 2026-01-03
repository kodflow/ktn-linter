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

// Test_runVar004 tests the runVar004 function.
func Test_runVar004(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		ruleEnabled bool
		expectCount int
	}{
		{
			name:        "enabled with violation",
			code:        `package test; var a = 1`,
			ruleEnabled: true,
			expectCount: 1,
		},
		{
			name:        "enabled without violation",
			code:        `package test; var ok = 1`,
			ruleEnabled: true,
			expectCount: 0,
		},
		{
			name:        "disabled",
			code:        `package test; var a = 1`,
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
						ruleCodeVar004: {Enabled: config.Bool(false)},
					},
				})
			}
			defer config.Reset()

			fset := token.NewFileSet()
			file, _ := parser.ParseFile(fset, "test.go", tt.code, 0)
			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
				Report:   func(_ analysis.Diagnostic) { reportCount++ },
			}

			_, _ = runVar004(pass)

			if reportCount != tt.expectCount {
				t.Errorf("runVar004() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar004PackageLevel tests the checkVar004PackageLevel function.
func Test_checkVar004PackageLevel(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectCount int
	}{
		{
			name:        "short name violation",
			code:        `package test; var a = 1`,
			expectCount: 1,
		},
		{
			name:        "valid name",
			code:        `package test; var ok = 1`,
			expectCount: 0,
		},
		{
			name:        "const not var",
			code:        `package test; const a = 1`,
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

			checkVar004PackageLevel(pass, insp, cfg)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar004PackageLevel() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar004LocalVars tests the checkVar004LocalVars function.
func Test_checkVar004LocalVars(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectCount int
	}{
		{
			name:        "short name violation",
			code:        `package test; func f() { x := 1; _ = x }`,
			expectCount: 1,
		},
		{
			name:        "idiomatic name",
			code:        `package test; func f() { i := 1; _ = i }`,
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
			if err != nil {
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

			checkVar004LocalVars(pass, insp, cfg)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar004LocalVars() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}

	// Test nil body avec AST manuel (car `func f()` ne parse pas)
	t.Run("nil body", func(t *testing.T) {
		config.Reset()
		fset := token.NewFileSet()

		// Construire AST manuellement pour fonction sans body
		file := &ast.File{
			Name: &ast.Ident{Name: "test"},
			Decls: []ast.Decl{
				&ast.FuncDecl{
					Name: &ast.Ident{Name: "f"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
					},
					Body: nil,
				},
			},
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
			t.Errorf("checkVar004LocalVars() reported %d issues, expected 0", reportCount)
		}
	})
}

// Test_checkVar004Node tests the checkVar004Node function.
func Test_checkVar004Node(t *testing.T) {
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
			name:        "return statement",
			node:        &ast.ReturnStmt{},
			expectCount: 0,
		},
		{
			name:        "assign with short name",
			node:        &ast.AssignStmt{Tok: token.DEFINE, Lhs: []ast.Expr{&ast.Ident{Name: "x"}}},
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

			checkVar004Node(pass, tt.node)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar004Node() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar004AssignStmt tests the checkVar004AssignStmt function.
func Test_checkVar004AssignStmt(t *testing.T) {
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
			name:        "define with short name",
			stmt:        &ast.AssignStmt{Tok: token.DEFINE, Lhs: []ast.Expr{&ast.Ident{Name: "x"}}},
			expectCount: 1,
		},
		{
			name:        "define with idiomatic name",
			stmt:        &ast.AssignStmt{Tok: token.DEFINE, Lhs: []ast.Expr{&ast.Ident{Name: "i"}}},
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

			checkVar004AssignStmt(pass, tt.stmt)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar004AssignStmt() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar004DeclStmt tests the checkVar004DeclStmt function.
func Test_checkVar004DeclStmt(t *testing.T) {
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
		{
			name:        "var decl short name",
			stmt:        &ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{&ast.ValueSpec{Names: []*ast.Ident{{Name: "x"}}}}}},
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

			checkVar004DeclStmt(pass, tt.stmt)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar004DeclStmt() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar004Spec tests the checkVar004Spec function.
func Test_checkVar004Spec(t *testing.T) {
	tests := []struct {
		name           string
		spec           *ast.ValueSpec
		isPackageLevel bool
		expectCount    int
	}{
		{
			name:           "short name package level",
			spec:           &ast.ValueSpec{Names: []*ast.Ident{{Name: "a"}}},
			isPackageLevel: true,
			expectCount:    1,
		},
		{
			name:           "short name func level",
			spec:           &ast.ValueSpec{Names: []*ast.Ident{{Name: "x"}}},
			isPackageLevel: false,
			expectCount:    1,
		},
		{
			name:           "valid name",
			spec:           &ast.ValueSpec{Names: []*ast.Ident{{Name: "ok"}}},
			isPackageLevel: false,
			expectCount:    0,
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

			checkVar004Spec(pass, tt.spec, tt.isPackageLevel)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar004Spec() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_checkVar004Name tests the checkVar004Name function.
func Test_checkVar004Name(t *testing.T) {
	tests := []struct {
		name           string
		ident          *ast.Ident
		isPackageLevel bool
		expectCount    int
	}{
		{
			name:           "blank identifier",
			ident:          &ast.Ident{Name: "_"},
			isPackageLevel: false,
			expectCount:    0,
		},
		{
			name:           "long enough name",
			ident:          &ast.Ident{Name: "ok"},
			isPackageLevel: false,
			expectCount:    0,
		},
		{
			name:           "short package level",
			ident:          &ast.Ident{Name: "a"},
			isPackageLevel: true,
			expectCount:    1,
		},
		{
			name:           "idiomatic func level",
			ident:          &ast.Ident{Name: "i"},
			isPackageLevel: false,
			expectCount:    0,
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

			checkVar004Name(pass, tt.ident, tt.isPackageLevel)

			if reportCount != tt.expectCount {
				t.Errorf("checkVar004Name() reported %d issues, expected %d", reportCount, tt.expectCount)
			}
		})
	}
}

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

// Test_checkVar004PackageLevel_nonValueSpec tests non-ValueSpec in package-level.
func Test_checkVar004PackageLevel_nonValueSpec(t *testing.T) {
	config.Reset()

	fset := token.NewFileSet()
	// Create a fake file with a var GenDecl containing non-ValueSpec
	file := &ast.File{
		Name: &ast.Ident{Name: "test"},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					// Import spec in var decl is invalid but tests defensive check
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
				},
			},
		},
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

	checkVar004PackageLevel(pass, insp, cfg)
	// Should not crash and should not report for non-ValueSpec
	if reportCount != 0 {
		t.Errorf("expected 0 reports for non-ValueSpec at package level, got %d", reportCount)
	}
}

// Test_checkVar004Name_packageLevelShort tests short name at package-level.
func Test_checkVar004Name_packageLevelShort(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// Single character name at package level (should report)
	ident := &ast.Ident{Name: "a"}
	checkVar004Name(pass, ident, true)

	if reportCount != 1 {
		t.Errorf("expected 1 report for short package-level name 'a', got %d", reportCount)
	}
}

// Test_checkVar004Name_functionLevelShortNonIdiomatic tests non-idiomatic short name.
func Test_checkVar004Name_functionLevelShortNonIdiomatic(t *testing.T) {
	config.Reset()

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// Single character 'x' is not in idiomaticOneChar004
	ident := &ast.Ident{Name: "x"}
	checkVar004Name(pass, ident, false)

	if reportCount != 1 {
		t.Errorf("expected 1 report for non-idiomatic 'x', got %d", reportCount)
	}
}

// Test_checkVar004Name_idiomaticShortBranch tests idiomaticShort004 branch.
func Test_checkVar004Name_idiomaticShortBranch(t *testing.T) {
	config.Reset()

	// Temporarily add a 1-char name to idiomaticShort004 to test the branch
	idiomaticShort004["q"] = true
	defer delete(idiomaticShort004, "q")

	reportCount := 0
	pass := &analysis.Pass{
		Fset:   token.NewFileSet(),
		Report: func(_ analysis.Diagnostic) { reportCount++ },
	}

	// "q" is now in idiomaticShort004 and len("q") == 1 < 2
	ident := &ast.Ident{Name: "q"}
	checkVar004Name(pass, ident, false)

	// Should not report for idiomatic short name
	if reportCount != 0 {
		t.Errorf("expected 0 reports for idiomatic 'q', got %d", reportCount)
	}
}
