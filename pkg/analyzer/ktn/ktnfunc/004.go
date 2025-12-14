// Analyzer 004 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeFunc004 is the rule code for this analyzer
	ruleCodeFunc004 string = "KTN-FUNC-004"
	// initialPrivateFuncsCap initial capacity for private funcs map
	initialPrivateFuncsCap int = 32
	// initialCalledFuncsCap initial capacity for called funcs map
	initialCalledFuncsCap int = 64
)

// Analyzer004 checks that all private functions are used in production.
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc004",
	Doc:      "KTN-FUNC-004: fonctions privées non utilisées dans le code de production (code mort)",
	Run:      runFunc004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// privateFuncInfo stocke les informations sur une fonction privée.
type privateFuncInfo struct {
	name         string
	pos          token.Pos
	receiverType string // vide si fonction, nom du type si méthode
}

// runFunc004 exécute l'analyse KTN-FUNC-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc004(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeFunc004) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter les fonctions privées et les appels
	privateFuncs := collectPrivateFunctions(pass, insp, cfg)
	calledInProduction := collectCalledFunctions(pass, insp, cfg)

	// Détecter les fonctions passées comme callbacks
	collectCallbackUsages(pass, insp, privateFuncs, calledInProduction, cfg)

	// Reporter les fonctions privées non utilisées
	reportUnusedPrivateFuncs(pass, privateFuncs, calledInProduction)

	// Retour succès de l'analyse
	return nil, nil
}

// collectPrivateFunctions collecte toutes les fonctions privées du code de production.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspector AST
//   - cfg: configuration
//
// Returns:
//   - map[string][]*privateFuncInfo: map des fonctions privées par nom
func collectPrivateFunctions(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) map[string][]*privateFuncInfo {
	privateFuncs := make(map[string][]*privateFuncInfo, initialPrivateFuncsCap)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeFunc004, filename) {
			// Fichier exclu
			return
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Retour pour ignorer le fichier de test
			return
		}

		// Vérifier et collecter la fonction privée
		info := extractPrivateFuncInfo(funcDecl)
		// Ajouter à la map si valide
		if info != nil {
			// Ajout de la fonction privée à la collection
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
		// Retour si nom invalide
		return nil
	}

	funcName := funcDecl.Name.Name

	// Ignorer les fonctions spéciales Go
	if funcName == "main" || funcName == "init" {
		// Retour si fonction spéciale
		return nil
	}

	// Vérifier si c'est privé (première lettre en minuscule)
	firstChar := rune(funcName[0])
	// Ignorer les fonctions publiques
	if firstChar < 'a' || firstChar > 'z' {
		// Retour si fonction publique
		return nil
	}

	info := &privateFuncInfo{name: funcName, pos: funcDecl.Pos()}

	// Si c'est une méthode, extraire le type du receiver
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		// Extraction du type du receiver
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
//   - cfg: configuration
//
// Returns:
//   - map[string]bool: map des fonctions appelées
func collectCalledFunctions(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) map[string]bool {
	calledInProduction := make(map[string]bool, initialCalledFuncsCap)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)
		filename := pass.Fset.Position(callExpr.Pos()).Filename

		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeFunc004, filename) {
			// Fichier exclu
			return
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Retour pour ignorer le fichier de test
			return
		}

		// Marquer la fonction comme appelée
		if funcName := extractCalledFuncName(callExpr); funcName != "" {
			// Marquage de la fonction comme appelée
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
//   - cfg: configuration
func collectCallbackUsages(pass *analysis.Pass, insp *inspector.Inspector, privateFuncs map[string][]*privateFuncInfo, calledInProduction map[string]bool, cfg *config.Config) {
	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.ValueSpec)(nil),
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename

		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeFunc004, filename) {
			// Fichier exclu
			return
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Retour pour ignorer le fichier de test
			return
		}

		// Traitement spécial pour les CallExpr (méthodes passées en arguments)
		if callExpr, ok := n.(*ast.CallExpr); ok {
			// Collection des callbacks dans les CallExpr
			collectCallExprCallbacks(callExpr, privateFuncs, calledInProduction)
			// Retour après traitement CallExpr
			return
		}

		// Parcourir tous les identifiants pour les autres types de noeuds
		collectIdentCallbacks(n, privateFuncs, calledInProduction)
	})
}

// collectCallExprCallbacks détecte les callbacks passés comme arguments à des appels.
//
// Params:
//   - callExpr: expression d'appel de fonction
//   - privateFuncs: map des fonctions privées
//   - calledInProduction: map des fonctions appelées (modifiée)
func collectCallExprCallbacks(callExpr *ast.CallExpr, privateFuncs map[string][]*privateFuncInfo, calledInProduction map[string]bool) {
	// Parcourir tous les arguments de l'appel
	for _, arg := range callExpr.Args {
		// Cas 1: fonction directe (ex: handler)
		if ident, ok := arg.(*ast.Ident); ok {
			// Marquage de la fonction identifiant
			markIfPrivateFunc(ident.Name, privateFuncs, calledInProduction)
			// Passage au prochain argument
			continue
		}

		// Cas 2: méthode (ex: a.handleLiveness, obj.Method)
		if selExpr, ok := arg.(*ast.SelectorExpr); ok {
			// Marquage de la méthode sélectionnée
			markIfPrivateFunc(selExpr.Sel.Name, privateFuncs, calledInProduction)
		}
	}
}

// collectIdentCallbacks parcourt un noeud pour détecter les identifiants de callbacks.
//
// Params:
//   - n: noeud AST à parcourir
//   - privateFuncs: map des fonctions privées
//   - calledInProduction: map des fonctions appelées (modifiée)
func collectIdentCallbacks(n ast.Node, privateFuncs map[string][]*privateFuncInfo, calledInProduction map[string]bool) {
	ast.Inspect(n, func(node ast.Node) bool {
		// Vérifier si c'est un identifiant
		if ident, ok := node.(*ast.Ident); ok {
			// Marquage si fonction privée
			markIfPrivateFunc(ident.Name, privateFuncs, calledInProduction)
		}
		// Continuer la traversée
		return true
	})
}

// markIfPrivateFunc marque une fonction comme appelée si elle est privée.
//
// Params:
//   - funcName: nom de la fonction
//   - privateFuncs: map des fonctions privées
//   - calledInProduction: map des fonctions appelées (modifiée)
func markIfPrivateFunc(funcName string, privateFuncs map[string][]*privateFuncInfo, calledInProduction map[string]bool) {
	// Vérifier si c'est une fonction privée connue
	if _, exists := privateFuncs[funcName]; exists {
		// Marquage de la fonction comme appelée
		calledInProduction[funcName] = true
	}
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
			// Passage à la fonction suivante si appelée
			continue
		}

		// Reporter toutes les fonctions avec ce nom
		for _, info := range infos {
			// Rapport de fonction non utilisée
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
	msg, _ := messages.Get(ruleCodeFunc004)
	// Vérifier si c'est une méthode
	if info.receiverType != "" {
		pass.Reportf(info.pos,
			"%s: %s",
			ruleCodeFunc004,
			msg.Format(config.Get().Verbose, info.receiverType+"."+info.name))
		// Fin du traitement pour les méthodes
		return
	}

	pass.Reportf(info.pos,
		"%s: %s",
		ruleCodeFunc004,
		msg.Format(config.Get().Verbose, info.name))
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
		// Retour du nom de l'identifiant
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
