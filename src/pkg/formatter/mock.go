//go:build test
// +build test

package formatter

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// MockFormatter est un mock de Formatter pour les tests.
type MockFormatter struct {
	FormatFunc func(fset *token.FileSet, diagnostics []analysis.Diagnostic)
}

// Format implémente Formatter.Format.
func (m *MockFormatter) Format(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	if m.FormatFunc != nil {
		m.FormatFunc(fset, diagnostics)
	}
}

// Vérification à la compilation que MockFormatter implémente Formatter
var _ Formatter = (*MockFormatter)(nil)
