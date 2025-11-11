package func013

import "context"

// Delete removes resource but doesn't use all params.
//
// Params:
//   - ctx: context
//   - req: request
//   - resp: response
func Delete(ctx context.Context, req string, resp string) {
	// Only ignores ctx, but req and resp are also unused!
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
func ProcessData(ctx context.Context, data string, options map[string]string) string {
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
func PartialIgnore(a string, b string, c string) string {
	_ = b
	// a et c non utilisés et non ignorés!
	return "result"
}
