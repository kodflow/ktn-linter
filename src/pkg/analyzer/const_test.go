package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestConstAnalyzer_KTN_CONST_001 teste ConstAnalyzer KTN CONST 001.
//
// Params:
//   - t: instance de test
func TestConstAnalyzer_KTN_CONST_001(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "ungrouped const should trigger KTN-CONST-001",
			code: `package test
const MaxConnections int = 100
`,
			wantDiag: true,
			wantMsg:  "KTN-CONST-001",
		},
		{
			name: "grouped const should not trigger KTN-CONST-001",
			code: `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`,
			wantDiag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
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

			if tt.wantDiag && len(diagnostics) == 0 {
				t.Errorf("expected diagnostic but got none")
			}
			if !tt.wantDiag && len(diagnostics) > 0 {
				t.Errorf("expected no diagnostic but got %d: %v", len(diagnostics), diagnostics)
			}
			if tt.wantDiag && len(diagnostics) > 0 {
				found := false
				for _, d := range diagnostics {
					if contains(d.Message, tt.wantMsg) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected message containing %q but got: %v", tt.wantMsg, diagnostics)
				}
			}
		})
	}
}

// TestConstAnalyzer_KTN_CONST_002 teste ConstAnalyzer KTN CONST 002.
//
// Params:
//   - t: instance de test
func TestConstAnalyzer_KTN_CONST_002(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "group without group comment should trigger KTN-CONST-002",
			code: `package test
const (
	MaxConnections int = 100
)
`,
			wantDiag: true,
			wantMsg:  "KTN-CONST-002",
		},
		{
			name: "group with group comment should not trigger KTN-CONST-002",
			code: `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`,
			wantDiag: false,
			wantMsg:  "KTN-CONST-002",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
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
				if contains(d.Message, tt.wantMsg) {
					hasExpectedDiag = true
					break
				}
			}

			if tt.wantDiag && !hasExpectedDiag {
				t.Errorf("expected diagnostic %q but got none", tt.wantMsg)
			}
			if !tt.wantDiag && hasExpectedDiag {
				t.Errorf("expected no diagnostic %q but got one", tt.wantMsg)
			}
		})
	}
}

// TestConstAnalyzer_KTN_CONST_003 teste ConstAnalyzer KTN CONST 003.
//
// Params:
//   - t: instance de test
func TestConstAnalyzer_KTN_CONST_003(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "const without individual comment should trigger KTN-CONST-003",
			code: `package test
// Connection limits
// These constants define connection limits
const (
	MaxConnections int = 100
)
`,
			wantDiag: true,
			wantMsg:  "KTN-CONST-003",
		},
		{
			name: "const with individual comment should not trigger KTN-CONST-003",
			code: `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`,
			wantDiag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
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
				if contains(d.Message, tt.wantMsg) {
					hasExpectedDiag = true
					break
				}
			}

			if tt.wantDiag && !hasExpectedDiag {
				t.Errorf("expected diagnostic %q but got none", tt.wantMsg)
			}
			if !tt.wantDiag && hasExpectedDiag {
				t.Errorf("expected no diagnostic %q but got one", tt.wantMsg)
			}
		})
	}
}

// TestConstAnalyzer_KTN_CONST_004 teste ConstAnalyzer KTN CONST 004.
//
// Params:
//   - t: instance de test
func TestConstAnalyzer_KTN_CONST_004(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "const without explicit type should trigger KTN-CONST-004",
			code: `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections = 100
)
`,
			wantDiag: true,
			wantMsg:  "KTN-CONST-004",
		},
		{
			name: "const with explicit type should not trigger KTN-CONST-004",
			code: `package test
// Connection limits
// These constants define connection limits
const (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)
`,
			wantDiag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
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
				if contains(d.Message, tt.wantMsg) {
					hasExpectedDiag = true
					break
				}
			}

			if tt.wantDiag && !hasExpectedDiag {
				t.Errorf("expected diagnostic %q but got none", tt.wantMsg)
			}
			if !tt.wantDiag && hasExpectedDiag {
				t.Errorf("expected no diagnostic %q but got one", tt.wantMsg)
			}
		})
	}
}

// TestConstAnalyzer_EdgeCases tests edge cases for better coverage
//
// Params:
//   - t: instance de test
func TestConstAnalyzer_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "multiple consts on one line",
			code: `package test
// Network
// Network constants
const (
	// Host and Port
	Host, Port string = "localhost", "8080"
)
`,
			wantDiag: false,
		},
		{
			name: "const with uint8 type",
			code: `package test
// Bytes
// Byte constants
const (
	// MaxByte is max byte value
	MaxByte uint8 = 255
)
`,
			wantDiag: false,
		},
		{
			name: "const with float32 type",
			code: `package test
// Floats
// Float constants
const (
	// DefaultRatio is the ratio
	DefaultRatio float32 = 1.5
)
`,
			wantDiag: false,
		},
		{
			name: "const with complex64 type",
			code: `package test
// Complex numbers
// Complex constants
const (
	// ImaginaryUnit is i
	ImaginaryUnit complex64 = 1i
)
`,
			wantDiag: false,
		},
		{
			name: "const with rune type",
			code: `package test
// Characters
// Character constants
const (
	// NewLine is newline char
	NewLine rune = '\n'
)
`,
			wantDiag: false,
		},
		{
			name: "const with byte type",
			code: `package test
// Bytes
// Byte constants
const (
	// NullByte is null
	NullByte byte = 0
)
`,
			wantDiag: false,
		},
		{
			name: "underscore const should be skipped",
			code: `package test
// Ignored
// Ignored constants
const (
	_ int = 999
)
`,
			wantDiag: false,
		},
		{
			name: "const with line comment only",
			code: `package test
// Config
// Configuration
const (
	MaxSize int = 1024 // Maximum size
)
`,
			wantDiag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
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
				if tt.wantMsg != "" && contains(d.Message, tt.wantMsg) {
					hasExpectedDiag = true
					break
				}
			}

			if tt.wantDiag && !hasExpectedDiag {
				t.Errorf("expected diagnostic %q but got none. Got: %v", tt.wantMsg, diagnostics)
			}
			if !tt.wantDiag && hasExpectedDiag {
				t.Errorf("expected no diagnostic %q but got one. Got: %v", tt.wantMsg, diagnostics)
			}
		})
	}
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
