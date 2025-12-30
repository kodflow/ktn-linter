package orchestrator_test

import (
	"os"
	"path/filepath"
	"strconv"
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
		pathCount      int // Number of times to duplicate the setup path (0 or 1 = single path)
	}{
		{
			name:          "empty input returns empty result",
			paths:         []string{},
			expectedCount: 0,
		},
		{
			name:          "current directory finds module",
			paths:         []string{}, // Will be set by setupFunc
			expectedCount: 1,
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0o644)
				if err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			useSetupResult: true,
		},
		{
			name:          "non-existent path returns empty",
			paths:         []string{filepath.Join(string(filepath.Separator), "nonexistent", "path")},
			expectedCount: 0,
		},
		{
			name:          "file path finds module",
			paths:         []string{}, // Will be set by setupFunc
			expectedCount: 1,
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				goModPath := filepath.Join(tmpDir, "go.mod")
				err := os.WriteFile(goModPath, []byte("module test\n"), 0o644)
				if err != nil {
					t.Fatal(err)
				}
				return goModPath // Return the go.mod file path
			},
			useSetupResult: true,
		},
		{
			name:          "duplicate paths are deduplicated",
			paths:         []string{}, // Will be set by setupFunc
			expectedCount: 1,
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0o644)
				if err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			useSetupResult: true,
			pathCount:      2, // Duplicate the path to test deduplication
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := orchestrator.NewModuleDiscovery()
			paths := tt.paths
			// Use setup result if configured
			if tt.setupFunc != nil && tt.useSetupResult {
				setupResult := tt.setupFunc(t)
				// Handle path count for duplicate testing
				count := tt.pathCount
				if count <= 1 {
					count = 1
				}
				paths = make([]string, count)
				for i := range paths {
					paths[i] = setupResult
				}
				// Cleanup if needed
				if tt.cleanupFunc != nil {
					defer tt.cleanupFunc(setupResult)
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
		patterns []string
		expected []string
	}{
		{
			name:     "recursive pattern preserved",
			patterns: []string{"./..."},
			expected: []string{"./..."},
		},
		{
			name:     "path defaults to recursive",
			patterns: []string{"/some/path"},
			expected: []string{"./..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use temp directory as rootPath
			rootPath := t.TempDir()
			d := orchestrator.NewModuleDiscovery()
			patterns := d.ResolvePatterns(rootPath, tt.patterns)
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
				suffix := strconv.Itoa(i)
				modDir := filepath.Join(tmpDir, "mod"+suffix)
				err := os.MkdirAll(modDir, 0o755)
				// Check error
				if err != nil {
					t.Fatal(err)
				}
				err = os.WriteFile(
					filepath.Join(modDir, "go.mod"),
					[]byte("module mod"+suffix+"\n"),
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
