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

// FuncInfo stocke les informations d'une fonction avec sa position.
// Peut être utilisé pour fonctions publiques, méthodes, etc.
type FuncInfo struct {
	Name     string
	Pos      token.Pos
	Filename string
}

// ASTNodeInfo représente un nœud AST générique avec métadonnées.
// Utilisable pour struct, interface, type alias, etc.
type ASTNodeInfo struct {
	Name string
	Node ast.Node
}

// MethodSignature représente la signature d'une méthode.
// Utilisé pour comparer les signatures entre struct et interface.
type MethodSignature struct {
	Name       string
	ParamsStr  string
	ResultsStr string
}
