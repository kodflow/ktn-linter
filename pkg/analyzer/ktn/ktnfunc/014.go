package ktnfunc

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer014 vérifie que toutes les fonctions privées sont utilisées dans le code de production.
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
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Map pour stocker toutes les fonctions privées déclarées (nom -> liste d'infos)
	privateFuncs := make(map[string][]*privateFuncInfo)

	// Map pour stocker les fonctions appelées dans le code de production
	calledInProduction := make(map[string]bool)

	// Étape 1: Collecter toutes les fonctions privées
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Récupérer le nom du fichier
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Fichier de test - ignorer
			return
		}

		// Vérifier si c'est une fonction privée (commence par minuscule)
		if funcDecl.Name == nil || len(funcDecl.Name.Name) == 0 {
			// Pas de nom
			return
		}

		funcName := funcDecl.Name.Name

		// Ignorer les fonctions spéciales Go (main et init sont appelées implicitement)
		if funcName == "main" || funcName == "init" {
			// Fonction spéciale - ignorer
			return
		}

		// Vérifier si c'est privé (première lettre en minuscule)
		firstChar := rune(funcName[0])
		if firstChar >= 'a' && firstChar <= 'z' {
			// Fonction privée
			info := &privateFuncInfo{
				name: funcName,
				pos:  funcDecl.Pos(),
			}

			// Si c'est une méthode, extraire le type du receiver pour le message d'erreur
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				info.receiverType = extractReceiverType(funcDecl.Recv.List[0].Type)
			}

			// Utiliser uniquement le nom comme clé (simplifié)
			// Cela signifie que si une méthode compute() est appelée,
			// TOUTES les méthodes compute() sont considérées comme utilisées
			key := info.name

			// Ajouter à la map de slices
			privateFuncs[key] = append(privateFuncs[key], info)
		}
	})

	// Étape 2: Parcourir tout le code de production pour trouver les appels
	nodeFilter2 := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter2, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Récupérer le nom du fichier
		filename := pass.Fset.Position(callExpr.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Fichier de test - ignorer
			return
		}

		// Extraire le nom de la fonction appelée
		funcName := extractCalledFuncName(callExpr)
		// Vérification du nom
		if funcName != "" {
			// Marquer comme utilisée
			calledInProduction[funcName] = true
		}
	})

	// Étape 2b: Détecter les fonctions passées comme valeurs (callbacks)
	nodeFilter3 := []ast.Node{
		(*ast.CompositeLit)(nil), // struct{Field: func}
		(*ast.AssignStmt)(nil),   // var = func
		(*ast.ValueSpec)(nil),    // var x = func
	}

	inspect.Preorder(nodeFilter3, func(n ast.Node) {
		// Récupérer le nom du fichier
		filename := pass.Fset.Position(n.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Fichier de test - ignorer
			return
		}

		// Parcourir tous les identifiants dans ce nœud
		ast.Inspect(n, func(node ast.Node) bool {
			// Si c'est un identifiant
			if ident, ok := node.(*ast.Ident); ok {
				// Vérifier si c'est une fonction privée qu'on connaît
				if _, exists := privateFuncs[ident.Name]; exists {
					// Marquer comme utilisée (passée comme valeur/callback)
					calledInProduction[ident.Name] = true
				}
			}
			// Continuer la traversée
			return true
		})
	})

	// Étape 3: Vérifier les fonctions privées non utilisées
	for key, infos := range privateFuncs {
		// Vérifier si la fonction est appelée dans le code de production
		if !calledInProduction[key] {
			// Toutes les fonctions avec ce nom sont non utilisées
			for _, info := range infos {
				// Fonction privée non utilisée
				if info.receiverType != "" {
					// Méthode privée
					pass.Reportf(
						info.pos,
						"KTN-FUNC-014: la méthode privée '%s.%s' n'est jamais appelée dans le code de production. Si elle n'est utilisée que dans les tests, c'est du code mort créé pour contourner les règles - supprimez-la",
						info.receiverType,
						info.name,
					)
				} else {
					// Fonction privée
					pass.Reportf(
						info.pos,
						"KTN-FUNC-014: la fonction privée '%s' n'est jamais appelée dans le code de production. Si elle n'est utilisée que dans les tests, c'est du code mort créé pour contourner les règles - supprimez-la",
						info.name,
					)
				}
			}
		}
	}

	// Retour de la fonction
	return nil, nil
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
		// Pour les méthodes, on retourne juste le nom de la méthode
		// Cela signifie que si TOUTE méthode avec ce nom est appelée,
		// on considère qu'elle est utilisée (moins précis mais évite les faux positifs)
		return fun.Sel.Name

	// Autres types d'appels
	default:
		// Non géré
		return ""
	}
}
