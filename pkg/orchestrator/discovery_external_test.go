package orchestrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
)

// TestNewModuleDiscovery tests the constructor.
func TestNewModuleDiscovery(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	// Check not nil
	if d == nil {
		t.Error("expected non-nil ModuleDiscovery")
	}
}

// TestFindModulesEmpty tests finding modules with empty input.
func TestFindModulesEmpty(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	modules, err := d.FindModules([]string{})
	// Check no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Check empty result
	if len(modules) != 0 {
		t.Errorf("expected 0 modules, got %d", len(modules))
	}
}

// TestFindModulesCurrentDir tests finding module in current directory.
func TestFindModulesCurrentDir(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	// Use workspace root which has go.mod
	modules, err := d.FindModules([]string{"/workspace"})
	// Check no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Check found module
	if len(modules) != 1 {
		t.Errorf("expected 1 module, got %d", len(modules))
	}
}

// TestFindModulesNonExistent tests finding modules in non-existent path.
func TestFindModulesNonExistent(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	modules, err := d.FindModules([]string{"/nonexistent/path"})
	// Check no error (returns empty)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Check empty result
	if len(modules) != 0 {
		t.Errorf("expected 0 modules, got %d", len(modules))
	}
}

// TestFindModulesFile tests finding module for a file.
func TestFindModulesFile(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	// Use a file in the workspace
	modules, err := d.FindModules([]string{"/workspace/go.mod"})
	// Check no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Check found module
	if len(modules) != 1 {
		t.Errorf("expected 1 module, got %d", len(modules))
	}
}

// TestFindModulesSkipsVendor tests that vendor directories are skipped.
func TestFindModulesSkipsVendor(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	// Create main go.mod
	mainMod := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(mainMod, []byte("module test\n"), 0o644)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	// Create vendor directory with go.mod
	vendorDir := filepath.Join(tmpDir, "vendor", "dep")
	err = os.MkdirAll(vendorDir, 0o755)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	vendorMod := filepath.Join(vendorDir, "go.mod")
	err = os.WriteFile(vendorMod, []byte("module dep\n"), 0o644)
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
	// Check only main module found (vendor skipped)
	if len(modules) != 1 {
		t.Errorf("expected 1 module (vendor skipped), got %d", len(modules))
	}
}

// TestFindModulesSkipsHidden tests that hidden directories are skipped.
func TestFindModulesSkipsHidden(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	// Create main go.mod
	mainMod := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(mainMod, []byte("module test\n"), 0o644)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	// Create hidden directory with go.mod
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	err = os.MkdirAll(hiddenDir, 0o755)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	hiddenMod := filepath.Join(hiddenDir, "go.mod")
	err = os.WriteFile(hiddenMod, []byte("module hidden\n"), 0o644)
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
	// Check only main module found (hidden skipped)
	if len(modules) != 1 {
		t.Errorf("expected 1 module (hidden skipped), got %d", len(modules))
	}
}

// TestFindModulesSkipsTestdata tests that testdata directories are skipped.
func TestFindModulesSkipsTestdata(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	// Create main go.mod
	mainMod := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(mainMod, []byte("module test\n"), 0o644)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	// Create testdata directory with go.mod
	testdataDir := filepath.Join(tmpDir, "testdata", "src", "example")
	err = os.MkdirAll(testdataDir, 0o755)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	testdataMod := filepath.Join(testdataDir, "go.mod")
	err = os.WriteFile(testdataMod, []byte("module example\n"), 0o644)
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
	// Check only main module found (testdata skipped)
	if len(modules) != 1 {
		t.Errorf("expected 1 module (testdata skipped), got %d", len(modules))
	}
}

// TestResolvePatternsRecursive tests pattern resolution for recursive.
func TestResolvePatternsRecursive(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	patterns := d.ResolvePatterns("/workspace", []string{"./..."})
	// Check result
	if len(patterns) != 1 || patterns[0] != "./..." {
		t.Errorf("expected [./...], got %v", patterns)
	}
}

// TestResolvePatternsDefault tests pattern resolution for default.
func TestResolvePatternsDefault(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	patterns := d.ResolvePatterns("/workspace", []string{"/some/path"})
	// Check result (default to recursive)
	if len(patterns) != 1 || patterns[0] != "./..." {
		t.Errorf("expected [./...], got %v", patterns)
	}
}

// TestFindModulesMultiple tests finding multiple modules.
func TestFindModulesMultiple(t *testing.T) {
	// Create temp directory structure with 2 modules
	tmpDir := t.TempDir()
	// Create first module
	mod1 := filepath.Join(tmpDir, "mod1")
	err := os.MkdirAll(mod1, 0o755)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(mod1, "go.mod"), []byte("module mod1\n"), 0o644)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	// Create second module
	mod2 := filepath.Join(tmpDir, "mod2")
	err = os.MkdirAll(mod2, 0o755)
	// Check error
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(mod2, "go.mod"), []byte("module mod2\n"), 0o644)
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
	// Check found both modules
	if len(modules) != 2 {
		t.Errorf("expected 2 modules, got %d", len(modules))
	}
}

// TestFindModulesDeduplicate tests that duplicate modules are handled.
func TestFindModulesDeduplicate(t *testing.T) {
	d := orchestrator.NewModuleDiscovery()
	// Pass same path twice
	modules, err := d.FindModules([]string{"/workspace", "/workspace"})
	// Check no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Check deduplicated
	if len(modules) != 1 {
		t.Errorf("expected 1 module (deduplicated), got %d", len(modules))
	}
}
