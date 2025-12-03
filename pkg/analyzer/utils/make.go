// Utility functions for make operations.
package utils

import (
	"go/ast"
)

// IsMakeCall vérifie si un CallExpr est un appel à make.
//
// Params:
//   - call: expression d'appel à vérifier
//
// Returns:
//   - bool: true si c'est un appel à make
func IsMakeCall(call *ast.CallExpr) bool {
	// Check if call is to 'make'
	return IsIdentCall(call, "make")
}

// IsMakeCallWithLength vérifie si c'est un appel à make avec une longueur spécifiée.
//
// Params:
//   - call: expression d'appel à vérifier
//   - minArgs: nombre minimum d'arguments attendus
//
// Returns:
//   - bool: true si c'est un appel à make avec au moins minArgs arguments
func IsMakeCallWithLength(call *ast.CallExpr, minArgs int) bool {
	// Check if it's a make call
	if !IsMakeCall(call) {
		// Not a make call
		return false
	}
	// Check argument count
	return len(call.Args) >= minArgs
}

// IsMakeSliceCall vérifie si c'est un appel à make pour créer un slice.
//
// Params:
//   - call: expression d'appel à vérifier
//
// Returns:
//   - bool: true si c'est make([]T, ...)
func IsMakeSliceCall(call *ast.CallExpr) bool {
	// Check if it's a make call
	if !IsMakeCall(call) {
		// Not a make call
		return false
	}
	// Check if first argument is a slice type
	if len(call.Args) > 0 {
		// Check first argument is slice type
		return IsSliceType(call.Args[0])
	}
	// No arguments
	return false
}

// IsMakeMapCall vérifie si c'est un appel à make pour créer une map.
//
// Params:
//   - call: expression d'appel à vérifier
//
// Returns:
//   - bool: true si c'est make(map[K]V, ...)
func IsMakeMapCall(call *ast.CallExpr) bool {
	// Check if it's a make call
	if !IsMakeCall(call) {
		// Not a make call
		return false
	}
	// Check if first argument is a map type
	if len(call.Args) > 0 {
		// Check first argument is map type
		return IsMapType(call.Args[0])
	}
	// No arguments
	return false
}

// IsMakeByteSliceCall vérifie si c'est un appel à make pour créer []byte.
//
// Params:
//   - call: expression d'appel à vérifier
//
// Returns:
//   - bool: true si c'est make([]byte, ...)
func IsMakeByteSliceCall(call *ast.CallExpr) bool {
	// Check if it's a make call
	if !IsMakeCall(call) {
		// Not a make call
		return false
	}
	// Check if first argument is []byte
	if len(call.Args) > 0 {
		// Check first argument is byte slice
		return IsByteSlice(call.Args[0])
	}
	// No arguments
	return false
}
