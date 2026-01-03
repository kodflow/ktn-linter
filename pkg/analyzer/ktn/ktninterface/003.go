// Package ktninterface provides analyzers for interface-related lint rules.
package ktninterface

import (
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeInterface003 code de la règle KTN-INTERFACE-003
	ruleCodeInterface003 string = "KTN-INTERFACE-003"
)

// Analyzer003 checks that single-method interfaces follow -er naming convention.
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktninterface003",
	Doc:      "KTN-INTERFACE-003: Les interfaces à une méthode doivent suivre la convention -er",
	Run:      runInterface003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runInterface003 exécute l'analyse KTN-INTERFACE-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runInterface003(pass *analysis.Pass) (any, error) {
	cfg := config.Get()
	if !cfg.IsRuleEnabled(ruleCodeInterface003) {
		return nil, nil
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.TypeSpec)(nil)}
	insp.Preorder(nodeFilter, func(n ast.Node) {
		analyzeTypeSpec003(pass, cfg, n.(*ast.TypeSpec))
	})
	return nil, nil
}

// analyzeTypeSpec003 analyse une déclaration de type pour la convention -er.
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration
//   - typeSpec: déclaration de type à analyser
func analyzeTypeSpec003(pass *analysis.Pass, cfg *config.Config, typeSpec *ast.TypeSpec) {
	filename := pass.Fset.Position(typeSpec.Pos()).Filename
	if cfg.IsFileExcluded(ruleCodeInterface003, filename) {
		return
	}
	ifaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok || !isSingleMethodInterface(ifaceType) {
		return
	}
	methodName := extractSingleMethodName(ifaceType)
	if methodName == "" {
		return
	}
	checkErNamingConvention(pass, typeSpec, methodName)
}

// isSingleMethodInterface vérifie si une interface a exactement une méthode.
// Ne compte que les méthodes réelles, pas les types embarqués.
//
// Params:
//   - ifaceType: type interface à vérifier
//
// Returns:
//   - bool: true si une seule méthode
func isSingleMethodInterface(ifaceType *ast.InterfaceType) bool {
	if ifaceType.Methods == nil {
		return false
	}
	methodCount := 0
	for _, m := range ifaceType.Methods.List {
		if _, isFunc := m.Type.(*ast.FuncType); isFunc {
			methodCount++
		}
	}
	return methodCount == 1
}

// extractSingleMethodName extrait le nom de la méthode unique.
//
// Params:
//   - ifaceType: type interface
//
// Returns:
//   - string: nom de la méthode ou vide
func extractSingleMethodName(ifaceType *ast.InterfaceType) string {
	for _, m := range ifaceType.Methods.List {
		if _, isFunc := m.Type.(*ast.FuncType); isFunc && len(m.Names) > 0 {
			return m.Names[0].Name
		}
	}
	return ""
}

// checkErNamingConvention vérifie la convention de nommage -er et reporte si non conforme.
//
// Params:
//   - pass: contexte d'analyse
//   - typeSpec: déclaration de type
//   - methodName: nom de la méthode
func checkErNamingConvention(pass *analysis.Pass, typeSpec *ast.TypeSpec, methodName string) {
	interfaceName := typeSpec.Name.Name
	if followsErConvention(interfaceName, methodName) {
		return
	}
	suggested := suggestErName(methodName)
	msg, _ := messages.Get(ruleCodeInterface003)
	pass.Reportf(
		typeSpec.Pos(),
		"%s: %s",
		ruleCodeInterface003,
		msg.Format(config.Get().Verbose, interfaceName, methodName, suggested),
	)
}

// followsErConvention vérifie si le nom d'interface suit la convention -er/-or.
//
// Params:
//   - interfaceName: nom de l'interface
//   - methodName: nom de la méthode
//
// Returns:
//   - bool: true si convention respectée
func followsErConvention(interfaceName, methodName string) bool {
	capitalized := capitalizeFirst(methodName)
	// Read → Reader
	if interfaceName == capitalized+"er" {
		return true
	}
	// Write → Writer (méthode finissant par 'e')
	if strings.HasSuffix(capitalized, "e") && interfaceName == capitalized+"r" {
		return true
	}
	// Validate → Validator (-ate → -ator)
	if strings.HasSuffix(capitalized, "ate") {
		base := strings.TrimSuffix(capitalized, "e")
		if interfaceName == base+"or" {
			return true
		}
	}
	// Handle → Handler (-or suffix avec méthode finissant par 'e')
	if strings.HasSuffix(interfaceName, "or") && interfaceName == capitalized+"r" {
		return true
	}
	return false
}

// capitalizeFirst met la première lettre en majuscule (UTF-8 safe).
//
// Params:
//   - s: chaîne à capitaliser
//
// Returns:
//   - string: chaîne capitalisée
func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// suggestErName génère un nom suggéré basé sur la méthode (UTF-8 safe).
// Gère les cas spéciaux: -ate → -ator (Validate → Validator).
//
// Params:
//   - methodName: nom de la méthode
//
// Returns:
//   - string: nom suggéré
func suggestErName(methodName string) string {
	if methodName == "" {
		return ""
	}
	capitalized := capitalizeFirst(methodName)
	// Verbes en -ate → -ator (Validate → Validator, Generate → Generator)
	if strings.HasSuffix(capitalized, "ate") {
		return strings.TrimSuffix(capitalized, "e") + "or"
	}
	// Verbes en -e → -er (Write → Writer)
	if strings.HasSuffix(capitalized, "e") {
		return capitalized + "r"
	}
	// Cas standard (Read → Reader)
	return capitalized + "er"
}
