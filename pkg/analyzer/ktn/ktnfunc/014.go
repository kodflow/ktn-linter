// Analyzer 014 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// INITIAL_PRIVATE_FUNCS_CAP initial capacity for private funcs map
	INITIAL_PRIVATE_FUNCS_CAP int = 32
	// INITIAL_CALLED_FUNCS_CAP initial capacity for called funcs map
	INITIAL_CALLED_FUNCS_CAP int = 64
)

// Analyzer014 checks that all private functions are used in production.
var Analyzer014 = &analysis.Analyzer{
	Name:     "ktnfunc014",
	Doc:      "KTN-FUNC-014: fonctions privées non utilisées dans le code de production (code mort)",
	Run:      runFunc014,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// privateFuncInfo stocke les informations sur une fonction privée.
type privateFuncInfo struct {
	name         string
	pos          token.Pos
	receiverType string // vide si fonction, nom du type si méthode
}

// runFunc014 exécute l'analyse KTN-FUNC-014.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc014(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter les fonctions privées et les appels
	privateFuncs := collectPrivateFunctions(pass, insp)
	calledInProduction := collectCalledFunctions(pass, insp)

	// Détecter les fonctions passées comme callbacks
	collectCallbackUsages(pass, insp, privateFuncs, calledInProduction)

	// Reporter les fonctions privées non utilisées
	reportUnusedPrivateFuncs(pass, privateFuncs, calledInProduction)

	return nil, nil
}

// collectPrivateFunctions collecte toutes les fonctions privées du code de production.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspector AST
//
// Returns:
//   - map[string][]*privateFuncInfo: map des fonctions privées par nom
func collectPrivateFunctions(pass *analysis.Pass, insp *inspector.Inspector) map[string][]*privateFuncInfo {
	privateFuncs := make(map[string][]*privateFuncInfo, INITIAL_PRIVATE_FUNCS_CAP)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			return
		}

		// Vérifier et collecter la fonction privée
		info := extractPrivateFuncInfo(funcDecl)
		// Ajouter à la map si valide
		if info != nil {
			privateFuncs[info.name] = append(privateFuncs[info.name], info)
		}
	})

	// Retour de la map des fonctions privées
	return privateFuncs
}

// extractPrivateFuncInfo extrait les infos d'une fonction privée si applicable.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - *privateFuncInfo: infos de la fonction ou nil si non applicable
func extractPrivateFuncInfo(funcDecl *ast.FuncDecl) *privateFuncInfo {
	// Vérifier le nom
	if funcDecl.Name == nil || len(funcDecl.Name.Name) == 0 {
		return nil
	}

	funcName := funcDecl.Name.Name

	// Ignorer les fonctions spéciales Go
	if funcName == "main" || funcName == "init" {
		return nil
	}

	// Vérifier si c'est privé (première lettre en minuscule)
	firstChar := rune(funcName[0])
	// Ignorer les fonctions publiques
	if firstChar < 'a' || firstChar > 'z' {
		return nil
	}

	info := &privateFuncInfo{name: funcName, pos: funcDecl.Pos()}

	// Si c'est une méthode, extraire le type du receiver
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		info.receiverType = extractReceiverType(funcDecl.Recv.List[0].Type)
	}

	// Retour des infos de la fonction privée
	return info
}

// collectCalledFunctions collecte les fonctions appelées dans le code de production.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspector AST
//
// Returns:
//   - map[string]bool: map des fonctions appelées
func collectCalledFunctions(pass *analysis.Pass, insp *inspector.Inspector) map[string]bool {
	calledInProduction := make(map[string]bool, INITIAL_CALLED_FUNCS_CAP)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)
		filename := pass.Fset.Position(callExpr.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			return
		}

		// Marquer la fonction comme appelée
		if funcName := extractCalledFuncName(callExpr); funcName != "" {
			calledInProduction[funcName] = true
		}
	})

	// Retour de la map des fonctions appelées
	return calledInProduction
}

// collectCallbackUsages détecte les fonctions passées comme callbacks.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspector AST
//   - privateFuncs: map des fonctions privées
//   - calledInProduction: map des fonctions appelées (modifiée)
func collectCallbackUsages(pass *analysis.Pass, insp *inspector.Inspector, privateFuncs map[string][]*privateFuncInfo, calledInProduction map[string]bool) {
	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.ValueSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			return
		}

		// Parcourir tous les identifiants
		ast.Inspect(n, func(node ast.Node) bool {
			// Vérifier si c'est un identifiant
			if ident, ok := node.(*ast.Ident); ok {
				// Vérifier si c'est une fonction privée connue
				if _, exists := privateFuncs[ident.Name]; exists {
					calledInProduction[ident.Name] = true
				}
			}
			// Continuer la traversée
			return true
		})
	})
}

// reportUnusedPrivateFuncs reporte les fonctions privées non utilisées.
//
// Params:
//   - pass: contexte d'analyse
//   - privateFuncs: map des fonctions privées
//   - calledInProduction: map des fonctions appelées
func reportUnusedPrivateFuncs(pass *analysis.Pass, privateFuncs map[string][]*privateFuncInfo, calledInProduction map[string]bool) {
	// Parcours des fonctions privées
	for key, infos := range privateFuncs {
		// Vérifier si la fonction est appelée
		if calledInProduction[key] {
			continue
		}

		// Reporter toutes les fonctions avec ce nom
		for _, info := range infos {
			reportUnusedFunc(pass, info)
		}
	}
}

// reportUnusedFunc reporte une fonction privée non utilisée.
//
// Params:
//   - pass: contexte d'analyse
//   - info: infos de la fonction
func reportUnusedFunc(pass *analysis.Pass, info *privateFuncInfo) {
	// Vérifier si c'est une méthode
	if info.receiverType != "" {
		pass.Reportf(info.pos,
			"KTN-FUNC-014: la méthode privée '%s.%s' n'est jamais appelée dans le code de production. Si elle n'est utilisée que dans les tests, c'est du code mort créé pour contourner les règles - supprimez-la",
			info.receiverType, info.name)
		// Fin du traitement pour les méthodes
		return
	}

	pass.Reportf(info.pos,
		"KTN-FUNC-014: la fonction privée '%s' n'est jamais appelée dans le code de production. Si elle n'est utilisée que dans les tests, c'est du code mort créé pour contourner les règles - supprimez-la",
		info.name)
}

// extractReceiverType extrait le nom du type du receiver.
//
// Params:
//   - expr: expression du type
//
// Returns:
//   - string: nom du type
func extractReceiverType(expr ast.Expr) string {
	// Gérer les pointeurs (*Type)
	if starExpr, ok := expr.(*ast.StarExpr); ok {
		// Récursion sur l'expression pointée
		return extractReceiverType(starExpr.X)
	}

	// Gérer les identifiants simples (Type)
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour du nom
		return ident.Name
	}

	// Type non géré
	return ""
}

// extractCalledFuncName extrait le nom de la fonction appelée.
//
// Params:
//   - callExpr: expression d'appel
//
// Returns:
//   - string: nom de la fonction
func extractCalledFuncName(callExpr *ast.CallExpr) string {
	// Vérifier le type de l'appelant
	switch fun := callExpr.Fun.(type) {
	// Appel direct de fonction (funcName())
	case *ast.Ident:
		// Retour du nom de la fonction
		return fun.Name

	// Appel de méthode (obj.method())
	case *ast.SelectorExpr:
		// For methods, return method name only
		// Any method with this name is considered used
		return fun.Sel.Name

	// Autres types d'appels
	default:
		// Non géré
		return ""
	}
}
