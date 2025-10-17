package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// DeclarationAnalyzer vérifie les déclarations (var, types prédéclarés, méthodes).
	DeclarationAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktndeclaration",
		Doc:  "Vérifie les déclarations de variables, types prédéclarés et méthodes",
		Run:  runDeclarationAnalyzer,
	}
)

// Liste des identifiants prédéclarés en Go
var predeclaredIdentifiers = map[string]bool{
	// Types
	"bool":       true,
	"byte":       true,
	"complex64":  true,
	"complex128": true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	// Constants
	"true":  true,
	"false": true,
	"iota":  true,
	// Zero value
	"nil": true,
	// Functions
	"append":  true,
	"cap":     true,
	"close":   true,
	"complex": true,
	"copy":    true,
	"delete":  true,
	"imag":    true,
	"len":     true,
	"make":    true,
	"new":     true,
	"panic":   true,
	"print":   true,
	"println": true,
	"real":    true,
	"recover": true,
}

// runDeclarationAnalyzer exécute l'analyseur de déclarations.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runDeclarationAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.GenDecl:
				if node.Tok == token.VAR {
					checkVarDeclaration(pass, node)
				}
				if node.Tok == token.TYPE || node.Tok == token.VAR || node.Tok == token.CONST {
					checkPredeclaredShadowing(pass, node)
				}
			case *ast.FuncDecl:
				if node.Recv != nil {
					checkMethodPointerReceiver(pass, node)
				}
			}
			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkVarDeclaration vérifie les déclarations var inutiles.
//
// Params:
//   - pass: la passe d'analyse
//   - decl: la déclaration générale
func checkVarDeclaration(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// Vérifier si c'est var x = value (avec initialisation)
		if len(valueSpec.Values) == 0 {
			// Pas d'initialisation, c'est ok
			continue
		}

		if valueSpec.Type != nil {
			// Type explicite, c'est ok
			continue
		}

		// C'est var x = value sans type explicite
		// On devrait utiliser := à la place
		for _, name := range valueSpec.Names {
			reportRedundantVarDecl(pass, valueSpec, name.Name)
		}
	}
}

// checkPredeclaredShadowing vérifie le shadowing d'identifiants prédéclarés.
//
// Params:
//   - pass: la passe d'analyse
//   - decl: la déclaration générale
func checkPredeclaredShadowing(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if predeclaredIdentifiers[s.Name.Name] {
				reportPredeclaredShadowing(pass, s.Name, s.Name.Name)
			}
		case *ast.ValueSpec:
			for _, name := range s.Names {
				if predeclaredIdentifiers[name.Name] {
					reportPredeclaredShadowing(pass, name, name.Name)
				}
			}
		}
	}
}

// checkMethodPointerReceiver vérifie les méthodes avec receiver non-pointeur.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction (méthode)
func checkMethodPointerReceiver(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		// Pas de receiver
		// Retourne
		return
	}

	recvField := funcDecl.Recv.List[0]
	recvType := recvField.Type

	// Vérifier si c'est un pointeur
	if _, isPointer := recvType.(*ast.StarExpr); isPointer {
		// C'est un pointeur, c'est ok
		// Retourne
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
		// Pas de body
		return false
	}

	receiverName := ""
	if len(funcDecl.Recv.List) > 0 && len(funcDecl.Recv.List[0].Names) > 0 {
		receiverName = funcDecl.Recv.List[0].Names[0].Name
	}

	if receiverName == "" {
		// Pas de nom de receiver
		return false
	}

	modifies := false

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			// Retourne true pour continuer
			return true
		}

		// Vérifier si Lhs contient le receiver ou un champ du receiver
		for _, lhs := range assignStmt.Lhs {
			if refersToReceiver(lhs, receiverName) {
				modifies = true
				// Retourne false pour arrêter
				return false
			}
		}

		// Retourne true pour continuer
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

// reportRedundantVarDecl rapporte une violation KTN-VAR-002.
//
// Params:
//   - pass: la passe d'analyse
//   - spec: la spécification de valeur
//   - varName: nom de la variable
func reportRedundantVarDecl(pass *analysis.Pass, spec *ast.ValueSpec, varName string) {
	pass.Reportf(spec.Pos(),
		"[KTN-VAR-002] Déclaration 'var %s = ...' redondante.\n"+
			"Quand une variable est initialisée, utilisez := au lieu de var.\n"+
			"C'est plus concis et idiomatique en Go.\n"+
			"Exemple:\n"+
			"  // ❌ NON-IDIOMATIQUE\n"+
			"  var name = \"John\"\n"+
			"\n"+
			"  // ✅ IDIOMATIQUE GO\n"+
			"  name := \"John\"\n"+
			"\n"+
			"Note: Utilisez var uniquement quand:\n"+
			"  - Pas d'initialisation: var count int\n"+
			"  - Type explicite différent: var x int64 = 42",
		varName)
}

// reportPredeclaredShadowing rapporte une violation KTN-PREDECL-002.
//
// Params:
//   - pass: la passe d'analyse
//   - ident: l'identifiant
//   - name: nom de l'identifiant prédéclaré
func reportPredeclaredShadowing(pass *analysis.Pass, ident *ast.Ident, name string) {
	pass.Reportf(ident.Pos(),
		"[KTN-PREDECL-002] Shadowing de l'identifiant prédéclaré '%s'.\n"+
			"Redéfinir un identifiant prédéclaré (type, fonction built-in, etc.) rend le code confus.\n"+
			"Cela cache l'identifiant original dans ce scope.\n"+
			"Choisissez un nom différent.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - shadow de 'string'\n"+
			"  var string = \"hello\"  // Cache le type string!\n"+
			"\n"+
			"  // ❌ MAUVAIS - shadow de 'len'\n"+
			"  len := 5  // Cache la fonction len()!\n"+
			"\n"+
			"  // ✅ CORRECT - noms différents\n"+
			"  var message = \"hello\"\n"+
			"  length := 5",
		name)
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
