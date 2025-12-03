// Shared utilities for types handling.
package shared

import (
	"go/ast"
	"go/token"
)

// DeclGroup représente un groupe de déclarations (const ou var).
// Utilisé par ktnconst/002.go et ktnvar/002.go pour vérifier le groupement.
type DeclGroup struct {
	Decl *ast.GenDecl
	Pos  token.Pos
}
