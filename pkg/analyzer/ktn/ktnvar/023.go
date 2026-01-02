// Package ktnvar implements KTN linter rules.
package ktnvar

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar023 is the rule code for this analyzer
	ruleCodeVar023 string = "KTN-VAR-023"
	// initialAliasMapCap is the initial capacity for alias map
	initialAliasMapCap int = 4
)

var (
	// securityKeywords contains keywords indicating security context.
	securityKeywords []string = []string{
		"key", "token", "secret", "password", "salt", "nonce",
		"crypt", "auth", "credential",
	}

	// Analyzer023 detects usage of math/rand in security contexts.
	//
	// Using math/rand for cryptographic purposes is insecure as it uses
	// a predictable PRNG. crypto/rand should be used instead.
	Analyzer023 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktnvar023",
		Doc:      "KTN-VAR-023: Détecte l'utilisation de math/rand dans un contexte sécurité",
		Run:      runVar023,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
)

// runVar023 executes the analysis for math/rand in security contexts.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: always nil
//   - error: potential error
func runVar023(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar023) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp, ok := inspAny.(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// Collecter les alias pour math/rand
	mathRandAliases := collectMathRandAliases(pass)

	// Si math/rand n'est pas importé, rien à vérifier
	if len(mathRandAliases) == 0 {
		// Pas d'import math/rand
		return nil, nil
	}

	// Analyser les appels et variables
	checkMathRandUsage(pass, insp, cfg, mathRandAliases)

	// Traitement
	return nil, nil
}

// collectMathRandAliases collects all aliases for math/rand import.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - map[string]bool: set of aliases for math/rand
func collectMathRandAliases(pass *analysis.Pass) map[string]bool {
	// Map pour stocker les alias
	aliases := make(map[string]bool, initialAliasMapCap)

	// Parcours des fichiers
	for _, file := range pass.Files {
		// Parcours des imports
		for _, imp := range file.Imports {
			// Vérification si c'est math/rand
			if imp.Path.Value == `"math/rand"` || imp.Path.Value == `"math/rand/v2"` {
				// Détermination de l'alias
				if imp.Name != nil {
					// Import avec alias
					aliases[imp.Name.Name] = true
				} else {
					// Import standard
					aliases["rand"] = true
				}
			}
		}
	}

	// Retour des alias
	return aliases
}

// checkMathRandUsage checks for math/rand usage in security contexts.
//
// Params:
//   - pass: analysis context
//   - insp: AST inspector
//   - cfg: configuration
//   - aliases: math/rand aliases
func checkMathRandUsage(
	pass *analysis.Pass,
	insp *inspector.Inspector,
	cfg *config.Config,
	aliases map[string]bool,
) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.GenDecl)(nil),
	}

	// Parcours des nœuds
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar023, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Traitement selon le type de nœud
		switch node := n.(type) {
		// Vérification des fonctions
		case *ast.FuncDecl:
			checkFuncForMathRand(pass, node, aliases)
		// Vérification des déclarations globales
		case *ast.GenDecl:
			checkGenDeclForMathRand(pass, node, aliases)
		}
	})
}

// checkFuncForMathRand checks a function for math/rand in security context.
//
// Params:
//   - pass: analysis context
//   - funcDecl: function declaration
//   - aliases: math/rand aliases
func checkFuncForMathRand(
	pass *analysis.Pass,
	funcDecl *ast.FuncDecl,
	aliases map[string]bool,
) {
	// Vérification si le nom de fonction indique un contexte sécurité
	funcName := funcDecl.Name.Name
	funcIsSecurityContext := isSecurityName(funcName)

	// Si pas de corps
	if funcDecl.Body == nil {
		// Pas de corps
		return
	}

	// Parcours du corps de la fonction
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Si fonction en contexte sécurité, vérifier tous les appels
		if funcIsSecurityContext {
			// Vérification de la condition
			if call, ok := n.(*ast.CallExpr); ok {
				// Vérification de la condition
				if isMathRandCall(call, aliases) {
					reportMathRandInSecurity(pass, call.Pos())
				}
			}
			// Continuer le parcours (pas besoin de checker var names)
			return true
		}

		// Fonction pas en contexte sécurité: vérifier les noms de variables
		// Vérification des assignations avec contexte variable
		if assign, ok := n.(*ast.AssignStmt); ok {
			checkAssignForMathRand(pass, assign, aliases, false)
		}

		// Vérification des déclarations de variables locales
		if decl, ok := n.(*ast.DeclStmt); ok {
			// Vérification de la condition
			if genDecl, ok := decl.Decl.(*ast.GenDecl); ok {
				checkLocalVarDecl(pass, genDecl, aliases, false)
			}
		}

		// Continuer le parcours
		return true
	})
}

// checkAssignForMathRand checks an assignment for math/rand in security context.
//
// Params:
//   - pass: analysis context
//   - assign: assignment statement
//   - aliases: math/rand aliases
//   - funcIsSecurityContext: whether function name is security-related
func checkAssignForMathRand(
	pass *analysis.Pass,
	assign *ast.AssignStmt,
	aliases map[string]bool,
	funcIsSecurityContext bool,
) {
	// Vérification des valeurs assignées
	for i, rhs := range assign.Rhs {
		// Vérification si c'est un appel math/rand
		call, ok := rhs.(*ast.CallExpr)
		// Vérification de la condition
		if !ok || !isMathRandCall(call, aliases) {
			// Pas un appel math/rand
			continue
		}

		// Vérification du contexte sécurité
		varIsSecurityContext := false
		// Vérification de la condition
		if i < len(assign.Lhs) {
			// Récupération du nom de la variable
			if ident, ok := assign.Lhs[i].(*ast.Ident); ok {
				varIsSecurityContext = isSecurityName(ident.Name)
			}
		}

		// Si contexte sécurité détecté
		if funcIsSecurityContext || varIsSecurityContext {
			reportMathRandInSecurity(pass, call.Pos())
		}
	}
}

// checkLocalVarDecl checks local variable declarations for math/rand.
//
// Params:
//   - pass: analysis context
//   - genDecl: general declaration
//   - aliases: math/rand aliases
//   - funcIsSecurityContext: whether function name is security-related
func checkLocalVarDecl(
	pass *analysis.Pass,
	genDecl *ast.GenDecl,
	aliases map[string]bool,
	funcIsSecurityContext bool,
) {
	// Parcours des specs
	for _, spec := range genDecl.Specs {
		// Process each spec
		processLocalVarSpec(pass, spec, aliases, funcIsSecurityContext)
	}
}

// processLocalVarSpec processes a single spec for math/rand usage.
//
// Params:
//   - pass: analysis context
//   - spec: spec to check
//   - aliases: math/rand aliases
//   - funcIsSecurityContext: whether function name is security-related
func processLocalVarSpec(
	pass *analysis.Pass,
	spec ast.Spec,
	aliases map[string]bool,
	funcIsSecurityContext bool,
) {
	// Vérification si c'est une ValueSpec
	valueSpec, ok := spec.(*ast.ValueSpec)
	// Vérification de la condition
	if !ok {
		// Pas une ValueSpec
		return
	}

	// Vérification des noms de variables
	varIsSecurityContext := hasSecurityVarName(valueSpec.Names)

	// Vérification des valeurs
	checkValuesForMathRand(pass, valueSpec.Values, aliases, funcIsSecurityContext || varIsSecurityContext)
}

// hasSecurityVarName checks if any variable name is security-related.
//
// Params:
//   - names: list of identifiers to check
//
// Returns:
//   - bool: true if any name is security-related
func hasSecurityVarName(names []*ast.Ident) bool {
	// Parcours des noms
	for _, name := range names {
		// Vérification de la condition
		if isSecurityName(name.Name) {
			// Contexte sécurité détecté
			return true
		}
	}
	// Pas de contexte sécurité
	return false
}

// checkValuesForMathRand checks values for math/rand calls.
//
// Params:
//   - pass: analysis context
//   - values: values to check
//   - aliases: math/rand aliases
//   - isSecurityContext: whether we're in security context
func checkValuesForMathRand(
	pass *analysis.Pass,
	values []ast.Expr,
	aliases map[string]bool,
	isSecurityContext bool,
) {
	// Vérification des valeurs
	for _, value := range values {
		// Vérification si c'est un appel math/rand
		call, ok := value.(*ast.CallExpr)
		// Vérification de la condition
		if !ok || !isMathRandCall(call, aliases) {
			// Pas un appel math/rand
			continue
		}

		// Si contexte sécurité détecté
		if isSecurityContext {
			reportMathRandInSecurity(pass, call.Pos())
		}
	}
}

// checkGenDeclForMathRand checks package-level declarations for math/rand.
//
// Params:
//   - pass: analysis context
//   - genDecl: general declaration
//   - aliases: math/rand aliases
func checkGenDeclForMathRand(
	pass *analysis.Pass,
	genDecl *ast.GenDecl,
	aliases map[string]bool,
) {
	// Parcours des specs
	for _, spec := range genDecl.Specs {
		// Process each spec
		processGenDeclSpec023(pass, spec, aliases)
	}
}

// processGenDeclSpec023 processes a single spec for math/rand in security context.
//
// Params:
//   - pass: analysis context
//   - spec: spec to check
//   - aliases: math/rand aliases
func processGenDeclSpec023(
	pass *analysis.Pass,
	spec ast.Spec,
	aliases map[string]bool,
) {
	// Vérification si c'est une ValueSpec
	valueSpec, ok := spec.(*ast.ValueSpec)
	// Vérification de la condition
	if !ok {
		// Pas une ValueSpec
		return
	}

	// Vérification des noms de variables via helper existant
	isSecurityContext := hasSecurityVarName(valueSpec.Names)

	// Si pas de contexte sécurité, ne rien faire
	if !isSecurityContext {
		// Pas de contexte sécurité
		return
	}

	// Vérification des valeurs via helper existant
	checkValuesForMathRand(pass, valueSpec.Values, aliases, true)
}

// isMathRandCall checks if a call expression is a math/rand call.
//
// Params:
//   - call: call expression
//   - aliases: math/rand aliases
//
// Returns:
//   - bool: true if math/rand call
func isMathRandCall(call *ast.CallExpr, aliases map[string]bool) bool {
	// Vérification si c'est un sélecteur (rand.XXX)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	// Vérification de la condition
	if !ok {
		// Pas un sélecteur
		return false
	}

	// Vérification de l'identifiant
	ident, ok := sel.X.(*ast.Ident)
	// Vérification de la condition
	if !ok {
		// Pas un identifiant
		return false
	}

	// Vérification si c'est un alias de math/rand
	return aliases[ident.Name]
}

// isSecurityName checks if a name indicates a security context.
//
// Params:
//   - name: name to check
//
// Returns:
//   - bool: true if security context
func isSecurityName(name string) bool {
	// Conversion en minuscules pour comparaison
	lower := strings.ToLower(name)

	// Parcours des mots-clés
	for _, keyword := range securityKeywords {
		// Vérification de la condition
		if strings.Contains(lower, keyword) {
			// Contexte sécurité détecté
			return true
		}
	}

	// Pas de contexte sécurité
	return false
}

// reportMathRandInSecurity reports math/rand usage in security context.
//
// Params:
//   - pass: analysis context
//   - pos: position to report
func reportMathRandInSecurity(pass *analysis.Pass, pos token.Pos) {
	// Récupération du message
	msg, _ := messages.Get(ruleCodeVar023)

	// Rapport de l'erreur
	pass.Reportf(
		pos,
		"%s: %s",
		ruleCodeVar023,
		msg.Format(config.Get().Verbose),
	)
}
