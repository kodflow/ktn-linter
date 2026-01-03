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
	cfg := config.Get()
	if !cfg.IsRuleEnabled(ruleCodeStruct008) {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	methodsByType := collectMethodReceiverTypes(pass, insp, cfg)
	checkReceiverTypeConsistency(pass, methodsByType)

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
		if cfg.IsFileExcluded(ruleCodeStruct008, filename) {
			return
		}
		info := extractMethodInfo(funcDecl)
		if info != nil {
			result[info.typeName] = append(result[info.typeName], info.receiver)
		}
	})

	return result
}

// methodInfoResult contient le résultat de l'extraction des infos d'une méthode.
type methodInfoResult struct {
	typeName string
	receiver methodReceiverInfo
}

// extractMethodInfo extrait les informations d'une méthode.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - *methodInfoResult: informations extraites ou nil si pas une méthode valide
func extractMethodInfo(funcDecl *ast.FuncDecl) *methodInfoResult {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return nil
	}

	recv := funcDecl.Recv.List[0]
	typeName, isPointer := extractReceiverTypeInfo(recv.Type)
	if typeName == "" {
		return nil
	}

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
	if star, ok := expr.(*ast.StarExpr); ok {
		isPointer = true
		expr = star.X
	}

	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name, isPointer
	}

	return "", false
}

// checkReceiverTypeConsistency vérifie la cohérence des types de receiver.
//
// Params:
//   - pass: contexte d'analyse
//   - methodsByType: map des méthodes par type
func checkReceiverTypeConsistency(pass *analysis.Pass, methodsByType map[string][]methodReceiverInfo) {
	for typeName, methods := range methodsByType {
		if len(methods) < 2 {
			continue
		}

		pointerCount, valueCount := countReceiverTypes(methods)
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
	for _, m := range methods {
		if m.isPointer {
			pointerCount++
		} else {
			valueCount++
		}
	}

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
	if !majorityIsPointer {
		majorityType, minorityType = "value", "pointer"
	}

	exampleMethod := findMajorityMethod(methods, majorityIsPointer)
	for _, m := range methods {
		if m.isPointer == majorityIsPointer {
			continue
		}
		reportReceiverInconsistency(pass, m, minorityType, exampleMethod, majorityType, typeName)
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
	for _, m := range methods {
		if m.isPointer == majorityIsPointer {
			return m.name
		}
	}

	return ""
}

// reportReceiverInconsistency signale une incohérence de receiver.
//
// Params:
//   - pass: contexte d'analyse
//   - m: méthode incohérente
//   - minorityType: type minoritaire
//   - exampleMethod: méthode exemple
//   - majorityType: type majoritaire
//   - typeName: nom du type
func reportReceiverInconsistency(pass *analysis.Pass, m methodReceiverInfo, minorityType, exampleMethod, majorityType, typeName string) {
	msg, _ := messages.Get(ruleCodeStruct008)
	pass.Reportf(
		m.funcDecl.Recv.List[0].Pos(),
		"%s: %s",
		ruleCodeStruct008,
		msg.Format(config.Get().Verbose, m.name, minorityType, exampleMethod, majorityType, typeName),
	)
}
