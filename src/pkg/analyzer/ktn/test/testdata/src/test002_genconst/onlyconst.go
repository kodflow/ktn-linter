package test002_genconst

// Fichier avec seulement const - pas d'éléments testables
// isTestableType002 retourne false à la ligne 141 car GenDecl.Tok != "type"
const (
	MaxRetries = 3
	Timeout    = 30
)
