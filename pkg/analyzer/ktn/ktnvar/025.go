// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar025 is the rule code for this analyzer
	ruleCodeVar025 string = "KTN-VAR-025"
)

// Analyzer025 checks for patterns that can be replaced with clear() built-in (Go 1.21+)
var Analyzer025 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar025",
	Doc:      "KTN-VAR-025: Utiliser clear() au lieu de boucles range delete/zero",
	Run:      runVar025,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar025 exécute l'analyse KTN-VAR-025.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar025(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar025) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp, ok := inspAny.(*inspector.Inspector)
	// Defensive: ensure inspector is available
	if !ok || insp == nil {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.RangeStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		rangeStmt, ok := n.(*ast.RangeStmt)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar025, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Vérification du pattern clear
		checkClearPattern(pass, rangeStmt)
	})

	// Retour de la fonction
	return nil, nil
}

// checkClearPattern vérifie si une boucle range peut être remplacée par clear().
//
// Params:
//   - pass: contexte d'analyse
//   - rangeStmt: boucle range à vérifier
func checkClearPattern(pass *analysis.Pass, rangeStmt *ast.RangeStmt) {
	// Le corps doit contenir exactement une instruction
	if rangeStmt.Body == nil || len(rangeStmt.Body.List) != 1 {
		// Corps vide ou plusieurs instructions
		return
	}

	// Vérification du pattern map delete
	if checkMapDeletePattern(pass, rangeStmt) {
		// Pattern détecté et signalé
		return
	}

	// Vérification du pattern slice zero
	checkSliceZeroPattern(pass, rangeStmt)
}

// checkMapDeletePattern vérifie le pattern for k := range m { delete(m, k) }.
//
// Params:
//   - pass: contexte d'analyse
//   - rangeStmt: boucle range à vérifier
//
// Returns:
//   - bool: true si pattern détecté
func checkMapDeletePattern(pass *analysis.Pass, rangeStmt *ast.RangeStmt) bool {
	// La clé doit être définie
	if rangeStmt.Key == nil {
		// Pas de clé dans le range
		return false
	}

	keyIdent, ok := rangeStmt.Key.(*ast.Ident)
	// Vérification de l'identifiant de clé
	if !ok {
		// Clé n'est pas un identifiant
		return false
	}

	// Récupération de l'identifiant de la collection
	rangeIdent := getRangeCollectionIdent(rangeStmt.X)
	// Vérification de la collection
	if rangeIdent == nil {
		// Collection n'est pas un identifiant
		return false
	}

	// L'unique instruction doit être un ExprStmt contenant un appel delete
	exprStmt, ok := rangeStmt.Body.List[0].(*ast.ExprStmt)
	// Vérification du type d'instruction
	if !ok {
		// Pas un statement expression
		return false
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	// Vérification de l'appel de fonction
	if !ok {
		// Pas un appel de fonction
		return false
	}

	// Vérification de l'appel delete
	return isDeleteCallWithKeyAndMap(callExpr, keyIdent, rangeIdent, pass)
}

// isDeleteCallWithKeyAndMap vérifie si l'appel est delete(map, key).
//
// Params:
//   - callExpr: expression d'appel
//   - keyIdent: identifiant de la clé de range
//   - mapIdent: identifiant de la map rangée
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si c'est un appel delete correct
func isDeleteCallWithKeyAndMap(
	callExpr *ast.CallExpr,
	keyIdent, mapIdent *ast.Ident,
	pass *analysis.Pass,
) bool {
	funIdent, ok := callExpr.Fun.(*ast.Ident)
	// Vérification du nom de la fonction
	if !ok || funIdent.Name != "delete" {
		// Pas la fonction delete
		return false
	}

	// delete doit avoir 2 arguments
	if len(callExpr.Args) != 2 {
		// Nombre d'arguments incorrect
		return false
	}

	// Premier argument: la map (doit correspondre à celle du range)
	argMapIdent := getRangeCollectionIdent(callExpr.Args[0])
	// Vérification de la map
	if argMapIdent == nil || argMapIdent.Name != mapIdent.Name {
		// Map différente
		return false
	}

	// Deuxième argument: la clé (doit correspondre à celle du range)
	argKeyIdent, ok := callExpr.Args[1].(*ast.Ident)
	// Vérification de la clé
	if !ok || argKeyIdent.Name != keyIdent.Name {
		// Clé différente
		return false
	}

	// Pattern détecté: signaler
	reportClearPattern(pass, callExpr, "map")
	// Pattern détecté
	return true
}

// checkSliceZeroPattern vérifie le pattern for i := range s { s[i] = zeroValue }.
//
// Params:
//   - pass: contexte d'analyse
//   - rangeStmt: boucle range à vérifier
func checkSliceZeroPattern(pass *analysis.Pass, rangeStmt *ast.RangeStmt) {
	// La clé doit être définie (index)
	if rangeStmt.Key == nil {
		// Pas d'index dans le range
		return
	}

	indexIdent, ok := rangeStmt.Key.(*ast.Ident)
	// Vérification de l'identifiant d'index
	if !ok {
		// Index n'est pas un identifiant
		return
	}

	// Récupération de l'identifiant de la slice
	sliceIdent := getRangeCollectionIdent(rangeStmt.X)
	// Vérification de la slice
	if sliceIdent == nil {
		// Collection n'est pas un identifiant
		return
	}

	// L'unique instruction doit être une affectation
	assignStmt, ok := rangeStmt.Body.List[0].(*ast.AssignStmt)
	// Vérification du type d'instruction
	if !ok || assignStmt.Tok.String() != "=" {
		// Pas une affectation simple
		return
	}

	// Vérification de l'affectation s[i] = zeroValue
	checkIndexAssignZero(pass, assignStmt, indexIdent, sliceIdent)
}

// checkIndexAssignZero vérifie si l'affectation est s[i] = zeroValue.
//
// Params:
//   - pass: contexte d'analyse
//   - assignStmt: affectation à vérifier
//   - indexIdent: identifiant de l'index de range
//   - sliceIdent: identifiant de la slice rangée
func checkIndexAssignZero(
	pass *analysis.Pass,
	assignStmt *ast.AssignStmt,
	indexIdent, sliceIdent *ast.Ident,
) {
	// Doit avoir exactement 1 Lhs et 1 Rhs
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		// Affectation multiple
		return
	}

	// Le Lhs doit être un IndexExpr: s[i]
	indexExpr, ok := assignStmt.Lhs[0].(*ast.IndexExpr)
	// Vérification de l'expression d'index
	if !ok {
		// Pas une expression d'index
		return
	}

	// Vérification de la slice indexée
	if !isMatchingSliceIndex(indexExpr, indexIdent, sliceIdent) {
		// Index ou slice ne correspond pas
		return
	}

	// Le Rhs doit être une valeur zéro
	if !isZeroValue(pass, assignStmt.Rhs[0]) {
		// Pas une valeur zéro
		return
	}

	// Pattern détecté: signaler
	reportClearPattern(pass, assignStmt, "slice")
}

// isMatchingSliceIndex vérifie si l'expression est slice[index].
//
// Params:
//   - indexExpr: expression d'index
//   - indexIdent: identifiant de l'index de range
//   - sliceIdent: identifiant de la slice rangée
//
// Returns:
//   - bool: true si correspondance
func isMatchingSliceIndex(indexExpr *ast.IndexExpr, indexIdent, sliceIdent *ast.Ident) bool {
	// Vérification de la slice
	sliceExprIdent := getRangeCollectionIdent(indexExpr.X)
	// La slice doit correspondre
	if sliceExprIdent == nil || sliceExprIdent.Name != sliceIdent.Name {
		// Slice différente
		return false
	}

	// Vérification de l'index
	idxIdent, ok := indexExpr.Index.(*ast.Ident)
	// L'index doit correspondre
	if !ok || idxIdent.Name != indexIdent.Name {
		// Index différent
		return false
	}

	// Correspondance trouvée
	return true
}

// isZeroValue vérifie si l'expression est une valeur zéro.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si valeur zéro
func isZeroValue(pass *analysis.Pass, expr ast.Expr) bool {
	// Vérification du type d'expression
	switch e := expr.(type) {
	// Littéral basique
	case *ast.BasicLit:
		// Zéro numérique ou chaîne vide
		return e.Value == "0" || e.Value == "0.0" || e.Value == `""` || e.Value == "``"
	// Identifiant (false, nil)
	case *ast.Ident:
		// nil ou false
		return e.Name == "nil" || e.Name == "false"
	// Littéral composite vide
	case *ast.CompositeLit:
		// Struct ou slice vide
		return len(e.Elts) == 0
	// Conversion de type
	case *ast.CallExpr:
		// Conversion vers type basique avec zéro
		return isZeroConversion(pass, e)
	}
	// Pas une valeur zéro reconnue
	return false
}

// isZeroConversion vérifie si c'est une conversion vers zéro (ex: string("")).
//
// Params:
//   - pass: contexte d'analyse
//   - callExpr: expression d'appel
//
// Returns:
//   - bool: true si conversion vers zéro
func isZeroConversion(pass *analysis.Pass, callExpr *ast.CallExpr) bool {
	// Doit avoir exactement 1 argument
	if len(callExpr.Args) != 1 {
		// Pas une conversion simple
		return false
	}

	// Le callee doit être un identifiant (type cible)
	funIdent, ok := callExpr.Fun.(*ast.Ident)
	// Doit être un identifiant de type
	if !ok {
		// Pas un identifiant
		return false
	}

	// Guard against nil type info
	if pass.TypesInfo == nil || pass.TypesInfo.Uses == nil {
		// No type information available to confirm conversion
		return false
	}

	// Vérifier que c'est bien une conversion de type (et pas un appel de fonction)
	obj := pass.TypesInfo.Uses[funIdent]
	// Doit être un type
	if obj == nil {
		// Type information missing
		return false
	}
	// Vérification du type de l'objet
	if _, isType := obj.(*types.TypeName); !isType {
		// Not a type conversion
		return false
	}

	// Maintenant seulement: l'argument doit être une valeur zéro
	return isZeroValue(pass, callExpr.Args[0])
}

// getRangeCollectionIdent récupère l'identifiant d'une collection rangée.
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - *ast.Ident: identifiant ou nil
func getRangeCollectionIdent(expr ast.Expr) *ast.Ident {
	// Vérification du type d'expression
	switch e := expr.(type) {
	// Identifiant direct
	case *ast.Ident:
		// Retourne l'identifiant
		return e
	// Expression parenthésée
	case *ast.ParenExpr:
		// Récursion sur l'expression interne
		return getRangeCollectionIdent(e.X)
	}
	// Pas un identifiant simple
	return nil
}

// reportClearPattern signale un pattern clear détecté.
//
// Params:
//   - pass: contexte d'analyse
//   - node: noeud à signaler
//   - collectionType: type de collection (map ou slice)
func reportClearPattern(pass *analysis.Pass, node ast.Node, collectionType string) {
	msg, ok := messages.Get(ruleCodeVar025)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(node.Pos(), "%s: utiliser clear() pour vider une %s",
			ruleCodeVar025, collectionType)
		return
	}
	pass.Reportf(
		node.Pos(),
		"%s: %s",
		ruleCodeVar025,
		msg.Format(config.Get().Verbose, collectionType),
	)
}
