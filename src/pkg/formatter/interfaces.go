package formatter

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Formatter gère le formatage des diagnostics pour différents modes de sortie.
type Formatter interface {
	// Format affiche les diagnostics de manière lisible selon le mode configuré
	//
	// Params:
	//   - fset: le FileSet contenant les informations de position
	//   - diagnostics: la liste des diagnostics à formater
	Format(fset *token.FileSet, diagnostics []analysis.Diagnostic)
}
