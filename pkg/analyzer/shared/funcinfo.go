// Package shared provides common utilities for static analysis.
package shared

import "go/token"

// FuncInfo stocke les informations d'une fonction avec sa position.
// Peut être utilisé pour fonctions publiques, méthodes, etc.
type FuncInfo struct {
	Name     string
	Pos      token.Pos
	Filename string
}
