package orchestrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
)

// TestNewModuleDiscovery tests the constructor.
func TestNewModuleDiscovery(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "returns non-nil instance"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := orchestrator.NewModuleDiscovery()
			// Check not nil
			if d == nil {
				t.Error("expected non-nil ModuleDiscovery")
			}
		})
	}
}

// TestFindModules tests the FindModules method with various inputs.
func TestFindModules(t *testing.T) {
	tests := []struct {
		name           string
		paths          []string
		expectedCount  int
		setupFunc      func(t *testing.T) string
		cleanupFunc    func(string)
		useSetupResult bool
	}{
		{
			name:          "empty input returns empty result",
			paths:         []string{},
			expectedCount: 0,
		},
		{
			name:          "current directory finds module",
			paths:         []string{"/workspace"},
			expectedCount: 1,
		},
		{
			name:          "non-existent path returns empty",
			paths:         []string{"/nonexistent/path"},
			expectedCount: 0,
		},
		{
			name:          "file path finds module",
			paths:         []string{"/workspace/go.mod"},
			expectedCount: 1,
		},
		{
			name:          "duplicate paths are deduplicated",
			paths:         []string{"/workspace", "/workspace"},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := orchestrator.NewModuleDiscovery()
			paths := tt.paths
			// Use setup result if configured
			if tt.setupFunc != nil && tt.useSetupResult {
				tmpDir := tt.setupFunc(t)
				paths = []string{tmpDir}
				// Cleanup if needed
				if tt.cleanupFunc != nil {
					defer tt.cleanupFunc(tmpDir)
				}
			}

			modules, err := d.FindModules(paths)
			// Check no error
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Check expected count
			if len(modules) != tt.expectedCount {
				t.Errorf("expected %d modules, got %d", tt.expectedCount, len(modules))
			}
		})
	}
}

// TestFindModulesSkipsDirectories tests that specific directories are skipped.
func TestFindModulesSkipsDirectories(t *testing.T) {
	tests := []struct {
		name       string
		setupDir   string
		setupMod   string
		moduleName string
	}{
		{
			name:       "skips vendor directories",
			setupDir:   "vendor/dep",
			setupMod:   "module dep\n",
			moduleName: "dep",
		},
		{
			name:       "skips hidden directories",
			setupDir:   ".hidden",
			setupMod:   "module hidden\n",
			moduleName: "hidden",
		},
		{
			name:       "skips testdata directories",
			setupDir:   "testdata/src/example",
			setupMod:   "module example\n",
			moduleName: "example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory structure
			tmpDir := t.TempDir()
			// Create main go.mod
			mainMod := filepath.Join(tmpDir, "go.mod")
			err := os.WriteFile(mainMod, []byte("module test\n"), 0o644)
			// Check error
			if err != nil {
				t.Fatal(err)
			}
			// Create skipped directory with go.mod
			skipDir := filepath.Join(tmpDir, tt.setupDir)
			err = os.MkdirAll(skipDir, 0o755)
			// Check error
			if err != nil {
				t.Fatal(err)
			}
			skipMod := filepath.Join(skipDir, "go.mod")
			err = os.WriteFile(skipMod, []byte(tt.setupMod), 0o644)
			// Check error
			if err != nil {
				t.Fatal(err)
			}

			d := orchestrator.NewModuleDiscovery()
			modules, err := d.FindModules([]string{tmpDir})
			// Check no error
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Check only main module found (directory skipped)
			if len(modules) != 1 {
				t.Errorf("expected 1 module (%s skipped), got %d", tt.setupDir, len(modules))
			}
		})
	}
}

// TestResolvePatterns tests pattern resolution.
func TestResolvePatterns(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
		patterns []string
		expected []string
	}{
		{
			name:     "recursive pattern preserved",
			rootPath: "/workspace",
			patterns: []string{"./..."},
			expected: []string{"./..."},
		},
		{
			name:     "path defaults to recursive",
			rootPath: "/workspace",
			patterns: []string{"/some/path"},
			expected: []string{"./..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := orchestrator.NewModuleDiscovery()
			patterns := d.ResolvePatterns(tt.rootPath, tt.patterns)
			// Check result
			if len(patterns) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, patterns)
				return
			}
			// Check each pattern
			for i, p := range patterns {
				// Verify pattern matches
				if p != tt.expected[i] {
					t.Errorf("expected pattern[%d]=%s, got %s", i, tt.expected[i], p)
				}
			}
		})
	}
}

// TestFindModulesMultiple tests finding multiple modules.
func TestFindModulesMultiple(t *testing.T) {
	tests := []struct {
		name          string
		moduleCount   int
		expectedCount int
	}{
		{
			name:          "finds two modules in subdirectories",
			moduleCount:   2,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory structure with modules
			tmpDir := t.TempDir()

			// Create modules
			for i := 1; i <= tt.moduleCount; i++ {
				modDir := filepath.Join(tmpDir, "mod"+string(rune('0'+i)))
				err := os.MkdirAll(modDir, 0o755)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				err = os.WriteFile(
					filepath.Join(modDir, "go.mod"),
					[]byte("module mod"+string(rune('0'+i))+"\n"),
					0o644,
				)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
			}

			d := orchestrator.NewModuleDiscovery()
			modules, err := d.FindModules([]string{tmpDir})
			// Check no error
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Check found expected modules
			if len(modules) != tt.expectedCount {
				t.Errorf("expected %d modules, got %d", tt.expectedCount, len(modules))
			}
		})
	}
}
