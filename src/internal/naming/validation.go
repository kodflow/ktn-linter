// Package naming fournit des utilitaires pour valider les noms selon les conventions Go
package naming

import (
	"strings"
	"unicode"
)

// IsAllCaps vérifie si une chaîne est entièrement en majuscules
// Retourne true si au moins une lettre est présente et toutes sont en majuscules
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

// IsValidInitialism vérifie si le nom est composé uniquement d'initialismes Go valides
// Exemples valides: HTTP, HTTPS, URL, HTTPOK, URLID, APIURL, HTTPSPort
// Exemples invalides: MAX_BUFFER, HTTP_OK (contiennent des underscores)
//
// Cette fonction suit les conventions Go pour les initialismes courants
// (voir Effective Go et le guide de style de la communauté Go)
func IsValidInitialism(name string) bool {
	// Liste des initialismes Go courants (voir Effective Go)
	initialisisms := []string{
		"HTTP", "HTTPS", "URL", "URI", "ID", "API", "JSON", "XML", "HTML",
		"SQL", "TLS", "SSL", "TCP", "UDP", "IP", "DNS", "SSH", "FTP",
		"OK", "EOF", "UID", "UUID", "ASCII", "UTF", "CPU", "RAM", "IO",
		"DB", "RPC", "CDN", "AWS", "GCP", "TTL", "ACL", "CORS", "CSRF",
	}

	// Si le nom contient un underscore, c'est invalide (viole KTN-VAR-008 / KTN-CONST-008)
	if strings.Contains(name, "_") {
		return false
	}

	// Essayer de décomposer le nom en initialismes connus
	remaining := name
	matched := false

	for len(remaining) > 0 {
		foundMatch := false
		// Essayer de matcher le début avec un initialisme
		for _, init := range initialisisms {
			if strings.HasPrefix(remaining, init) {
				remaining = remaining[len(init):]
				foundMatch = true
				matched = true
				break
			}
		}

		// Si on n'a pas trouvé de match et qu'il reste des caractères
		if !foundMatch {
			// Vérifier si le reste est en MixedCaps (ex: HTTPServer, URLParser)
			if remaining != "" && unicode.IsUpper(rune(remaining[0])) {
				// C'est peut-être une combinaison initialism+nom (HTTPOK, HTTPNotFound)
				// On accepte si au moins un initialisme a été trouvé
				return matched
			}
			return false
		}
	}

	return matched
}
