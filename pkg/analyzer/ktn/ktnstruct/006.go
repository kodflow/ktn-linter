package ktnstruct

import (
	"go/ast"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MAX_SIMPLE_FIELDS nombre maximum de champs pour une struct simple (config)
	MAX_SIMPLE_FIELDS int = 3
)

// Analyzer006 vérifie l'encapsulation des structs avec méthodes
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct006",
	Doc:      "KTN-STRUCT-006: Struct exportée avec méthodes (>3 champs) doit avoir champs privés + getters",
	Run:      runStruct006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct006 exécute l'analyse KTN-STRUCT-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct006(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Parcourir chaque fichier du package
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Continuer avec le fichier suivant
			continue
		}

		// Collecter les structs exportées avec méthodes
		structs := collectExportedStructsWithMethods(file, pass, insp)

		// Vérifier chaque struct
		for _, s := range structs {
			// Si pas de méthodes publiques, skip (DTO)
			if len(s.methods) == 0 {
				// Continuer avec la struct suivante
				continue
			}

			// Collecter les champs de la struct
			fields := collectStructFields(s.node)

			// Si ≤3 champs, pas de règle stricte (config simple)
			if len(fields) <= MAX_SIMPLE_FIELDS {
				// Continuer avec la struct suivante
				continue
			}

			// Vérifier l'encapsulation
			checkEncapsulation(pass, s, fields)
		}
	}

	// Retour de la fonction
	return nil, nil
}

// collectStructFields collecte les champs d'une struct.
//
// Params:
//   - typeSpec: spécification du type
//
// Returns:
//   - []structField: liste des champs
func collectStructFields(typeSpec *ast.TypeSpec) []structField {
	var fields []structField

	// Récupérer le type struct
	structType, ok := typeSpec.Type.(*ast.StructType)
	// Vérification du type
	if !ok || structType.Fields == nil {
		// Retour vide
		return fields
	}

	// Parcourir les champs
	for _, field := range structType.Fields.List {
		// Parcourir les noms de champs
		for _, name := range field.Names {
			fields = append(fields, structField{
				name:     name.Name,
				exported: ast.IsExported(name.Name),
			})
		}
	}

	// Retour de la liste
	return fields
}

// checkEncapsulation vérifie l'encapsulation d'une struct.
//
// Params:
//   - pass: contexte d'analyse
//   - s: struct avec méthodes
//   - fields: champs de la struct
func checkEncapsulation(pass *analysis.Pass, s structWithMethods, fields []structField) {
	// Compter les champs publics
	var publicFields []string
	// Collecter les champs privés
	var privateFields []string

	// Parcourir les champs
	for _, f := range fields {
		// Vérification si exporté
		if f.exported {
			publicFields = append(publicFields, f.name)
		} else {
			// Champ privé
			privateFields = append(privateFields, f.name)
		}
	}

	// Si des champs publics, violation
	if len(publicFields) > 0 {
		pass.Reportf(
			s.node.Pos(),
			"KTN-STRUCT-006: la struct exportée '%s' a des méthodes et %d champ(s) public(s) [%s]. Utiliser des champs privés avec getters",
			s.name,
			len(publicFields),
			strings.Join(publicFields, ", "),
		)
		// Retour
		return
	}

	// Si tous les champs sont privés, vérifier les getters
	checkGetters(pass, s, privateFields)
}

// checkGetters vérifie la présence des getters.
//
// Params:
//   - pass: contexte d'analyse
//   - s: struct avec méthodes
//   - privateFields: liste des champs privés
func checkGetters(pass *analysis.Pass, s structWithMethods, privateFields []string) {
	// Collecter les noms des getters existants
	var existingGetters []string
	// Parcourir les méthodes
	for _, method := range s.methods {
		// Vérifier si méthode commence par Get
		if strings.HasPrefix(method.Name, "Get") {
			existingGetters = append(existingGetters, method.Name)
		}
	}

	// Vérifier chaque champ privé
	var missingGetters []string
	// Parcourir les champs privés
	for _, field := range privateFields {
		// Construire le nom du getter attendu
		expectedGetter := "Get" + capitalize(field)

		// Vérifier si getter existe
		if !contains(existingGetters, expectedGetter) {
			missingGetters = append(missingGetters, expectedGetter)
		}
	}

	// Si des getters manquants, violation
	if len(missingGetters) > 0 {
		pass.Reportf(
			s.node.Pos(),
			"KTN-STRUCT-006: la struct exportée '%s' a des champs privés mais il manque %d getter(s): %s",
			s.name,
			len(missingGetters),
			strings.Join(missingGetters, ", "),
		)
	}
}

// capitalize met la première lettre en majuscule.
//
// Params:
//   - s: chaîne à capitaliser
//
// Returns:
//   - string: chaîne capitalisée
func capitalize(s string) string {
	// Vérification chaîne vide
	if s == "" {
		// Retour vide
		return ""
	}

	// Conversion en runes
	runes := []rune(s)
	// Mise en majuscule de la première rune
	runes[0] = unicode.ToUpper(runes[0])

	// Retour de la chaîne
	return string(runes)
}

// contains vérifie si une slice contient une valeur.
//
// Params:
//   - slice: slice à vérifier
//   - value: valeur recherchée
//
// Returns:
//   - bool: true si trouvé
func contains(slice []string, value string) bool {
	// Parcourir la slice
	for _, item := range slice {
		// Vérification si égal
		if item == value {
			// Trouvé
			return true
		}
	}

	// Non trouvé
	return false
}

// structField représente un champ de struct.
type structField struct {
	name     string
	exported bool
}
