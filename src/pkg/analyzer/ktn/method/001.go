package ktn_method

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie que les méthodes qui modifient le receiver utilisent un receiver pointeur.
//
// KTN-METHOD-001: Les méthodes qui modifient le receiver doivent utiliser un receiver pointeur.
// Un receiver non-pointeur reçoit une copie, donc les modifications ne sont pas visibles.
//
// Incorrect:
//   func (s MyStruct) SetValue(v int) {
//       s.value = v  // Modifie la copie!
//   }
//
// Correct:
//   func (s *MyStruct) SetValue(v int) {
//       s.value = v  // Modifie l'original
//   }
var Rule001 = &analysis.Analyzer{
	Name: "KTN_METHOD_001",
	Doc:  "Vérifie que les méthodes qui modifient le receiver utilisent un receiver pointeur",
	Run:  runRule001,
}

// runRule001 exécute la vérification KTN-METHOD-001.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Recv == nil {
				continue
			}

			checkMethodPointerReceiver(pass, funcDecl)
		}
	}

	return nil, nil
}

// checkMethodPointerReceiver vérifie les méthodes avec receiver non-pointeur.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction (méthode)
func checkMethodPointerReceiver(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return
	}

	recvField := funcDecl.Recv.List[0]
	recvType := recvField.Type

	// Vérifier si c'est un pointeur
	if _, isPointer := recvType.(*ast.StarExpr); isPointer {
		return
	}

	// Vérifier si la méthode modifie le receiver
	if methodModifiesReceiver(funcDecl) {
		reportNonPointerReceiver(pass, funcDecl, getReceiverName(recvField))
	}
}

// methodModifiesReceiver vérifie si une méthode modifie son receiver.
//
// Params:
//   - funcDecl: la déclaration de fonction
//
// Returns:
//   - bool: true si le receiver est modifié
func methodModifiesReceiver(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Body == nil {
		return false
	}

	receiverName := ""
	if len(funcDecl.Recv.List) > 0 && len(funcDecl.Recv.List[0].Names) > 0 {
		receiverName = funcDecl.Recv.List[0].Names[0].Name
	}

	if receiverName == "" {
		return false
	}

	modifies := false

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			// Vérifier si Lhs contient le receiver ou un champ du receiver
			for _, lhs := range stmt.Lhs {
				if refersToReceiver(lhs, receiverName) {
					modifies = true
					return false
				}
			}
		case *ast.IncDecStmt:
			// Vérifier si c'est une incrémentation/décrémentation du receiver
			if refersToReceiver(stmt.X, receiverName) {
				modifies = true
				return false
			}
		}

		return true
	})

	return modifies
}

// refersToReceiver vérifie si une expression fait référence au receiver.
//
// Params:
//   - expr: l'expression
//   - receiverName: nom du receiver
//
// Returns:
//   - bool: true si fait référence au receiver
func refersToReceiver(expr ast.Expr, receiverName string) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Assignation directe au receiver
		return e.Name == receiverName
	case *ast.SelectorExpr:
		// Assignation à un champ du receiver
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Name == receiverName
		}
	case *ast.IndexExpr:
		// Assignation à un élément d'un slice/map du receiver
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Name == receiverName
		}
	}
	return false
}

// getReceiverName extrait le nom du type receiver.
//
// Params:
//   - field: le champ receiver
//
// Returns:
//   - string: nom du type
func getReceiverName(field *ast.Field) string {
	switch t := field.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return "receiver"
}

// reportNonPointerReceiver rapporte une violation KTN-METHOD-001.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
//   - receiverType: type du receiver
func reportNonPointerReceiver(pass *analysis.Pass, funcDecl *ast.FuncDecl, receiverType string) {
	pass.Reportf(funcDecl.Pos(),
		"[KTN-METHOD-001] Méthode '%s' avec receiver non-pointeur mais modifie le receiver.\n"+
			"Un receiver non-pointeur reçoit une copie.\n"+
			"Les modifications ne sont pas visibles à l'appelant.\n"+
			"Utilisez un receiver pointeur pour les méthodes qui modifient l'état.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - modifications perdues\n"+
			"  func (s MyStruct) SetValue(v int) {\n"+
			"      s.value = v  // Modifie la copie!\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - receiver pointeur\n"+
			"  func (s *MyStruct) SetValue(v int) {\n"+
			"      s.value = v  // Modifie l'original\n"+
			"  }",
		funcDecl.Name.Name)
}
