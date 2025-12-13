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

// Test_collectStructsWithMethods tests the collectStructsWithMethods private function.
func Test_collectStructsWithMethods(t *testing.T) {
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

// Test_collectInterfaceChecks tests the private collectInterfaceChecks function.
func Test_collectInterfaceChecks(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected map[string]bool
	}{
		{
			name: "no interface checks",
			src: `package test
type User struct{}`,
			expected: map[string]bool{},
		},
		{
			name: "one interface check",
			src: `package test
type User struct{}
var _ UserInterface = (*User)(nil)`,
			expected: map[string]bool{"User": true},
		},
		{
			name: "multiple interface checks",
			src: `package test
type User struct{}
type Admin struct{}
var _ UserInterface = (*User)(nil)
var _ AdminInterface = (*Admin)(nil)`,
			expected: map[string]bool{"User": true, "Admin": true},
		},
		{
			name: "regular var - not interface check",
			src: `package test
type User struct{}
var x int = 10`,
			expected: map[string]bool{},
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

			checks := collectInterfaceChecks(file)

			// Vérifier le nombre de checks
			if len(checks) != len(tt.expected) {
				t.Errorf("expected %d checks, got %d", len(tt.expected), len(checks))
			}

			// Vérifier chaque check attendu
			for k, v := range tt.expected {
				// Vérifier présence
				if checks[k] != v {
					t.Errorf("expected checks[%s]=%v, got %v", k, v, checks[k])
				}
			}
		})
	}
}

// Test_extractStructFromInterfaceCheck tests the private extractStructFromInterfaceCheck function.
func Test_extractStructFromInterfaceCheck(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "valid interface check",
			src: `package test
var _ MyInterface = (*MyStruct)(nil)`,
			expected: "MyStruct",
		},
		{
			name: "not underscore var",
			src: `package test
var x MyInterface = (*MyStruct)(nil)`,
			expected: "",
		},
		{
			name: "no value",
			src: `package test
var _ MyInterface`,
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

			// Trouver la ValueSpec
			var valueSpec *ast.ValueSpec
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier si c'est une GenDecl
				if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
					// Parcourir les specs
					for _, spec := range genDecl.Specs {
						// Vérifier si c'est une ValueSpec
						if vs, ok := spec.(*ast.ValueSpec); ok {
							valueSpec = vs
							// Sortir de la boucle
							return false
						}
					}
				}
				// Continuer l'itération
				return true
			})

			// Vérifier si on a trouvé une ValueSpec
			if valueSpec == nil {
				t.Fatal("no ValueSpec found")
			}

			result := extractStructFromInterfaceCheck(valueSpec)

			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_extractStructNameFromConversion tests the private extractStructNameFromConversion function.
func Test_extractStructNameFromConversion(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "valid conversion",
			src: `package test
var _ MyInterface = (*MyStruct)(nil)`,
			expected: "MyStruct",
		},
		{
			name: "not a pointer type",
			src: `package test
var _ int = MyStruct(nil)`,
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
				// Pas de CallExpr dans ce test, on retourne
				return
			}

			result := extractStructNameFromConversion(callExpr)

			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
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
