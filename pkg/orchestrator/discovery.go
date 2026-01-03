// Package orchestrator coordinates the linting pipeline.
package orchestrator

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
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
			var emptySlice []string
			// Return error with empty slice
			return emptySlice, err
		}
		// Add found modules
		for _, m := range found {
			modules[m] = struct{}{}
		}
	}

	// Convert to slice using modern Go 1.23 approach
	result := slices.Collect(maps.Keys(modules))

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
		var emptySlice []string
		// Return error with empty slice
		return emptySlice, err
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	// Check for error
	if err != nil {
		var emptySlice []string
		// Return empty if not exists
		return emptySlice, nil
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
	var emptySlice []string
	// No module found - return empty slice
	return emptySlice, nil
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

	// Walk directory tree with callback
	err := filepath.WalkDir(dir, d.createWalkCallback(&modules))

	// Check for walk error
	if err != nil {
		var emptySlice []string
		// Return error with empty slice
		return emptySlice, err
	}

	// Return found modules
	return modules, nil
}

// createWalkCallback creates a WalkDir callback that populates modules slice.
//
// Params:
//   - modules: pointer to slice to populate with found modules
//
// Returns:
//   - fs.WalkDirFunc: callback function for WalkDir
func (d *ModuleDiscovery) createWalkCallback(modules *[]string) fs.WalkDirFunc {
	// Return callback that processes each entry
	return func(path string, entry os.DirEntry, walkErr error) error {
		// Skip inaccessible directories
		if walkErr != nil {
			// Continue walking - error handled by returning nil
			return nil
		}

		// Check if directory should be skipped
		if shouldSkipDirectory(entry) {
			// Skip this directory
			return filepath.SkipDir
		}

		// Check for go.mod file
		if entry.Name() == "go.mod" {
			*modules = append(*modules, filepath.Dir(path))
		}

		// Continue walking - success indicated by nil error
		return nil
	}
}

// shouldSkipDirectory checks if a directory entry should be skipped during walk.
//
// Params:
//   - entry: directory entry to check
//
// Returns:
//   - bool: true if directory should be skipped
func shouldSkipDirectory(entry os.DirEntry) bool {
	// Only check directories
	if !entry.IsDir() {
		return false
	}

	name := entry.Name()

	// Skip hidden directories (start with dot)
	if strings.HasPrefix(name, ".") {
		return true
	}

	// Skip vendor and testdata directories
	return name == "vendor" || name == "testdata"
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
