// Package naming fournit des utilitaires pour valider les noms selon les conventions Go
package naming

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
		// Retourne false car une chaîne vide n'est pas ALL_CAPS
		return false
	}
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				// Retourne false car une lettre minuscule a été trouvée
				return false
			}
		}
	}
	// Retourne true si au moins une lettre en majuscule a été trouvée
	return hasLetter
}

// IsMixedCaps vérifie si un nom suit la convention MixedCaps/mixedCaps.
//
// Params:
//   - name: le nom à valider
//
// Returns:
//   - bool: true si valide (pas de snake_case, pas de ALL_CAPS sauf initialismes). Exemples valides: ParseHTTPRequest, calculateTotal, HTTPServer, UserID. Exemples invalides: parse_http_request, Calculate_Total, PARSE_REQUEST
func IsMixedCaps(name string) bool {
	// Vide invalide
	if len(name) == 0 {
		// Retourne false car un nom vide est invalide
		return false
	}

	// Contient underscore = invalide (sauf initialismes sans underscore)
	if strings.Contains(name, "_") {
		// Retourne false car le nom contient un underscore (snake_case)
		return false
	}

	// Si tout est en majuscules, vérifier si c'est un initialisme valide
	if IsAllCaps(name) {
		// Retourne le résultat de la validation d'initialisme
		return IsValidInitialism(name)
	}

	// Sinon, c'est valide si ça ne contient pas d'underscore
	// Retourne true car le nom est en MixedCaps sans underscore
	return true
}

// HasGetterPrefix vérifie si un nom de fonction a un préfixe "Get" inutile.
//
// Params:
//   - name: le nom de fonction à vérifier
//
// Returns:
//   - bool: true si le nom commence par "Get" et ce n'est pas une exception. Exemples à signaler: GetUserName, GetEmail, GetHTTPClient. Exceptions acceptées: GetOrCreate, GetAndSet (verbe composé)
func HasGetterPrefix(name string) bool {
	if !strings.HasPrefix(name, "Get") {
		// Retourne false car le nom ne commence pas par Get
		return false
	}

	// Exceptions: verbes composés avec Get
	exceptions := []string{
		"GetOrCreate", "GetOrSet", "GetOrDefault",
		"GetAndSet", "GetAndUpdate", "GetAndDelete",
		"GetOrElse", "GetOrInsert",
	}

	for _, ex := range exceptions {
		if name == ex || strings.HasPrefix(name, ex) {
			// Retourne false car c'est une exception autorisée (verbe composé)
			return false
		}
	}

	// Si ça commence par Get suivi d'une majuscule, c'est un getter
	if len(name) > 3 && unicode.IsUpper(rune(name[3])) {
		// Retourne true car c'est un getter inutile (ex: GetUserName)
		return true
	}

	// Retourne false car Get n'est pas suivi d'une majuscule
	return false
}

// getInitialismsMap retourne le map des initialismes incorrects vers corrects.
//
// Returns:
//   - map[string]string: map des initialismes avec forme incorrecte -> forme correcte
func getInitialismsMap() map[string]string {
	// Retourne le map des initialismes incorrects vers leur forme correcte
	return map[string]string{
		"Http": "HTTP", "Https": "HTTPS", "Url": "URL", "Uri": "URI",
		"Id": "ID", "Api": "API", "Json": "JSON", "Xml": "XML",
		"Html": "HTML", "Sql": "SQL", "Tls": "TLS", "Ssl": "SSL",
		"Tcp": "TCP", "Udp": "UDP", "Ip": "IP", "Dns": "DNS",
		"Ssh": "SSH", "Ftp": "FTP", "Ok": "OK", "Eof": "EOF",
		"Uid": "UID", "Uuid": "UUID", "Ascii": "ASCII", "Utf": "UTF",
		"Cpu": "CPU", "Ram": "RAM", "Io": "IO", "Db": "DB",
		"Rpc": "RPC", "Cdn": "CDN", "Aws": "AWS", "Gcp": "GCP",
		"Ttl": "TTL", "Acl": "ACL", "Cors": "CORS", "Csrf": "CSRF",
	}
}

// FixInitialisms trouve les initialismes incorrects dans un nom.
//
// Params:
//   - name: le nom contenant potentiellement des initialismes incorrects
//
// Returns:
//   - []string: une liste de corrections suggérées (maximum 1 suggestion avec toutes les corrections). Exemples: "HttpServer" -> ["HTTPServer"], "UrlParser" -> ["URLParser"]
func FixInitialisms(name string) []string {
	fixed := name
	hasChanges := false

	for incorrect, correct := range getInitialismsMap() {
		if strings.Contains(fixed, incorrect) {
			fixed = strings.ReplaceAll(fixed, incorrect, correct)
			hasChanges = true
		}
	}

	if hasChanges && fixed != name {
		// Retourne une suggestion avec tous les initialismes corrigés
		return []string{fixed}
	}
	// Retourne une liste vide car aucune correction n'est nécessaire
	return []string{}
}

// getKnownInitialisms retourne la liste des initialismes Go courants.
//
// Returns:
//   - []string: liste des initialismes valides selon Effective Go
func getKnownInitialisms() []string {
	// Retourne la liste complète des initialismes reconnus par Go
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
//   - initialisms: la liste des initialismes à tester
//
// Returns:
//   - string: la chaîne restante après le match
//   - bool: true si un match a été trouvé
func tryMatchInitialismPrefix(remaining string, initialisms []string) (string, bool) {
	for _, init := range initialisms {
		if strings.HasPrefix(remaining, init) {
			// Retourne la chaîne restante et true car un initialisme a été trouvé
			return remaining[len(init):], true
		}
	}
	// Retourne la chaîne inchangée et false car aucun initialisme ne correspond
	return remaining, false
}

// IsValidInitialism vérifie si le nom est composé uniquement d'initialismes Go valides.
//
// Params:
//   - name: le nom à vérifier
//
// Returns:
//   - bool: true si composé uniquement d'initialismes valides. Exemples valides: HTTP, HTTPS, URL, HTTPOK, URLID, APIURL, HTTPSPort. Exemples invalides: MAX_BUFFER, HTTP_OK (contiennent des underscores). Cette fonction suit les conventions Go pour les initialismes courants (voir Effective Go et le guide de style de la communauté Go)
func IsValidInitialism(name string) bool {
	// Si le nom contient un underscore, c'est invalide (viole KTN-VAR-008 / KTN-CONST-008)
	if strings.Contains(name, "_") {
		// Retourne false car les underscores violent les conventions Go
		return false
	}

	// Essayer de décomposer le nom en initialismes connus
	initialisms := getKnownInitialisms()
	remaining := name
	matched := false

	for len(remaining) > 0 {
		newRemaining, foundMatch := tryMatchInitialismPrefix(remaining, initialisms)
		if foundMatch {
			remaining = newRemaining
			matched = true
		} else {
			// Vérifier si le reste est en MixedCaps (ex: HTTPServer, URLParser)
			if remaining != "" && unicode.IsUpper(rune(remaining[0])) {
				// Retourne true si au moins un initialisme a été trouvé avant le MixedCaps
				return matched
			}
			// Retourne false car le reste n'est pas un initialisme valide
			return false
		}
	}

	// Retourne true si tout le nom a été décomposé en initialismes valides
	return matched
}
