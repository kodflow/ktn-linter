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
// Tested via integration tests in 003_external_test.go
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runInterface003(pass *analysis.Pass) (any, error) {
	cfg := config.Get()
	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeInterface003) {
		// Règle désactivée - retour immédiat
		return nil, nil
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.TypeSpec)(nil)}
	// Parcourir les déclarations de type
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Analyser chaque TypeSpec
		analyzeTypeSpec003(pass, cfg, n.(*ast.TypeSpec))
	})
	// Retour succès
	return nil, nil
}

// analyzeTypeSpec003 analyse une déclaration de type pour la convention -er.
// Tested via integration tests in 003_external_test.go
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration
//   - typeSpec: déclaration de type à analyser
func analyzeTypeSpec003(pass *analysis.Pass, cfg *config.Config, typeSpec *ast.TypeSpec) {
	filename := pass.Fset.Position(typeSpec.Pos()).Filename
	// Vérifier si le fichier est exclu
	if cfg.IsFileExcluded(ruleCodeInterface003, filename) {
		// Fichier exclu - skip
		return
	}
	ifaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	// Vérifier si c'est une interface à une méthode
	if !ok || !isSingleMethodInterface(ifaceType) {
		// Pas une interface à une méthode - skip
		return
	}
	methodName := extractSingleMethodName(ifaceType)
	// Vérifier si nom de méthode trouvé
	if methodName == "" {
		// Pas de nom de méthode - skip
		return
	}
	// Vérifier la convention de nommage
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
	// Vérifier si méthodes présentes
	if ifaceType.Methods == nil {
		// Pas de méthodes
		return false
	}
	methodCount := 0
	// Compter les méthodes (pas les types embarqués)
	for _, m := range ifaceType.Methods.List {
		// Vérifier si c'est une fonction
		if _, isFunc := m.Type.(*ast.FuncType); isFunc {
			// Incrémenter le compteur
			methodCount++
		}
	}
	// Retourner si exactement une méthode
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
	// Parcourir les méthodes
	for _, m := range ifaceType.Methods.List {
		// Vérifier si c'est une fonction avec un nom
		if _, isFunc := m.Type.(*ast.FuncType); isFunc && len(m.Names) > 0 {
			// Retourner le nom de la première méthode
			return m.Names[0].Name
		}
	}
	// Aucune méthode trouvée
	return ""
}

// checkErNamingConvention vérifie la convention de nommage -er et reporte si non conforme.
// Tested via integration tests in 003_external_test.go
//
// Params:
//   - pass: contexte d'analyse
//   - typeSpec: déclaration de type
//   - methodName: nom de la méthode
func checkErNamingConvention(pass *analysis.Pass, typeSpec *ast.TypeSpec, methodName string) {
	interfaceName := typeSpec.Name.Name
	// Vérifier si la convention est respectée
	if followsErConvention(interfaceName, methodName) {
		// Convention respectée - skip
		return
	}
	// Générer un nom suggéré
	suggested := suggestErName(methodName)
	msg, _ := messages.Get(ruleCodeInterface003)
	// Reporter la non-conformité
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
		// Convention -er respectée
		return true
	}
	// Write → Writer (méthode finissant par 'e')
	if strings.HasSuffix(capitalized, "e") && interfaceName == capitalized+"r" {
		// Convention -r respectée
		return true
	}
	// Validate → Validator (-ate → -ator)
	if strings.HasSuffix(capitalized, "ate") {
		base := strings.TrimSuffix(capitalized, "e")
		// Vérifier si convention -ator respectée
		if interfaceName == base+"or" {
			// Convention -ator respectée
			return true
		}
	}
	// Handle → Handler (-or suffix avec méthode finissant par 'e')
	if strings.HasSuffix(interfaceName, "or") && interfaceName == capitalized+"r" {
		// Convention -or respectée
		return true
	}
	// Aucune convention respectée
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
	// Vérifier si chaîne vide
	if s == "" {
		// Retourner chaîne vide
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	// Retourner avec première lettre en majuscule
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
	// Vérifier si chaîne vide
	if methodName == "" {
		// Retourner chaîne vide
		return ""
	}
	capitalized := capitalizeFirst(methodName)
	// Verbes en -ate → -ator (Validate → Validator, Generate → Generator)
	if strings.HasSuffix(capitalized, "ate") {
		// Retourner avec suffixe -ator
		return strings.TrimSuffix(capitalized, "e") + "or"
	}
	// Verbes en -e → -er (Write → Writer)
	if strings.HasSuffix(capitalized, "e") {
		// Retourner avec suffixe -r
		return capitalized + "r"
	}
	// Cas standard (Read → Reader)
	return capitalized + "er"
}
