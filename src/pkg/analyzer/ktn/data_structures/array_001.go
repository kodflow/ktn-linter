package ktn_data_structures

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleArray001 = &analysis.Analyzer{
	Name: "KTN_ARRAY_002",
	Doc:  "Détecte les arrays avec taille incohérente",
	Run:  runRuleArray001,
}

func runRuleArray001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			// Vérifier si c'est un array (pas un slice)
			arrayType, ok := lit.Type.(*ast.ArrayType)
			if !ok || arrayType.Len == nil {
				return true
			}

			// Extraire la taille déclarée
			declaredLen := getArraySize(arrayType)
			if declaredLen == -1 {
				return true
			}

			// Compter les éléments
			actualLen := len(lit.Elts)

			// Vérifier l'incohérence
			if actualLen > declaredLen {
				pass.Reportf(lit.Pos(),
					"[KTN-ARRAY-001] Taille d'array incohérente: déclaré %d, mais %d éléments fournis.\n"+
						"Un array ne peut pas contenir plus d'éléments que sa taille déclarée.\n"+
						"Soit augmentez la taille, soit utilisez un slice.\n"+
						"Exemple:\n"+
						"  // ❌ MAUVAIS - trop d'éléments\n"+
						"  arr := [2]int{1, 2, 3}  // ERREUR\n"+
						"\n"+
						"  // ✅ CORRECT - bonne taille\n"+
						"  arr := [3]int{1, 2, 3}\n"+
						"\n"+
						"  // ✅ CORRECT - utiliser un slice\n"+
						"  arr := []int{1, 2, 3}",
					declaredLen, actualLen)
			}
			return true
		})
	}
	return nil, nil
}

func getArraySize(arrayType *ast.ArrayType) int {
	if arrayType.Len == nil {
		return -1
	}
	basicLit, ok := arrayType.Len.(*ast.BasicLit)
	if !ok {
		return -1
	}
	var size int
	_, err := fmt.Sscanf(basicLit.Value, "%d", &size)
	if err != nil {
		return -1
	}
	return size
}
