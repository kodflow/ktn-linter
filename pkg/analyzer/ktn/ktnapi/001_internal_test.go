// Internal tests for ktnapi Analyzer001.
package ktnapi

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
		{"allowed_package_strings", "strings", "strings.Builder", true},
		{"allowed_package_bytes", "bytes", "bytes.Buffer", true},
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
	// Create test packages
	httpPkg := types.NewPackage("net/http", "http")
	timePkg := types.NewPackage("time", "time")
	localPkg := types.NewPackage("myapp/internal", "internal")

	// Create named types
	clientTypeName := types.NewTypeName(0, httpPkg, "Client", nil)
	clientNamed := types.NewNamed(clientTypeName, types.NewStruct(nil, nil), nil)

	timeTypeName := types.NewTypeName(0, timePkg, "Time", nil)
	timeNamed := types.NewNamed(timeTypeName, types.NewStruct(nil, nil), nil)

	localTypeName := types.NewTypeName(0, localPkg, "MyService", nil)
	localNamed := types.NewNamed(localTypeName, types.NewStruct(nil, nil), nil)

	tests := []struct {
		name        string
		typ         types.Type
		expectNil   bool
		expectPkg   string
	}{
		{"basic_string_returns_nil", types.Typ[types.String], true, ""},
		{"basic_int_returns_nil", types.Typ[types.Int], true, ""},
		{"external_http_client", clientNamed, false, "net/http"},
		{"pointer_to_external", types.NewPointer(clientNamed), false, "net/http"},
		{"allowed_time_returns_nil", timeNamed, true, ""},
		{"internal_package_returns_nil", localNamed, true, ""},
		{"slice_returns_nil", types.NewSlice(types.Typ[types.Int]), true, ""},
		{"interface_returns_nil", types.NewInterfaceType(nil, nil), true, ""},
	}

	// Create pass with minimal info
	pass := &analysis.Pass{
		Pkg: localPkg,
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getExternalConcreteType(pass, tt.typ)
			// Verify nil expectation
			if tt.expectNil {
				// Should be nil
				if result != nil {
					t.Errorf("getExternalConcreteType() = %v, want nil", result)
				}
			} else {
				// Should not be nil
				if result == nil {
					t.Error("getExternalConcreteType() = nil, want non-nil")
				} else if result.Obj().Pkg().Path() != tt.expectPkg {
					t.Errorf("getExternalConcreteType().Pkg().Path() = %q, want %q",
						result.Obj().Pkg().Path(), tt.expectPkg)
				}
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
	// Create test packages
	httpPkg := types.NewPackage("net/http", "http")
	osPkg := types.NewPackage("os", "os")
	testPkg := types.NewPackage("test", "test")

	// Create named types for testing
	clientTypeName := types.NewTypeName(0, httpPkg, "Client", nil)
	clientNamed := types.NewNamed(clientTypeName, types.NewStruct(nil, nil), nil)

	fileTypeName := types.NewTypeName(0, osPkg, "File", nil)
	fileNamed := types.NewNamed(fileTypeName, types.NewStruct(nil, nil), nil)

	localTypeName := types.NewTypeName(0, testPkg, "MyType", nil)
	localNamed := types.NewNamed(localTypeName, types.NewStruct(nil, nil), nil)

	tests := []struct {
		name     string
		typ      types.Type
		expected string
	}{
		{"basic_string", types.Typ[types.String], "string"},
		{"basic_int", types.Typ[types.Int], "int"},
		{"basic_bool", types.Typ[types.Bool], "bool"},
		{"named_http_client", clientNamed, "http.Client"},
		{"named_os_file", fileNamed, "os.File"},
		{"named_local_type", localNamed, "test.MyType"},
		{"pointer_to_named", types.NewPointer(clientNamed), "*http.Client"},
		{"slice_of_int", types.NewSlice(types.Typ[types.Int]), "[]int"},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTypeName(tt.typ)
			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("formatTypeName() = %q, want %q", result, tt.expected)
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
		{"single_int", types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.Int])), "int"},
		{"single_string", types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])), "string"},
		{"two_types", types.NewTuple(
			types.NewVar(0, nil, "", types.Typ[types.Int]),
			types.NewVar(0, nil, "", types.Universe.Lookup("error").Type()),
		), "int, error"},
		{"three_types", types.NewTuple(
			types.NewVar(0, nil, "", types.Typ[types.String]),
			types.NewVar(0, nil, "", types.Typ[types.Int]),
			types.NewVar(0, nil, "", types.Typ[types.Bool]),
		), "string, int, bool"},
		{"slice_type", types.NewTuple(
			types.NewVar(0, nil, "", types.NewSlice(types.Typ[types.Byte])),
		), "[]uint8"}, // Note: byte is an alias for uint8 in Go type system
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
		expected   string
	}{
		{
			"no_params_no_results",
			"Close",
			types.NewSignatureType(nil, nil, nil, types.NewTuple(), types.NewTuple(), false),
			"Close()",
		},
		{
			"with_params_and_results",
			"Read",
			types.NewSignatureType(nil, nil, nil,
				types.NewTuple(types.NewVar(0, nil, "p", types.NewSlice(types.Typ[types.Byte]))),
				types.NewTuple(types.NewVar(0, nil, "n", types.Typ[types.Int]), types.NewVar(0, nil, "err", types.Universe.Lookup("error").Type())),
				false),
			"Read(p []uint8) (n int, err error)",
		},
		{
			"with_named_param_and_result",
			"Get",
			types.NewSignatureType(nil, nil, nil,
				types.NewTuple(types.NewVar(0, nil, "url", types.Typ[types.String])),
				types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])),
				false),
			"Get(url string) string",
		},
		{
			"single_result_no_parens",
			"String",
			types.NewSignatureType(nil, nil, nil,
				types.NewTuple(),
				types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])),
				false),
			"String() string",
		},
		{
			"only_error_result",
			"Validate",
			types.NewSignatureType(nil, nil, nil,
				types.NewTuple(),
				types.NewTuple(types.NewVar(0, nil, "", types.Universe.Lookup("error").Type())),
				false),
			"Validate() error",
		},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMethodSignature(tt.methodName, tt.sig)
			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("formatMethodSignature() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_buildInterfaceSignatures tests the buildInterfaceSignatures function.
func Test_buildInterfaceSignatures(t *testing.T) {
	// Create test package and type with methods
	pkg := types.NewPackage("test", "test")
	typeName := types.NewTypeName(0, pkg, "MyService", nil)
	named := types.NewNamed(typeName, types.NewStruct(nil, nil), nil)

	// Add methods to the named type
	closeSig := types.NewSignatureType(
		types.NewVar(0, pkg, "s", types.NewPointer(named)), nil, nil,
		types.NewTuple(),
		types.NewTuple(types.NewVar(0, nil, "", types.Universe.Lookup("error").Type())),
		false)
	closeFunc := types.NewFunc(0, pkg, "Close", closeSig)
	named.AddMethod(closeFunc)

	readSig := types.NewSignatureType(
		types.NewVar(0, pkg, "s", types.NewPointer(named)), nil, nil,
		types.NewTuple(types.NewVar(0, nil, "p", types.NewSlice(types.Typ[types.Byte]))),
		types.NewTuple(
			types.NewVar(0, nil, "n", types.Typ[types.Int]),
			types.NewVar(0, nil, "err", types.Universe.Lookup("error").Type())),
		false)
	readFunc := types.NewFunc(0, pkg, "Read", readSig)
	named.AddMethod(readFunc)

	stringSig := types.NewSignatureType(
		types.NewVar(0, pkg, "s", types.NewPointer(named)), nil, nil,
		types.NewTuple(),
		types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])),
		false)
	stringFunc := types.NewFunc(0, pkg, "String", stringSig)
	named.AddMethod(stringFunc)

	tests := []struct {
		name     string
		methods  []string
		expected string
	}{
		{"empty_methods", []string{}, ""},
		{"single_method_close", []string{"Close"}, "\tClose() error"},
		{"single_method_string", []string{"String"}, "\tString() string"},
		{"single_method_read", []string{"Read"}, "\tRead(p []uint8) (n int, err error)"},
		{"multiple_methods_read_close", []string{"Read", "Close"}, "\tRead(p []uint8) (n int, err error)\n\tClose() error"},
		{"nonexistent_method_placeholder", []string{"NonExistent"}, "\tNonExistent(...)"},
		{"mixed_existent_nonexistent", []string{"Close", "NonExistent"}, "\tClose() error\n\tNonExistent(...)"},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	// Create test package and type with methods
	pkg := types.NewPackage("test", "test")
	typeName := types.NewTypeName(0, pkg, "MyService", nil)
	named := types.NewNamed(typeName, types.NewStruct(nil, nil), nil)

	// Add Close method
	closeSig := types.NewSignatureType(
		types.NewVar(0, pkg, "s", types.NewPointer(named)), nil, nil,
		types.NewTuple(),
		types.NewTuple(types.NewVar(0, nil, "", types.Universe.Lookup("error").Type())),
		false)
	closeFunc := types.NewFunc(0, pkg, "Close", closeSig)
	named.AddMethod(closeFunc)

	// Add Read method with params and multiple returns
	readSig := types.NewSignatureType(
		types.NewVar(0, pkg, "s", types.NewPointer(named)), nil, nil,
		types.NewTuple(types.NewVar(0, nil, "p", types.NewSlice(types.Typ[types.Byte]))),
		types.NewTuple(
			types.NewVar(0, nil, "n", types.Typ[types.Int]),
			types.NewVar(0, nil, "err", types.Universe.Lookup("error").Type())),
		false)
	readFunc := types.NewFunc(0, pkg, "Read", readSig)
	named.AddMethod(readFunc)

	// Add String method with single return
	stringSig := types.NewSignatureType(
		types.NewVar(0, pkg, "s", types.NewPointer(named)), nil, nil,
		types.NewTuple(),
		types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])),
		false)
	stringFunc := types.NewFunc(0, pkg, "String", stringSig)
	named.AddMethod(stringFunc)

	tests := []struct {
		name       string
		methodName string
		expected   string
	}{
		{"nonexistent_method", "NonExistent", ""},
		{"close_method", "Close", "Close() error"},
		{"read_method", "Read", "Read(p []uint8) (n int, err error)"},
		{"string_method", "String", "String() string"},
	}

	// Iterate over tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMethodSignature(named, tt.methodName)
			// Verify result
			if result != tt.expected {
				t.Errorf("getMethodSignature() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_scanBodyForFieldAccess tests the scanBodyForFieldAccess function.
func Test_scanBodyForFieldAccess(t *testing.T) {
	tests := []struct {
		name               string
		code               string
		expectedFieldAccess bool
	}{
		{
			name: "no_field_access",
			code: `package test
import "bytes"
func DoSomething(buf *bytes.Buffer) string {
	return buf.String()
}`,
			expectedFieldAccess: false,
		},
		{
			name: "with_field_access",
			code: `package test
import "net/http"
func DoSomething(req *http.Request) string {
	return req.Method
}`,
			expectedFieldAccess: true,
		},
		{
			name: "mixed_field_and_method",
			code: `package test
import "net/http"
func DoSomething(req *http.Request) string {
	_ = req.Context()
	return req.Method
}`,
			expectedFieldAccess: true,
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
			conf := types.Config{Importer: importer.Default()}
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
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Name.Name == "DoSomething" {
					funcDecl = fd
					return false
				}
				// Continue traversal
				return true
			})

			// Skip if no function found
			if funcDecl == nil {
				t.Fatal("Function not found")
			}

			// Create paramInfo map
			params := make(map[*ast.Ident]*paramInfo)
			// Get param ident
			if len(funcDecl.Type.Params.List) > 0 {
				// Get first param
				for _, name := range funcDecl.Type.Params.List[0].Names {
					params[name] = &paramInfo{
						ident:         name,
						methodsCalled: make(map[string]bool),
					}
				}
			}

			// Run scanBodyForFieldAccess
			scanBodyForFieldAccess(pass, funcDecl.Body, params)

			// Check result
			hasFieldAccess := false
			// Check all params
			for _, info := range params {
				// Check field access flag
				if info.hasFieldAccess {
					hasFieldAccess = true
					break
				}
			}

			// Verify result
			if hasFieldAccess != tt.expectedFieldAccess {
				t.Errorf("scanBodyForFieldAccess() hasFieldAccess = %v, want %v", hasFieldAccess, tt.expectedFieldAccess)
			}
		})
	}
}
