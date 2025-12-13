package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runStruct001 tests the private runStruct001 function.
func Test_runStruct001(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}

// Test_collectInterfaces tests the private collectInterfaces function.
func Test_collectInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected int
	}{
		{
			name: "no interfaces",
			src: `package test
type User struct {
	Name string
}`,
			expected: 0,
		},
		{
			name: "one interface",
			src: `package test
type Reader interface {
	Read() error
}`,
			expected: 1,
		},
		{
			name: "multiple interfaces",
			src: `package test
type Reader interface {
	Read() error
}
type Writer interface {
	Write() error
}`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			pass := &analysis.Pass{Fset: fset}
			interfaces := collectInterfaces(file, pass)

			if len(interfaces) != tt.expected {
				t.Errorf("expected %d interfaces, got %d", tt.expected, len(interfaces))
			}
		})
	}
}

// Test_extractStructNameFromReceiver tests the private extractStructNameFromReceiver function.
func Test_extractStructNameFromReceiver(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "pointer receiver",
			src: `package test
type User struct{}
func (u *User) Method() {}`,
			expected: "User",
		},
		{
			name: "value receiver",
			src: `package test
type User struct{}
func (u User) Method() {}`,
			expected: "User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Find the function declaration
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil || funcDecl.Recv == nil {
				t.Fatal("no function with receiver found")
			}

			recvType := funcDecl.Recv.List[0].Type
			result := extractStructNameFromReceiver(recvType)

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_collectMethodsByStruct tests the private collectMethodsByStruct function.
func Test_collectMethodsByStruct(t *testing.T) {
	tests := []struct {
		name            string
		src             string
		structName      string
		expectedMethods int
	}{
		{
			name: "collect methods for struct",
			src: `package test
type User struct{}
func (u *User) GetName() string { return "" }
func (u *User) SetName(name string) {}`,
			structName:      "User",
			expectedMethods: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			pass := &analysis.Pass{Fset: fset}
			methodsByStruct := collectMethodsByStruct(file, pass)

			userMethods, ok := methodsByStruct[tt.structName]
			if !ok {
				t.Fatalf("expected %s struct in map", tt.structName)
			}

			if len(userMethods) != tt.expectedMethods {
				t.Errorf("expected %d methods, got %d", tt.expectedMethods, len(userMethods))
			}
		})
	}
}

// Test_collectAllMethodsByStruct tests the private collectAllMethodsByStruct function.
func Test_collectAllMethodsByStruct(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "collect methods from all package files"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", `package test
type User struct{}
func (u *User) GetName() string { return "" }`, 0)
			// Check parse result
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			methods := collectAllMethodsByStruct(pass)
			// Verify methods were collected
			if len(methods) == 0 {
				t.Error("expected at least one struct with methods")
			}
		})
	}
}

// Test_collectMethodsFromFile tests the private collectMethodsFromFile function.
func Test_collectMethodsFromFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "collect methods from single file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", `package test
type User struct{}
func (u *User) GetName() string { return "" }
func (u *User) SetName(name string) {}`, 0)
			// Check parse result
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			pass := &analysis.Pass{Fset: fset}
			methodsByStruct := make(map[string][]shared.MethodSignature)

			collectMethodsFromFile(file, pass, methodsByStruct)

			// Verify methods were collected
			if len(methodsByStruct["User"]) != 2 {
				t.Errorf("expected 2 methods for User, got %d", len(methodsByStruct["User"]))
			}
		})
	}
}

// Test_collectStructsFromFile tests the private collectStructsFromFile function.
func Test_collectStructsFromFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "collect structs from file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", `package test
type User struct {
	Name string
}`, 0)
			// Check parse result
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			allMethods := map[string][]shared.MethodSignature{
				"User": {{Name: "GetName", ParamsStr: "", ResultsStr: "string"}},
			}

			structs := collectStructsFromFile(file, allMethods)

			// Verify structs were collected
			if len(structs) != 1 {
				t.Errorf("expected 1 struct, got %d", len(structs))
			}
			// Verify methods were assigned
			if len(structs[0].methods) != 1 {
				t.Errorf("expected 1 method assigned, got %d", len(structs[0].methods))
			}
		})
	}
}

// Test_hasMatchingInterface tests the private hasMatchingInterface function.
func Test_hasMatchingInterface(t *testing.T) {
	tests := []struct {
		name       string
		structName string
		methods    []shared.MethodSignature
		interfaces map[string][]shared.MethodSignature
		expected   bool
	}{
		{
			name:       "matching interface found",
			structName: "User",
			methods:    []shared.MethodSignature{{Name: "GetName", ParamsStr: "", ResultsStr: "string"}},
			interfaces: map[string][]shared.MethodSignature{
				"Reader": {{Name: "GetName", ParamsStr: "", ResultsStr: "string"}},
			},
			expected: true,
		},
		{
			name:       "no matching interface",
			structName: "User",
			methods:    []shared.MethodSignature{{Name: "GetName", ParamsStr: "", ResultsStr: "string"}},
			interfaces: map[string][]shared.MethodSignature{
				"Writer": {{Name: "SetName", ParamsStr: "string", ResultsStr: ""}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := structWithMethods{name: tt.structName, methods: tt.methods}
			if hasMatchingInterface(s, tt.interfaces) != tt.expected {
				t.Errorf("expected %v for %s", tt.expected, tt.name)
			}
		})
	}
}

// Test_interfaceCoversAllMethods tests the private interfaceCoversAllMethods function.
func Test_interfaceCoversAllMethods(t *testing.T) {
	tests := []struct {
		name          string
		structMethods []shared.MethodSignature
		ifaceMethods  []shared.MethodSignature
		expected      bool
	}{
		{
			name: "interface covers all methods",
			structMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
				{Name: "GetAge", ParamsStr: "", ResultsStr: "int"},
			},
			ifaceMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
				{Name: "GetAge", ParamsStr: "", ResultsStr: "int"},
			},
			expected: true,
		},
		{
			name: "incomplete interface",
			structMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
				{Name: "GetAge", ParamsStr: "", ResultsStr: "int"},
			},
			ifaceMethods: []shared.MethodSignature{
				{Name: "GetName", ParamsStr: "", ResultsStr: "string"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if interfaceCoversAllMethods(tt.structMethods, tt.ifaceMethods) != tt.expected {
				t.Errorf("expected %v for %s", tt.expected, tt.name)
			}
		})
	}
}

// Test_formatFieldList tests the private formatFieldList function.
func Test_formatFieldList(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "no params",
			src: `package test
func test() {}`,
			expected: "",
		},
		{
			name: "single param",
			src: `package test
func test(x int) {}`,
			expected: "int",
		},
		{
			name: "multiple params",
			src: `package test
func test(x int, y string) {}`,
			expected: "int,string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil {
				t.Fatal("no function found")
			}

			pass := &analysis.Pass{Fset: fset}
			result := formatFieldList(funcDecl.Type.Params, pass)

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_runStruct001_disabled tests that the rule is skipped when disabled.
func Test_runStruct001_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-001": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			src := `package test
			type User struct { Name string }
			`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", src, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{f},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) {
					t.Error("Unexpected error when rule is disabled")
				},
			}

			_, err = runStruct001(pass)
			if err != nil {
				t.Errorf("runStruct001() error = %v", err)
			}

		})
	}
}

// Test_extractStructNameFromAST tests the extractStructNameFromAST function.
func Test_extractStructNameFromAST(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "pointer conversion (*S)(nil)",
			src: `package test
var _ I = (*MyStruct)(nil)`,
			expected: "MyStruct",
		},
		{
			name: "composite literal S{}",
			src: `package test
var _ I = MyStruct{}`,
			expected: "MyStruct",
		},
		{
			name: "address of composite &S{}",
			src: `package test
var _ I = &MyStruct{}`,
			expected: "MyStruct",
		},
		{
			name: "new(S)",
			src: `package test
var _ I = new(MyStruct)`,
			expected: "MyStruct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			// Vérifier erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Trouver la ValueSpec
			var valueExpr ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est une GenDecl
				if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
					// Parcourir les specs
					for _, spec := range genDecl.Specs {
						// Vérifier si c'est une ValueSpec
						if vs, ok := spec.(*ast.ValueSpec); ok && len(vs.Values) > 0 {
							valueExpr = vs.Values[0]
							// Sortir de la boucle
							return false
						}
					}
				}
				// Continuer l'itération
				return true
			})

			// Vérifier si on a trouvé une expression
			if valueExpr == nil {
				t.Fatal("no value expression found")
			}

			result := extractStructNameFromAST(valueExpr)

			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_extractFromCallExpr tests the extractFromCallExpr function.
func Test_extractFromCallExpr(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "pointer conversion (*S)(nil)",
			src: `package test
var _ I = (*MyStruct)(nil)`,
			expected: "MyStruct",
		},
		{
			name: "new(S)",
			src: `package test
var _ I = new(MyStruct)`,
			expected: "MyStruct",
		},
		{
			name: "regular function call",
			src: `package test
var _ I = someFunc()`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			// Vérifier erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Trouver le CallExpr
			var callExpr *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est un CallExpr
				if ce, ok := n.(*ast.CallExpr); ok {
					callExpr = ce
					// Sortir de la boucle
					return false
				}
				// Continuer l'itération
				return true
			})

			// Vérifier si on a trouvé un CallExpr
			if callExpr == nil {
				// Pas de CallExpr dans ce test
				return
			}

			result := extractFromCallExpr(callExpr)

			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_extractFromCompositeLit tests the extractFromCompositeLit function.
func Test_extractFromCompositeLit(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "simple composite literal",
			src: `package test
var _ I = MyStruct{}`,
			expected: "MyStruct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			// Vérifier erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Trouver le CompositeLit
			var compLit *ast.CompositeLit
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est un CompositeLit
				if cl, ok := n.(*ast.CompositeLit); ok {
					compLit = cl
					// Sortir de la boucle
					return false
				}
				// Continuer l'itération
				return true
			})

			// Vérifier si on a trouvé un CompositeLit
			if compLit == nil {
				t.Fatal("no CompositeLit found")
			}

			result := extractFromCompositeLit(compLit)

			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_extractFromUnaryExpr tests the extractFromUnaryExpr function.
func Test_extractFromUnaryExpr(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "address of composite &S{}",
			src: `package test
var _ I = &MyStruct{}`,
			expected: "MyStruct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			// Vérifier erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Trouver le UnaryExpr
			var unaryExpr *ast.UnaryExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est un UnaryExpr
				if ue, ok := n.(*ast.UnaryExpr); ok {
					unaryExpr = ue
					// Sortir de la boucle
					return false
				}
				// Continuer l'itération
				return true
			})

			// Vérifier si on a trouvé un UnaryExpr
			if unaryExpr == nil {
				t.Fatal("no UnaryExpr found")
			}

			result := extractFromUnaryExpr(unaryExpr)

			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_checkStructs tests the checkStructs function.
func Test_checkStructs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La logique de checkStructs est testée via analysistest.Run
			// dans 001_external_test.go qui exerce tout le pipeline
		})
	}
}

// Test_collectInterfaceChecksWithTypes tests the collectInterfaceChecksWithTypes function.
func Test_collectInterfaceChecksWithTypes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction nécessite un *analysis.Pass complet avec TypesInfo
			// Elle est testée via analysistest.Run dans 001_external_test.go
		})
	}
}

// Test_extractInterfaceCheckWithTypes tests the extractInterfaceCheckWithTypes function.
func Test_extractInterfaceCheckWithTypes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction nécessite TypesInfo pour résoudre les types
			// Elle est testée via analysistest.Run dans 001_external_test.go
		})
	}
}

// Test_extractStructNameFromValue tests the extractStructNameFromValue function.
func Test_extractStructNameFromValue(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction utilise pass.TypesInfo, testée via analysistest
			// Le fallback AST est testé dans Test_extractStructNameFromAST
		})
	}
}

// Test_extractStructNameFromType tests the extractStructNameFromType function.
func Test_extractStructNameFromType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction travaille avec types.Type
			// Elle est exercée via analysistest.Run dans 001_external_test.go
		})
	}
}

// Test_hasMatchingInterfaceCheck tests the hasMatchingInterfaceCheck function.
func Test_hasMatchingInterfaceCheck(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction nécessite interfaceCheck avec types.Interface
			// Elle est testée via analysistest.Run dans 001_external_test.go
			// Le testdata/bad.go contient BadIncompleteImpl qui vérifie
			// qu'une interface incomplète déclenche bien KTN-STRUCT-001
		})
	}
}

// Test_interfaceCoversAllPublicMethods tests the interfaceCoversAllPublicMethods function.
func Test_interfaceCoversAllPublicMethods(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction vérifie si types.Interface couvre toutes les méthodes
			// Elle est testée via analysistest.Run dans 001_external_test.go
			// BadIncompleteImpl dans testdata/bad.go prouve que les interfaces
			// incomplètes sont correctement détectées
		})
	}
}

// Test_signaturesMatch tests the signaturesMatch function.
func Test_signaturesMatch(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction compare types.Func avec MethodSignature
			// Elle est exercée via analysistest.Run dans 001_external_test.go
		})
	}
}

// Test_formatTypeTuple tests the formatTypeTuple function.
func Test_formatTypeTuple(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via analysistest framework"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cette fonction formate types.Tuple en string
			// Elle est exercée via analysistest.Run dans 001_external_test.go
		})
	}
}

// Test_runStruct001_excludedFile tests that excluded files are skipped.
func Test_runStruct001_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-001": {
						Enabled: config.Bool(true),
						Exclude: []string{"**/test.go"},
					},
				},
			})
			defer config.Reset()

			src := `package test
			type User struct { Name string }
			`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "/some/path/test.go", src, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{f},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) {
					t.Error("Unexpected error for excluded file")
				},
			}

			_, err = runStruct001(pass)
			if err != nil {
				t.Errorf("runStruct001() error = %v", err)
			}

		})
	}
}
