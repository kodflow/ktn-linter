// Package orchestrator coordinates the linting pipeline.
package orchestrator

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/tools/go/packages"
)

// PackageLoader handles loading Go packages for analysis.
// Configures the packages.Config and checks for loading errors.
type PackageLoader struct {
	stderr io.Writer
}

// NewPackageLoader creates a new PackageLoader.
//
// Params:
//   - stderr: writer for error output
//
// Returns:
//   - *PackageLoader: new loader instance
func NewPackageLoader(stderr io.Writer) *PackageLoader {
	// Return new loader instance
	return &PackageLoader{stderr: stderr}
}

// Load loads Go packages from the given patterns.
//
// Params:
//   - patterns: package patterns to load
//
// Returns:
//   - []*packages.Package: loaded packages
//   - error: loading error if any
func (l *PackageLoader) Load(patterns []string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests:      true,
		BuildFlags: []string{"-buildvcs=false"},
	}

	pkgs, err := packages.Load(cfg, patterns...)
	// Check for load error
	if err != nil {
		// Return error
		return []*packages.Package{}, fmt.Errorf("loading packages: %w", err)
	}

	// Check for package errors
	if err := l.checkErrors(pkgs); err != nil {
		// Return error
		return []*packages.Package{}, err
	}

	// Return loaded packages
	return pkgs, nil
}

// checkErrors checks for package loading errors.
//
// Params:
//   - pkgs: packages to check
//
// Returns:
//   - error: first non-VCS error found
func (l *PackageLoader) checkErrors(pkgs []*packages.Package) error {
	// Iterate over packages
	for _, pkg := range pkgs {
		// Iterate over errors
		for _, err := range pkg.Errors {
			// Skip VCS errors
			if strings.Contains(err.Error(), "VCS status") {
				continue
			}
			// Return first real error
			return fmt.Errorf("package %s: %v", pkg.PkgPath, err)
		}
	}
	// No errors found
	return nil
}
