// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import "go/ast"

// structFieldsInfo contient les infos sur les champs d'une struct.
type structFieldsInfo struct {
	name          string
	privateFields map[string]bool
	pos           ast.Node
}
