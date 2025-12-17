// Internal tests for the analysis runner.
package orchestrator

import (
	"bytes"
	"go/ast"
	"go/token"
	"sync"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// TestAnalysisRunner_analyzePackageParallel tests the analyzePackageParallel method.
func TestAnalysisRunner_analyzePackageParallel(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{
			name:    "analyze without verbose",
			verbose: false,
		},
		{
			name:    "analyze with verbose",
			verbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, tt.verbose)

			fset := token.NewFileSet()
			pkg := &packages.Package{
				PkgPath: "test/pkg",
				Fset:    fset,
				Syntax:  []*ast.File{},
			}

			results := make(map[*analysis.Analyzer]any)
			diagChan := make(chan DiagnosticResult, 10)

			// Should not panic
			runner.analyzePackageParallel(pkg, []*analysis.Analyzer{}, results, diagChan)
			close(diagChan)

			// Verify verbose output
			if tt.verbose && !bytes.Contains(buf.Bytes(), []byte("Analyzing package")) {
				t.Error("expected verbose output")
			}
		})
	}
}

// TestAnalysisRunner_createPassParallel tests the createPassParallel method.
func TestAnalysisRunner_createPassParallel(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create pass with empty files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, false)

			fset := token.NewFileSet()
			pkg := &packages.Package{
				PkgPath: "test/pkg",
				Fset:    fset,
				Syntax:  []*ast.File{},
			}

			analyzer := &analysis.Analyzer{
				Name: "test",
				Run:  func(*analysis.Pass) (any, error) { return nil, nil },
			}

			results := make(map[*analysis.Analyzer]any)
			diagChan := make(chan DiagnosticResult, 10)

			pass := runner.createPassParallel(analyzer, pkg, fset, diagChan, results)
			close(diagChan)

			// Verify pass created
			if pass == nil {
				t.Error("expected non-nil pass")
			}
			// Verify analyzer set
			if pass.Analyzer != analyzer {
				t.Error("expected analyzer to be set")
			}
			// Verify fset set
			if pass.Fset != fset {
				t.Error("expected fset to be set")
			}
		})
	}
}

// TestAnalysisRunner_worker tests the worker method.
func TestAnalysisRunner_worker(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "worker processes packages from channel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, false)

			fset := token.NewFileSet()
			pkg := &packages.Package{
				PkgPath: "test/pkg",
				Fset:    fset,
				Syntax:  []*ast.File{},
			}

			pkgChan := make(chan *packages.Package, 1)
			diagChan := make(chan DiagnosticResult, 10)
			var wg sync.WaitGroup

			pkgChan <- pkg
			close(pkgChan)

			wg.Add(1)
			// Worker will call wg.Done() via defer
			runner.worker([]*analysis.Analyzer{}, pkgChan, diagChan, &wg)
			close(diagChan)
			// Wait for worker to complete
			wg.Wait()
		})
	}
}

// TestAnalysisRunner_selectFiles tests the selectFiles method.
func TestAnalysisRunner_selectFiles(t *testing.T) {
	tests := []struct {
		name         string
		analyzerName string
		files        []*ast.File
	}{
		{
			name:         "ktntest analyzer gets all files",
			analyzerName: "ktntest001",
			files:        []*ast.File{},
		},
		{
			name:         "other analyzer filters test files",
			analyzerName: "ktnfunc001",
			files:        []*ast.File{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, false)

			fset := token.NewFileSet()
			pkg := &packages.Package{
				PkgPath: "test/pkg",
				Fset:    fset,
				Syntax:  tt.files,
			}

			analyzer := &analysis.Analyzer{
				Name: tt.analyzerName,
			}

			files := runner.selectFiles(analyzer, pkg, fset)

			// Verify files returned
			if files == nil {
				t.Error("expected non-nil files slice")
			}
		})
	}
}

// TestAnalysisRunner_filterTestFiles tests the filterTestFiles method.
func TestAnalysisRunner_filterTestFiles(t *testing.T) {
	tests := []struct {
		name    string
		files   []*ast.File
		wantLen int
	}{
		{
			name:    "empty files",
			files:   []*ast.File{},
			wantLen: 0,
		},
		{
			name:    "nil files",
			files:   nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, false)

			fset := token.NewFileSet()
			filtered := runner.filterTestFiles(tt.files, fset)

			// Verify result length
			if len(filtered) != tt.wantLen {
				t.Errorf("expected %d files, got %d", tt.wantLen, len(filtered))
			}
		})
	}
}

// TestAnalysisRunner_runRequired tests the runRequired method.
func TestAnalysisRunner_runRequired(t *testing.T) {
	tests := []struct {
		name     string
		requires []*analysis.Analyzer
	}{
		{
			name:     "no requirements",
			requires: nil,
		},
		{
			name:     "empty requirements",
			requires: []*analysis.Analyzer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, false)

			fset := token.NewFileSet()
			pkg := &packages.Package{
				PkgPath: "test/pkg",
				Fset:    fset,
				Syntax:  []*ast.File{},
			}

			analyzer := &analysis.Analyzer{
				Name:     "test",
				Requires: tt.requires,
			}

			results := make(map[*analysis.Analyzer]any)

			// Should not panic
			runner.runRequired(analyzer, []*ast.File{}, pkg, fset, results)
		})
	}
}
