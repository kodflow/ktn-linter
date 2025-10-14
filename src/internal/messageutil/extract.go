// Package messageutil fournit des utilitaires pour parser et extraire des informations
// des messages d'erreur du linter
package messageutil

import (
	"strings"
)

// ExtractCode extrait le code d'erreur d'un message (ex: "KTN-CONST-001").
//
// Params:
//   - message: le message d'erreur à analyser
//
// Returns:
//   - string: le code d'erreur extrait, ou "UNKNOWN" si aucun code n'est trouvé
func ExtractCode(message string) string {
	if idx := strings.Index(message, "[KTN-"); idx != -1 {
		if end := strings.Index(message[idx:], "]"); end != -1 {
			return message[idx+1 : idx+end]
		}
	}
	return "UNKNOWN"
}

// ExtractMessage extrait le message principal (première ligne après le code).
//
// Params:
//   - message: le message complet à analyser
//
// Returns:
//   - string: le message sans le code d'erreur et sans les lignes suivantes
func ExtractMessage(message string) string {
	if idx := strings.Index(message, "]"); idx != -1 {
		msg := strings.TrimSpace(message[idx+1:])
		// Extraire seulement la première ligne (le message principal)
		if newline := strings.Index(msg, "\n"); newline != -1 {
			return msg[:newline]
		}
		return msg
	}
	return message
}

// ExtractSuggestion extrait l'exemple de code du message.
//
// Params:
//   - message: le message contenant potentiellement un exemple
//
// Returns:
//   - string: l'exemple de code extrait, ou une chaîne vide si aucun exemple n'est trouvé
func ExtractSuggestion(message string) string {
	// Extraire l'exemple de code
	if idx := strings.Index(message, "Exemple:"); idx != -1 {
		suggestion := strings.TrimSpace(message[idx+8:])
		// Nettoyer les lignes vides au début et à la fin
		lines := strings.Split(suggestion, "\n")
		var cleaned []string
		for _, line := range lines {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				cleaned = append(cleaned, line)
			}
		}
		return strings.Join(cleaned, "\n")
	}
	return ""
}

// ExtractConstName extrait le nom de la constante/variable du message.
//
// Params:
//   - message: le message contenant un nom entre quotes simples (ex: "Constante 'MaxValue' sans type")
//
// Returns:
//   - string: le nom extrait, ou "MyConst" par défaut si aucun nom n'est trouvé
func ExtractConstName(message string) string {
	// Chercher entre quotes simples
	start := strings.Index(message, "'")
	if start == -1 {
		return "MyConst"
	}
	end := strings.Index(message[start+1:], "'")
	if end == -1 {
		return "MyConst"
	}
	return message[start+1 : start+1+end]
}

// ExtractType extrait le type d'une suggestion de code.
//
// Params:
//   - suggestion: la suggestion contenant potentiellement un type Go
//
// Returns:
//   - string: le type Go trouvé, ou "int" par défaut si aucun type n'est trouvé
func ExtractType(suggestion string) string {
	words := strings.Fields(suggestion)
	for i, word := range words {
		if isGoType(word) {
			return word
		}
		if word == "<type>" && i > 0 {
			return "int"
		}
	}
	return "int"
}

// isGoType vérifie si un mot est un type Go connu.
//
// Params:
//   - word: le mot à vérifier
//
// Returns:
//   - bool: true si c'est un type Go standard
func isGoType(word string) bool {
	goTypes := map[string]bool{
		"bool": true, "string": true, "int": true,
		"int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"float32": true, "float64": true,
		"byte": true, "rune": true,
		"complex64": true, "complex128": true,
	}
	return goTypes[word]
}
