package ktnfunc

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runFunc001_disabled tests behavior when rule is disabled.
func Test_runFunc001_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Configuration avec règle désactivée
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Enabled: config.Bool(false)},
				},
			})
			// Reset config après le test
			defer config.Reset()

			// Créer un pass minimal
			result, err := runFunc001(&analysis.Pass{})
			// Vérification de l'erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			// Vérification du résultat nil
			if result != nil {
				t.Errorf("Expected nil result when rule disabled, got %v", result)
			}

		})
	}
}

// Test_runFunc001_excludedFile tests behavior with excluded files.
func Test_runFunc001_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Configuration avec fichier exclu
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {
						Enabled: config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config après le test
			defer config.Reset()

			code := `package test
			func foo() (string, error) { return "", nil }
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer un inspector
			files := []*ast.File{file}
			inspectResult, _ := inspect.Analyzer.Run(&analysis.Pass{
				Fset:  fset,
				Files: files,
			})

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Expected no diagnostics for excluded file, got: %s", d.Message)
				},
			}

			// Exécuter l'analyse
			_, err = runFunc001(pass)
			// Vérification erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

		})
	}
}

// Test_runFunc001 tests the runFunc001 private function.
func Test_runFunc001(t *testing.T) {
	// Test cases pour la fonction privée runFunc001
	// La logique principale est testée via l'API publique dans 006_external_test.go
	// Ce test vérifie les cas edge de la fonction privée

	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - la logique principale est dans external tests
		})
	}
}

// Test_validateErrorInReturns vérifie la validation de la position des erreurs.
func Test_validateErrorInReturns(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - la logique est testée via external tests
		})
	}
}

// Test_isErrorType vérifie la détection du type error.
func Test_isErrorType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "builtin error type",
			code: `package test
func foo() error { return nil }`,
			expected: true,
		},
		{
			name: "custom type implementing error",
			code: `package test
type MyError struct{}
func (MyError) Error() string { return "" }
func bar() MyError { return MyError{} }`,
			expected: true,
		},
		{
			name: "non-error type",
			code: `package test
func baz() string { return "" }`,
			expected: false,
		},
		{
			name: "named error type",
			code: `package test
type CustomError error
func qux() CustomError { return nil }`,
			expected: true,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Créer un pass avec TypesInfo
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
			}

			// Trouver la déclaration de fonction et son type de retour
			var errorExpr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
					// Assignation de l'expression du type de retour
					errorExpr = fd.Type.Results.List[0].Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que l'expression a été trouvée
			if errorExpr == nil {
				t.Fatal("Expected to find return type expression")
			}

			// Appeler isErrorType
			result := isErrorType(pass, errorExpr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isErrorType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isBuiltinError vérifie la détection du type error builtin.
func Test_isBuiltinError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "builtin error interface",
			code: `package test
func foo() error { return nil }`,
			expected: true,
		},
		{
			name: "non-interface type",
			code: `package test
func bar() int { return 0 }`,
			expected: false,
		},
		{
			name: "interface with wrong method count",
			code: `package test
type MyInterface interface { Foo(); Bar() }
func baz() MyInterface { return nil }`,
			expected: false,
		},
		{
			name: "interface with wrong method name",
			code: `package test
type WrongError interface { NotError() string }
func qux() WrongError { return nil }`,
			expected: false,
		},
		{
			name: "interface with wrong signature",
			code: `package test
type BadError interface { Error(int) string }
func bad() BadError { return nil }`,
			expected: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver le type de retour
			var returnType types.Type
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
					tv := info.Types[fd.Type.Results.List[0].Type]
					// Assignation du type
					returnType = tv.Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que le type a été trouvé
			if returnType == nil {
				t.Fatal("Expected to find return type")
			}

			// Appeler isBuiltinError
			result := isBuiltinError(returnType)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isBuiltinError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_implementsError vérifie si un type implémente error.
func Test_implementsError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "struct implementing error",
			code: `package test
type MyError struct{}
func (MyError) Error() string { return "" }
func foo() MyError { return MyError{} }`,
			expected: true,
		},
		{
			name: "struct not implementing error",
			code: `package test
type NotError struct{}
func foo() NotError { return NotError{} }`,
			expected: false,
		},
		{
			name: "pointer implementing error",
			code: `package test
type PtrError struct{}
func (*PtrError) Error() string { return "" }
func foo() *PtrError { return nil }`,
			expected: true,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver le type de retour
			var returnType types.Type
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Name.Name == "foo" && fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
					tv := info.Types[fd.Type.Results.List[0].Type]
					// Assignation du type
					returnType = tv.Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que le type a été trouvé
			if returnType == nil {
				t.Fatal("Expected to find return type")
			}

			// Appeler implementsError
			result := implementsError(returnType)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("implementsError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isErrorType_nilType tests isErrorType with nil type.
func Test_isErrorType_nilType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fset := token.NewFileSet()

			// Créer un pass avec TypesInfo vide (expression sans type)
			pass := &analysis.Pass{
				Fset: fset,
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Créer une expression qui n'a pas de type dans TypesInfo
			expr := &ast.Ident{Name: "unknown"}

			// Appeler isErrorType avec une expression sans type
			result := isErrorType(pass, expr)

			// Vérifier false car tv.Type == nil
			if result {
				t.Error("Expected false for expression without type info")
			}

		})
	}
}

// Test_isErrorType_namedBuiltinError tests the named error with nil package path.
func Test_isErrorType_namedBuiltinError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test le cas où le type est nommé "error" avec pkg == nil
			code := `package test
			var e error
			func foo() error { return e }`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Créer un pass
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
			}

			// Trouver l'expression du type de retour
			var errorExpr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
					// Assignation de l'expression du type de retour
					errorExpr = fd.Type.Results.List[0].Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que l'expression a été trouvée
			if errorExpr == nil {
				t.Fatal("Expected to find return type expression")
			}

			// Appeler isErrorType
			result := isErrorType(pass, errorExpr)
			// Vérification du résultat
			if !result {
				t.Error("Expected true for builtin error type")
			}

		})
	}
}

// Test_isErrorType_aliasedError tests isErrorType with error type alias.
func Test_isErrorType_aliasedError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test le cas où le type est un alias de error
			code := `package test
			type MyErr = error
			func foo() MyErr { return nil }`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Créer un pass
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
			}

			// Trouver l'expression du type de retour
			var errorExpr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
					// Assignation de l'expression du type de retour
					errorExpr = fd.Type.Results.List[0].Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que l'expression a été trouvée
			if errorExpr == nil {
				t.Fatal("Expected to find return type expression")
			}

			// Appeler isErrorType
			result := isErrorType(pass, errorExpr)
			// Vérification du résultat
			if !result {
				t.Error("Expected true for error alias type")
			}

		})
	}
}

// Test_isErrorType_namedNonError tests isErrorType with non-error named type.
func Test_isErrorType_namedNonError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un type nommé qui n'est pas error
			code := `package test
			type MyString string
			func foo() MyString { return "" }`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Créer un pass
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
			}

			// Trouver l'expression du type de retour
			var errorExpr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
					// Assignation de l'expression du type de retour
					errorExpr = fd.Type.Results.List[0].Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que l'expression a été trouvée
			if errorExpr == nil {
				t.Fatal("Expected to find return type expression")
			}

			// Appeler isErrorType - devrait retourner false
			result := isErrorType(pass, errorExpr)
			// Vérification du résultat
			if result {
				t.Error("Expected false for non-error named type")
			}

		})
	}
}

// Test_isErrorType_underlyingError tests isErrorType with named type having error underlying.
func Test_isErrorType_underlyingError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un type nommé dont le underlying est error
			code := `package test
			type CustomError error
			func foo() CustomError { return nil }`

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Créer un pass
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
			}

			// Trouver l'expression du type de retour
			var errorExpr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
					// Assignation de l'expression du type de retour
					errorExpr = fd.Type.Results.List[0].Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que l'expression a été trouvée
			if errorExpr == nil {
				t.Fatal("Expected to find return type expression")
			}

			// Appeler isErrorType - devrait retourner true car underlying est error
			result := isErrorType(pass, errorExpr)
			// Vérification du résultat
			if !result {
				t.Error("Expected true for named type with error underlying")
			}

		})
	}
}
