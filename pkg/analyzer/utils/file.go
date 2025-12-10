// Utility functions for file detection.
package utils

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// IsTestFile checks if a file is a test file (*_test.go).
//
// Params:
//   - filename: the file name to check
//
// Returns:
//   - bool: true if the file is a test file
func IsTestFile(filename string) bool {
	// Vérification du suffixe _test.go
	return strings.HasSuffix(filename, "_test.go")
}

// IsGeneratedFile checks if a file is generated (contains "Code generated" comment).
//
// Params:
//   - file: the AST file to check
//
// Returns:
//   - bool: true if the file is generated
func IsGeneratedFile(file *ast.File) bool {
	// Parcours des groupes de commentaires
	for _, cg := range file.Comments {
		// Parcours des commentaires individuels
		for _, c := range cg.List {
			// Vérification du marqueur Code generated
			if strings.Contains(c.Text, "Code generated") {
				// Fichier généré détecté
				return true
			}
			// Vérification du marqueur DO NOT EDIT
			if strings.Contains(c.Text, "DO NOT EDIT") {
				// Fichier généré détecté
				return true
			}
		}
	}
	// Fichier non généré
	return false
}

// ShouldSkipFile checks if a file should be skipped for analysis.
// Skips test files and generated files.
//
// Params:
//   - pass: analysis pass
//   - file: the AST file to check
//
// Returns:
//   - bool: true if the file should be skipped
func ShouldSkipFile(pass *analysis.Pass, file *ast.File) bool {
	// Récupération du nom de fichier
	filename := pass.Fset.Position(file.Pos()).Filename

	// Vérification fichier de test
	if IsTestFile(filename) {
		// Fichier de test à ignorer
		return true
	}

	// Vérification fichier généré
	if IsGeneratedFile(file) {
		// Fichier généré à ignorer
		return true
	}

	// Fichier à analyser
	return false
}

// ShouldSkipTestFile checks if a test file should be skipped.
// Only skips test files, not generated files.
//
// Params:
//   - pass: analysis pass
//   - file: the AST file to check
//
// Returns:
//   - bool: true if the file is a test file
func ShouldSkipTestFile(pass *analysis.Pass, file *ast.File) bool {
	filename := pass.Fset.Position(file.Pos()).Filename
	// Retour du résultat de la vérification
	return IsTestFile(filename)
}

// ShouldSkipGeneratedFile checks if a generated file should be skipped.
// Only skips generated files, not test files.
//
// Params:
//   - file: the AST file to check
//
// Returns:
//   - bool: true if the file is generated
func ShouldSkipGeneratedFile(file *ast.File) bool {
	// Retour du résultat de la vérification
	return IsGeneratedFile(file)
}
