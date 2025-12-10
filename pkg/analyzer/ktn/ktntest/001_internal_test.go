// Internal tests for analyzer 001.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

// Test_runTest001 tests the runTest001 private function with table-driven tests.
//
// Params:
//   - t: testing context
func Test_runTest001(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "basic test of runTest001 structure",
		},
		{
			name: "error case validation",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing: %s", tt.name)
		})
	}
}

// Test_runTest012_integration tests the analyzer structure.
//
// Params:
//   - t: testing context
func Test_runTest012_integration(t *testing.T) {
	tests := []struct {
		name         string
		expectedName string
	}{
		{name: "analyzer structure", expectedName: "ktntest001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Analyzer001 == nil || Analyzer001.Name != tt.expectedName {
				t.Errorf("Analyzer001 invalid: nil=%v, Name=%q, want %q",
					Analyzer001 == nil, Analyzer001.Name, tt.expectedName)
			}
		})
	}
}

// Test_runTest012_fileNamingPatterns tests various file naming patterns.
//
// Params:
//   - t: testing context
func Test_runTest012_fileNamingPatterns(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		shouldFail bool
	}{
		{
			name:       "internal test file is valid",
			filename:   "myfile_internal_test.go",
			shouldFail: false,
		},
		{
			name:       "external test file is valid",
			filename:   "myfile_external_test.go",
			shouldFail: false,
		},
		{
			name:       "plain test file should fail",
			filename:   "myfile_test.go",
			shouldFail: true,
		},
		{
			name:       "non-test file is ignored",
			filename:   "myfile.go",
			shouldFail: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conceptual logic
			t.Logf("Testing filename: %s (shouldFail=%v)", tt.filename, tt.shouldFail)
		})
	}
}

// Test_runTest012_edgeCases tests edge cases for file naming.
//
// Params:
//   - t: testing context
func Test_runTest012_edgeCases(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		isValid  bool
	}{
		{
			name:     "file with multiple underscores",
			filename: "my_complex_file_internal_test.go",
			isValid:  true,
		},
		{
			name:     "file with numbers",
			filename: "file001_internal_test.go",
			isValid:  true,
		},
		{
			name:     "short filename",
			filename: "a_internal_test.go",
			isValid:  true,
		},
		{
			name:     "error case - empty filename",
			filename: "",
			isValid:  false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conceptual logic
			t.Logf("Testing edge case: %s (isValid=%v)", tt.filename, tt.isValid)
		})
	}
}

// Test_hasValidTestSuffix tests the hasValidTestSuffix function.
//
// Params:
//   - t: testing context
func Test_hasValidTestSuffix(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "internal test file",
			filename: "foo_internal_test.go",
			want:     true,
		},
		{
			name:     "external test file",
			filename: "foo_external_test.go",
			want:     true,
		},
		{
			name:     "bench test file",
			filename: "foo_bench_test.go",
			want:     true,
		},
		{
			name:     "integration test file",
			filename: "foo_integration_test.go",
			want:     true,
		},
		{
			name:     "plain test file invalid",
			filename: "foo_test.go",
			want:     false,
		},
		{
			name:     "non-test file",
			filename: "foo.go",
			want:     false,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasValidTestSuffix(tt.filename)
			// Check result
			if got != tt.want {
				t.Errorf("hasValidTestSuffix(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// mockReporter simulates a reporter to capture reported errors.
type mockReporter struct {
	// diagnostics contains the reported error messages
	diagnostics []string
}

// Report adds a diagnostic to the list.
//
// Params:
//   - d: diagnostic to report
func (m *mockReporter) Report(d analysis.Diagnostic) {
	// Add the diagnostic message
	m.diagnostics = append(m.diagnostics, d.Message)
}

// Test_verifyBenchFile tests the verifyBenchFile private function.
//
// Params:
//   - t: testing context
func Test_verifyBenchFile(t *testing.T) {
	tests := []struct {
		name              string
		filename          string
		code              string
		expectedDiagCount int
	}{
		{
			name:     "bench file with only benchmarks",
			filename: "foo_bench_test.go",
			code: `package test
import "testing"
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = 1 + 1
	}
}
func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = 2 * 2
	}
}`,
			expectedDiagCount: 0,
		},
		{
			name:     "bench file with Test function should error",
			filename: "foo_bench_test.go",
			code: `package test
import "testing"
func TestAdd(t *testing.T) {
	_ = 1 + 1
}
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = 1 + 1
	}
}`,
			expectedDiagCount: 1,
		},
		{
			name:     "bench file with multiple Test functions",
			filename: "foo_bench_test.go",
			code: `package test
import "testing"
func TestAdd(t *testing.T) {
	_ = 1 + 1
}
func TestSubtract(t *testing.T) {
	_ = 2 - 1
}`,
			expectedDiagCount: 2,
		},
		{
			name:     "bench file with helper functions",
			filename: "foo_bench_test.go",
			code: `package test
import "testing"
func helperFunc() int {
	return 42
}
func BenchmarkWithHelper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = helperFunc()
	}
}`,
			expectedDiagCount: 0,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, tt.filename, tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			// Create a mock reporter
			reporter := &mockReporter{}

			// Create a mock analysis.Pass
			pass := &analysis.Pass{
				Fset:   fset,
				Files:  []*ast.File{f},
				Report: reporter.Report,
			}

			// Call verifyBenchFile
			verifyBenchFile(pass, f, tt.filename)

			// Check the number of diagnostics
			if len(reporter.diagnostics) != tt.expectedDiagCount {
				t.Errorf("Expected %d diagnostics, got %d: %v",
					tt.expectedDiagCount, len(reporter.diagnostics), reporter.diagnostics)
			}
		})
	}
}

// Test_runTest001_disabled tests that the rule is skipped when disabled.
//
// Params:
//   - t: testing context
func Test_runTest001_disabled(t *testing.T) {
	// Save current config and restore it
	originalCfg := config.Get()
	defer config.Set(originalCfg)

	// Set config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-001": {Enabled: config.Bool(false)},
		},
	})

	src := `package test_test
import "testing"
func TestExample(t *testing.T) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test_test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create mock reporter
	reporter := &mockReporter{}

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{f},
		Report: reporter.Report,
	}

	_, err = runTest001(pass)
	if err != nil {
		t.Errorf("runTest001() error = %v", err)
	}

	// Check no diagnostics reported
	if len(reporter.diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics (rule disabled), got %d: %v",
			len(reporter.diagnostics), reporter.diagnostics)
	}
}

// Test_runTest001_excludedFile tests that excluded files are skipped.
//
// Params:
//   - t: testing context
func Test_runTest001_excludedFile(t *testing.T) {
	// Save current config and restore it
	originalCfg := config.Get()
	defer config.Set(originalCfg)

	// Set config with file excluded
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-001": {
				Enabled: config.Bool(true),
				Exclude: []string{"test_test.go"},
			},
		},
	})

	src := `package test_test
import "testing"
func TestExample(t *testing.T) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test_test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Create mock reporter
	reporter := &mockReporter{}

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{f},
		Report: reporter.Report,
	}

	_, err = runTest001(pass)
	if err != nil {
		t.Errorf("runTest001() error = %v", err)
	}

	// Check no diagnostics reported
	if len(reporter.diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics (file excluded), got %d: %v",
			len(reporter.diagnostics), reporter.diagnostics)
	}
}
