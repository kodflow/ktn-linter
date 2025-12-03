// Shared utilities for astnodeinfo handling.
package shared

import "go/ast"

// ASTNodeInfo représente un nœud AST générique avec métadonnées.
// Utilisable pour struct, interface, type alias, etc.
type ASTNodeInfo struct {
	Name string
	Node ast.Node
}
