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
	if len(s) == 0 {
		return false
	}
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
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
	if len(name) == 0 {
		return false
	}

	if strings.Contains(name, "_") {
		return false
	}

	if IsAllCaps(name) {
		return IsValidInitialism(name)
	}

	return true
}

// getKnownInitialisms retourne la liste des initialismes Go courants.
//
// Returns:
//   - []string: liste des initialismes valides
func getKnownInitialisms() []string {
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
	for _, init := range initialisms {
		if strings.HasPrefix(remaining, init) {
			return remaining[len(init):], true
		}
	}
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
	if strings.Contains(name, "_") {
		return false
	}

	initialisms := getKnownInitialisms()
	remaining := name
	matched := false

	for len(remaining) > 0 {
		newRemaining, foundMatch := tryMatchInitialismPrefix(remaining, initialisms)
		if foundMatch {
			remaining = newRemaining
			matched = true
		} else {
			if remaining != "" && unicode.IsUpper(rune(remaining[0])) {
				return matched
			}
			return false
		}
	}

	return matched
}
