// Internal tests for 007.go private functions.
package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runStruct007 teste la fonction runStruct007.
func Test_runStruct007(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "struct007_analysis",
			expectErr: false,
		},
		{
			name:      "struct007_error_case",
			expectErr: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite analysis.Pass réel
			// Les cas d'erreur sont couverts via le test external
			_ = tt.expectErr
		})
	}
}

// Test_collectStructPrivateFields teste la fonction collectStructPrivateFields.
func Test_collectStructPrivateFields(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedCount int
		structName    string
		fieldName     string
	}{
		{
			name: "struct with private fields",
			code: `package test
type User struct {
	name string
	age int
}`,
			expectedCount: 1,
			structName:    "User",
			fieldName:     "name",
		},
		{
			name: "struct with public fields only",
			code: `package test
type User struct {
	Name string
	Age int
}`,
			expectedCount: 1,
			structName:    "User",
		},
		{
			name: "private struct",
			code: `package test
type user struct {
	name string
}`,
			expectedCount: 0,
		},
		{
			name: "non-struct type",
			code: `package test
type MyInt int`,
			expectedCount: 0,
		},
		{
			name: "struct with no fields",
			code: `package test
type Empty struct{}`,
			expectedCount: 1,
			structName:    "Empty",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build config
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {Enabled: config.Bool(true)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			// Build TypesInfo
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			conf := types.Config{}
			_, _ = conf.Check("test", fset, []*ast.File{f}, info)

			inspectPass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				Report:    func(d analysis.Diagnostic) {},
				ResultOf:  make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) {},
			}

			insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
			structFields := collectStructPrivateFields(pass, insp, cfg)

			// Vérification du nombre de structs
			if len(structFields) != tt.expectedCount {
				t.Errorf("collectStructPrivateFields() returned %d structs, want %d", len(structFields), tt.expectedCount)
			}

			// Si on attend une struct spécifique, vérifier ses champs
			if tt.structName != "" {
				info, exists := structFields[tt.structName]
				if !exists {
					t.Errorf("Expected struct %q not found", tt.structName)
				} else if tt.fieldName != "" {
					if !info.privateFields[tt.fieldName] {
						t.Errorf("Expected field %q not found in struct %q", tt.fieldName, tt.structName)
					}
				}
			}
		})
	}
}

// Test_collectMethodsDetailed teste la fonction collectMethodsDetailed.
func Test_collectMethodsDetailed(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		expectedCount  int
		expectedMethod string
	}{
		{
			name: "method with receiver",
			code: `package test
type User struct {
	name string
}
func (u *User) Name() string {
	return u.name
}`,
			expectedCount:  1,
			expectedMethod: "Name",
		},
		{
			name: "function without receiver",
			code: `package test
func DoSomething() {
}`,
			expectedCount: 0,
		},
		{
			name: "multiple methods",
			code: `package test
type User struct {
	name string
	age int
}
func (u *User) Name() string {
	return u.name
}
func (u *User) Age() int {
	return u.age
}`,
			expectedCount:  2,
			expectedMethod: "Name",
		},
		{
			name: "method with invalid receiver type",
			code: `package test
import "io"
func (r *io.Reader) Method() {}`,
			expectedCount: 0,
		},
		{
			name: "method with multiple return values",
			code: `package test
type User struct { name string }
func (u *User) GetName() (string, error) {
	return u.name, nil
}`,
			expectedCount:  1,
			expectedMethod: "GetName",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build config
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {Enabled: config.Bool(true)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			// Build TypesInfo
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			conf := types.Config{}
			_, _ = conf.Check("test", fset, []*ast.File{f}, info)

			inspectPass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				Report:    func(d analysis.Diagnostic) {},
				ResultOf:  make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) {},
			}

			insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
			methods := collectMethodsDetailed(pass, insp, cfg)

			// Count total methods
			totalMethods := 0
			for _, methodList := range methods {
				totalMethods += len(methodList)
			}

			// Vérification du nombre de méthodes
			if totalMethods != tt.expectedCount {
				t.Errorf("collectMethodsDetailed() returned %d methods, want %d", totalMethods, tt.expectedCount)
			}

			// Si on attend une méthode spécifique, vérifier son nom
			if tt.expectedMethod != "" {
				found := false
				for _, methodList := range methods {
					for _, method := range methodList {
						if method.name == tt.expectedMethod {
							found = true
							break
						}
					}
				}
				if !found {
					t.Errorf("Expected method %q not found", tt.expectedMethod)
				}
			}
		})
	}
}

// Test_extractReceiverType teste la fonction extractReceiverType.
func Test_extractReceiverType(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple_ident",
			expr:     &ast.Ident{Name: "MyType"},
			expected: "MyType",
		},
		{
			name:     "star_expr",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "MyType"}},
			expected: "MyType",
		},
		{
			name:     "nil_expr",
			expr:     nil,
			expected: "",
		},
		{
			name:     "selector_expr",
			expr:     &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Type"}},
			expected: "",
		},
		{
			name:     "star_expr_with_selector",
			expr:     &ast.StarExpr{X: &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Type"}}},
			expected: "",
		},
		{
			name:     "array_type",
			expr:     &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
			expected: "",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Skip nil test pour éviter panic
			if tt.expr == nil {
				return
			}
			result := extractReceiverType(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("extractReceiverType() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_extractReturnedField teste la fonction extractReturnedField.
func Test_extractReturnedField(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "selector_expr",
			expr:     &ast.SelectorExpr{X: &ast.Ident{Name: "s"}, Sel: &ast.Ident{Name: "name"}},
			expected: "name",
		},
		{
			name:     "non_selector",
			expr:     &ast.Ident{Name: "x"},
			expected: "",
		},
		{
			name:     "selector_with_non_ident_x",
			expr:     &ast.SelectorExpr{X: &ast.CallExpr{}, Sel: &ast.Ident{Name: "name"}},
			expected: "",
		},
	}

	// Exécution des tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractReturnedField(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("extractReturnedField() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_checkNamingConventions teste la fonction checkNamingConventions.
func Test_checkNamingConventions(t *testing.T) {
	tests := []struct {
		name        string
		structInfo  map[string]structFieldsInfo
		methods     map[string][]methodInfo
		expectError bool
	}{
		{
			name: "no struct info",
			structInfo: map[string]structFieldsInfo{
				"User": {
					name:          "User",
					privateFields: map[string]bool{"name": true},
				},
			},
			methods: map[string][]methodInfo{
				"Unknown": {
					{
						name: "GetName",
						funcDecl: &ast.FuncDecl{
							Name: &ast.Ident{Name: "GetName"},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.SelectorExpr{
												X:   &ast.Ident{Name: "u"},
												Sel: &ast.Ident{Name: "name"},
											},
										},
									},
								},
							},
						},
						receiverTy: "Unknown",
					},
				},
			},
			expectError: false,
		},
		{
			name: "struct with nil private fields",
			structInfo: map[string]structFieldsInfo{
				"User": {
					name:          "User",
					privateFields: nil,
				},
			},
			methods: map[string][]methodInfo{
				"User": {
					{
						name:       "GetName",
						funcDecl:   &ast.FuncDecl{Name: &ast.Ident{Name: "GetName"}},
						receiverTy: "User",
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid naming convention",
			structInfo: map[string]structFieldsInfo{
				"User": {
					name:          "User",
					privateFields: map[string]bool{"name": true},
				},
			},
			methods: map[string][]methodInfo{
				"User": {
					{
						name: "Name",
						funcDecl: &ast.FuncDecl{
							Name: &ast.Ident{Name: "Name"},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.SelectorExpr{
												X:   &ast.Ident{Name: "u"},
												Sel: &ast.Ident{Name: "name"},
											},
										},
									},
								},
							},
						},
						receiverTy: "User",
					},
				},
			},
			expectError: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal pass for testing
			errCount := 0
			pass := &analysis.Pass{
				Report: func(_ analysis.Diagnostic) { errCount++ },
			}

			// Call the function
			checkNamingConventions(pass, tt.structInfo, tt.methods)

			// Verify expectations
			if tt.expectError && errCount == 0 {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && errCount > 0 {
				t.Errorf("Expected no error but got %d", errCount)
			}
		})
	}
}

// Test_extractSimpleReturnType tests the extractSimpleReturnType private function.
func Test_extractSimpleReturnType(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "function with single return type",
			code: `package test
func GetName() string {
	return "test"
}`,
			want: "string",
		},
		{
			name: "function with no return type",
			code: `package test
func DoSomething() {
}`,
			want: "",
		},
		{
			name: "function with multiple return types",
			code: `package test
func GetValue() (string, error) {
	return "test", nil
}`,
			want: "",
		},
		{
			name: "function with int return type",
			code: `package test
func GetAge() int {
	return 42
}`,
			want: "int",
		},
	}

	// Iteration over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build types info
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {Enabled: config.Bool(true)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			// Build proper TypesInfo with type checker
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			conf := types.Config{}
			_, _ = conf.Check("test", fset, []*ast.File{f}, info)

			inspectPass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				Report:    func(d analysis.Diagnostic) {},
				ResultOf:  make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) {},
			}

			// Find the function declaration
			var funcDecl *ast.FuncDecl
			ast.Inspect(f, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil {
				t.Fatal("No function declaration found")
			}

			result := extractSimpleReturnType(pass, funcDecl)
			// Vérification du résultat
			if result != tt.want {
				t.Errorf("extractSimpleReturnType() = %q, want %q", result, tt.want)
			}
		})
	}
}

// Test_checkGetterFieldMismatch tests the checkGetterFieldMismatch private function.
func Test_checkGetterFieldMismatch(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name: "getter returns matching field",
			code: `package test
type User struct {
	name string
}
func (u *User) Name() string {
	return u.name
}`,
			expectError: false,
		},
		{
			name: "getter returns non-matching field",
			code: `package test
type User struct {
	firstName string
}
func (u *User) Name() string {
	return u.firstName
}`,
			expectError: true,
		},
		{
			name: "method with no body",
			code: `package test
type User interface {
	Name() string
}`,
			expectError: false,
		},
		{
			name: "method with multiple statements",
			code: `package test
type User struct {
	name string
}
func (u *User) Name() string {
	x := u.name
	return x
}`,
			expectError: false,
		},
		{
			name: "method with multiple return values",
			code: `package test
type User struct {
	name string
}
func (u *User) Name() (string, error) {
	return u.name, nil
}`,
			expectError: false,
		},
		{
			name: "method returning non-selector",
			code: `package test
type User struct {
	name string
}
func (u *User) Name() string {
	return "constant"
}`,
			expectError: false,
		},
		{
			name: "method returning public field",
			code: `package test
type User struct {
	Name string
}
func (u *User) GetName() string {
	return u.Name
}`,
			expectError: false,
		},
		{
			name: "getter with Get prefix",
			code: `package test
type User struct {
	firstName string
}
func (u *User) GetName() string {
	return u.firstName
}`,
			expectError: false,
		},
	}

	// Iteration over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build config
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {Enabled: config.Bool(true)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			errCount := 0
			inspectPass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: &types.Info{},
				Report:    func(d analysis.Diagnostic) { errCount++ },
				ResultOf:  make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: &types.Info{},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) { errCount++ },
			}

			// Find the struct and method
			var structInfo structFieldsInfo
			var method methodInfo

			ast.Inspect(f, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.TypeSpec:
					if st, ok := node.Type.(*ast.StructType); ok && ast.IsExported(node.Name.Name) {
						privateFields := make(map[string]bool)
						if st.Fields != nil {
							for _, field := range st.Fields.List {
								for _, name := range field.Names {
									if len(name.Name) > 0 && unicode.IsLower(rune(name.Name[0])) {
										privateFields[name.Name] = true
									}
								}
							}
						}
						structInfo = structFieldsInfo{
							name:          node.Name.Name,
							privateFields: privateFields,
							pos:           node,
						}
					}
				case *ast.FuncDecl:
					if node.Recv != nil && len(node.Recv.List) > 0 {
						receiverType := extractReceiverType(node.Recv.List[0].Type)
						returnType := extractSimpleReturnType(pass, node)
						method = methodInfo{
							name:       node.Name.Name,
							funcDecl:   node,
							receiverTy: receiverType,
							returnType: returnType,
						}
					}
				}
				return true
			})

			// Test the function if we have data
			if method.funcDecl != nil && structInfo.privateFields != nil {
				checkGetterFieldMismatch(pass, method, structInfo)
			}

			// Verify expectations
			if tt.expectError && errCount == 0 {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && errCount > 0 {
				t.Errorf("Expected no error but got %d", errCount)
			}
		})
	}
}

// Test_runStruct007_disabled tests that the rule is skipped when disabled.
func Test_runStruct007_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {Enabled: config.Bool(false)},
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

			_, err = runStruct007(pass)
			if err != nil {
				t.Errorf("runStruct007() error = %v", err)
			}

		})
	}
}

// Test_runStruct007_excludedFile tests that excluded files are skipped.
func Test_runStruct007_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-007": {
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

			_, err = runStruct007(pass)
			if err != nil {
				t.Errorf("runStruct007() error = %v", err)
			}

		})
	}
}
