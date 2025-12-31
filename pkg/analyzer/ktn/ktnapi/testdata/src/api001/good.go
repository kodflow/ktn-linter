// Package api001 contains test cases for KTN-API-001.
package api001

import (
	"net/http"
)

// httpDoer is a minimal consumer-side interface.
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// GoodWithInterface uses an interface parameter - no warning.
//
// Params:
//   - d: client HTTP
//
// Returns:
//   - *http.Response: réponse HTTP
//   - error: erreur éventuelle
func GoodWithInterface(d httpDoer) (*http.Response, error) {
	// Création de la requête
	req, err := http.NewRequest("GET", "http://example.com", nil)
	// Vérification de l'erreur de création
	if err != nil {
		// Retour de l'erreur
		return nil, err
	}

	// Exécution via interface
	return d.Do(req)
}

// GoodWithFieldAccess accesses a field on external type.
// No warning because no methods are called on the parameter directly.
//
// Params:
//   - req: requête HTTP
//
// Returns:
//   - string: méthode HTTP
func GoodWithFieldAccess(req *http.Request) string {
	// Only field access, no method calls on req
	// Retour de la méthode HTTP
	return req.Method
}

// init utilise les fonctions privées
func init() {
	// Appel GoodWithFieldAccess
	_ = GoodWithFieldAccess(&http.Request{Method: "GET"})
}
