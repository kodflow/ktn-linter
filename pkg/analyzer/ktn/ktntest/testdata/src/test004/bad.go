// Package test005 provides test cases for table-driven tests.
package test005

import "strings"

// StringLength retourne la longueur d'une string.
//
// Params:
//   - s: string à mesurer
//
// Returns:
//   - int: longueur de la string
func StringLength(s string) int {
	// Retour de la longueur
	return len(s)
}

// IsEmpty vérifie si une string est vide.
//
// Params:
//   - s: string à vérifier
//
// Returns:
//   - bool: true si vide, false sinon
func IsEmpty(s string) bool {
	// Retour du résultat
	return s == ""
}

// ToUpper convertit une string en majuscules.
//
// Params:
//   - s: string à convertir
//
// Returns:
//   - string: string en majuscules
func ToUpper(s string) string {
	// Retour de la conversion
	return strings.ToUpper(s)
}

// Contains vérifie si une string contient une sous-string.
//
// Params:
//   - s: string principale
//   - substr: sous-string à chercher
//
// Returns:
//   - bool: true si substr est dans s
func Contains(s, substr string) bool {
	// Retour du résultat
	return strings.Contains(s, substr)
}

// CountWords compte le nombre de mots dans une string.
//
// Params:
//   - s: string à analyser
//
// Returns:
//   - int: nombre de mots
func CountWords(s string) int {
	// Vérification string vide
	if s == "" {
		// Retour 0
		return 0
	}
	words := strings.Fields(s)
	// Retour du nombre de mots
	return len(words)
}
