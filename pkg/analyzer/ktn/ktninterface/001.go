// Analyzer 001 for the ktninterface package.
//
// KTN-INTERFACE-001: Détecte les interfaces privées non utilisées.
//
// Comportement par défaut:
//   - Interfaces exportées (majuscule): ignorées car API publique (usage externe possible)
//   - Interfaces privées (minuscule): reportées si non utilisées dans le package
//
// Une interface est considérée "utilisée" si elle apparaît comme type dans:
//   - Champ de struct
//   - Paramètre de fonction/méthode
//   - Retour de fonction/méthode
//   - Déclaration de variable (var x MyInterface)
//   - Compile-time check (var _ MyInterface = (*S)(nil))
//
// La détection supporte les types imbriqués: *T, []T, map[K]V, chan T, etc.
package ktninterface

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	// initialCap capacité initiale des maps
	initialCap int = 16
	// ruleCode code de la règle
	ruleCode string = "KTN-INTERFACE-001"
)

// Analyzer001 détecte les interfaces privées non utilisées.
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktninterface001",
	Doc:      "KTN-INTERFACE-001: interface privée non utilisée",
	Run:      runInterface001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runInterface001 exécute l'analyse.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: erreur éventuelle
func runInterface001(pass *analysis.Pass) (any, error) {
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCode) {
		// Règle désactivée
		return nil, nil
	}

	// Pass 1: Collecter toutes les interfaces du package
	interfaces := collectInterfaces(pass, cfg)

	// Pass 2: Trouver les usages d'interfaces dans l'AST
	usedInterfaces := findUsages(pass, interfaces)

	// Report: Interfaces privées non utilisées
	reportUnused(pass, interfaces, usedInterfaces)

	// Retour succès
	return nil, nil
}

// collectInterfaces collecte toutes les interfaces déclarées.
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration
//
// Returns:
//   - map[string]*ast.TypeSpec: map nom -> TypeSpec
func collectInterfaces(pass *analysis.Pass, cfg *config.Config) map[string]*ast.TypeSpec {
	interfaces := make(map[string]*ast.TypeSpec, initialCap)

	// Parcourir les fichiers
	for _, file := range pass.Files {
		// Vérifier exclusion fichier
		filename := pass.Fset.Position(file.Pos()).Filename

		// Fichier exclu
		if cfg.IsFileExcluded(ruleCode, filename) {
			// Passer au suivant
			continue
		}

		// Parcourir les déclarations du fichier
		collectInterfacesFromFile(file, interfaces)
	}

	// Retourner les interfaces collectées
	return interfaces
}

// collectInterfacesFromFile collecte les interfaces d'un fichier AST.
//
// Params:
//   - file: fichier AST à analyser
//   - interfaces: map où stocker les interfaces trouvées
func collectInterfacesFromFile(file *ast.File, interfaces map[string]*ast.TypeSpec) {
	// Parcourir l'AST
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier si c'est un TypeSpec
		typeSpec, ok := n.(*ast.TypeSpec)

		// Pas un TypeSpec
		if !ok {
			// Continuer l'inspection
			return true
		}

		// Vérifier si c'est une interface
		if _, isIface := typeSpec.Type.(*ast.InterfaceType); isIface {
			// Ajouter à la map
			interfaces[typeSpec.Name.Name] = typeSpec
		}

		// Continuer l'inspection
		return true
	})
}

// findUsages trouve les usages d'interfaces dans l'AST du package.
//
// Params:
//   - pass: contexte d'analyse
//   - interfaces: interfaces à chercher
//
// Returns:
//   - map[string]bool: interfaces utilisées
func findUsages(pass *analysis.Pass, interfaces map[string]*ast.TypeSpec) map[string]bool {
	used := make(map[string]bool, initialCap)

	// Parcourir les fichiers
	for _, file := range pass.Files {
		// Chercher les usages dans ce fichier
		findUsagesInFile(file, interfaces, used)
	}

	// Retourner les usages trouvés
	return used
}

// findUsagesInFile cherche les usages d'interfaces dans un fichier.
//
// Params:
//   - file: fichier AST
//   - interfaces: interfaces connues
//   - used: map des interfaces utilisées (modifiée in-place)
func findUsagesInFile(file *ast.File, interfaces map[string]*ast.TypeSpec, used map[string]bool) {
	// Parcourir l'AST
	ast.Inspect(file, func(n ast.Node) bool {
		// Extraire les types utilisés selon le nœud
		types := extractTypesFromNode(n)

		// Marquer comme utilisés si dans notre liste d'interfaces
		for _, typeName := range types {
			// Vérifier si c'est une interface connue
			if _, exists := interfaces[typeName]; exists {
				// Marquer comme utilisée
				used[typeName] = true
			}
		}

		// Continuer l'inspection
		return true
	})
}

// extractTypesFromNode extrait les noms de types d'un nœud AST.
//
// Params:
//   - n: nœud AST
//
// Returns:
//   - []string: noms de types trouvés
func extractTypesFromNode(n ast.Node) []string {
	// Analyser le type de nœud
	switch node := n.(type) {
	// Champ de struct / paramètre / retour
	case *ast.Field:
		// Extraire types du champ
		return extractTypeIdents(node.Type)

	// Déclaration de variable (var x Type)
	case *ast.ValueSpec:
		// Vérifier si type explicite
		if node.Type != nil {
			// Extraire types de la déclaration
			return extractTypeIdents(node.Type)
		}

	// Type assertion (x.(Type))
	case *ast.TypeAssertExpr:
		// Vérifier si type présent
		if node.Type != nil {
			// Extraire types de l'assertion
			return extractTypeIdents(node.Type)
		}

	// Type switch cases
	case *ast.CaseClause:
		// Extraire types des cases
		return extractCaseClauseTypes(node)

	// Interface embedding
	case *ast.InterfaceType:
		// Extraire types embarqués
		return extractEmbeddedTypes(node)
	}

	// Nœud sans types
	return []string{}
}

// extractCaseClauseTypes extrait les types d'un case clause.
//
// Params:
//   - node: case clause à analyser
//
// Returns:
//   - []string: types trouvés
func extractCaseClauseTypes(node *ast.CaseClause) []string {
	var types []string

	// Parcourir les expressions du case
	for _, expr := range node.List {
		// Extraire et ajouter les types
		types = append(types, extractTypeIdents(expr)...)
	}

	// Retourner les types trouvés
	return types
}

// extractEmbeddedTypes extrait les types embarqués d'une interface.
//
// Params:
//   - node: interface à analyser
//
// Returns:
//   - []string: types embarqués trouvés
func extractEmbeddedTypes(node *ast.InterfaceType) []string {
	// Vérifier si méthodes présentes
	if node.Methods == nil {
		// Pas de méthodes
		return []string{}
	}

	var types []string

	// Parcourir les méthodes/types embarqués
	for _, m := range node.Methods.List {
		// Vérifier si type présent
		if m.Type != nil {
			// Extraire et ajouter les types
			types = append(types, extractTypeIdents(m.Type)...)
		}
	}

	// Retourner les types embarqués
	return types
}

// extractTypeIdents extrait récursivement les identifiants de type.
// Supporte: *T, []T, [N]T, map[K]V, chan T, func(...), Foo[T], etc.
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - []string: identifiants trouvés
func extractTypeIdents(expr ast.Expr) []string {
	// Vérifier nil
	if expr == nil {
		// Expression nulle
		return []string{}
	}

	// Essayer extraction simple (ident, selector)
	if result := extractSimpleType(expr); len(result) > 0 {
		// Type simple extrait
		return result
	}

	// Essayer extraction récursive (pointer, array, chan, ellipsis, paren)
	if result := extractRecursiveType(expr); len(result) > 0 {
		// Type récursif extrait
		return result
	}

	// Essayer extraction composite (map, func, generics)
	if result := extractCompositeType(expr); len(result) > 0 {
		// Type composite extrait
		return result
	}

	// Type non reconnu
	return []string{}
}

// extractSimpleType extrait les types simples (ident, selector).
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - []string: types extraits ou nil si non applicable
func extractSimpleType(expr ast.Expr) []string {
	// Analyser le type
	switch t := expr.(type) {
	// Identifiant simple: MyInterface
	case *ast.Ident:
		// Retourner le nom
		return []string{t.Name}

	// Sélecteur: pkg.MyInterface
	case *ast.SelectorExpr:
		// Retourner le sélecteur
		return []string{t.Sel.Name}
	}

	// Pas un type simple
	return []string{}
}

// extractRecursiveType extrait les types récursifs (pointer, array, etc).
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - []string: types extraits ou nil si non applicable
func extractRecursiveType(expr ast.Expr) []string {
	// Analyser le type
	switch t := expr.(type) {
	// Pointeur: *T
	case *ast.StarExpr:
		// Extraire le type pointé
		return extractTypeIdents(t.X)

	// Slice/Array: []T ou [N]T
	case *ast.ArrayType:
		// Extraire le type élément
		return extractTypeIdents(t.Elt)

	// Channel: chan T
	case *ast.ChanType:
		// Extraire le type du channel
		return extractTypeIdents(t.Value)

	// Ellipsis: ...T
	case *ast.Ellipsis:
		// Extraire le type variadique
		return extractTypeIdents(t.Elt)

	// Parenthèses: (T)
	case *ast.ParenExpr:
		// Extraire le contenu
		return extractTypeIdents(t.X)
	}

	// Pas un type récursif
	return []string{}
}

// extractCompositeType extrait les types composites (map, func, generics).
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - []string: types extraits ou nil si non applicable
func extractCompositeType(expr ast.Expr) []string {
	// Analyser le type
	switch t := expr.(type) {
	// Map: map[K]V
	case *ast.MapType:
		// Extraire clé et valeur
		return extractMapTypes(t)

	// Function type: func(T) R
	case *ast.FuncType:
		// Extraire params et retours
		return extractFuncTypes(t)

	// Generic: Foo[T]
	case *ast.IndexExpr:
		// Extraire type générique simple
		return extractIndexTypes(t)

	// Generic multiple: Foo[T, U]
	case *ast.IndexListExpr:
		// Extraire types génériques multiples
		return extractIndexListTypes(t)
	}

	// Pas un type composite
	return []string{}
}

// extractMapTypes extrait les types d'une expression map.
//
// Params:
//   - t: expression map
//
// Returns:
//   - []string: types clé et valeur
func extractMapTypes(t *ast.MapType) []string {
	var result []string

	// Extraire type clé
	result = append(result, extractTypeIdents(t.Key)...)

	// Extraire type valeur
	result = append(result, extractTypeIdents(t.Value)...)

	// Retourner résultat
	return result
}

// extractFuncTypes extrait les types d'une expression fonction.
//
// Params:
//   - t: expression fonction
//
// Returns:
//   - []string: types params et retours
func extractFuncTypes(t *ast.FuncType) []string {
	var result []string

	// Extraire paramètres
	if t.Params != nil {
		// Parcourir les paramètres
		for _, f := range t.Params.List {
			// Ajouter les types
			result = append(result, extractTypeIdents(f.Type)...)
		}
	}

	// Extraire retours
	if t.Results != nil {
		// Parcourir les retours
		for _, f := range t.Results.List {
			// Ajouter les types
			result = append(result, extractTypeIdents(f.Type)...)
		}
	}

	// Retourner résultat
	return result
}

// extractIndexTypes extrait les types d'une expression générique simple.
//
// Params:
//   - t: expression générique
//
// Returns:
//   - []string: types trouvés
func extractIndexTypes(t *ast.IndexExpr) []string {
	var result []string

	// Extraire type de base
	result = append(result, extractTypeIdents(t.X)...)

	// Extraire type paramètre
	result = append(result, extractTypeIdents(t.Index)...)

	// Retourner résultat
	return result
}

// extractIndexListTypes extrait les types d'une expression générique multiple.
//
// Params:
//   - t: expression générique multiple
//
// Returns:
//   - []string: types trouvés
func extractIndexListTypes(t *ast.IndexListExpr) []string {
	var result []string

	// Extraire type de base
	result = append(result, extractTypeIdents(t.X)...)

	// Parcourir les indices
	for _, idx := range t.Indices {
		// Ajouter chaque type
		result = append(result, extractTypeIdents(idx)...)
	}

	// Retourner résultat
	return result
}

// reportUnused reporte les interfaces privées non utilisées.
//
// Params:
//   - pass: contexte d'analyse
//   - interfaces: toutes les interfaces
//   - used: interfaces utilisées
func reportUnused(pass *analysis.Pass, interfaces map[string]*ast.TypeSpec, used map[string]bool) {
	// Parcourir les interfaces
	for name, typeSpec := range interfaces {
		// Skip interfaces exportées (API publique)
		if ast.IsExported(name) {
			// Interface publique
			continue
		}

		// Skip interfaces utilisées
		if used[name] {
			// Interface utilisée
			continue
		}

		// Report interface privée non utilisée
		msg, _ := messages.Get(ruleCode)
		pass.Reportf(
			typeSpec.Pos(),
			"%s: %s",
			ruleCode,
			msg.Format(config.Get().Verbose, name),
		)
	}
}
