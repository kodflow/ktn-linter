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
	// minMethodsForConsistency nombre minimum de méthodes pour vérifier la cohérence
	minMethodsForConsistency int = 2
)

// Analyzer008 checks receiver type consistency (pointer vs value).
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct008",
	Doc:      "KTN-STRUCT-008: Le type de receiver doit être cohérent (pointer ou value)",
	Run:      runStruct008,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
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
	cfg := config.Get()
	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct008) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	methodsByType := collectMethodReceiverTypes(pass, insp, cfg)
	checkReceiverTypeConsistency(pass, methodsByType)

	// Fin de l'analyse
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
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(n.Pos()).Filename
		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeStruct008, filename) {
			// Fichier exclu
			return
		}
		info := extractMethodInfo(funcDecl)
		// Ajouter la méthode si info extraite
		if info != nil {
			result[info.typeName] = append(result[info.typeName], info.receiver)
		}
	})

	// Retour de la map des méthodes
	return result
}

// extractMethodInfo extrait les informations d'une méthode.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - *methodInfoResult: informations extraites ou nil si pas une méthode valide
func extractMethodInfo(funcDecl *ast.FuncDecl) *methodInfoResult {
	// Vérifier si c'est une méthode
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		// Pas une méthode
		return nil
	}

	recv := funcDecl.Recv.List[0]
	typeName, isPointer := extractReceiverTypeInfo(recv.Type)
	// Vérifier si le type est valide
	if typeName == "" {
		// Type invalide
		return nil
	}

	// Retour des informations extraites
	return &methodInfoResult{
		typeName: typeName,
		receiver: methodReceiverInfo{
			name:      funcDecl.Name.Name,
			isPointer: isPointer,
			funcDecl:  funcDecl,
		},
	}
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
	isPointer := false
	// Gérer le cas *Type
	if star, ok := expr.(*ast.StarExpr); ok {
		isPointer = true
		expr = star.X
	}

	// Extraire l'identifiant
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour du nom du type
		return ident.Name, isPointer
	}

	// Type invalide
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
		// Ignorer si moins de minMethodsForConsistency méthodes
		if len(methods) < minMethodsForConsistency {
			continue
		}

		pointerCount, valueCount := countReceiverTypes(methods)
		// Vérifier s'il y a un mélange pointer/value
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
	pointerCount, valueCount := 0, 0
	// Parcourir les méthodes
	for _, m := range methods {
		// Compter pointer vs value
		if m.isPointer {
			pointerCount++
		} else {
			// Value receiver
			valueCount++
		}
	}

	// Retour des compteurs
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
	majorityIsPointer := pointerCount >= valueCount
	majorityType, minorityType := "pointer", "value"
	// Inverser si value est majoritaire
	if !majorityIsPointer {
		majorityType, minorityType = "value", "pointer"
	}

	exampleMethod := findMajorityMethod(methods, majorityIsPointer)
	report := inconsistencyReport{
		typeName:      typeName,
		minorityType:  minorityType,
		majorityType:  majorityType,
		exampleMethod: exampleMethod,
	}
	// Parcourir et signaler les méthodes minoritaires
	for _, m := range methods {
		// Ignorer les méthodes du type majoritaire
		if m.isPointer == majorityIsPointer {
			continue
		}
		reportReceiverInconsistency(pass, m, report)
	}
}

// findMajorityMethod trouve une méthode du type majoritaire.
//
// Params:
//   - methods: liste des méthodes
//   - majorityIsPointer: true si le type majoritaire est pointer
//
// Returns:
//   - string: nom de la méthode exemple
func findMajorityMethod(methods []methodReceiverInfo, majorityIsPointer bool) string {
	// Chercher une méthode du type majoritaire
	for _, m := range methods {
		// Trouver la première méthode correspondante
		if m.isPointer == majorityIsPointer {
			// Retour du nom de la méthode
			return m.name
		}
	}

	// Aucune méthode trouvée
	return ""
}

// reportReceiverInconsistency signale une incohérence de receiver.
//
// Params:
//   - pass: contexte d'analyse
//   - m: méthode incohérente
//   - report: informations du rapport
func reportReceiverInconsistency(pass *analysis.Pass, m methodReceiverInfo, report inconsistencyReport) {
	msg, _ := messages.Get(ruleCodeStruct008)
	pass.Reportf(
		m.funcDecl.Recv.List[0].Pos(),
		"%s: %s",
		ruleCodeStruct008,
		msg.Format(config.Get().Verbose, m.name, report.minorityType, report.exampleMethod, report.majorityType, report.typeName),
	)
}
