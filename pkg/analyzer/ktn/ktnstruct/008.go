// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct008 code de la règle KTN-STRUCT-008
	ruleCodeStruct008 string = "KTN-STRUCT-008"
	// methodsMapCap008 capacité initiale pour la map des méthodes
	methodsMapCap008 int = 16
)

// Analyzer008 checks receiver type consistency (pointer vs value).
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct008",
	Doc:      "KTN-STRUCT-008: Le type de receiver doit être cohérent (pointer ou value)",
	Run:      runStruct008,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// methodReceiverInfo contient les informations sur le receiver d'une méthode.
type methodReceiverInfo struct {
	name      string
	isPointer bool
	funcDecl  *ast.FuncDecl
}

// runStruct008 exécute l'analyse KTN-STRUCT-008.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct008(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct008) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter les méthodes par type avec info pointer/value
	methodsByType := collectMethodReceiverTypes(pass, insp, cfg)

	// Vérifier la cohérence des types de receiver
	checkReceiverTypeConsistency(pass, methodsByType)

	// Retour de la fonction
	return nil, nil
}

// collectMethodReceiverTypes collecte les méthodes par type.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
//
// Returns:
//   - map[string][]methodReceiverInfo: map type -> liste de méthodes
func collectMethodReceiverTypes(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) map[string][]methodReceiverInfo {
	result := make(map[string][]methodReceiverInfo, methodsMapCap008)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct008, filename) {
			// Fichier exclu
			return
		}

		// Vérifier si c'est une méthode
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			// Pas une méthode
			return
		}

		// Extraire les infos du receiver
		recv := funcDecl.Recv.List[0]
		typeName, isPointer := extractReceiverTypeInfo(recv.Type)
		// Si pas de type valide, ignorer
		if typeName == "" {
			// Retour anticipé
			return
		}

		// Ajouter à la map
		result[typeName] = append(result[typeName], methodReceiverInfo{
			name:      funcDecl.Name.Name,
			isPointer: isPointer,
			funcDecl:  funcDecl,
		})
	})

	// Retour de la map
	return result
}

// extractReceiverTypeInfo extrait le nom du type et si c'est un pointer.
//
// Params:
//   - expr: expression du type
//
// Returns:
//   - string: nom du type
//   - bool: true si pointer receiver
func extractReceiverTypeInfo(expr ast.Expr) (string, bool) {
	// Vérifier si c'est un pointer
	isPointer := false
	if star, ok := expr.(*ast.StarExpr); ok {
		isPointer = true
		expr = star.X
	}

	// Gérer le cas Type
	if ident, ok := expr.(*ast.Ident); ok {
		// Retourner le nom et si pointer
		return ident.Name, isPointer
	}

	// Type non reconnu
	return "", false
}

// checkReceiverTypeConsistency vérifie la cohérence des types de receiver.
//
// Params:
//   - pass: contexte d'analyse
//   - methodsByType: map des méthodes par type
func checkReceiverTypeConsistency(pass *analysis.Pass, methodsByType map[string][]methodReceiverInfo) {
	// Parcourir les types
	for typeName, methods := range methodsByType {
		// Si moins de 2 méthodes, pas de vérification
		if len(methods) < 2 {
			// Continuer au type suivant
			continue
		}

		// Compter les pointer et value receivers
		pointerCount, valueCount := countReceiverTypes(methods)

		// Si mix de pointer et value, reporter les inconsistants
		if pointerCount > 0 && valueCount > 0 {
			reportInconsistentReceivers(pass, typeName, methods, pointerCount, valueCount)
		}
	}
}

// countReceiverTypes compte les pointer et value receivers.
//
// Params:
//   - methods: liste des méthodes
//
// Returns:
//   - int: nombre de pointer receivers
//   - int: nombre de value receivers
func countReceiverTypes(methods []methodReceiverInfo) (int, int) {
	pointerCount := 0
	valueCount := 0

	// Parcourir les méthodes
	for _, m := range methods {
		if m.isPointer {
			pointerCount++
		} else {
			valueCount++
		}
	}

	// Retourner les compteurs
	return pointerCount, valueCount
}

// reportInconsistentReceivers signale les receivers inconsistants.
//
// Params:
//   - pass: contexte d'analyse
//   - typeName: nom du type
//   - methods: liste des méthodes
//   - pointerCount: nombre de pointer receivers
//   - valueCount: nombre de value receivers
func reportInconsistentReceivers(pass *analysis.Pass, typeName string, methods []methodReceiverInfo, pointerCount, valueCount int) {
	// Déterminer le type majoritaire
	majorityIsPointer := pointerCount >= valueCount
	majorityType := "pointer"
	minorityType := "value"
	// Inverser si value est majoritaire
	if !majorityIsPointer {
		majorityType = "value"
		minorityType = "pointer"
	}

	// Reporter les méthodes minoritaires
	for _, m := range methods {
		// Si le type correspond au majoritaire, ignorer
		if m.isPointer == majorityIsPointer {
			// Méthode conforme
			continue
		}

		// Trouver une méthode majoritaire pour le message
		var exampleMethod string
		for _, other := range methods {
			if other.isPointer == majorityIsPointer {
				exampleMethod = other.name
				break
			}
		}

		msg, _ := messages.Get(ruleCodeStruct008)
		pass.Reportf(
			m.funcDecl.Recv.List[0].Pos(),
			"%s: %s",
			ruleCodeStruct008,
			msg.Format(config.Get().Verbose, m.name, minorityType, exampleMethod, majorityType, typeName),
		)
	}
}
