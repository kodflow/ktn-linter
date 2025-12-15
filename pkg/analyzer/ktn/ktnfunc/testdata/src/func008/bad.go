// Package func008 contains test cases for KTN rules.
package func008

import "context"

// Delete removes resource but doesn't use all params.
//
// Params:
//   - ctx: context
//   - req: request
//   - resp: response
func Delete(ctx context.Context, req string, resp string) { // want "ctx" "req" "resp"
	// Uses _ = ctx pattern which is now flagged as bad practice
	_ = ctx
	// No-op: intentionally empty
}

// ProcessData doesn't use some params.
//
// Params:
//   - ctx: context
//   - data: données
//   - options: options
//
// Returns:
//   - string: résultat
func ProcessData(ctx context.Context, data string, options map[string]string) string { // want "ctx" "options"
	// Utilise seulement data, ctx et options non utilisés
	return data
}

// PartialIgnore ignore only one unused param.
//
// Params:
//   - a: premier
//   - b: deuxième
//   - c: troisième
//
// Returns:
//   - string: résultat
func PartialIgnore(a string, b string, c string) string { // want "a" "b" "c"
	// Uses _ = b pattern which is now flagged
	_ = b
	// a et c non utilisés et non ignorés!
	return "result"
}

// badHandler représente un gestionnaire qui n'utilise pas ses paramètres.
// Ce type est un exemple de mauvaise pratique pour KTN-FUNC-008.
type badHandler struct{}

// handle implémente une méthode qui n'utilise pas le préfixe _.
// Ceci est un mauvais exemple car les params non utilisés
// doivent utiliser le préfixe _ pour être explicites.
//
// Params:
//   - ctx: context (non utilisé)
//   - data: données (non utilisé)
//
// Returns:
//   - error: nil
func (h *badHandler) handle(ctx context.Context, data string) error { // want "ctx" "data"
	// Ne fait rien mais n'utilise pas le préfixe _ pour les params non utilisés
	return nil
}

// init appelle les méthodes pour éviter dead code.
func init() {
	// Crée une instance et appelle handle
	h := &badHandler{}
	_ = h.handle(nil, "")
}
