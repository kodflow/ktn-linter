// Internal tests for the analysis runner.
package orchestrator

import (
	"bytes"
	"go/ast"
	"go/token"
	"strings"
	"sync"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
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
		tt := tt // Capture range variable
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
		{
			name: "create pass and test Report function",
		},
		{
			name: "create pass and test ReadFile function",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
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

			// Test Report function
			if tt.name == "create pass and test Report function" {
				pos := fset.AddFile("test.go", -1, 100).Pos(0)
				pass.Report(analysis.Diagnostic{
					Pos:     pos,
					Message: "test diagnostic",
				})
				// Verify diagnostic sent to channel
				select {
				case diag := <-diagChan:
					// Verify diagnostic received
					if diag.AnalyzerName != analyzer.Name {
						t.Errorf("expected analyzer name %s, got %s", analyzer.Name, diag.AnalyzerName)
					}
				default:
					t.Error("expected diagnostic in channel")
				}
			}

			// Test ReadFile function
			if tt.name == "create pass and test ReadFile function" {
				// Try to read a non-existent file
				_, err := pass.ReadFile("/nonexistent/file.go")
				// Verify error returned
				if err == nil {
					t.Error("expected error reading non-existent file")
				}
			}

			close(diagChan)
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
		tt := tt // Capture range variable
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
			runner.worker([]*analysis.Analyzer{}, pkgChan, diagChan, &wg, 0)
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
		tt := tt // Capture range variable
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

// TestAnalysisRunner_selectFilesWithConfig tests selectFiles with config.
func TestAnalysisRunner_selectFilesWithConfig(t *testing.T) {
	tests := []struct {
		name         string
		analyzerName string
		forceAll     bool
		wantAll      bool
	}{
		{
			name:         "force all rules on tests enabled",
			analyzerName: "ktnfunc001",
			forceAll:     true,
			wantAll:      true,
		},
		{
			name:         "force all rules on tests disabled",
			analyzerName: "ktnfunc001",
			forceAll:     false,
			wantAll:      false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore config
			oldCfg := config.Get()
			defer config.Set(oldCfg)

			// Set test config
			cfg := config.DefaultConfig()
			cfg.ForceAllRulesOnTests = tt.forceAll
			config.Set(cfg)

			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, false)

			fset := token.NewFileSet()
			// Create test file
			testFile := fset.AddFile("test_test.go", -1, 100)
			normalFile := fset.AddFile("normal.go", -1, 100)

			files := []*ast.File{
				{
					Package: testFile.Pos(0),
				},
				{
					Package: normalFile.Pos(0),
				},
			}

			pkg := &packages.Package{
				PkgPath: "test/pkg",
				Fset:    fset,
				Syntax:  files,
			}

			analyzer := &analysis.Analyzer{
				Name: tt.analyzerName,
			}

			selected := runner.selectFiles(analyzer, pkg, fset)

			// Verify files returned
			if selected == nil {
				t.Error("expected non-nil files slice")
			}

			// Verify all files returned when force is enabled
			if tt.wantAll && len(selected) != len(files) {
				t.Errorf("expected all %d files with ForceAllRulesOnTests, got %d", len(files), len(selected))
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
		tt := tt // Capture range variable
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

// TestAnalysisRunner_filterTestFilesWithActualFiles tests filtering with actual test files.
func TestAnalysisRunner_filterTestFilesWithActualFiles(t *testing.T) {
	tests := []struct {
		name      string
		filenames []string
		wantLen   int
	}{
		{
			name:      "mix of test and non-test files",
			filenames: []string{"main.go", "main_test.go", "util.go"},
			wantLen:   2, // Only main.go and util.go
		},
		{
			name:      "only test files",
			filenames: []string{"main_test.go", "util_test.go"},
			wantLen:   0,
		},
		{
			name:      "only non-test files",
			filenames: []string{"main.go", "util.go"},
			wantLen:   2,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, false)

			fset := token.NewFileSet()
			files := make([]*ast.File, len(tt.filenames))
			for i, filename := range tt.filenames {
				f := fset.AddFile(filename, -1, 100)
				files[i] = &ast.File{
					Package: f.Pos(0),
				}
			}

			filtered := runner.filterTestFiles(files, fset)

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
		wantRun  bool
	}{
		{
			name:     "no requirements",
			requires: nil,
			wantRun:  false,
		},
		{
			name:     "empty requirements",
			requires: []*analysis.Analyzer{},
			wantRun:  false,
		},
		{
			name: "with required analyzer",
			requires: []*analysis.Analyzer{
				{
					Name: "required",
					Run: func(pass *analysis.Pass) (any, error) {
						return "required result", nil
					},
				},
			},
			wantRun: true,
		},
		{
			name: "with already cached required analyzer",
			requires: []*analysis.Analyzer{
				{
					Name: "cached",
					Run:  func(pass *analysis.Pass) (any, error) { return nil, nil },
				},
			},
			wantRun: false,
		},
		{
			name: "with required analyzer that uses ReadFile",
			requires: []*analysis.Analyzer{
				{
					Name: "readfile",
					Run: func(pass *analysis.Pass) (any, error) {
						// Test ReadFile function
						_, err := pass.ReadFile("/nonexistent/file.go")
						// Return result even if ReadFile fails
						return "readfile result", err
					},
				},
			},
			wantRun: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
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

			// Pre-cache if needed
			if tt.name == "with already cached required analyzer" && len(tt.requires) > 0 {
				results[tt.requires[0]] = "cached result"
			}

			// Should not panic
			runner.runRequired(analyzer, []*ast.File{}, pkg, fset, results)

			// Verify required analyzer was run
			if tt.wantRun && len(tt.requires) > 0 {
				if _, exists := results[tt.requires[0]]; !exists {
					t.Error("expected required analyzer to be run")
				}
			}
		})
	}
}

// TestAnalysisRunner_filterExcludedFiles tests the filterExcludedFiles method.
func TestAnalysisRunner_filterExcludedFiles(t *testing.T) {
	tests := []struct {
		name           string
		excludePattern []string
		files          []string
		wantCount      int
		verbose        bool
	}{
		{
			name:           "no exclusions returns all files",
			excludePattern: []string{},
			files:          []string{"file1.go", "file2.go"},
			wantCount:      2,
			verbose:        false,
		},
		{
			name:           "excludes matching files",
			excludePattern: []string{"**/gen/**"},
			files:          []string{"src/gen/file.go", "src/main.go"},
			wantCount:      1,
			verbose:        false,
		},
		{
			name:           "excludes pb.go files",
			excludePattern: []string{"*.pb.go"},
			files:          []string{"service.pb.go", "main.go"},
			wantCount:      1,
			verbose:        false,
		},
		{
			name:           "verbose logs excluded files",
			excludePattern: []string{"*.pb.go"},
			files:          []string{"service.pb.go"},
			wantCount:      0,
			verbose:        true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Setup config
			config.Set(&config.Config{
				Exclude: tt.excludePattern,
			})
			defer config.Reset()

			var buf bytes.Buffer
			runner := NewAnalysisRunner(&buf, tt.verbose)

			// Create fset and files
			fset := token.NewFileSet()
			var astFiles []*ast.File
			// Add files to fset
			for _, filename := range tt.files {
				f := fset.AddFile(filename, -1, 100)
				astFile := &ast.File{
					Package: f.Pos(0),
				}
				astFiles = append(astFiles, astFile)
			}

			// Call filterExcludedFiles
			result := runner.filterExcludedFiles(astFiles, fset)

			// Verify count
			if len(result) != tt.wantCount {
				t.Errorf("filterExcludedFiles() got %d files, want %d", len(result), tt.wantCount)
			}

			// Verify verbose output
			if tt.verbose && len(tt.excludePattern) > 0 && len(tt.files) > 0 {
				output := buf.String()
				// Check if "Excluding file" is in output
				if len(result) < len(astFiles) && !strings.Contains(output, "Excluding file") {
					t.Errorf("expected verbose exclusion output, got: %s", output)
				}
			}
		})
	}
}
