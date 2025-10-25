package shared

// MethodSignature représente la signature d'une méthode.
// Utilisé pour comparer les signatures entre struct et interface.
type MethodSignature struct {
	Name       string
	ParamsStr  string
	ResultsStr string
}
