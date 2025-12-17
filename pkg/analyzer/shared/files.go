// Package shared provides common utilities for static analysis.
package shared

import (
	"go/ast"
	"strings"
)

// IsTestFile vérifie si un nom de fichier est un fichier de test.
//
// Params:
//   - filename: nom du fichier à vérifier
//
// Returns:
//   - bool: true si le fichier se termine par "_test.go"
func IsTestFile(filename string) bool {
	// Retour de la vérification
	return strings.HasSuffix(filename, "_test.go")
}

// IsTestFunction vérifie si une fonction est une fonction de test.
//
// Params:
//   - funcDecl: déclaration de fonction à vérifier
//
// Returns:
//   - bool: true si c'est une fonction de test (Test, Benchmark, Example, Fuzz)
func IsTestFunction(funcDecl *ast.FuncDecl) bool {
	// Vérification de la condition
	if funcDecl == nil || funcDecl.Name == nil {
		// Pas une fonction valide
		return false
	}

	name := funcDecl.Name.Name
	// Vérifier si le nom commence par Test, Benchmark, Example ou Fuzz
	return strings.HasPrefix(name, "Test") ||
		strings.HasPrefix(name, "Benchmark") ||
		strings.HasPrefix(name, "Example") ||
		strings.HasPrefix(name, "Fuzz")
}

// IsUnitTestFunction vérifie si une fonction est une fonction de test unitaire (Test*).
//
// Params:
//   - funcDecl: déclaration de fonction à vérifier
//
// Returns:
//   - bool: true si c'est une fonction Test* uniquement (pas Benchmark/Example/Fuzz)
func IsUnitTestFunction(funcDecl *ast.FuncDecl) bool {
	// Vérification de la condition
	if funcDecl == nil || funcDecl.Name == nil {
		// Pas une fonction valide
		return false
	}

	// Vérifier uniquement le préfixe "Test"
	return strings.HasPrefix(funcDecl.Name.Name, "Test")
}

// IsExportedFunction vérifie si une fonction est exportée (publique).
//
// Params:
//   - funcDecl: déclaration de fonction à vérifier
//
// Returns:
//   - bool: true si la fonction est exportée
func IsExportedFunction(funcDecl *ast.FuncDecl) bool {
	// Vérification de la condition
	if funcDecl == nil || funcDecl.Name == nil {
		// Pas une fonction valide
		return false
	}

	// Utiliser ast.IsExported pour vérifier
	return ast.IsExported(funcDecl.Name.Name)
}
