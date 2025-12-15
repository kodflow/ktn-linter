// Package func008 provides good test cases.
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

// MixedApproach utilise le préfixe _ pour tous les non-utilisés.
//
// Params:
//   - _ctx: context (non utilisé)
//   - _unused1: non utilisé
//   - used: utilisé
//   - _unused2: non utilisé
//
// Returns:
//   - string: résultat
func MixedApproach(_ctx context.Context, _unused1 string, used string, _unused2 int) string {
	// Retourne used
	return used
}

// Handler interface pour les gestionnaires de requêtes.
type Handler interface {
	Handle(ctx context.Context, data string) error
}

// ProcessHandler traite un handler.
//
// Params:
//   - h: handler à traiter
//
// Returns:
//   - error: erreur éventuelle
func ProcessHandler(h Handler) error {
	// Retourne nil si succès
	return h.Handle(context.Background(), "")
}

// NoOpHandler implémente Handler mais n'utilise pas tous les params.
// C'est un cas valide car les paramètres sont requis par l'interface.
type NoOpHandler struct{}

// NewNoOpHandler crée une nouvelle instance de NoOpHandler.
//
// Returns:
//   - *NoOpHandler: instance qui implémente l'interface Handler
func NewNoOpHandler() *NoOpHandler {
	// Retourne une nouvelle instance
	return &NoOpHandler{}
}

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
