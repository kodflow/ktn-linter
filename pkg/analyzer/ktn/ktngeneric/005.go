// Package ktngeneric implements KTN linter rules for generic functions.
package ktngeneric

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeGeneric005 est le code de la regle KTN-GENERIC-005.
	ruleCodeGeneric005 string = "KTN-GENERIC-005"
)

// predeclaredIdentifiers contient tous les identifiants predeclares de Go.
// Types (22): bool, byte, complex64, complex128, error, float32, float64,
//
//	int, int8, int16, int32, int64, rune, string, uint, uint8, uint16,
//	uint32, uint64, uintptr, any, comparable
//
// Constantes (4): true, false, iota, nil
// Fonctions (18): append, cap, clear, close, complex, copy, delete, imag,
//
//	len, make, max, min, new, panic, print, println, real, recover
var predeclaredIdentifiers map[string]bool = map[string]bool{
	// Types
	"bool": true, "byte": true, "complex64": true, "complex128": true,
	"error": true, "float32": true, "float64": true, "int": true,
	"int8": true, "int16": true, "int32": true, "int64": true,
	"rune": true, "string": true, "uint": true, "uint8": true,
	"uint16": true, "uint32": true, "uint64": true, "uintptr": true,
	"any": true, "comparable": true,
	// Constants
	"true": true, "false": true, "iota": true, "nil": true,
	// Functions
	"append": true, "cap": true, "clear": true, "close": true,
	"complex": true, "copy": true, "delete": true, "imag": true,
	"len": true, "make": true, "max": true, "min": true,
	"new": true, "panic": true, "print": true, "println": true,
	"real": true, "recover": true,
}

// Analyzer005 checks that type parameters do not shadow predeclared identifiers.
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktngeneric005",
	Doc:      "KTN-GENERIC-005: Type parameters must not shadow predeclared identifiers",
	Run:      runGeneric005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runGeneric005 execute l'analyse KTN-GENERIC-005.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: resultat de l'analyse
//   - error: erreur eventuelle
func runGeneric005(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeGeneric005) {
		// Regle desactivee
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeGeneric005, filename) {
			// File is excluded
			return
		}

		// Verifier les type parameters selon le type de noeud
		switch node := n.(type) {
		// Cas d'une declaration de fonction
		case *ast.FuncDecl:
			// Analyser les type parameters de la fonction
			checkFuncTypeParams(pass, node)
		// Cas d'une specification de type
		case *ast.TypeSpec:
			// Analyser les type parameters du type
			checkTypeSpecTypeParams(pass, node)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// checkFuncTypeParams verifie les type parameters d'une fonction.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction a analyser
func checkFuncTypeParams(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Verifier si la fonction a des parametres de type
	if funcDecl.Type.TypeParams == nil {
		// Pas une fonction generique
		return
	}

	// Verifier chaque type parameter
	checkTypeParamList(pass, funcDecl.Type.TypeParams)
}

// checkTypeSpecTypeParams verifie les type parameters d'un type.
//
// Params:
//   - pass: contexte d'analyse
//   - typeSpec: specification de type a analyser
func checkTypeSpecTypeParams(pass *analysis.Pass, typeSpec *ast.TypeSpec) {
	// Verifier si le type a des parametres de type
	if typeSpec.TypeParams == nil {
		// Pas un type generique
		return
	}

	// Verifier chaque type parameter
	checkTypeParamList(pass, typeSpec.TypeParams)
}

// checkTypeParamList verifie une liste de parametres de type.
//
// Params:
//   - pass: contexte d'analyse
//   - typeParams: liste des parametres de type
func checkTypeParamList(pass *analysis.Pass, typeParams *ast.FieldList) {
	// Parcourir les type parameters
	for _, field := range typeParams.List {
		// Parcourir les noms de chaque field
		for _, name := range field.Names {
			// Verifier si le nom est un identifiant predeclare
			if predeclaredIdentifiers[name.Name] {
				// Reporter l'erreur
				reportShadowing(pass, name)
			}
		}
	}
}

// reportShadowing reporte une erreur de shadowing.
//
// Params:
//   - pass: contexte d'analyse
//   - name: identifiant du type parameter
func reportShadowing(pass *analysis.Pass, name *ast.Ident) {
	// Guard contre nil (pour tests unitaires)
	if pass == nil {
		// Pas de contexte pour reporter
		return
	}

	// Recuperer le message
	cfg := config.Get()
	msg, _ := messages.Get(ruleCodeGeneric005)

	// Reporter l'erreur
	pass.Reportf(
		name.Pos(),
		"%s: %s",
		ruleCodeGeneric005,
		msg.Format(cfg.Verbose, name.Name),
	)
}
