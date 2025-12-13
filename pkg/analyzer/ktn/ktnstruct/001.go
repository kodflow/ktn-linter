// Analyzer 001 for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct001 code de la règle KTN-STRUCT-001
	ruleCodeStruct001 string = "KTN-STRUCT-001"
	// interfaceChecksCap capacité initiale pour la map des interface checks
	interfaceChecksCap int = 8
)

// Analyzer001 vérifie qu'une interface existe pour chaque struct avec méthodes publiques
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct001",
	Doc:      "KTN-STRUCT-001: Chaque struct doit avoir une interface reprenant 100% de ses méthodes publiques",
	Run:      runStruct001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// structWithMethods stocke une struct et ses méthodes publiques
type structWithMethods struct {
	name       string
	node       *ast.TypeSpec
	structType *ast.StructType
	methods    []shared.MethodSignature
}

// runStruct001 exécute l'analyse KTN-STRUCT-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct001(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct001) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Parcourir chaque fichier du package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct001, filename) {
			// Fichier exclu
			continue
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Continuer avec le fichier suivant
			continue
		}

		// Collecter les interfaces et leurs méthodes
		interfaces := collectInterfaces(file, pass)

		// Collecter les compile-time interface checks (var _ Interface = (*Struct)(nil))
		interfaceChecks := collectInterfaceChecks(file)

		// Collecter les structs et leurs méthodes
		structs := collectStructsWithMethods(file, pass, insp)

		// Vérifier chaque struct
		for _, s := range structs {
			// Si la struct n'a pas de méthodes publiques, skip
			if len(s.methods) == 0 {
				// Continuer avec la struct suivante
				continue
			}

			// Exception: les DTOs n'ont pas besoin d'interface
			if shared.IsSerializableStruct(s.structType, s.name) {
				// DTO - pas besoin d'interface
				continue
			}

			// Exception: struct avec compile-time interface check (DDD pattern)
			if interfaceChecks[s.name] {
				// Interface existe dans un autre package
				continue
			}

			// Trouver une interface qui couvre toutes les méthodes
			if !hasMatchingInterface(s, interfaces) {
				pass.Reportf(
					s.node.Pos(),
					"KTN-STRUCT-001: la struct '%s' a %d méthode(s) publique(s) mais aucune interface ne les reprend toutes. Créer une interface complète dans le même fichier",
					s.name,
					len(s.methods),
				)
			}
		}
	}

	// Retour de la fonction
	return nil, nil
}

// collectInterfaces collecte toutes les interfaces et leurs méthodes.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string][]shared.MethodSignature: map nom interface -> signatures méthodes
func collectInterfaces(file *ast.File, pass *analysis.Pass) map[string][]shared.MethodSignature {
	interfaces := make(map[string][]shared.MethodSignature, 0)

	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier si c'est une TypeSpec
		typeSpec, ok := n.(*ast.TypeSpec)
		// Si ce n'est pas une TypeSpec, continuer
		if !ok {
			// Continue traversal
			return true
		}

		// Vérifier si c'est une interface
		ifaceType, isInterface := typeSpec.Type.(*ast.InterfaceType)
		// Si ce n'est pas une interface, continuer
		if !isInterface {
			// Continue traversal
			return true
		}

		// Collecter les méthodes de l'interface
		var methods []shared.MethodSignature
		// Parcourir les méthodes de l'interface
		for _, method := range ifaceType.Methods.List {
			// Vérifier si c'est une méthode (pas un embedded interface)
			funcType, isFunc := method.Type.(*ast.FuncType)
			// Si ce n'est pas une fonction, continuer
			if !isFunc {
				// Continue with next method
				continue
			}

			// Extraire le nom de la méthode
			for _, name := range method.Names {
				methods = append(methods, shared.MethodSignature{
					Name:       name.Name,
					ParamsStr:  formatFieldList(funcType.Params, pass),
					ResultsStr: formatFieldList(funcType.Results, pass),
				})
			}
		}

		interfaces[typeSpec.Name.Name] = methods
		// Continue traversal
		return true
	})

	// Retour de la map
	return interfaces
}

// collectInterfaceChecks collecte les compile-time interface checks.
// Pattern: var _ InterfaceName = (*StructName)(nil)
// Cela prouve que la struct implémente une interface (même cross-package).
//
// Params:
//   - file: fichier AST
//
// Returns:
//   - map[string]bool: noms des structs avec interface check
func collectInterfaceChecks(file *ast.File) map[string]bool {
	checks := make(map[string]bool, interfaceChecksCap)

	// Parcourir les déclarations
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		// Vérifier si c'est une déclaration var
		if !ok || genDecl.Tok != token.VAR {
			// Continuer l'itération
			continue
		}

		// Parcourir les specs
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			// Vérifier si c'est une ValueSpec
			if !ok {
				// Continuer l'itération
				continue
			}

			// Vérifier le pattern: var _ Interface = (*Struct)(nil)
			structName := extractStructFromInterfaceCheck(valueSpec)
			// Ajouter à la map si trouvé
			if structName != "" {
				checks[structName] = true
			}
		}
	}

	// Retour de la map
	return checks
}

// extractStructFromInterfaceCheck extrait le nom de struct d'un interface check.
// Pattern: var _ Interface = (*Struct)(nil)
//
// Params:
//   - spec: ValueSpec à analyser
//
// Returns:
//   - string: nom de la struct ou vide si pas trouvé
func extractStructFromInterfaceCheck(spec *ast.ValueSpec) string {
	// Vérifier que le nom est "_"
	if len(spec.Names) != 1 || spec.Names[0].Name != "_" {
		// Pas le pattern attendu
		return ""
	}

	// Vérifier qu'il y a une valeur
	if len(spec.Values) != 1 {
		// Pas de valeur
		return ""
	}

	// La valeur doit être une conversion: (*Struct)(nil)
	callExpr, ok := spec.Values[0].(*ast.CallExpr)
	// Vérifier si c'est un CallExpr
	if !ok {
		// Pas une conversion
		return ""
	}

	// Extraire le type de la conversion
	return extractStructNameFromConversion(callExpr)
}

// extractStructNameFromConversion extrait le nom de struct d'une conversion.
// Pattern: (*Struct)(nil)
//
// Params:
//   - callExpr: expression de conversion
//
// Returns:
//   - string: nom de la struct ou vide si pas trouvé
func extractStructNameFromConversion(callExpr *ast.CallExpr) string {
	// Le Fun doit être un type pointeur
	parenExpr, ok := callExpr.Fun.(*ast.ParenExpr)
	// Vérifier si c'est une expression parenthésée
	if !ok {
		// Pas le pattern attendu
		return ""
	}

	// Le type dans les parenthèses doit être *Struct
	starExpr, ok := parenExpr.X.(*ast.StarExpr)
	// Vérifier si c'est un pointeur
	if !ok {
		// Pas un pointeur
		return ""
	}

	// Extraire le nom de la struct
	ident, ok := starExpr.X.(*ast.Ident)
	// Vérifier si c'est un identifiant
	if ok {
		// Retour du nom
		return ident.Name
	}

	// Retour vide
	return ""
}

// extractStructNameFromReceiver extrait le nom de la struct depuis le receiver.
//
// Params:
//   - recvType: type du receiver
//
// Returns:
//   - string: nom de la struct
func extractStructNameFromReceiver(recvType ast.Expr) string {
	var structName string
	// Gérer les receivers de type *T ou T
	switch t := recvType.(type) {
	// Traitement
	case *ast.StarExpr:
		// Receiver de type *T
		if ident, ok := t.X.(*ast.Ident); ok {
			structName = ident.Name
		}
	// Traitement
	case *ast.Ident:
		// Receiver de type T
		structName = t.Name
	}
	// Retour du nom
	return structName
}

// collectMethodsByStruct collecte les méthodes publiques pour chaque struct.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string][]shared.MethodSignature: map des méthodes par struct
func collectMethodsByStruct(file *ast.File, pass *analysis.Pass) map[string][]shared.MethodSignature {
	methodsByStruct := make(map[string][]shared.MethodSignature, 0)

	// Collecter les méthodes du fichier
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier FuncDecl
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification si FuncDecl
		if !ok {
			// Continue traversal
			return true
		}

		// Vérifier receiver
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			// Pas de receiver
			return true
		}

		// Vérifier méthode publique
		if !ast.IsExported(funcDecl.Name.Name) {
			// Méthode privée
			return true
		}

		// Extraire nom de struct
		recvType := funcDecl.Recv.List[0].Type
		structName := extractStructNameFromReceiver(recvType)

		// Ajouter méthode
		if structName != "" {
			methodsByStruct[structName] = append(methodsByStruct[structName], shared.MethodSignature{
				Name:       funcDecl.Name.Name,
				ParamsStr:  formatFieldList(funcDecl.Type.Params, pass),
				ResultsStr: formatFieldList(funcDecl.Type.Results, pass),
			})
		}

		// Continue traversal
		return true
	})

	// Retour de la map
	return methodsByStruct
}

// collectStructsWithMethods collecte les structs et leurs méthodes publiques.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//   - insp: inspector
//
// Returns:
//   - []structWithMethods: liste des structs avec méthodes
func collectStructsWithMethods(file *ast.File, pass *analysis.Pass, _insp *inspector.Inspector) []structWithMethods {
	// Collecter les méthodes
	methodsByStruct := collectMethodsByStruct(file, pass)

	// Collecter les structs du fichier
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
		// Si c'est une struct
		if isStruct {
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

// hasMatchingInterface vérifie si une interface couvre toutes les méthodes.
//
// Params:
//   - s: struct avec méthodes
//   - interfaces: map des interfaces
//
// Returns:
//   - bool: true si une interface matching existe
func hasMatchingInterface(s structWithMethods, interfaces map[string][]shared.MethodSignature) bool {
	// Parcourir toutes les interfaces
	for _, ifaceMethods := range interfaces {
		// Vérifier si cette interface couvre toutes les méthodes de la struct
		if interfaceCoversAllMethods(s.methods, ifaceMethods) {
			// Interface trouvée
			return true
		}
	}

	// Aucune interface ne couvre toutes les méthodes
	return false
}

// interfaceCoversAllMethods vérifie si l'interface couvre toutes les méthodes.
//
// Params:
//   - structMethods: méthodes de la struct
//   - ifaceMethods: méthodes de l'interface
//
// Returns:
//   - bool: true si toutes les méthodes sont couvertes
func interfaceCoversAllMethods(structMethods []shared.MethodSignature, ifaceMethods []shared.MethodSignature) bool {
	// Chaque méthode de la struct doit être dans l'interface
	for _, sm := range structMethods {
		found := false
		// Chercher la méthode dans l'interface
		for _, im := range ifaceMethods {
			// Comparer nom et signatures
			if sm.Name == im.Name && sm.ParamsStr == im.ParamsStr && sm.ResultsStr == im.ResultsStr {
				found = true
				// Sortir de la boucle
				break
			}
		}

		// Si une méthode n'est pas trouvée, l'interface ne couvre pas tout
		if !found {
			// Retour false
			return false
		}
	}

	// Toutes les méthodes sont couvertes
	return true
}

// formatFieldList formate une liste de champs en string.
//
// Params:
//   - fields: liste de champs
//   - pass: contexte d'analyse
//
// Returns:
//   - string: représentation string
func formatFieldList(fields *ast.FieldList, _pass *analysis.Pass) string {
	// Si pas de champs
	if fields == nil {
		// Retour vide
		return ""
	}

	var parts []string
	// Parcourir les champs
	for _, field := range fields.List {
		typeStr := types.ExprString(field.Type)
		parts = append(parts, typeStr)
	}

	// Retour de la string jointe
	return strings.Join(parts, ",")
}
