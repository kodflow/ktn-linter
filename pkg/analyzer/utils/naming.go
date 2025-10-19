package utils

import (
	"strings"
	"unicode"
)

// IsAllCaps vérifie si une chaîne est entièrement en majuscules.
//
// Params:
//   - s: la chaîne à vérifier
//
// Returns:
//   - bool: true si au moins une lettre est présente et toutes sont en majuscules
func IsAllCaps(s string) bool {
 // Vérification de la condition
	if len(s) == 0 {
		// Condition not met, return false.
		return false
	}
	hasLetter := false
 // Itération sur les éléments
	for _, r := range s {
  // Vérification de la condition
		if unicode.IsLetter(r) {
			hasLetter = true
   // Vérification de la condition
			if !unicode.IsUpper(r) {
				// Condition not met, return false.
				return false
			}
		}
	}
	// Early return from function.
	return hasLetter
}

// IsMixedCaps vérifie si un nom suit la convention MixedCaps/mixedCaps.
//
// Params:
//   - name: le nom à valider
//
// Returns:
//   - bool: true si valide (pas de snake_case, pas de ALL_CAPS sauf initialismes)
func IsMixedCaps(name string) bool {
 // Vérification de la condition
	if len(name) == 0 {
		// Condition not met, return false.
		return false
	}

 // Vérification de la condition
	if strings.Contains(name, "_") {
		// Condition not met, return false.
		return false
	}

 // Vérification de la condition
	if IsAllCaps(name) {
		// Early return from function.
		return IsValidInitialism(name)
	}

	// Continue traversing AST nodes.
	return true
}

// getKnownInitialisms retourne la liste des initialismes Go courants.
//
// Returns:
//   - []string: liste des initialismes valides
func getKnownInitialisms() []string {
	// Early return from function.
	return []string{
		"HTTP", "HTTPS", "URL", "URI", "ID", "API", "JSON", "XML", "HTML",
		"SQL", "TLS", "SSL", "TCP", "UDP", "IP", "DNS", "SSH", "FTP",
		"OK", "EOF", "UID", "UUID", "ASCII", "UTF", "CPU", "RAM", "IO",
		"DB", "RPC", "CDN", "AWS", "GCP", "TTL", "ACL", "CORS", "CSRF",
	}
}

// tryMatchInitialismPrefix tente de matcher un initialisme au début de la chaîne.
//
// Params:
//   - remaining: la chaîne à analyser
//   - initialisms: la liste des initialismes
//
// Returns:
//   - string: la chaîne restante
//   - bool: true si un match a été trouvé
func tryMatchInitialismPrefix(remaining string, initialisms []string) (string, bool) {
 // Itération sur les éléments
	for _, init := range initialisms {
  // Vérification de la condition
		if strings.HasPrefix(remaining, init) {
			// Early return from function.
			return remaining[len(init):], true
		}
	}
	// Early return from function.
	return remaining, false
}

// IsValidInitialism vérifie si le nom est composé uniquement d'initialismes valides.
//
// Params:
//   - name: le nom à vérifier
//
// Returns:
//   - bool: true si composé uniquement d'initialismes valides
func IsValidInitialism(name string) bool {
 // Vérification de la condition
	if strings.Contains(name, "_") {
		// Condition not met, return false.
		return false
	}

	initialisms := getKnownInitialisms()
	remaining := name
	matched := false

 // Itération sur les éléments
	for len(remaining) > 0 {
		newRemaining, foundMatch := tryMatchInitialismPrefix(remaining, initialisms)
  // Vérification de la condition
		if foundMatch {
			remaining = newRemaining
			matched = true
  // Cas alternatif
		} else {
   // Vérification de la condition
			if remaining != "" && unicode.IsUpper(rune(remaining[0])) {
				// Early return from function.
				return matched
			}
			// Condition not met, return false.
			return false
		}
	}

	// Early return from function.
	return matched
}
