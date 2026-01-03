// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import "go/ast"

// methodInfo contient les informations sur une méthode.
type methodInfo struct {
	name       string
	funcDecl   *ast.FuncDecl
	receiverTy string
	returnType string
}
