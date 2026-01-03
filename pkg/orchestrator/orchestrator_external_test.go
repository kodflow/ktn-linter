// External tests for the orchestrator package.
package orchestrator_test

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
)

// TestNewOrchestrator tests the NewOrchestrator function.
func TestNewOrchestrator(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{
			name:    "create orchestrator without verbose",
			verbose: false,
		},
		{
			name:    "create orchestrator with verbose",
			verbose: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, tt.verbose)

			// Verify orchestrator created
			if orch == nil {
				t.Error("expected non-nil orchestrator")
			}
		})
	}
}

// TestOrchestrator_LoadPackages tests the LoadPackages method.
func TestOrchestrator_LoadPackages(t *testing.T) {
	tests := []struct {
		name        string
		patterns    []string
		expectError bool
		minPackages int
	}{
		{
			name:        "load valid package",
			patterns:    []string{"../../pkg/formatter"},
			expectError: false,
			minPackages: 1,
		},
		{
			name:        "load invalid pattern",
			patterns:    []string{"./nonexistent/package"},
			expectError: true,
			minPackages: 0,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			pkgs, err := orch.LoadPackages(tt.patterns)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Verify package count
			if !tt.expectError && len(pkgs) < tt.minPackages {
				t.Errorf("expected at least %d packages, got %d", tt.minPackages, len(pkgs))
			}
		})
	}
}

// TestOrchestrator_SelectAnalyzers tests the SelectAnalyzers method.
func TestOrchestrator_SelectAnalyzers(t *testing.T) {
	tests := []struct {
		name        string
		opts        orchestrator.Options
		expectError bool
		minCount    int
	}{
		{
			name:        "select all analyzers",
			opts:        orchestrator.Options{},
			expectError: false,
			minCount:    1,
		},
		{
			name:        "select by category func",
			opts:        orchestrator.Options{Category: "func"},
			expectError: false,
			minCount:    1,
		},
		{
			name:        "select single rule",
			opts:        orchestrator.Options{OnlyRule: "KTN-FUNC-001"},
			expectError: false,
			minCount:    1,
		},
		{
			name:        "unknown category returns error",
			opts:        orchestrator.Options{Category: "nonexistent"},
			expectError: true,
			minCount:    0,
		},
		{
			name:        "unknown rule returns error",
			opts:        orchestrator.Options{OnlyRule: "KTN-INVALID-999"},
			expectError: true,
			minCount:    0,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			analyzers, err := orch.SelectAnalyzers(tt.opts)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Verify analyzer count
			if !tt.expectError && len(analyzers) < tt.minCount {
				t.Errorf("expected at least %d analyzers, got %d", tt.minCount, len(analyzers))
			}
		})
	}
}

// TestOrchestrator_RunAnalyzers tests the RunAnalyzers method.
func TestOrchestrator_RunAnalyzers(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		opts     orchestrator.Options
	}{
		{
			name:     "run analyzers on valid package",
			patterns: []string{"../../pkg/formatter"},
			opts:     orchestrator.Options{OnlyRule: "KTN-FUNC-001"},
		},
		{
			name:     "run analyzers with empty packages",
			patterns: []string{},
			opts:     orchestrator.Options{},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			// Handle empty patterns case
			if len(tt.patterns) == 0 {
				diags := orch.RunAnalyzers(nil, nil)
				// Verify empty result
				if len(diags) != 0 {
					t.Errorf("expected empty diagnostics, got %d", len(diags))
				}
				return
			}

			pkgs, err := orch.LoadPackages(tt.patterns)
			// Verify load succeeded
			if err != nil {
				t.Fatalf("load error: %v", err)
			}

			analyzers, err := orch.SelectAnalyzers(tt.opts)
			// Verify select succeeded
			if err != nil {
				t.Fatalf("select error: %v", err)
			}

			diags := orch.RunAnalyzers(pkgs, analyzers)

			// RunAnalyzers should not panic - diags can be nil or empty
			_ = diags
		})
	}
}

// TestOrchestrator_FilterDiagnostics tests the FilterDiagnostics method.
func TestOrchestrator_FilterDiagnostics(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics []orchestrator.DiagnosticResult
		wantLen     int
	}{
		{
			name:        "filter nil diagnostics",
			diagnostics: nil,
			wantLen:     0,
		},
		{
			name:        "filter empty diagnostics",
			diagnostics: []orchestrator.DiagnosticResult{},
			wantLen:     0,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			filtered := orch.FilterDiagnostics(tt.diagnostics)

			// Verify result length
			if len(filtered) != tt.wantLen {
				t.Errorf("expected %d diagnostics, got %d", tt.wantLen, len(filtered))
			}
		})
	}
}

// TestOrchestrator_ExtractDiagnostics tests the ExtractDiagnostics method.
func TestOrchestrator_ExtractDiagnostics(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics []orchestrator.DiagnosticResult
		wantLen     int
	}{
		{
			name:        "extract nil diagnostics",
			diagnostics: nil,
			wantLen:     0,
		},
		{
			name:        "extract empty diagnostics",
			diagnostics: []orchestrator.DiagnosticResult{},
			wantLen:     0,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			diags := orch.ExtractDiagnostics(tt.diagnostics)

			// Verify result length
			if len(diags) != tt.wantLen {
				t.Errorf("expected %d diagnostics, got %d", tt.wantLen, len(diags))
			}
		})
	}
}

// TestOrchestrator_Run tests the Run method.
func TestOrchestrator_Run(t *testing.T) {
	tests := []struct {
		name        string
		patterns    []string
		opts        orchestrator.Options
		expectError bool
	}{
		{
			name:        "run on valid package",
			patterns:    []string{"../../pkg/formatter"},
			opts:        orchestrator.Options{OnlyRule: "KTN-FUNC-001"},
			expectError: false,
		},
		{
			name:        "run with invalid category",
			patterns:    []string{"../../pkg/formatter"},
			opts:        orchestrator.Options{Category: "nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			diags, err := orch.Run(tt.patterns, tt.opts)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Verify diagnostics not nil on success
			if !tt.expectError && diags == nil {
				t.Error("expected non-nil diagnostics")
			}
		})
	}
}

// TestGetFirstFset tests the GetFirstFset function.
func TestGetFirstFset(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics []orchestrator.DiagnosticResult
		wantNil     bool
	}{
		{
			name:        "nil diagnostics returns nil",
			diagnostics: nil,
			wantNil:     true,
		},
		{
			name:        "empty diagnostics returns nil",
			diagnostics: []orchestrator.DiagnosticResult{},
			wantNil:     true,
		},
		{
			name: "non-empty diagnostics returns fset",
			diagnostics: []orchestrator.DiagnosticResult{
				{
					Fset: token.NewFileSet(),
				},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := orchestrator.GetFirstFset(tt.diagnostics)

			// Verify nil expectation
			if tt.wantNil && result != nil {
				t.Error("expected nil result")
			}
			// Verify non-nil expectation
			if !tt.wantNil && result == nil {
				t.Error("expected non-nil result")
			}
		})
	}
}

// TestOptionsDefaults tests Options default values.
func TestOptionsDefaults(t *testing.T) {
	tests := []struct {
		name  string
		field string
		check func(opts orchestrator.Options) bool
	}{
		{
			name:  "Verbose defaults to false",
			field: "Verbose",
			check: func(opts orchestrator.Options) bool { return !opts.Verbose },
		},
		{
			name:  "Category defaults to empty",
			field: "Category",
			check: func(opts orchestrator.Options) bool { return opts.Category == "" },
		},
		{
			name:  "OnlyRule defaults to empty",
			field: "OnlyRule",
			check: func(opts orchestrator.Options) bool { return opts.OnlyRule == "" },
		},
		{
			name:  "ConfigPath defaults to empty",
			field: "ConfigPath",
			check: func(opts orchestrator.Options) bool { return opts.ConfigPath == "" },
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			opts := orchestrator.Options{}

			// Verify default value
			if !tt.check(opts) {
				t.Errorf("%s does not have expected default value", tt.field)
			}
		})
	}
}

// TestOrchestrator_DiscoverModules tests the DiscoverModules method.
func TestOrchestrator_DiscoverModules(t *testing.T) {
	tests := []struct {
		name        string
		paths       []string
		expectError bool
	}{
		{
			name:        "discover in empty path list",
			paths:       []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			_, err := orch.DiscoverModules(tt.paths)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestOrchestrator_LoadPackagesFromDir tests the LoadPackagesFromDir method.
func TestOrchestrator_LoadPackagesFromDir(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		patterns    []string
		expectError bool
	}{
		{
			name:        "load from invalid directory",
			dir:         "/nonexistent/directory",
			patterns:    []string{"./..."},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			_, err := orch.LoadPackagesFromDir(tt.dir, tt.patterns)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestOrchestrator_RunMultiModule tests the RunMultiModule method.
func TestOrchestrator_RunMultiModule(t *testing.T) {
	tests := []struct {
		name        string
		paths       []string
		opts        orchestrator.Options
		expectError bool
	}{
		{
			name:        "run on empty path list",
			paths:       []string{},
			opts:        orchestrator.Options{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := orchestrator.NewOrchestrator(&buf, false)

			_, err := orch.RunMultiModule(tt.paths, tt.opts)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
