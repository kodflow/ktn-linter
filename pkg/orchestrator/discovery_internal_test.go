// Internal tests for the module discovery.
package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
)

// TestModuleDiscovery_findInPath tests the findInPath method.
func TestModuleDiscovery_findInPath(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantCount int
		wantErr   bool
	}{
		{
			name: "directory with go.mod",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "non-existent path",
			setupFunc: func(t *testing.T) string {
				return "/nonexistent/path/that/does/not/exist"
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			d := NewModuleDiscovery()
			path := tt.setupFunc(t)

			result, err := d.findInPath(path)

			// Verify error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("findInPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Verify result count
			if len(result) != tt.wantCount {
				t.Errorf("findInPath() got %d modules, want %d", len(result), tt.wantCount)
			}
		})
	}
}

// TestModuleDiscovery_findModuleForFile tests the findModuleForFile method.
func TestModuleDiscovery_findModuleForFile(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantCount int
		wantErr   bool
	}{
		{
			name: "file in module directory",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				filePath := filepath.Join(dir, "main.go")
				err = os.WriteFile(filePath, []byte("package main\n"), 0o644)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "file without module",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				filePath := filepath.Join(dir, "main.go")
				err := os.WriteFile(filePath, []byte("package main\n"), 0o644)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			d := NewModuleDiscovery()
			path := tt.setupFunc(t)

			result, err := d.findModuleForFile(path)

			// Verify error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("findModuleForFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Verify result count
			if len(result) != tt.wantCount {
				t.Errorf("findModuleForFile() got %d modules, want %d", len(result), tt.wantCount)
			}
		})
	}
}

// TestModuleDiscovery_searchDirectory tests the searchDirectory method.
func TestModuleDiscovery_searchDirectory(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantCount int
		wantErr   bool
	}{
		{
			name: "directory with single module",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "empty directory",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "directory with nested modules",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				// Create main module
				err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				// Create nested module
				subDir := filepath.Join(dir, "submodule")
				err = os.MkdirAll(subDir, 0o755)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				err = os.WriteFile(filepath.Join(subDir, "go.mod"), []byte("module test/submodule\n"), 0o644)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			d := NewModuleDiscovery()
			path := tt.setupFunc(t)

			result, err := d.searchDirectory(path)

			// Verify error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("searchDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Verify result count
			if len(result) != tt.wantCount {
				t.Errorf("searchDirectory() got %d modules, want %d", len(result), tt.wantCount)
			}
		})
	}
}
