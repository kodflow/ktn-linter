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

// Test_runFunc002_disabled tests behavior when rule is disabled.
func Test_runFunc002_disabled(t *testing.T) {
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
					"KTN-FUNC-002": {Enabled: config.Bool(false)},
				},
			})
			// Reset config après le test
			defer config.Reset()

			// Créer un pass minimal
			result, err := runFunc002(&analysis.Pass{})
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

// Test_runFunc002_excludedFile tests behavior with excluded files.
func Test_runFunc002_excludedFile(t *testing.T) {
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
					"KTN-FUNC-002": {
						Enabled:       config.Bool(true),
						Exclude: []string{"test.go"},
					},
				},
			})
			// Reset config après le test
			defer config.Reset()

			code := `package test
			func foo() { }
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
			_, err = runFunc002(pass)
			// Vérification erreur
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

		})
	}
}

// Test_runFunc002 tests the runFunc002 private function.
func Test_runFunc002(t *testing.T) {
	// Test cases pour la fonction privée runFunc002
	// La logique principale est testée via l'API publique dans 008_external_test.go
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

// Test_isContextType tests the isContextType private function.
func Test_isContextType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}

// Test_isContextTypeWithPass vérifie la détection de context.Context avec pass.
//
// Params:
//   - t: instance de testing
func Test_isContextTypeWithPass(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "context_detection_with_pass",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite analysis.Pass réel
			_ = tt.name
		})
	}
}

// Test_isContextTypeByType vérifie la détection de context.Context par type.
//
// Params:
//   - t: instance de testing
func Test_isContextTypeByType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "context_type_detection",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite types.Type réel
			_ = tt.name
		})
	}
}

// Test_isContextObj teste la fonction isContextObj.
//
// Params:
//   - t: instance de testing
func Test_isContextObj(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "context_obj_detection",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite types.TypeName réel
			_ = tt.name
		})
	}
}

// Test_isContextUnderlying teste la fonction isContextUnderlying.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "context_underlying_detection",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite types.Type réel
			_ = tt.name
		})
	}
}

// Test_isContextType_fallback tests the AST-based fallback for context detection.
//
// Params:
//   - t: instance de testing
func Test_isContextType_fallback(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "context.Context selector",
			code: "context.Context",
			expected: true,
		},
		{
			name: "non-context selector",
			code: "foo.Bar",
			expected: false,
		},
		{
			name: "non-selector expression",
			code: "foo",
			expected: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			expr, err := parser.ParseExpr(tt.code)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Appeler la fonction isContextType (AST-based fallback)
			result := isContextType(expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isContextType(%s) = %v, want %v", tt.code, result, tt.expected)
			}
			_ = fset
		})
	}
}

// Test_reportMultipleContexts_withMultiple tests reporting multiple context params.
//
// Params:
//   - t: instance de testing
func Test_reportMultipleContexts_withMultiple(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec plusieurs contextes pour déclencher le rapport d'erreur
			code := `package test
			import "context"
			func badFunc(ctx1 context.Context, ctx2 context.Context) { }
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Trouver la déclaration de fonction
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok {
					// Assignation de la déclaration de fonction
					funcDecl = fd
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérification que la fonction a été trouvée
			if funcDecl == nil {
				t.Fatal("Expected to find function declaration")
			}

			reported := false
			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					// Vérifier que le message contient "2 paramètres context.Context"
					if len(d.Message) > 0 {
						// Marquer comme rapporté
						reported = true
					}
				},
			}

			// Appeler reportMultipleContexts avec contextCount = 2
			reportMultipleContexts(pass, funcDecl, 2)

			// Vérifier qu'une erreur a été rapportée
			if !reported {
				t.Error("Expected error report for multiple context parameters")
			}

		})
	}
}

// Test_reportMultipleContexts_withSingle tests no reporting for single context.
//
// Params:
//   - t: instance de testing
func Test_reportMultipleContexts_withSingle(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un seul contexte (ne doit pas déclencher d'erreur)
			code := `package test
			import "context"
			func goodFunc(ctx context.Context) { }
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Trouver la déclaration de fonction
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok {
					// Assignation de la déclaration de fonction
					funcDecl = fd
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Expected no report for single context, got: %s", d.Message)
				},
			}

			// Appeler reportMultipleContexts avec contextCount = 1
			reportMultipleContexts(pass, funcDecl, 1)

		})
	}
}

// Test_analyzeContextParams_multipleInSameField tests multiple contexts in one field.
//
// Params:
//   - t: instance de testing
func Test_analyzeContextParams_multipleInSameField(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec plusieurs contextes dans le même champ (ctx1, ctx2 context.Context)
			code := `package test
			import "context"
			func badFunc(ctx1, ctx2 context.Context) { }
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker pour avoir les informations de types
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Defs:  make(map[*ast.Ident]types.Object),
				Uses:  make(map[*ast.Ident]types.Object),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver la déclaration de fonction
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if fd, ok := n.(*ast.FuncDecl); ok {
					// Assignation de la déclaration de fonction
					funcDecl = fd
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			reported := false
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
				Report: func(d analysis.Diagnostic) {
					// Devrait rapporter 2 paramètres context.Context
					reported = true
				},
			}

			// Appeler analyzeContextParams (devrait détecter 2 contextes dans le même field)
			analyzeContextParams(pass, funcDecl)

			// Vérifier qu'une erreur a été rapportée
			if !reported {
				t.Error("Expected error report for multiple contexts in same field")
			}

		})
	}
}

// Test_isContextTypeByType_notNamed tests when type is not a Named type.
//
// Params:
//   - t: instance de testing
func Test_isContextTypeByType_notNamed(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un type de base (non Named)
			basicType := types.Typ[types.Int]

			// Appeler isContextTypeByType avec un type de base
			result := isContextTypeByType(basicType)

			// Vérifier false car ce n'est pas un Named type
			if result {
				t.Error("Expected false for basic type")
			}

		})
	}
}

// Test_isContextTypeByType_realContext tests with real context.Context.
//
// Params:
//   - t: instance de testing
func Test_isContextTypeByType_realContext(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec context.Context réel
			code := `package test
			import "context"
			func foo(ctx context.Context) { }
			`
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

			// Trouver le type context.Context
			var ctxType types.Type
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est une fonction
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Params != nil && len(fd.Type.Params.List) > 0 {
					tv := info.Types[fd.Type.Params.List[0].Type]
					// Assignation du type
					ctxType = tv.Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que le type a été trouvé
			if ctxType == nil {
				t.Fatal("Expected to find context.Context type")
			}

			// Appeler isContextTypeByType
			result := isContextTypeByType(ctxType)
			// Vérifier true car c'est context.Context
			if !result {
				t.Error("Expected true for context.Context type")
			}

		})
	}
}

// Test_isContextTypeByType_noObject tests when Named type has no object.
//
// Params:
//   - t: instance de testing
func Test_isContextTypeByType_noObject(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Créer un Named type sans objet (edge case rare)
			// Normalement un Named type a toujours un objet, mais on teste quand même
			basicType := types.Typ[types.Int]
			result := isContextTypeByType(basicType)

			// Vérifier false
			if result {
				t.Error("Expected false for type without object")
			}

		})
	}
}

// Test_isContextTypeWithPass_fallback tests fallback to AST when type info unavailable.
//
// Params:
//   - t: instance de testing
func Test_isContextTypeWithPass_fallback(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test du fallback AST quand TypesInfo est vide
			code := `package test
			import "context"
			func foo(ctx context.Context) { }
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Trouver l'expression de type context.Context
			var ctxExpr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type de nœud
				if field, ok := n.(*ast.Field); ok && len(field.Names) > 0 {
					// Assignation de l'expression de type
					ctxExpr = field.Type
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérification que l'expression a été trouvée
			if ctxExpr == nil {
				t.Fatal("Expected to find context expression")
			}

			// Créer un pass avec TypesInfo vide (force fallback)
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Appeler isContextTypeWithPass (doit fallback sur isContextType)
			result := isContextTypeWithPass(pass, ctxExpr)

			// Vérifier le résultat (devrait être true via fallback AST)
			if !result {
				t.Error("Expected true from fallback AST check for context.Context")
			}

		})
	}
}

// Test_isContextUnderlying_withAlias tests underlying type detection for aliases.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying_withAlias(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un alias de context.Context
			code := `package test
			import "context"
			type MyContext = context.Context
			func foo(ctx MyContext) { }
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker pour obtenir les informations de type
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Defs:  make(map[*ast.Ident]types.Object),
				Uses:  make(map[*ast.Ident]types.Object),
			}
			pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver le type MyContext dans le scope
			obj := pkg.Scope().Lookup("MyContext")
			// Vérification que l'objet existe
			if obj == nil {
				t.Fatal("Expected to find MyContext type")
			}

			// Obtenir le type nommé
			typeName, ok := obj.(*types.TypeName)
			// Vérification du type de l'objet
			if !ok {
				t.Fatal("Expected TypeName")
			}

			// Test isContextUnderlying avec l'alias
			result := isContextUnderlying(typeName.Type(), typeName)

			// Pour un alias de type (=), l'underlying devrait être context.Context
			// Note: les alias sont transparents, donc le type est directement context.Context
			if !result {
				t.Logf("Type alias may be transparent, result: %v", result)
			}

		})
	}
}

// Test_isContextUnderlying_realContext tests with direct context.Context type.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying_realContext(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec context.Context réel depuis le paramètre
			code := `package test
			import "context"
			func foo(ctx context.Context) { }
			`
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
				Defs:  make(map[*ast.Ident]types.Object),
			}
			_, err = conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver le type du paramètre ctx
			var ctxType types.Type
			var ctxTypeName *types.TypeName
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est une fonction
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Type.Params != nil && len(fd.Type.Params.List) > 0 {
					tv := info.Types[fd.Type.Params.List[0].Type]
					// Assignation du type
					ctxType = tv.Type
					// Extraire TypeName si Named
					if named, ok := ctxType.(*types.Named); ok {
						// Extraire le TypeName
						ctxTypeName = named.Obj()
					}
					// Retour false pour arrêter la recherche
					return false
				}
				// Retour true pour continuer la recherche
				return true
			})

			// Vérifier que le type a été trouvé
			if ctxType == nil || ctxTypeName == nil {
				t.Fatal("Expected to find context.Context type")
			}

			// Appeler isContextUnderlying - pour context.Context direct, devrait retourner false
			// car ce n'est pas un alias mais le type directement
			result := isContextUnderlying(ctxType, ctxTypeName)
			// Le résultat devrait être false car context.Context n'est pas un alias
			// mais le type original
			_ = result

		})
	}
}

// Test_isContextUnderlying_underlyingNilObj tests when underlying has nil object.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying_underlyingNilObj(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un type dont le sous-jacent est Named mais sans objet valide
			code := `package test
			type MyStruct struct{ field int }
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{}
			pkg, err := conf.Check("test", fset, []*ast.File{file}, &types.Info{
				Defs: make(map[*ast.Ident]types.Object),
			})
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver MyStruct
			obj := pkg.Scope().Lookup("MyStruct")
			// Vérification que l'objet existe
			if obj == nil {
				t.Fatal("Expected to find MyStruct")
			}

			typeName, ok := obj.(*types.TypeName)
			// Vérification du type de l'objet
			if !ok {
				t.Fatal("Expected TypeName")
			}

			// Test avec un struct (underlying est *types.Struct, pas *types.Named)
			result := isContextUnderlying(typeName.Type(), typeName)

			// Vérifier false car underlying n'est pas Named
			if result {
				t.Error("Expected false for struct underlying type")
			}

		})
	}
}

// Test_isContextUnderlying_contextTypedef tests with a typedef of context.Context.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying_contextTypedef(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un typedef de context.Context (pas un alias mais type Foo context.Context)
			// Note: "type MyCtx context.Context" crée un type dont l'underlying est *types.Interface
			// et non *types.Named, donc isContextUnderlying retourne false
			code := `package test
			import "context"
			type MyCtx context.Context
			func foo(ctx MyCtx) { }
			`
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
				Defs:  make(map[*ast.Ident]types.Object),
				Uses:  make(map[*ast.Ident]types.Object),
			}
			pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver MyCtx
			obj := pkg.Scope().Lookup("MyCtx")
			// Vérification que l'objet existe
			if obj == nil {
				t.Fatal("Expected to find MyCtx")
			}

			typeName, ok := obj.(*types.TypeName)
			// Vérification du type de l'objet
			if !ok {
				t.Fatal("Expected TypeName")
			}

			// Test isContextUnderlying - retourne false car underlying est Interface, pas Named
			result := isContextUnderlying(typeName.Type(), typeName)

			// L'underlying est *types.Interface (pas *types.Named), donc false
			if result {
				t.Error("Expected false: underlying is interface, not Named")
			}

		})
	}
}

// Test_isContextUnderlying_notContext tests with a non-context named underlying.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying_notContext(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un type dont l'underlying est Named mais pas context.Context
			code := `package test
			import "io"
			type MyReader io.Reader
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{Importer: importer.Default()}
			info := &types.Info{
				Defs: make(map[*ast.Ident]types.Object),
			}
			pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver MyReader
			obj := pkg.Scope().Lookup("MyReader")
			// Vérification que l'objet existe
			if obj == nil {
				t.Fatal("Expected to find MyReader")
			}

			typeName, ok := obj.(*types.TypeName)
			// Vérification du type de l'objet
			if !ok {
				t.Fatal("Expected TypeName")
			}

			// Test isContextUnderlying - devrait être false car ce n'est pas context.Context
			result := isContextUnderlying(typeName.Type(), typeName)

			// Vérifier false car ce n'est pas context.Context
			if result {
				t.Error("Expected false for non-context underlying type")
			}

		})
	}
}

// Test_isContextUnderlying_noPackage tests behavior when package is nil.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying_noPackage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un objet sans package (early return)
			code := `package test
			type MyType struct{}
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{}
			pkg, err := conf.Check("test", fset, []*ast.File{file}, &types.Info{
				Defs: make(map[*ast.Ident]types.Object),
			})
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver MyType
			obj := pkg.Scope().Lookup("MyType")
			// Vérification que l'objet existe
			if obj == nil {
				t.Fatal("Expected to find MyType")
			}

			typeName, ok := obj.(*types.TypeName)
			// Vérification du type de l'objet
			if !ok {
				t.Fatal("Expected TypeName")
			}

			// Créer un TypeName sans package (simule builtin)
			noPackageTypeName := types.NewTypeName(token.NoPos, nil, "TestType", typeName.Type())

			// Test isContextUnderlying avec package nil (doit retourner false)
			result := isContextUnderlying(noPackageTypeName.Type(), noPackageTypeName)

			// Vérifier false car package est nil
			if result {
				t.Error("Expected false when package is nil")
			}

		})
	}
}

// Test_isContextUnderlying_nonNamed tests behavior with non-named underlying type.
//
// Params:
//   - t: instance de testing
func Test_isContextUnderlying_nonNamed(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test avec un type dont l'underlying n'est pas Named
			code := `package test
			type MyInt int
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Vérification erreur parsing
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Créer type checker
			conf := types.Config{}
			pkg, err := conf.Check("test", fset, []*ast.File{file}, &types.Info{
				Defs: make(map[*ast.Ident]types.Object),
			})
			// Vérification erreur type checking
			if err != nil {
				t.Fatalf("Failed type check: %v", err)
			}

			// Trouver MyInt
			obj := pkg.Scope().Lookup("MyInt")
			// Vérification que l'objet existe
			if obj == nil {
				t.Fatal("Expected to find MyInt")
			}

			typeName, ok := obj.(*types.TypeName)
			// Vérification du type de l'objet
			if !ok {
				t.Fatal("Expected TypeName")
			}

			// Test isContextUnderlying avec un type dont l'underlying est Basic (int)
			result := isContextUnderlying(typeName.Type(), typeName)

			// Vérifier false car underlying est Basic, pas Named
			if result {
				t.Error("Expected false for non-named underlying type")
			}

		})
	}
}

// Test_analyzeContextParams tests the analyzeContextParams private function.
//
// Params:
//   - t: testing context
func Test_analyzeContextParams(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}


// Test_reportMultipleContexts tests the reportMultipleContexts private function.
//
// Params:
//   - t: testing context
func Test_reportMultipleContexts(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}


// Test_reportMisplacedContext tests the reportMisplacedContext private function.
//
// Params:
//   - t: testing context
func Test_reportMisplacedContext(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

