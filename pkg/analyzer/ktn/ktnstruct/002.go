// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct002 code de la règle KTN-STRUCT-002
	ruleCodeStruct002 string = "KTN-STRUCT-002"
)

// Analyzer002 vérifie que les structs exportées avec méthodes ont un constructeur
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct002",
	Doc:      "KTN-STRUCT-002: Struct exportée avec méthodes doit avoir un constructeur NewX()",
	Run:      runStruct002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct002 exécute l'analyse KTN-STRUCT-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct002(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct002) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Parcourir chaque fichier du package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct002, filename) {
			// Fichier exclu
			continue
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Continuer avec le fichier suivant
			continue
		}

		// Collecter les structs exportées et leurs méthodes
		structs := collectExportedStructsWithMethods(file, pass, insp)

		// Collecter les constructeurs disponibles
		constructors := collectConstructors(file)

		// Vérifier chaque struct
		for _, s := range structs {
			// Si la struct n'a pas de méthodes publiques, skip
			if len(s.methods) == 0 {
				// Continuer avec la struct suivante
				continue
			}

			// Exception: les DTOs n'ont pas besoin de constructeur
			if shared.IsSerializableStruct(s.structType, s.name) {
				// DTO - pas besoin de constructeur
				continue
			}

			// Chercher un constructeur pour cette struct
			expectedName := "New" + s.name
			// Vérification si constructeur trouvé
			if !hasConstructor(constructors, expectedName, s.name) {
				msg, _ := messages.Get(ruleCodeStruct002)
				pass.Reportf(
					s.node.Pos(),
					"%s: %s",
					ruleCodeStruct002,
					msg.Format(config.Get().Verbose, s.name, expectedName),
				)
			}
		}
	}

	// Retour de la fonction
	return nil, nil
}

// collectExportedStructsWithMethods collecte les structs exportées et leurs méthodes.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//   - insp: inspector
//
// Returns:
//   - []structWithMethods: liste des structs avec méthodes
func collectExportedStructsWithMethods(file *ast.File, pass *analysis.Pass, _insp *inspector.Inspector) []structWithMethods {
	// Collecter les méthodes
	methodsByStruct := shared.CollectMethodsByStruct(file, pass)

	// Collecter les structs exportées du fichier
	var structs []structWithMethods
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier si c'est une TypeSpec
		typeSpec, ok := n.(*ast.TypeSpec)
		// Si ce n'est pas une TypeSpec, continuer
		if !ok {
			// Continue traversal
			return true
		}

		// Vérifier si c'est une struct
		structType, isStruct := typeSpec.Type.(*ast.StructType)
		// Si c'est une struct ET exportée
		if isStruct && ast.IsExported(typeSpec.Name.Name) {
			structs = append(structs, structWithMethods{
				name:       typeSpec.Name.Name,
				node:       typeSpec,
				structType: structType,
				methods:    methodsByStruct[typeSpec.Name.Name],
			})
		}

		// Continue traversal
		return true
	})

	// Retour de la liste
	return structs
}

// constructorInfo stocke les informations d'un constructeur.
type constructorInfo struct {
	name       string
	returnType string
}

// collectConstructors collecte tous les constructeurs du fichier.
//
// Params:
//   - file: fichier AST
//
// Returns:
//   - []constructorInfo: liste des constructeurs
func collectConstructors(file *ast.File) []constructorInfo {
	var constructors []constructorInfo

	// Parcourir les fonctions du fichier
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier FuncDecl
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification si FuncDecl
		if !ok {
			// Continue traversal
			return true
		}

		// Ignorer les méthodes (avec receiver)
		if funcDecl.Recv != nil {
			// Continue traversal
			return true
		}

		// Vérifier que le nom commence par "New"
		if !strings.HasPrefix(funcDecl.Name.Name, "New") {
			// Continue traversal
			return true
		}

		// Vérifier qu'il y a un type de retour
		if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) == 0 {
			// Continue traversal
			return true
		}

		// Extraire le type de retour
		returnType := extractReturnTypeName(funcDecl.Type.Results)
		// Si type de retour valide
		if returnType != "" {
			constructors = append(constructors, constructorInfo{
				name:       funcDecl.Name.Name,
				returnType: returnType,
			})
		}

		// Continue traversal
		return true
	})

	// Retour de la liste
	return constructors
}

// extractReturnTypeName extrait le nom du type de retour.
//
// Params:
//   - results: liste des résultats
//
// Returns:
//   - string: nom du type
func extractReturnTypeName(results *ast.FieldList) string {
	// Si pas de résultats
	if results == nil || len(results.List) == 0 {
		// Retour vide
		return ""
	}

	// Prendre le premier type de retour
	firstResult := results.List[0].Type

	// Gérer les différents types
	switch t := firstResult.(type) {
	// Traitement du pointeur
	case *ast.StarExpr:
		// Retour de type pointeur (*T)
		if ident, ok := t.X.(*ast.Ident); ok {
			// Retour du nom du type extrait
			return ident.Name
		}
	// Traitement de l'identifiant
	case *ast.Ident:
		// Retour de type direct (T)
		return t.Name
	}

	// Type non géré
	return ""
}

// hasConstructor vérifie si un constructeur existe pour la struct.
//
// Params:
//   - constructors: liste des constructeurs
//   - expectedName: nom attendu du constructeur
//   - structName: nom de la struct
//
// Returns:
//   - bool: true si constructeur trouvé
func hasConstructor(constructors []constructorInfo, expectedName string, structName string) bool {
	// Parcourir les constructeurs
	for _, c := range constructors {
		// Vérifier le nom ET le type de retour
		if c.name == expectedName && c.returnType == structName {
			// Constructeur trouvé
			return true
		}
	}

	// Constructeur non trouvé
	return false
}
