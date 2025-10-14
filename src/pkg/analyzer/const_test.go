package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestConstAnalyzerKTNCONST001 teste ConstAnalyzer KTN CONST 001.
//
// Params:
//   - t: instance de test
func TestConstAnalyzerKTNCONST001(t *testing.T) {
	runConstTest(t, "ungrouped", `package test
const MaxConnections int = 100
`, true, "KTN-CONST-001")

	runConstTest(t, "grouped", `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`, false, "")
}

// TestConstAnalyzerKTNCONST002 teste ConstAnalyzer KTN CONST 002.
//
// Params:
//   - t: instance de test
func TestConstAnalyzerKTNCONST002(t *testing.T) {
	runConstTest(t, "no group comment", `package test
const (
	MaxConnections int = 100
)
`, true, "KTN-CONST-002")

	runConstTest(t, "with group comment", `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`, false, "")
}

// TestConstAnalyzerKTNCONST003 teste ConstAnalyzer KTN CONST 003.
//
// Params:
//   - t: instance de test
func TestConstAnalyzerKTNCONST003(t *testing.T) {
	runConstTest(t, "no individual comment", `package test
// Connection limits
// These constants define connection limits
const (
	MaxConnections int = 100
)
`, true, "KTN-CONST-003")

	runConstTest(t, "with individual comment", `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`, false, "")
}

// TestConstAnalyzerKTNCONST004 teste ConstAnalyzer KTN CONST 004.
//
// Params:
//   - t: instance de test
func TestConstAnalyzerKTNCONST004(t *testing.T) {
	runConstTest(t, "no explicit type", `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections = 100
)
`, true, "KTN-CONST-004")

	runConstTest(t, "with explicit type", `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`, false, "")
}

// runConstTest exécute un test pour le ConstAnalyzer.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - wantDiag: true si on attend un diagnostic
//   - wantMsg: message attendu dans le diagnostic
func runConstTest(t *testing.T, name, code string, wantDiag bool, wantMsg string) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: analyzer.ConstAnalyzer,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(diag analysis.Diagnostic) {
			diagnostics = append(diagnostics, diag)
		},
	}

	_, err = analyzer.ConstAnalyzer.Run(pass)
	if err != nil {
		t.Fatalf("analyzer failed: %v", err)
	}

	hasExpectedDiag := false
	for _, d := range diagnostics {
		if wantMsg != "" && contains(d.Message, wantMsg) {
			hasExpectedDiag = true
			break
		}
	}

	if wantDiag && !hasExpectedDiag {
		t.Errorf("expected diagnostic %q but got none", wantMsg)
	}
	if !wantDiag && len(diagnostics) > 0 {
		t.Errorf("expected no diagnostic but got: %v", diagnostics)
	}
}

// TestConstAnalyzerNumericTypes teste les types numériques.
//
// Params:
//   - t: instance de test
func TestConstAnalyzerNumericTypes(t *testing.T) {
	runConstTest(t, "uint8", `package test
// Bytes
// Byte constants
const (
	// MaxByte is max byte value
	MaxByte uint8 = 255
)`, false, "")

	runConstTest(t, "float32", `package test
// Floats
// Float constants
const (
	// DefaultRatio is the ratio
	DefaultRatio float32 = 1.5
)`, false, "")

	runConstTest(t, "complex64", `package test
// Complex numbers
// Complex constants
const (
	// ImaginaryUnit is i
	ImaginaryUnit complex64 = 1i
)`, false, "")

	runConstTest(t, "rune", `package test
// Characters
// Character constants
const (
	// NewLine is newline char
	NewLine rune = '\n'
)`, false, "")

	runConstTest(t, "byte", `package test
// Bytes
// Byte constants
const (
	// NullByte is null
	NullByte byte = 0
)`, false, "")
}

// TestConstAnalyzerSpecialCases teste les cas spéciaux.
//
// Params:
//   - t: instance de test
func TestConstAnalyzerSpecialCases(t *testing.T) {
	runConstTest(t, "multiple consts", `package test
// Network
// Network constants
const (
	// Host and Port
	Host, Port string = "localhost", "8080"
)`, false, "")

	runConstTest(t, "underscore const", `package test
// Ignored
// Ignored constants
const (
	_ int = 999
)`, false, "")

	runConstTest(t, "line comment only", `package test
// Config
// Configuration
const (
	MaxSize int = 1024 // Maximum size
)`, false, "")
}

// TestPlugin tests the plugin interface functions
//
// Params:
//   - t: instance de test
func TestPlugin(t *testing.T) {
	// Test New function
	plugin, err := analyzer.New(nil)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if plugin == nil {
		t.Fatal("New() returned nil plugin")
	}

	// Test BuildAnalyzers
	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers() returned error: %v", err)
	}
	if len(analyzers) == 0 {
		t.Fatal("BuildAnalyzers() returned empty list")
	}

	// Verify ConstAnalyzer is included
	found := false
	for _, a := range analyzers {
		if a.Name == "ktnconst" {
			found = true
			break
		}
	}
	if !found {
		t.Error("BuildAnalyzers() did not include ConstAnalyzer")
	}

	// Test GetLoadMode
	loadMode := plugin.GetLoadMode()
	if loadMode == "" {
		t.Error("GetLoadMode() returned empty string")
	}
}

// contains vérifie si une chaîne contient une sous-chaîne.
//
// Params:
//   - s: la chaîne à analyser
//   - substr: la sous-chaîne recherchée
//
// Returns:
//   - bool: true si substr est trouvé dans s
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		 stringContains(s, substr)))
}

// stringContains est une fonction helper pour rechercher une sous-chaîne.
//
// Params:
//   - s: la chaîne à analyser
//   - substr: la sous-chaîne recherchée
//
// Returns:
//   - bool: true si substr est trouvé dans s
func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
