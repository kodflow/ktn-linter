// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct009 code de la règle KTN-STRUCT-009
	ruleCodeStruct009 string = "KTN-STRUCT-009"
	// receiverMapCap capacité initiale pour la map des receivers
	receiverMapCap int = 16
)

// badReceiverNames noms de receiver génériques à éviter.
var badReceiverNames = map[string]bool{
	"me":   true,
	"this": true,
	"self": true,
}

// Analyzer009 checks receiver name consistency across methods.
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct009",
	Doc:      "KTN-STRUCT-009: Les noms de receiver doivent être cohérents entre méthodes",
	Run:      runStruct009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// receiverInfo contient les informations sur un receiver.
type receiverInfo struct {
	name     string
	funcDecl *ast.FuncDecl
}

// runStruct009 exécute l'analyse KTN-STRUCT-009.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct009(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct009) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter les receivers par type
	receiversByType := collectReceivers(pass, insp, cfg)

	// Vérifier la cohérence des noms
	checkReceiverConsistency(pass, receiversByType)

	// Vérifier les noms génériques
	checkGenericReceiverNames(pass, receiversByType)

	// Retour de la fonction
	return nil, nil
}

// collectReceivers collecte les receivers par type.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
//
// Returns:
//   - map[string][]receiverInfo: map type -> liste de receivers
func collectReceivers(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) map[string][]receiverInfo {
	result := make(map[string][]receiverInfo, receiverMapCap)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct009, filename) {
			// Fichier exclu
			return
		}

		// Vérifier si c'est une méthode
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			// Pas une méthode
			return
		}

		// Extraire le nom du type receiver
		recv := funcDecl.Recv.List[0]
		typeName := extractReceiverTypeName(recv.Type)
		// Si pas de type valide, ignorer
		if typeName == "" {
			// Retour anticipé
			return
		}

		// Extraire le nom du receiver
		receiverName := extractReceiverName(recv)
		// Si pas de nom, ignorer
		if receiverName == "" {
			// Retour anticipé
			return
		}

		// Ajouter à la map
		result[typeName] = append(result[typeName], receiverInfo{
			name:     receiverName,
			funcDecl: funcDecl,
		})
	})

	// Retour de la map
	return result
}

// extractReceiverTypeName extrait le nom du type du receiver.
//
// Params:
//   - expr: expression du type
//
// Returns:
//   - string: nom du type
func extractReceiverTypeName(expr ast.Expr) string {
	// Gérer le cas *Type
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}

	// Gérer le cas Type
	if ident, ok := expr.(*ast.Ident); ok {
		// Retourner le nom
		return ident.Name
	}

	// Type non reconnu
	return ""
}

// extractReceiverName extrait le nom du receiver.
//
// Params:
//   - field: champ du receiver
//
// Returns:
//   - string: nom du receiver
func extractReceiverName(field *ast.Field) string {
	// Vérifier si le receiver a un nom
	if len(field.Names) == 0 {
		// Pas de nom explicite
		return ""
	}

	// Retourner le premier nom
	return field.Names[0].Name
}

// checkReceiverConsistency vérifie la cohérence des noms de receiver.
//
// Params:
//   - pass: contexte d'analyse
//   - receiversByType: map des receivers par type
func checkReceiverConsistency(pass *analysis.Pass, receiversByType map[string][]receiverInfo) {
	// Parcourir les types
	for typeName, receivers := range receiversByType {
		// Si moins de 2 receivers, pas de vérification
		if len(receivers) < 2 {
			// Continuer au type suivant
			continue
		}

		// Trouver le nom le plus courant
		firstName := receivers[0].name

		// Vérifier la cohérence
		for i := 1; i < len(receivers); i++ {
			recv := receivers[i]
			// Si nom différent, reporter
			if recv.name != firstName {
				msg, _ := messages.Get(ruleCodeStruct009)
				pass.Reportf(
					recv.funcDecl.Recv.List[0].Pos(),
					"%s: %s",
					ruleCodeStruct009,
					msg.Format(config.Get().Verbose, recv.name, firstName, typeName),
				)
			}
		}
	}
}

// checkGenericReceiverNames vérifie les noms de receiver génériques.
//
// Params:
//   - pass: contexte d'analyse
//   - receiversByType: map des receivers par type
func checkGenericReceiverNames(pass *analysis.Pass, receiversByType map[string][]receiverInfo) {
	// Parcourir les types
	for typeName, receivers := range receiversByType {
		// Parcourir les receivers
		for _, recv := range receivers {
			// Vérifier si nom générique
			if badReceiverNames[recv.name] {
				msg, _ := messages.Get(ruleCodeStruct009)
				pass.Reportf(
					recv.funcDecl.Recv.List[0].Pos(),
					"%s: %s",
					ruleCodeStruct009,
					msg.Format(config.Get().Verbose, recv.name, suggestReceiverName(typeName), typeName),
				)
			}
		}
	}
}

// suggestReceiverName suggère un nom de receiver basé sur le type.
//
// Params:
//   - typeName: nom du type
//
// Returns:
//   - string: nom suggéré (1-2 lettres)
func suggestReceiverName(typeName string) string {
	// Vérifier si le type est vide
	if len(typeName) == 0 {
		// Retourner valeur par défaut
		return "v"
	}

	// Retourner la première lettre en minuscule
	first := typeName[0]
	// Convertir en minuscule si nécessaire
	if first >= 'A' && first <= 'Z' {
		first = first + 32
	}

	// Retourner le nom suggéré
	return string(first)
}
