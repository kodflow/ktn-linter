package astutil

import (
	"go/ast"
	"strings"
)

// IsConstCompatibleType vérifie si le type est compatible avec const
// Retourne true pour les types de base : bool, string, int*, uint*, float*, complex*, byte, rune
func IsConstCompatibleType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Types de base compatibles avec const
		switch e.Name {
		case "bool", "string",
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64",
			"complex64", "complex128",
			"byte", "rune":
			return true
		}
	}
	return false
}

// IsLiteralValue vérifie si l'expression est une valeur littérale
// Retourne true pour les BasicLit (nombres, strings) et Ident (true, false, nil)
func IsLiteralValue(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.BasicLit: // true, false, 123, "string", 3.14, etc.
		return true
	case *ast.Ident: // true, false, nil
		return true
	}
	return false
}

// LooksLikeConstantName vérifie si le nom ressemble à une constante mathématique ou scientifique
// Ne signale que les cas évidents (Pi, E, Euler, etc.) pour éviter les faux positifs
func LooksLikeConstantName(name string) bool {
	// Liste de noms connus de constantes mathématiques/scientifiques
	knownConstants := map[string]bool{
		"Pi":             true,
		"E":              true,
		"Euler":          true,
		"EulerNumber":    true,
		"GoldenRatio":    true,
		"Phi":            true,
		"Tau":            true,
		"SpeedOfLight":   true,
		"PlanckConstant": true,
		"AvogadroNumber": true,
		"BoltzmannConst": true,
		"GravityConst":   true,
	}

	// Vérifier si c'est un nom connu
	if knownConstants[name] {
		return true
	}

	// Vérifier si le nom contient des indicateurs de constante mathématique
	nameLower := strings.ToLower(name)
	if strings.Contains(nameLower, "constant") ||
		strings.Contains(nameLower, "ratio") && strings.Contains(nameLower, "golden") {
		return true
	}

	return false
}
