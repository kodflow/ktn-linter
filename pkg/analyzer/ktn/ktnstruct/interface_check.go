// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import "go/types"

// interfaceCheck représente un compile-time check var _ I = S
type interfaceCheck struct {
	structName    string
	interfaceType *types.Interface
}
