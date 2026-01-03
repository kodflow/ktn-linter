// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import (
	"go/ast"
	"unicode"
	"unicode/utf8"

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
	// minReceiversForConsistency nombre minimum de receivers pour vérifier la cohérence
	minReceiversForConsistency int = 2
)

var (
	// badReceiverNames noms de receiver génériques à éviter.
	badReceiverNames map[string]bool = map[string]bool{
		"me":   true,
		"this": true,
		"self": true,
	}

	// Analyzer009 checks receiver name consistency across methods.
	Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktnstruct009",
		Doc:      "KTN-STRUCT-009: Les noms de receiver doivent être cohérents entre méthodes",
		Run:      runStruct009,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
)

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
	cfg := config.Get()
	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct009) {
		// Règle désactivée
		return nil, nil
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	receiversByType := collectReceivers(pass, insp, cfg)
	checkReceiverConsistency(pass, receiversByType)
	checkGenericReceiverNames(pass, receiversByType)
	// Fin de l'analyse
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
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	insp.Preorder(nodeFilter, func(n ast.Node) {
		processMethodDecl009(pass, cfg, n.(*ast.FuncDecl), result)
	})
	// Retour de la map des receivers
	return result
}

// processMethodDecl009 traite une déclaration de méthode pour collecte des receivers.
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration
//   - funcDecl: déclaration de fonction
//   - result: map à remplir
func processMethodDecl009(pass *analysis.Pass, cfg *config.Config, funcDecl *ast.FuncDecl, result map[string][]receiverInfo) {
	filename := pass.Fset.Position(funcDecl.Pos()).Filename
	// Vérifier si le fichier est exclu
	if cfg.IsFileExcluded(ruleCodeStruct009, filename) {
		// Fichier exclu
		return
	}
	// Vérifier si c'est une méthode
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		// Pas une méthode
		return
	}
	recv := funcDecl.Recv.List[0]
	typeName := extractReceiverTypeName(recv.Type)
	// Vérifier si le type est valide
	if typeName == "" {
		// Type invalide
		return
	}
	receiverName := extractReceiverName(recv)
	// Vérifier si le nom est valide
	if receiverName == "" {
		// Nom invalide
		return
	}
	result[typeName] = append(result[typeName], receiverInfo{
		name:     receiverName,
		funcDecl: funcDecl,
	})
}

// extractReceiverTypeName extrait le nom du type du receiver.
// Gère les pointeurs (*Type) et les types simples (Type).
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
	// Extraire l'identifiant
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour du nom du type
		return ident.Name
	}
	// Type invalide
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
	// Vérifier si le champ a un nom
	if len(field.Names) == 0 {
		// Pas de nom
		return ""
	}
	// Retour du premier nom
	return field.Names[0].Name
}

// checkReceiverConsistency vérifie la cohérence des noms de receiver.
// Reporte si différents noms sont utilisés pour le même type.
//
// Params:
//   - pass: contexte d'analyse
//   - receiversByType: map des receivers par type
func checkReceiverConsistency(pass *analysis.Pass, receiversByType map[string][]receiverInfo) {
	// Parcourir les types
	for typeName, receivers := range receiversByType {
		// Ignorer si moins de minReceiversForConsistency receivers
		if len(receivers) < minReceiversForConsistency {
			continue
		}
		firstName := receivers[0].name
		// Comparer avec les autres receivers
		for i := 1; i < len(receivers); i++ {
			recv := receivers[i]
			// Vérifier si le nom est différent
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

// checkGenericReceiverNames vérifie les noms de receiver génériques (this, self, me).
//
// Params:
//   - pass: contexte d'analyse
//   - receiversByType: map des receivers par type
func checkGenericReceiverNames(pass *analysis.Pass, receiversByType map[string][]receiverInfo) {
	// Parcourir les types
	for typeName, receivers := range receiversByType {
		// Parcourir les receivers
		for _, recv := range receivers {
			// Vérifier si le nom est générique
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

// suggestReceiverName suggère un nom de receiver basé sur le type (UTF-8 safe).
//
// Params:
//   - typeName: nom du type
//
// Returns:
//   - string: nom suggéré (1-2 lettres)
func suggestReceiverName(typeName string) string {
	// Vérifier si le nom du type est vide
	if typeName == "" {
		// Nom par défaut
		return "v"
	}
	r, _ := utf8.DecodeRuneInString(typeName)
	// Première lettre en minuscule
	return string(unicode.ToLower(r))
}
