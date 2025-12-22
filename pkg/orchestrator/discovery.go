// Package orchestrator coordinates the linting pipeline.
package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
)

// ModuleDiscovery handles finding Go modules recursively.
// Searches directories for go.mod files and returns module root paths.
type ModuleDiscovery struct{}

// NewModuleDiscovery creates a new ModuleDiscovery.
//
// Returns:
//   - *ModuleDiscovery: new discovery instance
func NewModuleDiscovery() *ModuleDiscovery {
	// Return new instance
	return &ModuleDiscovery{}
}

// FindModules finds all go.mod files in paths.
//
// Params:
//   - paths: paths to search (files or directories)
//
// Returns:
//   - []string: module root directories
//   - error: search error if any
func (d *ModuleDiscovery) FindModules(paths []string) ([]string, error) {
	modules := make(map[string]struct{}, len(paths))

	// Process each path
	for _, p := range paths {
		// Find modules in path
		found, err := d.findInPath(p)
		// Check for error
		if err != nil {
			// Return error with empty slice
			return []string{}, err
		}
		// Add found modules
		for _, m := range found {
			modules[m] = struct{}{}
		}
	}

	// Convert to slice
	result := make([]string, 0, len(modules))
	// Add all modules
	for m := range modules {
		result = append(result, m)
	}

	// Return found modules
	return result, nil
}

// findInPath finds go.mod files in a single path.
//
// Params:
//   - path: path to search
//
// Returns:
//   - []string: module roots found
//   - error: search error if any
func (d *ModuleDiscovery) findInPath(path string) ([]string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	// Check for error
	if err != nil {
		// Return error with empty slice
		return []string{}, err
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	// Check for error
	if err != nil {
		// Return empty if not exists
		return []string{}, nil
	}

	// Check if file
	if !info.IsDir() {
		// Find module for file
		return d.findModuleForFile(absPath)
	}

	// Search directory recursively
	return d.searchDirectory(absPath)
}

// findModuleForFile finds the module containing a file.
//
// Params:
//   - filePath: path to file
//
// Returns:
//   - []string: module root if found
//   - error: search error if any
func (d *ModuleDiscovery) findModuleForFile(filePath string) ([]string, error) {
	dir := filepath.Dir(filePath)
	// Walk up directory tree
	for dir != "/" && dir != "." {
		goMod := filepath.Join(dir, "go.mod")
		// Check if go.mod exists
		if _, err := os.Stat(goMod); err == nil {
			// Return module root
			return []string{dir}, nil
		}
		// Move up
		dir = filepath.Dir(dir)
	}
	// No module found - return empty slice
	return []string{}, nil
}

// searchDirectory searches for go.mod files recursively.
//
// Params:
//   - dir: directory to search
//
// Returns:
//   - []string: module roots found
//   - error: search error if any
func (d *ModuleDiscovery) searchDirectory(dir string) ([]string, error) {
	var modules []string

	// Walk directory tree
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
		// Skip inaccessible directories
		if walkErr != nil {
			// Continue walking
			return nil
		}

		// Skip hidden directories
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") {
			// Skip hidden directory
			return filepath.SkipDir
		}

		// Skip vendor directories
		if entry.IsDir() && entry.Name() == "vendor" {
			// Skip vendor directory
			return filepath.SkipDir
		}

		// Check for go.mod file
		if entry.Name() == "go.mod" {
			modules = append(modules, filepath.Dir(path))
		}

		// Continue walking
		return nil
	})

	// Check for walk error
	if err != nil {
		// Return error with empty slice
		return []string{}, err
	}

	// Return found modules
	return modules, nil
}

// ResolvePatterns resolves patterns for a module.
//
// Params:
//   - _moduleRoot: root directory of module (reserved for future use)
//   - originalPatterns: original patterns from CLI
//
// Returns:
//   - []string: resolved patterns for this module
func (d *ModuleDiscovery) ResolvePatterns(_moduleRoot string, originalPatterns []string) []string {
	// Check for ./... pattern
	for _, p := range originalPatterns {
		// Check if recursive pattern
		if p == "./..." || strings.HasSuffix(p, "/...") {
			// Return recursive pattern for module
			return []string{"./..."}
		}
	}

	// Default to recursive
	return []string{"./..."}
}
