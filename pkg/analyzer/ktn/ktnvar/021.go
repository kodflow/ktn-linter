// Package ktnvar implements KTN linter rules.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar021 is the rule code for this analyzer
	ruleCodeVar021 string = "KTN-VAR-021"
	// initialMethodsCap initial capacity for methods map
	initialMethodsCap int = 10
)

// receiverInfo stocke les informations sur un receiver.
type receiverInfo struct {
	isPointer bool
	pos       ast.Node
}

// Analyzer021 détecte les incohérences de receiver (pointeur vs valeur).
//
// Toutes les méthodes d'un type doivent utiliser le même type de receiver.
var Analyzer021 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar021",
	Doc:      "KTN-VAR-021: Vérifie la cohérence des receivers (tous pointeur ou tous valeur)",
	Run:      runVar021,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar021 exécute l'analyse de cohérence des receivers.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: erreur éventuelle
func runVar021(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar021) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp := inspAny.(*inspector.Inspector)
	// Defensive: avoid nil dereference when resolving positions
	if pass.Fset == nil {
		return nil, nil
	}

	// Collecte des receivers par type
	typeReceivers := collectReceivers(pass, insp, cfg)

	// Vérification de la cohérence
	checkReceiverConsistency(pass, typeReceivers)

	// Traitement terminé
	return nil, nil
}

// collectReceivers collecte les receivers de toutes les méthodes par type.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
//
// Returns:
//   - map[string][]receiverInfo: map des receivers par type
func collectReceivers(
	pass *analysis.Pass,
	insp *inspector.Inspector,
	cfg *config.Config,
) map[string][]receiverInfo {
	// Map pour stocker les receivers par type
	typeReceivers := make(map[string][]receiverInfo, initialMethodsCap)

	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcours des fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en fonction
		funcDecl := n.(*ast.FuncDecl)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar021, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Vérification si c'est une méthode
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			// Pas une méthode
			return
		}

		// Extraction du receiver
		recv := funcDecl.Recv.List[0]
		typeName, isPointer := extractReceiverType(recv.Type)

		// Vérification du nom de type
		if typeName == "" {
			// Type non identifiable
			return
		}

		// Ajout à la map
		typeReceivers[typeName] = append(typeReceivers[typeName], receiverInfo{
			isPointer: isPointer,
			pos:       recv,
		})
	})

	// Retour de la map
	return typeReceivers
}

// extractReceiverType extrait le nom du type et si c'est un pointeur.
//
// Params:
//   - expr: expression du type receiver
//
// Returns:
//   - string: nom du type
//   - bool: true si pointeur
func extractReceiverType(expr ast.Expr) (string, bool) {
	// Vérification si c'est un pointeur
	if star, ok := expr.(*ast.StarExpr); ok {
		// Extraction du nom sous le pointeur
		if ident, ok := star.X.(*ast.Ident); ok {
			// Retour du nom et flag pointeur
			return ident.Name, true
		}
		// Type non identifiable
		return "", false
	}

	// Vérification si c'est un identifiant direct
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour du nom sans pointeur
		return ident.Name, false
	}

	// Type non identifiable
	return "", false
}

// checkReceiverConsistency vérifie la cohérence des receivers par type.
//
// Params:
//   - pass: contexte d'analyse
//   - typeReceivers: map des receivers par type
func checkReceiverConsistency(pass *analysis.Pass, typeReceivers map[string][]receiverInfo) {
	// Parcours des types
	for typeName, receivers := range typeReceivers {
		// Vérification si le type a au moins 2 méthodes
		if len(receivers) < 2 {
			// Pas assez de méthodes pour détecter une incohérence
			continue
		}

		// Détermination du type dominant (premier receiver)
		dominantIsPointer := receivers[0].isPointer

		// Vérification de la cohérence
		for _, recv := range receivers[1:] {
			// Vérification si le receiver est différent du dominant
			if recv.isPointer != dominantIsPointer {
				// Incohérence détectée
				reportInconsistency(pass, recv.pos, typeName, dominantIsPointer)
			}
		}
	}
}

// reportInconsistency signale une incohérence de receiver.
//
// Params:
//   - pass: contexte d'analyse
//   - pos: position du receiver incohérent
//   - typeName: nom du type
//   - dominantIsPointer: si le type dominant est pointeur
func reportInconsistency(pass *analysis.Pass, pos ast.Node, typeName string, dominantIsPointer bool) {
	// Détermination du type attendu
	expected := "valeur"
	// Vérification si pointeur attendu
	if dominantIsPointer {
		expected = "pointeur"
	}

	// Récupération du message
	msg, ok := messages.Get(ruleCodeVar021)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(pos.Pos(), "%s: receiver incohérent pour %s, attendu %s", ruleCodeVar021, typeName, expected)
		return
	}

	// Rapport d'erreur
	pass.Reportf(
		pos.Pos(),
		"%s: %s",
		ruleCodeVar021,
		msg.Format(config.Get().Verbose, typeName, expected),
	)
}
