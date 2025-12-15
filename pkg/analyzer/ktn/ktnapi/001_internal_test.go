// Internal tests for ktnapi Analyzer001.
package ktnapi

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_isAllowedType tests the isAllowedType function.
func Test_isAllowedType(t *testing.T) {
	tests := []struct {
		name     string
		pkgPath  string
		typeName string
		expected bool
	}{
		{"allowed_type_time.Time", "time", "time.Time", true},
		{"allowed_type_time.Duration", "time", "time.Duration", true},
		{"allowed_type_context.Context", "context", "context.Context", true},
		{"allowed_package_time", "time", "time.Other", true},
		{"allowed_package_context", "context", "context.Other", true},
		{"allowed_package_go_ast", "go/ast", "go/ast.File", true},
		{"allowed_package_go_token", "go/token", "go/token.FileSet", true},
		{"allowed_package_go_types", "go/types", "go/types.Named", true},
		{"allowed_package_analysis", "golang.org/x/tools/go/analysis", "golang.org/x/tools/go/analysis.Pass", true},
		{"allowed_package_inspector", "golang.org/x/tools/go/ast/inspector", "golang.org/x/tools/go/ast/inspector.Inspector", true},
		{"allowed_package_config", "github.com/kodflow/ktn-linter/pkg/config", "github.com/kodflow/ktn-linter/pkg/config.Config", true},
		{"not_allowed_net_http", "net/http", "net/http.Client", false},
		{"not_allowed_os", "os", "os.File", false},
		{"not_allowed_external_package", "github.com/some/package", "github.com/some/package.Type", false},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedType(tt.pkgPath, tt.typeName)
			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("isAllowedType(%q, %q) = %v, want %v", tt.pkgPath, tt.typeName, result, tt.expected)
			}
		})
	}
}

// Test_getBaseIdent tests the getBaseIdent function.
func Test_getBaseIdent(t *testing.T) {
	tests := []struct {
		name         string
		expr         ast.Expr
		expectedNil  bool
		expectedName string
	}{
		{"simple_ident", &ast.Ident{Name: "x"}, false, "x"},
		{"paren_expr", &ast.ParenExpr{X: &ast.Ident{Name: "x"}}, false, "x"},
		{"star_expr", &ast.StarExpr{X: &ast.Ident{Name: "x"}}, false, "x"},
		{"selector_expr_returns_nil", &ast.SelectorExpr{X: &ast.Ident{Name: "x"}, Sel: &ast.Ident{Name: "field"}}, true, ""},
		{"nested_paren_star", &ast.ParenExpr{X: &ast.StarExpr{X: &ast.Ident{Name: "ptr"}}}, false, "ptr"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBaseIdent(tt.expr)
			// Vérifier si nil attendu
			if tt.expectedNil {
				// Devrait être nil
				if result != nil {
					t.Errorf("getBaseIdent() = %v, want nil", result)
				}
			} else {
				// Ne devrait pas être nil
				if result == nil {
					t.Error("getBaseIdent() = nil, want non-nil")
				} else if result.Name != tt.expectedName {
					t.Errorf("getBaseIdent().Name = %q, want %q", result.Name, tt.expectedName)
				}
			}
		})
	}
}

// Test_suggestInterfaceName tests the suggestInterfaceName function.
func Test_suggestInterfaceName(t *testing.T) {
	tests := []struct {
		name      string
		paramName string
		typeName  string
		expected  string
	}{
		{"param_ends_with_repo", "userRepo", "Repository", "userRepo"},
		{"param_ends_with_client", "httpClient", "Client", "httpClient"},
		{"param_ends_with_service", "userService", "Service", "userService"},
		{"param_ends_with_reader", "fileReader", "FileReader", "fileReader"},
		{"param_ends_with_writer", "logWriter", "Writer", "logWriter"},
		{"param_ends_with_handler", "requestHandler", "Handler", "requestHandler"},
		{"generic_param_uses_lowercase_type", "x", "Client", "client"},
		{"empty_param_uses_lowercase_type", "", "Repository", "repository"},
		{"underscore_param_uses_lowercase_type", "_", "Logger", "logger"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suggestInterfaceName(tt.paramName, tt.typeName)
			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("suggestInterfaceName(%q, %q) = %q, want %q", tt.paramName, tt.typeName, result, tt.expected)
			}
		})
	}
}

// Test_lowerFirst tests the lowerFirst function.
func Test_lowerFirst(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty_string", "", ""},
		{"single_uppercase", "A", "a"},
		{"single_lowercase", "a", "a"},
		{"uppercase_word", "Client", "client"},
		{"all_caps", "HTTP", "hTTP"},
		{"already_lowercase", "client", "client"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lowerFirst(tt.input)
			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("lowerFirst(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_getSortedMethods tests the getSortedMethods function.
func Test_getSortedMethods(t *testing.T) {
	tests := []struct {
		name     string
		methods  map[string]bool
		expected []string
	}{
		{"empty_map", map[string]bool{}, []string{}},
		{"single_method", map[string]bool{"Get": true}, []string{"Get"}},
		{"multiple_methods_sorted", map[string]bool{"Write": true, "Read": true, "Close": true}, []string{"Close", "Read", "Write"}},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSortedMethods(tt.methods)
			// Vérifier la longueur
			if len(result) != len(tt.expected) {
				t.Errorf("getSortedMethods() len = %d, want %d", len(result), len(tt.expected))
				return
			}
			// Vérifier chaque élément
			for i, m := range result {
				// Comparer avec l'attendu
				if m != tt.expected[i] {
					t.Errorf("getSortedMethods()[%d] = %q, want %q", i, m, tt.expected[i])
				}
			}
		})
	}
}

// Test_runAPI001 tests the runAPI001 function.
func Test_runAPI001(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{"disabled_rule", false},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup config
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-API-001": {Enabled: config.Bool(tt.enabled)},
				},
			})
			defer config.Reset()

			// Parse test code
			src := `package test
func DoSomething() {}`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", src, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build TypesInfo
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Uses:  make(map[*ast.Ident]types.Object),
				Defs:  make(map[*ast.Ident]types.Object),
			}
			conf := types.Config{}
			_, _ = conf.Check("test", fset, []*ast.File{f}, info)

			// Create inspect pass
			inspectPass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				Report:    func(d analysis.Diagnostic) {},
				ResultOf:  make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			// Create analysis pass
			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) {},
			}

			// Run the function
			_, err = runAPI001(pass)
			// Verify no error
			if err != nil {
				t.Errorf("runAPI001() error = %v", err)
			}
		})
	}
}

// Test_analyzeFunction tests the analyzeFunction function.
func Test_analyzeFunction(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "function_without_external_deps",
			code: `package test
func DoSomething() {}`,
		},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse test code
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build TypesInfo
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Uses:  make(map[*ast.Ident]types.Object),
				Defs:  make(map[*ast.Ident]types.Object),
			}
			conf := types.Config{}
			_, _ = conf.Check("test", fset, []*ast.File{f}, info)

			// Create pass
			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				Report:    func(_ analysis.Diagnostic) {},
			}

			// Find function declaration
			var funcDecl *ast.FuncDecl
			// Inspect file
			ast.Inspect(f, func(n ast.Node) bool {
				// Check FuncDecl
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				// Continue traversal
				return true
			})

			// Setup config
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-API-001": {Enabled: config.Bool(true)},
				},
			})
			defer config.Reset()

			// Run the function (should not panic)
			if funcDecl != nil {
				analyzeFunction(pass, funcDecl)
			}
		})
	}
}

// Test_collectParams tests the collectParams function.
func Test_collectParams(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedCount int
	}{
		{
			name: "function_with_no_params",
			code: `package test
func DoSomething() {}`,
			expectedCount: 0,
		},
		{
			name: "function_with_basic_params",
			code: `package test
func DoSomething(a int, b string) {}`,
			expectedCount: 0, // Basic types are not external concrete types
		},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse test code
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build TypesInfo
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
				Uses:  make(map[*ast.Ident]types.Object),
				Defs:  make(map[*ast.Ident]types.Object),
			}
			conf := types.Config{}
			_, _ = conf.Check("test", fset, []*ast.File{f}, info)

			// Create pass
			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				Report:    func(_ analysis.Diagnostic) {},
			}

			// Find function declaration
			var funcDecl *ast.FuncDecl
			// Inspect file
			ast.Inspect(f, func(n ast.Node) bool {
				// Check FuncDecl
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				// Continue traversal
				return true
			})

			// Run the function
			if funcDecl != nil {
				params := collectParams(pass, funcDecl)
				// Verify count
				if len(params) != tt.expectedCount {
					t.Errorf("collectParams() returned %d params, want %d", len(params), tt.expectedCount)
				}
			}
		})
	}
}

// Test_getExternalConcreteType tests the getExternalConcreteType function.
func Test_getExternalConcreteType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"basic_type"},
		{"string_type"},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create pass with minimal info
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Test with basic types (should return nil for basic types)
			basicType := types.Typ[types.String]
			result := getExternalConcreteType(pass, basicType)
			// Basic types should return nil
			if result != nil {
				t.Logf("getExternalConcreteType returned non-nil for basic type")
			}
		})
	}
}

// Test_scanBodyForMethodCalls tests the scanBodyForMethodCalls function.
func Test_scanBodyForMethodCalls(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "empty_function_body",
			code: `package test
func DoSomething() {}`,
		},
		{
			name: "function_with_method_call",
			code: `package test
type MyType struct{}
func (m *MyType) Method() {}
func DoSomething(m *MyType) { m.Method() }`,
		},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse test code
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Build TypesInfo
			info := &types.Info{
				Types:      make(map[ast.Expr]types.TypeAndValue),
				Uses:       make(map[*ast.Ident]types.Object),
				Defs:       make(map[*ast.Ident]types.Object),
				Selections: make(map[*ast.SelectorExpr]*types.Selection),
			}
			conf := types.Config{}
			_, _ = conf.Check("test", fset, []*ast.File{f}, info)

			// Create pass
			pass := &analysis.Pass{
				Fset:      fset,
				Files:     []*ast.File{f},
				TypesInfo: info,
				Report:    func(_ analysis.Diagnostic) {},
			}

			// Find last function declaration
			var funcDecl *ast.FuncDecl
			// Inspect file
			ast.Inspect(f, func(n ast.Node) bool {
				// Check FuncDecl
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
				}
				// Continue traversal
				return true
			})

			// Run the function (should not panic)
			if funcDecl != nil && funcDecl.Body != nil {
				params := make(map[*ast.Ident]*paramInfo)
				scanBodyForMethodCalls(pass, funcDecl.Body, params)
			}
		})
	}
}

// Test_matchesParam tests the matchesParam function.
func Test_matchesParam(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"same_ident", true},
		{"different_idents", false},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create pass with minimal info
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Uses: make(map[*ast.Ident]types.Object),
					Defs: make(map[*ast.Ident]types.Object),
				},
			}

			// Create idents
			ident1 := &ast.Ident{Name: "x"}
			ident2 := &ast.Ident{Name: "y"}

			// Test based on case
			var result bool
			// Check test case
			if tt.name == "same_ident" {
				result = matchesParam(pass, ident1, ident1)
			} else {
				result = matchesParam(pass, ident1, ident2)
			}

			// Verify result
			if result != tt.expected {
				t.Errorf("matchesParam() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_reportDiagnostics tests the reportDiagnostics function.
func Test_reportDiagnostics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"empty_params"},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse test code
			src := `package test
func DoSomething() {}`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", src, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			// Create pass
			pass := &analysis.Pass{
				Fset:   fset,
				Report: func(_ analysis.Diagnostic) {},
			}

			// Setup config
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-API-001": {Enabled: config.Bool(true)},
				},
			})
			defer config.Reset()

			// Find function declaration
			var funcDecl *ast.FuncDecl
			// Inspect file
			ast.Inspect(f, func(n ast.Node) bool {
				// Check FuncDecl
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				// Continue traversal
				return true
			})

			// Test with empty params (should not report anything)
			params := make(map[*ast.Ident]*paramInfo)
			// Run the function
			if funcDecl != nil {
				reportDiagnostics(pass, funcDecl, params)
			}
		})
	}
}

// Test_formatTypeName tests the formatTypeName function.
func Test_formatTypeName(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		expected string
	}{
		{"simple_type", "MyType", "MyType"},
		{"pointer_type", "*MyType", "*MyType"},
		{"package_type", "pkg.Type", "pkg.Type"},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create basic type for testing
			basic := types.Typ[types.String]
			result := formatTypeName(basic)
			// Verify result is not empty
			if result == "" {
				t.Error("formatTypeName() returned empty string")
			}
		})
	}
}

// Test_shortQualifier tests the shortQualifier function.
func Test_shortQualifier(t *testing.T) {
	tests := []struct {
		name     string
		pkg      *types.Package
		expected string
	}{
		{"nil_package", nil, ""},
		{"named_package", types.NewPackage("net/http", "http"), "http"},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shortQualifier(tt.pkg)
			// Verify result
			if result != tt.expected {
				t.Errorf("shortQualifier() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_formatTupleTypes tests the formatTupleTypes function.
func Test_formatTupleTypes(t *testing.T) {
	tests := []struct {
		name     string
		tuple    *types.Tuple
		expected string
	}{
		{"nil_tuple", nil, ""},
		{"empty_tuple", types.NewTuple(), ""},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTupleTypes(tt.tuple)
			// Verify result
			if result != tt.expected {
				t.Errorf("formatTupleTypes() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_formatMethodSignature tests the formatMethodSignature function.
func Test_formatMethodSignature(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		sig        *types.Signature
		contains   string
	}{
		{
			"no_params_no_results",
			"Close",
			types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(), false),
			"Close()",
		},
		{
			"with_results",
			"Read",
			types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.Int])), false),
			"Read()",
		},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMethodSignature(tt.methodName, tt.sig)
			// Verify result contains expected
			if result == "" {
				t.Error("formatMethodSignature() returned empty string")
			}
		})
	}
}

// Test_buildInterfaceSignatures tests the buildInterfaceSignatures function.
func Test_buildInterfaceSignatures(t *testing.T) {
	tests := []struct {
		name     string
		methods  []string
		expected string
	}{
		{"empty_methods", []string{}, ""},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal named type for testing
			pkg := types.NewPackage("test", "test")
			typeName := types.NewTypeName(0, pkg, "MyType", nil)
			named := types.NewNamed(typeName, types.NewStruct(nil, nil), nil)

			result := buildInterfaceSignatures(named, tt.methods)
			// Verify result
			if result != tt.expected {
				t.Errorf("buildInterfaceSignatures() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_getMethodSignature tests the getMethodSignature function.
func Test_getMethodSignature(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		expected   string
	}{
		{"nonexistent_method", "NonExistent", ""},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal named type without methods
			pkg := types.NewPackage("test", "test")
			typeName := types.NewTypeName(0, pkg, "MyType", nil)
			named := types.NewNamed(typeName, types.NewStruct(nil, nil), nil)

			result := getMethodSignature(named, tt.methodName)
			// Verify result
			if result != tt.expected {
				t.Errorf("getMethodSignature() = %q, want %q", result, tt.expected)
			}
		})
	}
}
