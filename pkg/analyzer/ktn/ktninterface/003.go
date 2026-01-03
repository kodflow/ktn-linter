// Package ktninterface provides analyzers for interface-related lint rules.
package ktninterface

import (
	"go/ast"
	"strings"
	"unicode"

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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeInterface003) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		typeSpec := n.(*ast.TypeSpec)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeInterface003, filename) {
			// Fichier exclu
			return
		}

		// Vérifier si c'est une interface
		ifaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		// Si pas une interface, ignorer
		if !ok {
			// Retour anticipé
			return
		}

		// Vérifier si c'est une interface à une seule méthode
		if !isSingleMethodInterface(ifaceType) {
			// Plus d'une méthode ou zéro
			return
		}

		// Extraire le nom de la méthode
		methodName := extractSingleMethodName(ifaceType)
		// Si pas de nom trouvé, ignorer
		if methodName == "" {
			// Retour anticipé
			return
		}

		// Vérifier la convention de nommage -er
		checkErNamingConvention(pass, typeSpec, methodName)
	})

	// Retour de la fonction
	return nil, nil
}

// isSingleMethodInterface vérifie si une interface a exactement une méthode.
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

	// Compter les méthodes (pas les types embarqués)
	methodCount := 0
	// Parcourir les champs
	for _, m := range ifaceType.Methods.List {
		// Une méthode a un type FuncType
		if _, isFunc := m.Type.(*ast.FuncType); isFunc {
			methodCount++
		}
	}

	// Une seule méthode
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
	// Parcourir les champs
	for _, m := range ifaceType.Methods.List {
		// Une méthode a un type FuncType
		if _, isFunc := m.Type.(*ast.FuncType); isFunc {
			// Vérifier si nom présent
			if len(m.Names) > 0 {
				// Retourner le nom
				return m.Names[0].Name
			}
		}
	}

	// Pas de nom trouvé
	return ""
}

// checkErNamingConvention vérifie la convention de nommage -er.
//
// Params:
//   - pass: contexte d'analyse
//   - typeSpec: déclaration de type
//   - methodName: nom de la méthode
func checkErNamingConvention(pass *analysis.Pass, typeSpec *ast.TypeSpec, methodName string) {
	interfaceName := typeSpec.Name.Name

	// Générer le nom suggéré (méthode + "er")
	suggested := suggestErName(methodName)

	// Vérifier si le nom suit la convention
	if followsErConvention(interfaceName, methodName) {
		// Convention respectée
		return
	}

	// Reporter le problème
	msg, _ := messages.Get(ruleCodeInterface003)
	pass.Reportf(
		typeSpec.Pos(),
		"%s: %s",
		ruleCodeInterface003,
		msg.Format(config.Get().Verbose, interfaceName, methodName, suggested),
	)
}

// followsErConvention vérifie si le nom d'interface suit la convention -er.
//
// Params:
//   - interfaceName: nom de l'interface
//   - methodName: nom de la méthode
//
// Returns:
//   - bool: true si convention respectée
func followsErConvention(interfaceName, methodName string) bool {
	// Capitaliser le nom de méthode pour comparaison
	capitalizedMethod := capitalizeFirst(methodName)

	// Vérifier si l'interface finit par "er" et correspond à la méthode
	if strings.HasSuffix(interfaceName, "er") {
		// Vérifier cohérence avec la méthode: Reader = Read + er
		expectedName := capitalizedMethod + "er"
		if interfaceName == expectedName {
			// Convention parfaite: Read -> Reader
			return true
		}
		// Aussi accepter Write -> Writer (le 'e' est déjà là)
		if strings.HasSuffix(capitalizedMethod, "e") {
			altName := capitalizedMethod + "r"
			if interfaceName == altName {
				// Convention acceptée: Write -> Writer
				return true
			}
		}
	}

	// Vérifier les cas spéciaux comme "or" (Handler, Validator, etc.)
	if strings.HasSuffix(interfaceName, "or") {
		// Vérifier cohérence: Handle -> Handler
		expectedName := capitalizedMethod + "r"
		if interfaceName == expectedName {
			// Convention acceptée
			return true
		}
	}

	// Non conforme
	return false
}

// capitalizeFirst met la première lettre en majuscule.
//
// Params:
//   - s: chaîne à capitaliser
//
// Returns:
//   - string: chaîne capitalisée
func capitalizeFirst(s string) string {
	// Vérifier si vide
	if len(s) == 0 {
		// Retourner vide
		return ""
	}

	// Capitaliser la première lettre
	first := unicode.ToUpper(rune(s[0]))

	// Retourner le résultat
	return string(first) + s[1:]
}

// suggestErName génère un nom suggéré basé sur la méthode.
//
// Params:
//   - methodName: nom de la méthode
//
// Returns:
//   - string: nom suggéré (méthode + "er")
func suggestErName(methodName string) string {
	// Vérifier si le nom est vide
	if len(methodName) == 0 {
		// Retourner vide
		return ""
	}

	// Capitaliser la première lettre
	first := unicode.ToUpper(rune(methodName[0]))

	// Ajouter "er" à la fin
	return string(first) + methodName[1:] + "er"
}
