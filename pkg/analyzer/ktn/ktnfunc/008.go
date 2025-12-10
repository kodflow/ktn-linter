// Analyzer 008 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeFunc008 is the rule code for this analyzer
	ruleCodeFunc008 string = "KTN-FUNC-008"
	// initialParamsCap initial capacity for params map
	initialParamsCap int = 8
	// initialUsedVarsCap initial capacity for used vars map
	initialUsedVarsCap int = 16
)

// paramCheckContext contains context for parameter checking.
type paramCheckContext struct {
	pass        *analysis.Pass
	usedVars    map[string]bool
	ignoredVars map[string]bool
	ifaceName   string
}

// Analyzer008 vérifie que les paramètres non utilisés sont explicitement ignorés.
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeFunc008) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filtrer uniquement les déclarations de fonction
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeFunc008, filename) {
			// Fichier exclu
			return
		}

		// Analyser la fonction
		analyzeFunc008(pass, funcDecl)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeFunc008 analyse une déclaration de fonction pour FUNC-008.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction à analyser
func analyzeFunc008(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Vérification de la présence d'un corps de fonction
	if funcDecl.Body == nil {
		// Retour si fonction sans corps (interface method signature)
		return
	}

	// Collecter tous les paramètres
	params := collectFunctionParams008(funcDecl)

	// Collecter les variables utilisées dans le corps
	usedVars := collectUsedVariables008(funcDecl.Body)

	// Collecter les variables explicitement ignorées avec _ = var
	ignoredVars := collectIgnoredVariables008(funcDecl.Body)

	// Déterminer si c'est une implémentation d'interface
	ifaceName := findImplementedInterface(pass, funcDecl)

	// Créer le contexte de vérification
	ctx := &paramCheckContext{
		pass:        pass,
		usedVars:    usedVars,
		ignoredVars: ignoredVars,
		ifaceName:   ifaceName,
	}

	// Vérifier chaque paramètre
	for paramName, paramPos := range params {
		// Vérifier le paramètre
		ctx.checkParam008(paramName, paramPos)
	}
}

// checkParam008 vérifie un paramètre individuel.
//
// Params:
//   - paramName: nom du paramètre
//   - paramPos: position du paramètre
func (ctx *paramCheckContext) checkParam008(paramName string, paramPos token.Pos) {
	// Ignorer les paramètres déjà préfixés par _
	if len(paramName) > 0 && paramName[0] == '_' {
		// Paramètre déjà marqué comme ignoré
		return
	}

	// Vérifier si le paramètre est utilisé
	if ctx.usedVars[paramName] {
		// Retour si paramètre utilisé
		return
	}

	// Générer le message approprié selon le contexte
	if ctx.ignoredVars[paramName] {
		// Pattern _ = param détecté
		reportUnusedWithBypass(ctx.pass, paramPos, paramName, ctx.ifaceName)
		// Retour après signalement
		return
	}

	// Paramètre non utilisé et non ignoré
	reportUnusedParam(ctx.pass, paramPos, paramName, ctx.ifaceName)
}

// reportUnusedWithBypass signale un paramètre avec contournement _ = param.
//
// Params:
//   - pass: contexte d'analyse
//   - pos: position du paramètre
//   - name: nom du paramètre
//   - ifaceName: nom de l'interface (vide si fonction native)
func reportUnusedWithBypass(pass *analysis.Pass, pos token.Pos, name string, ifaceName string) {
	// Vérifier si implémentation d'interface
	if ifaceName != "" {
		// Implémentation d'interface - toléré avec préfixe _
		pass.Reportf(
			pos,
			"KTN-FUNC-008: le paramètre '%s' utilise '_ = %s' pour contourner le compilateur. "+
				"Cette méthode implémente l'interface '%s'. "+
				"Utilisez le préfixe _ (ex: _%s) pour indiquer explicitement que ce paramètre est imposé par l'interface",
			name, name, ifaceName, name,
		)
		// Fin du traitement
		return
	}

	// Fonction native - doit supprimer le paramètre
	pass.Reportf(
		pos,
		"KTN-FUNC-008: le paramètre '%s' utilise '_ = %s' pour contourner le compilateur. "+
			"SUPPRIMEZ ce paramètre inutilisé. "+
			"Si ce paramètre est requis par une interface, vérifiez que le type implémente bien cette interface. "+
			"Sinon, ce paramètre indique peut-être une fonctionnalité non implémentée",
		name, name,
	)
}

// reportUnusedParam signale un paramètre non utilisé.
//
// Params:
//   - pass: contexte d'analyse
//   - pos: position du paramètre
//   - name: nom du paramètre
//   - ifaceName: nom de l'interface (vide si fonction native)
func reportUnusedParam(pass *analysis.Pass, pos token.Pos, name string, ifaceName string) {
	// Vérifier si implémentation d'interface
	if ifaceName != "" {
		// Implémentation d'interface - toléré avec préfixe _
		pass.Reportf(
			pos,
			"KTN-FUNC-008: le paramètre '%s' n'est pas utilisé. "+
				"Cette méthode implémente l'interface '%s'. "+
				"Utilisez le préfixe _ (ex: _%s) pour indiquer explicitement que ce paramètre est imposé par l'interface",
			name, ifaceName, name,
		)
		// Fin du traitement
		return
	}

	// Fonction native - doit supprimer le paramètre
	pass.Reportf(
		pos,
		"KTN-FUNC-008: le paramètre '%s' n'est pas utilisé. "+
			"SUPPRIMEZ ce paramètre inutilisé. "+
			"Si ce paramètre est requis par une interface, vérifiez que le type implémente bien cette interface. "+
			"Sinon, ce paramètre indique peut-être une fonctionnalité non implémentée",
		name,
	)
}

// findImplementedInterface détecte si une méthode implémente une interface.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - string: nom de l'interface implémentée (vide si aucune)
func findImplementedInterface(pass *analysis.Pass, funcDecl *ast.FuncDecl) string {
	// Vérifier si c'est une méthode (a un receiver)
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		// Retour chaîne vide si fonction sans receiver
		return ""
	}

	// Obtenir le type du receiver
	recvType := getReceiverType(pass, funcDecl.Recv.List[0].Type)
	// Vérifier si le type est valide
	if recvType == nil {
		// Retour chaîne vide si type non résolu
		return ""
	}

	// Obtenir le nom de la méthode
	methodName := funcDecl.Name.Name

	// Chercher dans le scope du package les interfaces
	return findInterfaceForMethod(pass, recvType, methodName)
}

// getReceiverType obtient le type sous-jacent du receiver.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression du type du receiver
//
// Returns:
//   - types.Type: type du receiver (nil si non résolu)
func getReceiverType(pass *analysis.Pass, expr ast.Expr) types.Type {
	// Gérer le cas du pointeur (*Type)
	if starExpr, isStar := expr.(*ast.StarExpr); isStar {
		// Extraction du type pointé
		expr = starExpr.X
	}

	// Obtenir le type via TypesInfo
	typeInfo := pass.TypesInfo.TypeOf(expr)
	// Vérifier si le type est valide
	if typeInfo == nil {
		// Retour nil si type non résolu
		return nil
	}

	// Retourner le type
	return typeInfo
}

// findInterfaceForMethod cherche une interface implémentée par le type pour la méthode.
//
// Params:
//   - pass: contexte d'analyse
//   - recvType: type du receiver
//   - methodName: nom de la méthode
//
// Returns:
//   - string: nom de l'interface trouvée (vide si aucune)
func findInterfaceForMethod(pass *analysis.Pass, recvType types.Type, methodName string) string {
	// Parcourir tous les packages importés et le package courant
	pkgScope := pass.Pkg.Scope()

	// Chercher dans le scope du package courant
	ifaceName := searchInterfaceInScope(pkgScope, recvType, methodName)
	// Vérifier si trouvé
	if ifaceName != "" {
		// Retour du nom de l'interface trouvée
		return ifaceName
	}

	// Chercher dans les imports
	for _, imp := range pass.Pkg.Imports() {
		// Chercher dans le scope de l'import
		ifaceName = searchInterfaceInScope(imp.Scope(), recvType, methodName)
		// Vérifier si trouvé
		if ifaceName != "" {
			// Retour du nom de l'interface trouvée
			return ifaceName
		}
	}

	// Aucune interface trouvée
	return ""
}

// searchInterfaceInScope cherche une interface dans un scope.
//
// Params:
//   - scope: scope à parcourir
//   - recvType: type du receiver
//   - methodName: nom de la méthode
//
// Returns:
//   - string: nom de l'interface trouvée (vide si aucune)
func searchInterfaceInScope(scope *types.Scope, recvType types.Type, methodName string) string {
	// Parcourir tous les noms du scope
	for _, name := range scope.Names() {
		// Récupérer l'objet
		obj := scope.Lookup(name)
		// Vérifier si c'est un type
		typeName, isTypeName := obj.(*types.TypeName)
		// Si pas un type, continuer
		if !isTypeName {
			// Continuer la recherche
			continue
		}

		// Vérifier si c'est une interface
		iface, isIface := typeName.Type().Underlying().(*types.Interface)
		// Si pas une interface, continuer
		if !isIface {
			// Continuer la recherche
			continue
		}

		// Vérifier si l'interface a cette méthode
		if !interfaceHasMethod(iface, methodName) {
			// Passage à l'interface suivante
			continue
		}

		// Vérifier si le type implémente l'interface
		if types.Implements(recvType, iface) || implementsWithPointer(recvType, iface) {
			// Retour du nom de l'interface trouvée
			return name
		}
	}

	// Aucune interface trouvée
	return ""
}

// interfaceHasMethod vérifie si une interface a une méthode donnée.
//
// Params:
//   - iface: interface à vérifier
//   - methodName: nom de la méthode
//
// Returns:
//   - bool: true si l'interface a la méthode
func interfaceHasMethod(iface *types.Interface, methodName string) bool {
	// Parcourir les méthodes de l'interface avec l'itérateur
	for method := range iface.Methods() {
		// Vérifier le nom
		if method.Name() == methodName {
			// Retour true si méthode trouvée
			return true
		}
	}
	// Méthode non trouvée
	return false
}

// implementsWithPointer vérifie si *T implémente l'interface.
//
// Params:
//   - t: type à vérifier
//   - iface: interface à implémenter
//
// Returns:
//   - bool: true si *T implémente l'interface
func implementsWithPointer(t types.Type, iface *types.Interface) bool {
	// Créer le type pointeur
	ptrType := types.NewPointer(t)
	// Vérifier l'implémentation
	return types.Implements(ptrType, iface)
}

// collectFunctionParams008 collecte tous les paramètres d'une fonction.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - map[string]token.Pos: map des noms de paramètres vers leurs positions
func collectFunctionParams008(funcDecl *ast.FuncDecl) map[string]token.Pos {
	params := make(map[string]token.Pos, initialParamsCap)

	// Vérification de la présence de paramètres
	if funcDecl.Type.Params == nil {
		// Pas de paramètres
		return params
	}

	// Parcourir tous les paramètres
	for _, field := range funcDecl.Type.Params.List {
		// Parcourir tous les noms dans ce champ
		for _, name := range field.Names {
			// Vérification du nom
			if name != nil && name.Name != "_" {
				// Ajout du paramètre à la collection
				params[name.Name] = name.Pos()
			}
		}
	}

	// Retour de la map des paramètres
	return params
}

// collectUsedVariables008 collecte toutes les variables utilisées dans le corps.
//
// Params:
//   - body: corps de fonction
//
// Returns:
//   - map[string]bool: map des variables utilisées
func collectUsedVariables008(body *ast.BlockStmt) map[string]bool {
	used := make(map[string]bool, initialUsedVarsCap)

	ast.Inspect(body, func(n ast.Node) bool {
		ident, isIdent := n.(*ast.Ident)
		// Vérifier si c'est un identifiant
		if isIdent {
			inBlank, found := findParentAssignToBlank008(body, ident)
			// Vérifier si dans une assignation à _
			if found && inBlank {
				// Retour pour ne pas compter comme utilisé
				return true
			}
			// Ajout comme variable utilisée
			used[ident.Name] = true
		}
		// Continuer la traversée
		return true
	})

	// Retour de la map des variables utilisées
	return used
}

// collectIgnoredVariables008 collecte les variables ignorées avec _ = var.
//
// Params:
//   - body: corps de fonction
//
// Returns:
//   - map[string]bool: map des variables ignorées
func collectIgnoredVariables008(body *ast.BlockStmt) map[string]bool {
	ignored := make(map[string]bool, initialParamsCap)

	ast.Inspect(body, func(n ast.Node) bool {
		assign, isAssign := n.(*ast.AssignStmt)
		// Vérifier si c'est une assignation
		if !isAssign {
			// Continuer la traversée
			return true
		}
		// Vérifier la structure
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
			// Ajout comme variable ignorée
			ignored[rhsIdent.Name] = true
		}
		// Continuer la traversée
		return true
	})

	// Retour de la map des variables ignorées
	return ignored
}

// findParentAssignToBlank008 vérifie si un identifiant est dans une assignation à _.
//
// Params:
//   - body: corps de fonction
//   - target: identifiant cible
//
// Returns:
//   - bool: true si dans une assignation à _
//   - bool: true si trouvé
func findParentAssignToBlank008(body *ast.BlockStmt, target *ast.Ident) (bool, bool) {
	found := false
	inAssignToBlank := false

	ast.Inspect(body, func(n ast.Node) bool {
		assign, isAssign := n.(*ast.AssignStmt)
		// Vérifier si c'est une assignation
		if !isAssign {
			// Continuer la traversée
			return true
		}
		// Vérifier la structure
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
			// Marquage comme trouvé et assigné à blank
			found = true
			inAssignToBlank = true
			// Retour false pour arrêter la recherche
			return false
		}
		// Continuer la traversée
		return true
	})

	// Retour du résultat
	return inAssignToBlank, found
}
