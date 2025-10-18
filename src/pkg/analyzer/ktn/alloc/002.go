package ktn_alloc

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
)

// Rule002 vérifie le pattern make([]T, 0) suivi d'append.
// KTN-ALLOC-002: make([]T, 0) + append() force des réallocations coûteuses.
var Rule002 = &analysis.Analyzer{
	Name: "KTN_ALLOC_002",
	Doc:  "Prévenir make([]T, 0) suivi d'append (préallouer la capacité)",
	Run:  runRule002,
}

// runRule002 exécute l'analyse pour la règle ALLOC-002.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule002(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Body == nil {
				return true
			}

			checkMakeAppendPattern(pass, funcDecl)
			return true
		})
	}
	return nil, nil
}

// checkMakeAppendPattern vérifie le pattern make(slice, 0) suivi d'append.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction à analyser
func checkMakeAppendPattern(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	makeSliceVars := make(map[string]token.Pos)

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		trackMakeSliceZeroAssignment(n, makeSliceVars)
		checkAppendOnTrackedSlice(pass, n, makeSliceVars)
		return true
	})
}

// trackMakeSliceZeroAssignment détecte et enregistre les assignations make([]T, 0).
//
// Params:
//   - n: le nœud AST à analyser
//   - makeSliceVars: map des variables trackées avec leur position
func trackMakeSliceZeroAssignment(n ast.Node, makeSliceVars map[string]token.Pos) {
	assignStmt, ok := n.(*ast.AssignStmt)
	if !ok || len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return
	}

	if !utils.IsMakeSliceZero(assignStmt.Rhs[0]) {
		return
	}

	if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
		makeSliceVars[ident.Name] = assignStmt.Pos()
	}
}

// checkAppendOnTrackedSlice vérifie si append() est appelé sur un slice tracké.
//
// Params:
//   - pass: la passe d'analyse
//   - n: le nœud AST à analyser
//   - makeSliceVars: map des variables trackées
func checkAppendOnTrackedSlice(pass *analysis.Pass, n ast.Node, makeSliceVars map[string]token.Pos) {
	callExpr, ok := n.(*ast.CallExpr)
	if !ok {
		return
	}

	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok || ident.Name != "append" {
		return
	}

	if len(callExpr.Args) == 0 {
		return
	}

	firstArg, ok := callExpr.Args[0].(*ast.Ident)
	if !ok {
		return
	}

	makePos, found := makeSliceVars[firstArg.Name]
	if !found {
		return
	}

	reportMakeAppendViolation(pass, makePos, firstArg.Name)
	delete(makeSliceVars, firstArg.Name)
}

// reportMakeAppendViolation rapporte une violation KTN-ALLOC-002.
//
// Params:
//   - pass: la passe d'analyse
//   - pos: position du make()
//   - sliceName: nom du slice
func reportMakeAppendViolation(pass *analysis.Pass, pos token.Pos, sliceName string) {
	pass.Reportf(pos,
		"[KTN-ALLOC-002] Slice '%s' créé avec make([]T, 0) puis utilisé avec append().\n"+
			"Cela force des réallocations coûteuses à chaque append.\n"+
			"Si la taille est connue, spécifiez la capacité.\n"+
			"Exemple:\n"+
			"  // ❌ INEFFICACE\n"+
			"  items := make([]Item, 0)\n"+
			"  for _, v := range source {\n"+
			"      items = append(items, v)  // Réallocation O(log n)\n"+
			"  }\n"+
			"\n"+
			"  // ✅ OPTIMISÉ\n"+
			"  items := make([]Item, 0, len(source))  // Préallocation\n"+
			"  for _, v := range source {\n"+
			"      items = append(items, v)  // Pas de réallocation\n"+
			"  }",
		sliceName)
}
