// Good examples for the func010 test case.
package func008

import "context"

// AllParamsUsed utilise tous les paramètres.
//
// Params:
//   - ctx: context
//   - name: nom
//   - value: valeur
//
// Returns:
//   - string: résultat
func AllParamsUsed(ctx context.Context, name string, value int) string {
	// Utilise ctx
	_ = ctx.Done()
	// Retourne avec name et value
	return name + string(rune(value))
}

// UnusedWithUnderscore préfixe les params non utilisés.
//
// Params:
//   - _ctx: context (non utilisé)
//   - name: nom
//   - _value: valeur (non utilisée)
//
// Returns:
//   - string: nom
func UnusedWithUnderscore(_ctx context.Context, name string, _value int) string {
	// Retourne uniquement name
	return name
}

// UnusedWithBlankAssign assigne à _ les params non utilisés.
//
// Params:
//   - ctx: context
//   - req: requête
//   - resp: réponse
//
// Returns:
//   - string: résultat
func UnusedWithBlankAssign(ctx context.Context, req string, resp string) string {
	// Ignore explicitement ctx et resp
	_ = ctx
	_ = resp
	// Retourne uniquement req
	return req
}

// MixedApproach mélange les deux approches.
//
// Params:
//   - ctx: context
//   - _unused1: non utilisé
//   - used: utilisé
//   - unused2: non utilisé
//
// Returns:
//   - string: résultat
func MixedApproach(ctx context.Context, _unused1 string, used string, unused2 int) string {
	_ = ctx
	_ = unused2
	// Retourne used
	return used
}

// Handler interface pour les gestionnaires de requêtes.
type Handler interface {
	Handle(ctx context.Context, data string) error
}

// NoOpHandler implémente Handler mais n'utilise pas tous les params.
// C'est un cas valide car les paramètres sont requis par l'interface.
type NoOpHandler struct{}

// Handle implémente Handler.Handle.
// Les paramètres ctx et data sont requis par l'interface mais non utilisés ici.
//
// Params:
//   - _ctx: context (requis par interface)
//   - _data: données (requis par interface)
//
// Returns:
//   - error: nil
func (h *NoOpHandler) Handle(_ctx context.Context, _data string) error {
	// Ne fait rien
	return nil
}
