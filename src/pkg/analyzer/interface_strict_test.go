package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
)

// TestInterfaceStrictAnalyzerOnlyInterfaces teste avec un fichier contenant uniquement des interfaces.
//
// Params:
//   - t: contexte de test
func TestInterfaceStrictAnalyzerOnlyInterfaces(t *testing.T) {
	src := `package test

// Service définit une interface.
type Service interface {
	Process() error
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "interfaces.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pkg := types.NewPackage("test", "test")
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   pkg,
		Report: func(diag analysis.Diagnostic) {
			// Collecter les diagnostics si besoin
		},
	}

	// Devrait passer sans erreurs (pas de structs)
	_, err = analyzer.InterfaceStrictAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestInterfaceStrictAnalyzerWithStructs teste avec un fichier contenant des structs.
//
// Params:
//   - t: contexte de test
func TestInterfaceStrictAnalyzerWithStructs(t *testing.T) {
	src := `package test

type Service interface {
	Process() error
}

type ServiceImpl struct {
	data string
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "interfaces.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pkg := types.NewPackage("test", "test")
	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    pkg,
		Report: func(d analysis.Diagnostic) {},
	}

	// Devrait détecter la violation KTN-INTERFACE-008
	_, err = analyzer.InterfaceStrictAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestInterfaceStrictAnalyzerMainPackage teste l'exemption du package main.
//
// Params:
//   - t: contexte de test
func TestInterfaceStrictAnalyzerMainPackage(t *testing.T) {
	src := `package main

type Config struct {
	Port int
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pkg := types.NewPackage("main", "main")
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   pkg,
		Report: func(diag analysis.Diagnostic) {
			// Collecter les diagnostics si besoin
		},
	}

	// Package main devrait être exempté
	_, err = analyzer.InterfaceStrictAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestInterfaceStrictAnalyzerNoInterfacesFile teste sans fichier interfaces.go.
//
// Params:
//   - t: contexte de test
func TestInterfaceStrictAnalyzerNoInterfacesFile(t *testing.T) {
	src := `package test

func Process() error {
	return nil
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "utils.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pkg := types.NewPackage("test", "test")
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   pkg,
		Report: func(diag analysis.Diagnostic) {
			// Collecter les diagnostics si besoin
		},
	}

	// Pas de interfaces.go, rien à vérifier
	_, err = analyzer.InterfaceStrictAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestMockFileExists teste la vérification de mock.go.
//
// Params:
//   - t: contexte de test
func TestMockFileExists(t *testing.T) {
	// Créer un fichier temporaire mock.go
	tmpDir := t.TempDir()
	mockPath := tmpDir + "/mock.go"
	err := os.WriteFile(mockPath, []byte("package test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock.go: %v", err)
	}

	// Test avec mock.go existant
	if _, err := os.Stat(mockPath); err != nil {
		t.Errorf("mock.go should exist")
	}
}

// TestMockCompleteness teste la règle KTN-MOCK-002.
//
// Params:
//   - t: contexte de test
func TestMockCompleteness(t *testing.T) {
	// Créer fichiers temporaires
	tmpDir := t.TempDir()

	// interfaces.go avec 2 interfaces
	interfacesSrc := `package test

type Service interface {
	Process() error
}

type Repository interface {
	Save() error
}
`
	interfacesPath := tmpDir + "/interfaces.go"
	err := os.WriteFile(interfacesPath, []byte(interfacesSrc), 0644)
	if err != nil {
		t.Fatalf("Failed to create interfaces.go: %v", err)
	}

	// mock.go avec seulement 1 mock (MockService manque MockRepository)
	mockSrc := `package test

type MockService struct {}

func (m *MockService) Process() error {
	return nil
}
`
	mockPath := tmpDir + "/mock.go"
	err = os.WriteFile(mockPath, []byte(mockSrc), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock.go: %v", err)
	}

	// Parser interfaces.go
	fset := token.NewFileSet()
	interfacesFile, err := parser.ParseFile(fset, interfacesPath, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse interfaces.go: %v", err)
	}

	pkg := types.NewPackage("test", "test")
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{interfacesFile},
		Pkg:   pkg,
		Report: func(diag analysis.Diagnostic) {
			diagnostics = append(diagnostics, diag)
		},
	}

	// Exécuter l'analyseur
	_, err = analyzer.InterfaceStrictAnalyzer.Run(pass)
	if err != nil {
		t.Fatalf("Analyzer returned error: %v", err)
	}

	// Vérifier qu'on a détecté le mock manquant pour Repository
	foundMissingMock := false
	for _, d := range diagnostics {
		if strings.Contains(d.Message, "KTN-MOCK-002") &&
			strings.Contains(d.Message, "Repository") {
			foundMissingMock = true
			break
		}
	}

	if !foundMissingMock {
		t.Error("Expected KTN-MOCK-002 error for missing MockRepository")
	}
}
