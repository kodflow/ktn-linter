// Internal types for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
)

// structWithMethods stocke une struct et ses méthodes.
type structWithMethods struct {
	name       string
	node       *ast.TypeSpec
	structType *ast.StructType
	methods    []string
}
