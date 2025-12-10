// Package ktnstruct implements KTN linter rules.
package ktnstruct

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

const (
	// ruleCodeStruct004 code de la règle KTN-STRUCT-004
	ruleCodeStruct004 string = "KTN-STRUCT-004"
)

// Analyzer004 vérifie qu'il n'y a qu'une seule struct par fichier Go
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktnstruct004",
	Doc:  "KTN-STRUCT-004: Un fichier Go ne doit contenir qu'une seule struct (évite les fichiers de 10000 lignes)",
	Run:  runStruct004,
}

// structInfo stocke les informations d'une struct trouvée
type structInfo struct {
	name       string
	node       *ast.TypeSpec
	structType *ast.StructType
}

// runStruct004 exécute l'analyse KTN-STRUCT-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct004(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct004) {
		// Règle désactivée
		return nil, nil
	}

	// Parcourir chaque fichier du package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeStruct004, filename) {
			// Fichier exclu
			continue
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Continuer avec le fichier suivant
			continue
		}

		// Collecter toutes les structs du fichier
		structs := collectStructs(file)

		// If more than one struct, check if they are all DTOs
		if len(structs) > 1 {
			// Exception: si toutes les structs sont des DTOs/serializable liés
			if allStructsAreSerializable(structs) {
				// DTOs groupés sont autorisés
				continue
			}

			// Itération sur les structs (à partir de la 2ème)
			for i := 1; i < len(structs); i++ {
				s := structs[i]
				pass.Reportf(
					s.node.Pos(),
					"KTN-STRUCT-004: le fichier contient plusieurs structs (%d au total). Déplacer '%s' dans un fichier séparé pour respecter le principe 'une struct par fichier'",
					len(structs),
					s.name,
				)
			}
		}
	}

	// Retour de la fonction
	return nil, nil
}

// allStructsAreSerializable vérifie si toutes les structs sont des DTOs.
//
// Params:
//   - structs: liste des structs à vérifier
//
// Returns:
//   - bool: true si toutes sont des DTOs
func allStructsAreSerializable(structs []structInfo) bool {
	// Parcourir les structs
	for _, s := range structs {
		// Vérifier si c'est un DTO
		if !shared.IsSerializableStruct(s.structType, s.name) {
			// Une struct n'est pas un DTO
			return false
		}
	}
	// Toutes sont des DTOs
	return true
}

// collectStructs collecte toutes les déclarations de struct d'un fichier.
//
// Params:
//   - file: fichier AST à analyser
//
// Returns:
//   - []structInfo: liste des structs trouvées
func collectStructs(file *ast.File) []structInfo {
	var structs []structInfo

	// Parcourir l'AST du fichier
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier si c'est une déclaration de type générale
		genDecl, ok := n.(*ast.GenDecl)
		// Si ce n'est pas une GenDecl, continuer
		if !ok {
			// Continue traversal
			return true
		}

		// Parcourir les specs de la déclaration
		for _, spec := range genDecl.Specs {
			typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
			// Si ce n'est pas une TypeSpec, continuer
			if !isTypeSpec {
				// Continue with next spec
				continue
			}

			// Vérifier si le type est une struct
			structType, isStruct := typeSpec.Type.(*ast.StructType)
			// Si c'est une struct, l'ajouter à la liste
			if isStruct {
				structs = append(structs, structInfo{
					name:       typeSpec.Name.Name,
					node:       typeSpec,
					structType: structType,
				})
			}
		}

		// Continue traversal
		return true
	})

	// Retour de la liste des structs
	return structs
}
