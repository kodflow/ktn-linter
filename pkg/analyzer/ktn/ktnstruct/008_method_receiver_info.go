// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import "go/ast"

// methodReceiverInfo contient les informations sur le receiver d'une méthode.
type methodReceiverInfo struct {
	name      string
	isPointer bool
	funcDecl  *ast.FuncDecl
}
