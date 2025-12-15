// Shared utilities for methods handling.
package shared

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	// defaultMethodsMapCap est la capacité initiale pour la map des méthodes
	defaultMethodsMapCap int = 16
)

// CollectMethodsByStruct collecte les méthodes groupées par struct.
//
// Params:
//   - file: fichier AST
//   - pass: contexte d'analyse
//
// Returns:
//   - map[string][]string: méthodes par nom de struct
func CollectMethodsByStruct(file *ast.File, _pass *analysis.Pass) map[string][]string {
	methods := make(map[string][]string, defaultMethodsMapCap)

	// Parcourir le fichier
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier FuncDecl
		funcDecl, ok := n.(*ast.FuncDecl)
		// Si ce n'est pas une fonction
		if !ok {
			// Continue traversal
			return true
		}

		// Ignorer les fonctions (pas de receiver)
		if funcDecl.Recv == nil {
			// Continue traversal
			return true
		}

		// Ignorer les méthodes non exportées
		if !ast.IsExported(funcDecl.Name.Name) {
			// Continue traversal
			return true
		}

		// Extraire le nom du receiver
		receiverName := ExtractReceiverName(funcDecl.Recv)
		// Si receiver valide
		if receiverName != "" {
			methods[receiverName] = append(methods[receiverName], funcDecl.Name.Name)
		}

		// Continue traversal
		return true
	})

	// Retour des méthodes
	return methods
}

// ExtractReceiverName extrait le nom du type receiver.
//
// Params:
//   - recv: liste des receivers
//
// Returns:
//   - string: nom du type receiver
func ExtractReceiverName(recv *ast.FieldList) string {
	// Vérifier si la liste est vide
	if recv == nil || len(recv.List) == 0 {
		// Retour vide
		return ""
	}

	// Prendre le premier receiver
	recvType := recv.List[0].Type

	// Gérer pointeur ou valeur
	switch t := recvType.(type) {
	// Traitement du pointeur
	case *ast.StarExpr:
		// Receiver de type *T
		if ident, ok := t.X.(*ast.Ident); ok {
			// Retour du nom
			return ident.Name
		}
	// Traitement de l'identifiant
	case *ast.Ident:
		// Receiver de type T
		return t.Name
	}

	// Type non géré
	return ""
}
