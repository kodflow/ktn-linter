// Analyzer 002 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer002 checks that package-level variables have explicit types
var Analyzer002 = &analysis.Analyzer{
	Name:     "ktnvar002",
	Doc:      "KTN-VAR-002: Vérifie que les variables de package ont un type explicite",
	Run:      runVar002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar002 exécute l'analyse KTN-VAR-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar002(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for File nodes to access package-level declarations
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Check package-level declarations only
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Skip if not a GenDecl
			if !ok {
				// Not a general declaration
				continue
			}

			// Only check var declarations
			if genDecl.Tok != token.VAR {
				// Continue traversing AST nodes.
				continue
			}

			// Itération sur les spécifications
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				// Vérifier si le type est explicite ou visible dans la valeur
				checkVarSpec(pass, valueSpec)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}

// checkVarSpec vérifie une spécification de variable.
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
func checkVarSpec(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	hasExplicitType := valueSpec.Type != nil
	hasVisible := hasVisibleType(valueSpec.Values)

	// Cas 1: Type explicite + type visible = redondant
	if hasExplicitType && hasVisible {
		// Parcourir les noms
		for _, name := range valueSpec.Names {
			pass.Reportf(
				valueSpec.Type.Pos(),
				"KTN-VAR-002: la variable '%s' a un type redondant (le type est déjà visible dans la valeur)",
				name.Name,
			)
		}
		return
	}

	// Cas 2: Pas de type explicite + pas de type visible = erreur
	if !hasExplicitType && !hasVisible {
		// Parcourir les noms
		for _, name := range valueSpec.Names {
			pass.Reportf(
				name.Pos(),
				"KTN-VAR-002: la variable '%s' doit avoir un type explicite",
				name.Name,
			)
		}
	}
	// Cas 3: Type explicite sans type visible = OK
	// Cas 4: Pas de type explicite mais type visible = OK
}

// hasVisibleType vérifie si le type est visible dans les expressions.
// Exemples de types visibles:
//   - Composite literal: []string{}, map[K]V{}, Struct{}
//   - Unary &: &Struct{}
//   - Call make/new: make([]int, 10), new(Struct)
//
// Params:
//   - values: expressions de valeurs
//
// Returns:
//   - bool: true si type visible
func hasVisibleType(values []ast.Expr) bool {
	// Pas de valeurs = type non visible
	if len(values) == 0 {
		return false
	}

	// Vérifier chaque valeur
	for _, val := range values {
		// Vérifier si le type est visible
		if !isTypeVisible(val) {
			return false
		}
	}
	// Toutes les valeurs ont un type visible
	return true
}

// isTypeVisible vérifie si le type est visible dans une expression.
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si type visible
func isTypeVisible(expr ast.Expr) bool {
	// Switch sur le type d'expression
	switch e := expr.(type) {
	// Composite literal: []string{}, map[K]V{}, Struct{}
	case *ast.CompositeLit:
		return true
	// Unary &: &Struct{}, &value
	case *ast.UnaryExpr:
		// &CompositeLit = type visible
		if e.Op == token.AND {
			_, ok := e.X.(*ast.CompositeLit)
			// Retour du résultat
			return ok
		}
		// Pas un & sur composite literal
		return false
	// Appel de fonction: make(), new(), Type()
	case *ast.CallExpr:
		// Déléguer à isTypedCall
		return isTypedCall(e)
	// Pas un type visible
	default:
		return false
	}
}

// isTypedCall vérifie si un appel a un type visible.
// Détecte: make(T), new(T), Type(x), pkg.Type(x)
//
// Params:
//   - call: expression d'appel
//
// Returns:
//   - bool: true si type visible
func isTypedCall(call *ast.CallExpr) bool {
	// Vérifier le type de fonction
	switch fn := call.Fun.(type) {
	// make(T, ...) ou new(T) ou type conversion: int(x), string(y)
	case *ast.Ident:
		// Vérifier si c'est make, new ou un type de base
		return isBuiltinOrTypeConversion(fn.Name)
	// Type conversion: []byte(x), map[K]V(x), chan T(x)
	case *ast.ArrayType, *ast.MapType, *ast.ChanType, *ast.StructType:
		// Conversion vers type composite
		return true
	// pkg.Type(x) ou Type(x) via selector
	case *ast.SelectorExpr:
		// Appel qualifié = type visible
		return true
	// Pas un appel typé connu
	default:
		// Type inconnu
		return false
	}
}

// isBuiltinOrTypeConversion vérifie si un identifiant est make/new ou un type.
//
// Params:
//   - name: nom de l'identifiant
//
// Returns:
//   - bool: true si c'est make, new ou un type de base
func isBuiltinOrTypeConversion(name string) bool {
	// Liste des builtins et types de base
	switch name {
	// Builtins qui créent des types
	case "make", "new":
		return true
	// Types de base Go (conversions)
	case "bool", "byte", "rune", "string",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "complex64", "complex128":
		return true
	// Pas un builtin ou type connu
	default:
		return false
	}
}
