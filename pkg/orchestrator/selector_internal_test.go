// Internal tests for the analyzer selector.
package orchestrator

import (
	"bytes"
	"testing"
)

// TestAnalyzerSelector_selectSingleRule tests the selectSingleRule method.
func TestAnalyzerSelector_selectSingleRule(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		verbose     bool
		expectError bool
		wantCount   int
	}{
		{
			name:        "valid rule code",
			code:        "KTN-FUNC-001",
			verbose:     false,
			expectError: false,
			wantCount:   1,
		},
		{
			name:        "valid rule code with verbose",
			code:        "KTN-FUNC-001",
			verbose:     true,
			expectError: false,
			wantCount:   1,
		},
		{
			name:        "invalid rule code",
			code:        "KTN-INVALID-999",
			verbose:     false,
			expectError: true,
			wantCount:   0,
		},
		{
			name:        "empty rule code",
			code:        "",
			verbose:     false,
			expectError: true,
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			selector := NewAnalyzerSelector(&buf, tt.verbose)

			analyzers, err := selector.selectSingleRule(tt.code)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Verify count
			if len(analyzers) != tt.wantCount {
				t.Errorf("expected %d analyzers, got %d", tt.wantCount, len(analyzers))
			}
			// Verify verbose output
			if tt.verbose && !tt.expectError && !bytes.Contains(buf.Bytes(), []byte("Running only rule")) {
				t.Error("expected verbose output")
			}
		})
	}
}

// TestAnalyzerSelector_selectByCategory tests the selectByCategory method.
func TestAnalyzerSelector_selectByCategory(t *testing.T) {
	tests := []struct {
		name        string
		category    string
		verbose     bool
		expectError bool
		minCount    int
	}{
		{
			name:        "valid category func",
			category:    "func",
			verbose:     false,
			expectError: false,
			minCount:    1,
		},
		{
			name:        "valid category const with verbose",
			category:    "const",
			verbose:     true,
			expectError: false,
			minCount:    1,
		},
		{
			name:        "invalid category",
			category:    "nonexistent",
			verbose:     false,
			expectError: true,
			minCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			selector := NewAnalyzerSelector(&buf, tt.verbose)

			analyzers, err := selector.selectByCategory(tt.category)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Verify count
			if !tt.expectError && len(analyzers) < tt.minCount {
				t.Errorf("expected at least %d analyzers, got %d", tt.minCount, len(analyzers))
			}
			// Verify verbose output
			if tt.verbose && !tt.expectError && !bytes.Contains(buf.Bytes(), []byte("Running")) {
				t.Error("expected verbose output")
			}
		})
	}
}

// TestAnalyzerSelector_selectAll tests the selectAll method.
func TestAnalyzerSelector_selectAll(t *testing.T) {
	tests := []struct {
		name     string
		verbose  bool
		minCount int
	}{
		{
			name:     "select all without verbose",
			verbose:  false,
			minCount: 1,
		},
		{
			name:     "select all with verbose",
			verbose:  true,
			minCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			selector := NewAnalyzerSelector(&buf, tt.verbose)

			analyzers, err := selector.selectAll()

			// Verify no error
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Verify count
			if len(analyzers) < tt.minCount {
				t.Errorf("expected at least %d analyzers, got %d", tt.minCount, len(analyzers))
			}
		})
	}
}
