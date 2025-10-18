package ktn_error_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_error "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/error"
)

// TestRule001_Unit_NilPassTypesInfo teste le comportement avec pass.TypesInfo == nil
func TestRule001_Unit_NilPassTypesInfo(t *testing.T) {
	src := `package test
func badUnwrappedError() error {
	err := errors.New("base error")
	if err != nil {
		return err
	}
	return nil
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Pass avec TypesInfo == nil
	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: nil, // Ceci force les branches alternatives sans TypesInfo
	}

	// Exécuter l'analyse
	_, err = ktn_error.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRule001_Unit_FunctionWithoutError teste une fonction qui ne retourne pas error
func TestRule001_Unit_FunctionWithoutError(t *testing.T) {
	src := `package test
func noError() string {
	return "hello"
}

func noReturnType() {
	println("test")
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object),
			Defs:  make(map[*ast.Ident]types.Object),
		},
	}

	// Exécuter l'analyse - ne devrait pas rapporter d'erreur
	_, err = ktn_error.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRule001_Unit_NonErrorReturn teste des returns qui ne sont pas des erreurs
func TestRule001_Unit_NonErrorReturn(t *testing.T) {
	src := `package test
import "errors"

func mixedReturn() (string, error) {
	msg := "hello"
	err := errors.New("error")
	if err != nil {
		return msg, err
	}
	return msg, nil
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Configuration minimale pour TypesInfo
	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object),
			Defs:  make(map[*ast.Ident]types.Object),
		},
	}

	_, err = ktn_error.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRule001_Unit_FunctionWithoutResults teste une fonction sans résultats
func TestRule001_Unit_FunctionWithoutResults(t *testing.T) {
	src := `package test
func noResults() {
	println("test")
}

type MyType struct{}

func (m *MyType) method() {
	println("method")
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object),
			Defs:  make(map[*ast.Ident]types.Object),
		},
	}

	_, err = ktn_error.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRule001_Unit_ReturnLiterals teste des returns avec des littéraux
func TestRule001_Unit_ReturnLiterals(t *testing.T) {
	src := `package test
import "errors"

func returnLiteral() error {
	return errors.New("literal error")
}

func returnFunctionCall() error {
	return doSomething()
}

func doSomething() error {
	return errors.New("error")
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object),
			Defs:  make(map[*ast.Ident]types.Object),
		},
	}

	_, err = ktn_error.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRule001_Unit_ComplexScenarios teste des scénarios complexes
func TestRule001_Unit_ComplexScenarios(t *testing.T) {
	src := `package test
import "errors"

// Fonction sans body (interface)
type ErrorHandler interface {
	HandleError() error
}

// Fonction avec multiple returns
func multipleReturns(flag bool) (string, int, error) {
	if flag {
		return "ok", 42, nil
	}
	err := errors.New("error")
	value := "fail"
	code := 0
	return value, code, err // err should be detected
}

// Fonction avec return vide
func earlyReturn() error {
	if true {
		return nil
	}
	err := errors.New("unreachable")
	return err
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object),
			Defs:  make(map[*ast.Ident]types.Object),
		},
	}

	_, err = ktn_error.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRule001_Unit_UndefinedIdentifier teste le cas d'un identifiant sans objet
func TestRule001_Unit_UndefinedIdentifier(t *testing.T) {
	src := `package test

func undefinedReturn() error {
	return undefinedVar
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// TypesInfo avec Uses vide pour simuler un identifiant non résolu
	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object), // vide - pas d'objet pour undefinedVar
			Defs:  make(map[*ast.Ident]types.Object),
		},
	}

	_, err = ktn_error.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestIsErrorVariable_NilTypesInfo teste isErrorVariable avec TypesInfo == nil
func TestIsErrorVariable_NilTypesInfo(t *testing.T) {
	src := `package test
func test() error {
	err := error(nil)
	return err
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: nil, // TypesInfo is nil
	}

	// Trouver l'identifiant err dans le return
	var errIdent *ast.Ident
	ast.Inspect(file, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			if len(ret.Results) > 0 {
				if ident, ok := ret.Results[0].(*ast.Ident); ok {
					errIdent = ident
					return false
				}
			}
		}
		return true
	})

	if errIdent == nil {
		t.Fatal("Could not find err identifier")
	}

	// Tester isErrorVariable avec TypesInfo == nil
	result := ktn_error.ExportedIsErrorVariable(pass, errIdent)
	if result {
		t.Error("Expected false when TypesInfo is nil")
	}
}

// TestIsErrorVariable_NilObject teste isErrorVariable avec obj == nil
func TestIsErrorVariable_NilObject(t *testing.T) {
	src := `package test
func test() error {
	err := error(nil)
	return err
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// TypesInfo avec Uses vide (obj sera nil)
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object), // Vide - obj sera nil
			Defs:  make(map[*ast.Ident]types.Object),
		},
	}

	// Trouver l'identifiant err dans le return
	var errIdent *ast.Ident
	ast.Inspect(file, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			if len(ret.Results) > 0 {
				if ident, ok := ret.Results[0].(*ast.Ident); ok {
					errIdent = ident
					return false
				}
			}
		}
		return true
	})

	if errIdent == nil {
		t.Fatal("Could not find err identifier")
	}

	// Tester isErrorVariable avec obj == nil
	result := ktn_error.ExportedIsErrorVariable(pass, errIdent)
	if result {
		t.Error("Expected false when obj is nil")
	}
}
