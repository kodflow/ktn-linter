package rules_var

import "os"

// Variables du package.
var (
	// urlFromEnv lit une URL depuis l'environnement (OK car appel fonction).
	urlFromEnv string = os.Getenv("SERVICE_URL")

	// computedValue est calculé via une fonction (OK car appel fonction).
	computedValue int = len("hello")
)

// Constantes du package.
const (
	// ValidHTTPCode utilise l'initialisme HTTP correctement.
	ValidHTTPCode int = 200

	// MaxHTTPRetries utilise HTTP en majuscules (initialisme valide).
	MaxHTTPRetries int = 5

	// MaxRetriesLiteral est une constante (corrigé depuis var).
	MaxRetriesLiteral int = 3

	// MaxTimeout timeout maximum en secondes.
	MaxTimeout int = 30
)
