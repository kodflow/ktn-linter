// Analyzer 008 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// INITIAL_PARAMS_CAP initial capacity for params map
	INITIAL_PARAMS_CAP int = 8
	// INITIAL_USED_VARS_CAP initial capacity for used vars map
	INITIAL_USED_VARS_CAP int = 16
)

// Analyzer008 vérifie que les paramètres non utilisés sont explicitement ignorés.
var Analyzer008 = &analysis.Analyzer{
	Name:     "ktnfunc008",
	Doc:      "KTN-FUNC-008: paramètres non utilisés doivent être préfixés par _ ou assignés à _",
	Run:      runFunc008,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc008 exécute l'analyse KTN-FUNC-008.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc008(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filtrer uniquement les déclarations de fonction
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Vérification de la présence d'un corps de fonction
		if funcDecl.Body == nil {
			// Fonction sans corps (interface method signature)
			return
		}

		// Collecter tous les paramètres
		params := collectFunctionParams(funcDecl)

		// Collecter les variables utilisées dans le corps
		usedVars := collectUsedVariables(funcDecl.Body)

		// Collecter les variables explicitement ignorées avec _ = var
		ignoredVars := collectIgnoredVariables(funcDecl.Body)

		// Vérifier chaque paramètre
		for paramName, paramPos := range params {
			// Ignorer les paramètres déjà préfixés par _
			if len(paramName) > 0 && paramName[0] == '_' {
				// Paramètre déjà marqué comme ignoré
				continue
			}

			// Vérifier si le paramètre est utilisé
			if usedVars[paramName] {
				// Paramètre utilisé
				continue
			}

			// Vérifier si le paramètre est explicitement ignoré
			if ignoredVars[paramName] {
				// Paramètre explicitement ignoré
				continue
			}

			// Paramètre non utilisé et non ignoré - reporter l'erreur
			pass.Reportf(
				paramPos,
				"KTN-FUNC-008: le paramètre '%s' n'est pas utilisé. Préfixez-le par _ (ex: _%s) ou ajoutez '_ = %s' dans le corps de la fonction",
				paramName,
				paramName,
				paramName,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// collectFunctionParams collecte tous les paramètres d'une fonction.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - map[string]token.Pos: map des noms de paramètres vers leurs positions
func collectFunctionParams(funcDecl *ast.FuncDecl) map[string]token.Pos {
	params := make(map[string]token.Pos, INITIAL_PARAMS_CAP)

	// Vérification de la présence de paramètres
	if funcDecl.Type.Params == nil {
		// Pas de paramètres
		return params
	}

	// Parcourir tous les paramètres
	for _, field := range funcDecl.Type.Params.List {
		// Parcourir tous les noms dans ce champ (peut y avoir plusieurs: x, y int)
		for _, name := range field.Names {
			// Vérification du nom
			if name != nil && name.Name != "_" {
				// Ajouter le paramètre
				params[name.Name] = name.Pos()
			}
		}
	}

	// Retour de la map des paramètres
	return params
}

// collectUsedVariables collecte toutes les variables utilisées dans le corps.
//
// Params:
//   - body: corps de fonction
//
// Returns:
//   - map[string]bool: map des variables utilisées
func collectUsedVariables(body *ast.BlockStmt) map[string]bool {
	used := make(map[string]bool, INITIAL_USED_VARS_CAP)

	ast.Inspect(body, func(n ast.Node) bool {
		ident, isIdent := n.(*ast.Ident)
		// Vérifier si c'est un identifiant
		if isIdent {
			parent, found := findParentAssignToBlank(body, ident)
			// Vérifier si dans une assignation à _
			if found && parent {
				// C'est dans une assignation à _ - ne pas compter comme utilisé
				return true
			}
			// Ajouter comme variable utilisée
			used[ident.Name] = true
		}
		// Continuer la traversée
		return true
	})

	// Retour de la map des variables utilisées
	return used
}

// collectIgnoredVariables collecte les variables explicitement ignorées avec _ = var.
//
// Params:
//   - body: corps de fonction
//
// Returns:
//   - map[string]bool: map des variables ignorées
func collectIgnoredVariables(body *ast.BlockStmt) map[string]bool {
	ignored := make(map[string]bool, INITIAL_PARAMS_CAP)

	ast.Inspect(body, func(n ast.Node) bool {
		assign, isAssign := n.(*ast.AssignStmt)
		// Vérifier si c'est une assignation
		if !isAssign {
			// Continuer la traversée
			return true
		}
		// Vérifier si le côté gauche est _
		if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			// Continuer la traversée
			return true
		}
		lhsIdent, isLhsIdent := assign.Lhs[0].(*ast.Ident)
		// Vérification du côté gauche
		if !isLhsIdent || lhsIdent.Name != "_" {
			// Continuer la traversée
			return true
		}
		rhsIdent, isRhsIdent := assign.Rhs[0].(*ast.Ident)
		// Vérification du côté droit
		if isRhsIdent {
			// Ajouter comme variable ignorée
			ignored[rhsIdent.Name] = true
		}
		// Continuer la traversée
		return true
	})

	// Retour de la map des variables ignorées
	return ignored
}

// findParentAssignToBlank vérifie si un identifiant est dans une assignation à _.
//
// Params:
//   - body: corps de fonction
//   - target: identifiant cible
//
// Returns:
//   - bool: true si dans une assignation à _
//   - bool: true si trouvé
func findParentAssignToBlank(body *ast.BlockStmt, target *ast.Ident) (bool, bool) {
	found := false
	inAssignToBlank := false

	ast.Inspect(body, func(n ast.Node) bool {
		assign, isAssign := n.(*ast.AssignStmt)
		// Vérifier si c'est une assignation
		if !isAssign {
			// Continuer la traversée
			return true
		}
		// Vérifier la structure de l'assignation
		if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			// Continuer la traversée
			return true
		}
		lhsIdent, isLhsIdent := assign.Lhs[0].(*ast.Ident)
		// Vérification du côté gauche
		if !isLhsIdent || lhsIdent.Name != "_" {
			// Continuer la traversée
			return true
		}
		rhsIdent, isRhsIdent := assign.Rhs[0].(*ast.Ident)
		// Vérification du côté droit
		if !isRhsIdent {
			// Continuer la traversée
			return true
		}
		// Vérification si c'est notre target
		if rhsIdent.Pos() == target.Pos() {
			found = true
			inAssignToBlank = true
			// Arrêter la recherche
			return false
		}
		// Continuer la traversée
		return true
	})

	// Retour du résultat
	return inAssignToBlank, found
}
