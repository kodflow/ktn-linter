// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar022 is the rule code for this analyzer
	ruleCodeVar022 string = "KTN-VAR-022"
)

// Analyzer022 détecte les pointeurs vers interfaces.
//
// Une interface est déjà un fat pointer (type + data), un pointeur dessus
// est rarement utile et souvent signe d'une erreur de conception.
var Analyzer022 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar022",
	Doc:      "KTN-VAR-022: Détecte les pointeurs vers interfaces (*io.Reader, *interface{}, *any)",
	Run:      runVar022,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar022 exécute l'analyse de détection des pointeurs vers interfaces.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: erreur éventuelle
func runVar022(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar022) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp := inspAny.(*inspector.Inspector)
	// Defensive: avoid nil dereference when resolving positions
	if pass.Fset == nil {
		// Cannot analyze without file set
		return nil, nil
	}
	// Defensive: avoid nil dereference when resolving types
	if pass.TypesInfo == nil {
		// Cannot analyze without type information
		return nil, nil
	}

	// Analyse des déclarations de fonctions (paramètres et retours)
	checkFuncDecls(pass, insp, cfg)

	// Analyse des déclarations de variables
	checkVarDecls(pass, insp, cfg)

	// Analyse des champs de struct
	checkStructFields(pass, insp, cfg)

	// Traitement terminé
	return nil, nil
}

// checkFuncDecls vérifie les paramètres et retours de fonctions.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
func checkFuncDecls(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcours des fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en fonction
		funcDecl := n.(*ast.FuncDecl)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar022, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Vérification des paramètres
		if funcDecl.Type.Params != nil {
			// Parcours des paramètres
			checkFieldList(pass, funcDecl.Type.Params.List)
		}

		// Vérification des retours
		if funcDecl.Type.Results != nil {
			// Parcours des retours
			checkFieldList(pass, funcDecl.Type.Results.List)
		}
	})
}

// checkVarDecls vérifie les déclarations de variables.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
func checkVarDecls(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	// Parcours des déclarations
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en déclaration générale
		genDecl := n.(*ast.GenDecl)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar022, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Vérification que c'est une déclaration var
		if genDecl.Tok.String() != "var" {
			// Pas une déclaration var
			return
		}

		// Parcours des spécifications
		for _, spec := range genDecl.Specs {
			// Cast en value spec
			valueSpec, ok := spec.(*ast.ValueSpec)
			// Vérification du cast
			if !ok {
				// Pas une value spec
				continue
			}

			// Vérification du type explicite
			if valueSpec.Type != nil {
				// Vérification si c'est un pointeur vers interface
				checkPointerToInterface(pass, valueSpec.Type)
			}
		}
	})
}

// checkStructFields vérifie les champs de structures.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
func checkStructFields(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	// Parcours des types
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en type spec
		typeSpec := n.(*ast.TypeSpec)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar022, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Vérification que c'est une struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		// Vérification du cast
		if !ok {
			// Pas une struct
			return
		}

		// Vérification des champs
		if structType.Fields != nil {
			// Parcours des champs
			checkFieldList(pass, structType.Fields.List)
		}
	})
}

// checkFieldList vérifie une liste de champs.
//
// Params:
//   - pass: contexte d'analyse
//   - fields: liste de champs
func checkFieldList(pass *analysis.Pass, fields []*ast.Field) {
	// Parcours des champs
	for _, field := range fields {
		// Vérification du type
		checkPointerToInterface(pass, field.Type)
	}
}

// checkPointerToInterface vérifie si un type est un pointeur vers interface.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression de type
func checkPointerToInterface(pass *analysis.Pass, expr ast.Expr) {
	// Vérification si c'est un pointeur
	starExpr, ok := expr.(*ast.StarExpr)
	// Vérification du cast
	if !ok {
		// Pas un pointeur
		return
	}

	// Récupération du type sous-jacent
	underlyingType := pass.TypesInfo.TypeOf(starExpr.X)
	// Vérification du type
	if underlyingType == nil {
		// Type non trouvé
		return
	}

	// Vérification si c'est une interface
	if isInterfaceType(underlyingType) {
		// Rapport d'erreur
		msg, ok := messages.Get(ruleCodeVar022)
		// Defensive: avoid panic if message is missing
		if !ok {
			// Fallback message
			pass.Reportf(expr.Pos(), "%s: éviter pointeur vers interface %s", ruleCodeVar022, underlyingType.String())

			// Stop after reporting
			return
		}
		// Report avec message formatte
		pass.Reportf(
			expr.Pos(),
			"%s: %s",
			ruleCodeVar022,
			msg.Format(config.Get().Verbose, underlyingType.String()),
		)
	}
}

// isInterfaceType vérifie si un type est une interface.
//
// Params:
//   - t: type à vérifier
//
// Returns:
//   - bool: true si c'est une interface
func isInterfaceType(t types.Type) bool {
	// Récupération du type sous-jacent
	underlying := t.Underlying()
	// Vérification si c'est une interface
	_, ok := underlying.(*types.Interface)
	// Retour du résultat
	return ok
}
