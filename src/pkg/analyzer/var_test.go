package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestVarAnalyzerKTNVAR001 teste VarAnalyzer KTN VAR 001.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR001(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "ungrouped var should trigger KTN-VAR-001",
			code: `package test
var MaxConnections int = 100
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-001",
		},
		{
			name: "grouped var should not trigger KTN-VAR-001",
			code: `package test
// Connection limits
// These variables define connection limits
var (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)

func updateConfig() {
	MaxConnections = 200
}
`,
			wantDiag: false,
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR002 teste VarAnalyzer KTN VAR 002.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR002(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "group without group comment should trigger KTN-VAR-002",
			code: `package test
var (
	MaxConnections int = 100
)
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-002",
		},
		{
			name: "group with group comment should not trigger KTN-VAR-002",
			code: `package test
// Connection limits
// These variables define connection limits
var (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)

func updateConfig() {
	MaxConnections = 200
}
`,
			wantDiag: false,
			wantMsg:  "KTN-VAR-002",
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR003 teste VarAnalyzer KTN VAR 003.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR003(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "var without individual comment should trigger KTN-VAR-003",
			code: `package test
// Connection limits
// These variables define connection limits
var (
	MaxConnections int = 100
)
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-003",
		},
		{
			name: "var with individual comment should not trigger KTN-VAR-003",
			code: `package test
// Connection limits
// These variables define connection limits
var (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)

func updateConfig() {
	MaxConnections = 200
}
`,
			wantDiag: false,
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR004 teste VarAnalyzer KTN VAR 004.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR004(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "var without explicit type should trigger KTN-VAR-004",
			code: `package test
// Connection limits
// These variables define connection limits
var (
	// MaxConnections defines the maximum number of connections
	MaxConnections = 100
)
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-004",
		},
		{
			name: "var with explicit type should not trigger KTN-VAR-004",
			code: `package test
// Connection limits
// These variables define connection limits
var (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)

func updateConfig() {
	MaxConnections = 200
}
`,
			wantDiag: false,
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR005 teste VarAnalyzer KTN VAR 005.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR005(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "mathematical constant Pi should trigger KTN-VAR-005",
			code: `package test
// Mathematical value
// Pi is the mathematical constant
var (
	// Pi represents the value of pi
	Pi float64 = 3.14159265358979323846
)
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-005",
		},
		{
			name: "regular variable MaxConnections should not trigger KTN-VAR-005",
			code: `package test
// Connection limits
// These variables define connection limits
var (
	// MaxConnections defines the maximum number of connections
	MaxConnections int = 100
)

func updateConfig() {
	MaxConnections = 200
}
`,
			wantDiag: false,
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR006 teste VarAnalyzer KTN VAR 006.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR006(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "multiple variables on one line should trigger KTN-VAR-006",
			code: `package test
// Network settings
// These variables configure network
var (
	// HostName and Port are network settings
	HostName, Port = "localhost", 8080
)
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-006",
		},
		{
			name: "one variable per line should not trigger KTN-VAR-006",
			code: `package test
// Network settings
// These variables configure network
var (
	// HostName is the hostname
	HostName string = "localhost"
	// Port is the port number
	Port int = 8080
)

func updateNetwork() {
	HostName = "example.com"
	Port = 9090
}
`,
			wantDiag: false,
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR007 teste VarAnalyzer KTN VAR 007.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR007(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "channel without buffer size should trigger KTN-VAR-007",
			code: `package test
// Message channels
// These variables handle messages
var (
	// MessageQueue is the message channel
	MessageQueue chan string = make(chan string)
)
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-007",
		},
		{
			name: "channel with unbuffered comment should not trigger KTN-VAR-007",
			code: `package test
// Message channels
// These variables handle messages
var (
	// DoneSignal is unbuffered intentionally
	DoneSignal chan bool = make(chan bool)
)
`,
			wantDiag: false,
		},
		{
			name: "channel with buffer size should not trigger KTN-VAR-007",
			code: `package test
// Message channels
// These variables handle messages
var (
	// MessageQueue has buffer of 100
	MessageQueue chan string = make(chan string, 100)
)
`,
			wantDiag: false,
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR008 teste VarAnalyzer KTN VAR 008.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR008(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantDiag bool
		wantMsg  string
	}{
		{
			name: "variable with underscore should trigger KTN-VAR-008",
			code: `package test
// HTTP codes
// These variables contain HTTP codes
var (
	// HTTP_OK represents code 200
	HTTP_OK int = 200
)
`,
			wantDiag: true,
			wantMsg:  "KTN-VAR-008",
		},
		{
			name: "variable with MixedCaps should not trigger KTN-VAR-008",
			code: `package test
// HTTP codes
// These variables contain HTTP codes
var (
	// HTTPOK represents code 200
	HTTPOK int = 200
)

func updateHTTPCode() {
	HTTPOK = 201
}
`,
			wantDiag: false,
		},
	}

	runVarTests(t, tests)
}

// TestVarAnalyzerKTNVAR009 teste VarAnalyzer KTN VAR 009.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerKTNVAR009(t *testing.T) {
	runVarTest(t, "ALL_CAPS", `package test
// Buffer configuration
// These variables configure buffers
var (
	// MAX_BUFFER_SIZE is the max buffer size
	MAX_BUFFER_SIZE int = 1024
)
`, true, "KTN-VAR-009")

	runVarTest(t, "initialism", `package test
// HTTP codes
// These variables contain HTTP codes
var (
	// HTTPOK represents code 200
	HTTPOK int = 200
)

func updateHTTPCode() {
	HTTPOK = 201
}
`, false, "")

	runVarTest(t, "MixedCaps", `package test
// Buffer configuration
// These variables configure buffers
var (
	// MaxBufferSize is the max buffer size
	MaxBufferSize int = 1024
)

func updateBuffer() {
	MaxBufferSize = 2048
}
`, false, "")
}

// TestVarAnalyzerSpecialValues teste les valeurs spéciales.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerSpecialValues(t *testing.T) {
	runVarTest(t, "underscore", `package test
// Utilities
// Utility variables
var (
	_ int = 123
)
`, false, "")

	runVarTest(t, "composite literal", `package test
// Config values
// Configuration
var (
	// Config is the configuration
	Config []string = []string{"a", "b"}
)
`, false, "")

	runVarTest(t, "function call", `package test
// Init values
// Initialization
var (
	// Value from function
	Value int = getValue()
)
`, false, "")

	runVarTest(t, "constant-like unsupported type", `package test
// Constants
// Math constants
var (
	// Pi is mathematical constant
	Pi []int = []int{3}
)
`, false, "")
}

// TestVarAnalyzerNamingEdgeCases teste les cas limites du nommage.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerNamingEdgeCases(t *testing.T) {
	runVarTest(t, "short MixedCaps", `package test
// Test values
// Testing
var (
	// Ax is a test var
	Ax int = 1
)

func updateAx() {
	Ax = 2
}
`, false, "")

	runVarTest(t, "initialism start", `package test
// HTTP variables
// HTTP related
var (
	// HTTPServer is the server
	HTTPServer string = "localhost"
)

func updateHTTPServer() {
	HTTPServer = "example.com"
}
`, false, "")

	runVarTest(t, "mixed initialism", `package test
// API variables
// API configuration
var (
	// APIEndpoint is the endpoint
	APIEndpoint string = "https://api.example.com"
)

func updateAPIEndpoint() {
	APIEndpoint = "https://api.new.com"
}
`, false, "")
}

// TestVarAnalyzerChannelComments teste les commentaires de canaux.
//
// Params:
//   - t: instance de test
func TestVarAnalyzerChannelComments(t *testing.T) {
	runVarTest(t, "line comment", `package test
// Channels
// Channel variables
var (
	Ch chan int = make(chan int, 10) // buffered channel
)
`, false, "")
}

// runVarTest exécute un test pour le VarAnalyzer.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - wantDiag: true si on attend un diagnostic
//   - wantMsg: message attendu dans le diagnostic
func runVarTest(t *testing.T, name, code string, wantDiag bool, wantMsg string) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: analyzer.VarAnalyzer,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(diag analysis.Diagnostic) {
			diagnostics = append(diagnostics, diag)
		},
	}

	_, err = analyzer.VarAnalyzer.Run(pass)
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

// runVarTests est une fonction helper pour exécuter les tests du VarAnalyzer.
//
// Params:
//   - t: instance de test
//   - tests: slice de cas de test à exécuter
func runVarTests(t *testing.T, tests []struct {
	name     string
	code     string
	wantDiag bool
	wantMsg  string
}) {
	for _, tt := range tests {
		runVarTest(t, tt.name, tt.code, tt.wantDiag, tt.wantMsg)
	}
}
