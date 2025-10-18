package formatter

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Formatter définit l'interface pour formater et afficher les diagnostics.
type Formatter interface {
	// Format affiche les diagnostics de manière lisible
	//
	// Params:
	//   - fset: le FileSet contenant les informations de position
	//   - diagnostics: la liste des diagnostics à formater
	Format(fset *token.FileSet, diagnostics []analysis.Diagnostic)
}
