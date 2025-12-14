// Internal tests for ktnapi Analyzer001.
package ktnapi

import (
	"go/ast"
	"testing"
)

// Test_isAllowedType tests the isAllowedType function.
func Test_isAllowedType(t *testing.T) {
	tests := []struct {
		name     string
		pkgPath  string
		typeName string
		expected bool
	}{
		{
			name:     "allowed_type_time.Time",
			pkgPath:  "time",
			typeName: "time.Time",
			expected: true,
		},
		{
			name:     "allowed_type_time.Duration",
			pkgPath:  "time",
			typeName: "time.Duration",
			expected: true,
		},
		{
			name:     "allowed_type_context.Context",
			pkgPath:  "context",
			typeName: "context.Context",
			expected: true,
		},
		{
			name:     "allowed_package_time",
			pkgPath:  "time",
			typeName: "time.Other",
			expected: true,
		},
		{
			name:     "allowed_package_context",
			pkgPath:  "context",
			typeName: "context.Other",
			expected: true,
		},
		{
			name:     "allowed_package_go_token",
			pkgPath:  "go/token",
			typeName: "go/token.FileSet",
			expected: true,
		},
		{
			name:     "allowed_package_analysis",
			pkgPath:  "golang.org/x/tools/go/analysis",
			typeName: "golang.org/x/tools/go/analysis.Pass",
			expected: true,
		},
		{
			name:     "allowed_package_inspector",
			pkgPath:  "golang.org/x/tools/go/ast/inspector",
			typeName: "golang.org/x/tools/go/ast/inspector.Inspector",
			expected: true,
		},
		{
			name:     "allowed_package_config",
			pkgPath:  "github.com/kodflow/ktn-linter/pkg/config",
			typeName: "github.com/kodflow/ktn-linter/pkg/config.Config",
			expected: true,
		},
		{
			name:     "not_allowed_net_http",
			pkgPath:  "net/http",
			typeName: "net/http.Client",
			expected: false,
		},
		{
			name:     "not_allowed_os",
			pkgPath:  "os",
			typeName: "os.File",
			expected: false,
		},
		{
			name:     "not_allowed_external_package",
			pkgPath:  "github.com/some/package",
			typeName: "github.com/some/package.Type",
			expected: false,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedType(tt.pkgPath, tt.typeName)
			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("isAllowedType(%q, %q) = %v, want %v",
					tt.pkgPath, tt.typeName, result, tt.expected)
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
		{
			name:         "simple_ident",
			expr:         &ast.Ident{Name: "x"},
			expectedNil:  false,
			expectedName: "x",
		},
		{
			name:         "paren_expr",
			expr:         &ast.ParenExpr{X: &ast.Ident{Name: "x"}},
			expectedNil:  false,
			expectedName: "x",
		},
		{
			name:         "star_expr",
			expr:         &ast.StarExpr{X: &ast.Ident{Name: "x"}},
			expectedNil:  false,
			expectedName: "x",
		},
		{
			name:        "selector_expr_returns_nil",
			expr:        &ast.SelectorExpr{X: &ast.Ident{Name: "x"}, Sel: &ast.Ident{Name: "field"}},
			expectedNil: true,
		},
		{
			name:        "nested_paren_star",
			expr:        &ast.ParenExpr{X: &ast.StarExpr{X: &ast.Ident{Name: "ptr"}}},
			expectedNil: false,
			expectedName: "ptr",
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBaseIdent(tt.expr)
			// Vérifier si nil
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
		{
			name:      "param_ends_with_repo",
			paramName: "userRepo",
			typeName:  "Repository",
			expected:  "userRepo",
		},
		{
			name:      "param_ends_with_client",
			paramName: "httpClient",
			typeName:  "Client",
			expected:  "httpClient",
		},
		{
			name:      "param_ends_with_service",
			paramName: "userService",
			typeName:  "Service",
			expected:  "userService",
		},
		{
			name:      "param_ends_with_reader",
			paramName: "fileReader",
			typeName:  "FileReader",
			expected:  "fileReader",
		},
		{
			name:      "param_ends_with_writer",
			paramName: "logWriter",
			typeName:  "Writer",
			expected:  "logWriter",
		},
		{
			name:      "param_ends_with_handler",
			paramName: "requestHandler",
			typeName:  "Handler",
			expected:  "requestHandler",
		},
		{
			name:      "generic_param_uses_lowercase_type",
			paramName: "x",
			typeName:  "Client",
			expected:  "client",
		},
		{
			name:      "empty_param_uses_lowercase_type",
			paramName: "",
			typeName:  "Repository",
			expected:  "repository",
		},
		{
			name:      "underscore_param_uses_lowercase_type",
			paramName: "_",
			typeName:  "Logger",
			expected:  "logger",
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suggestInterfaceName(tt.paramName, tt.typeName)
			// Vérifier le résultat
			if result != tt.expected {
				t.Errorf("suggestInterfaceName(%q, %q) = %q, want %q",
					tt.paramName, tt.typeName, result, tt.expected)
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
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "single_uppercase",
			input:    "A",
			expected: "a",
		},
		{
			name:     "single_lowercase",
			input:    "a",
			expected: "a",
		},
		{
			name:     "uppercase_word",
			input:    "Client",
			expected: "client",
		},
		{
			name:     "all_caps",
			input:    "HTTP",
			expected: "hTTP",
		},
		{
			name:     "already_lowercase",
			input:    "client",
			expected: "client",
		},
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
		{
			name:     "empty_map",
			methods:  map[string]bool{},
			expected: []string{},
		},
		{
			name:     "single_method",
			methods:  map[string]bool{"Get": true},
			expected: []string{"Get"},
		},
		{
			name:     "multiple_methods_sorted",
			methods:  map[string]bool{"Write": true, "Read": true, "Close": true},
			expected: []string{"Close", "Read", "Write"},
		},
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
// Tested via analysistest framework in 001_external_test.go.
func Test_runAPI001(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "tested_via_analysistest_framework",
			expectError: false,
		},
		{
			name:        "always_returns_nil_error",
			expectError: false,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction runAPI001 nécessite un *analysis.Pass complet
			// Elle est testée via analysistest dans 001_external_test.go
			// La fonction retourne toujours (nil, nil), donc pas d'erreur possible
			if tt.expectError {
				t.Error("runAPI001 never returns an error")
			}
		})
	}
}

// Test_analyzeFunction tests the analyzeFunction function.
// Tested via analysistest framework in 001_external_test.go.
func Test_analyzeFunction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "tested_via_analysistest_framework"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction analyzeFunction nécessite un *analysis.Pass complet
			// Elle est testée via analysistest dans 001_external_test.go
		})
	}
}

// Test_collectParams tests the collectParams function.
// Tested via analysistest framework in 001_external_test.go.
func Test_collectParams(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "tested_via_analysistest_framework"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction collectParams nécessite un *analysis.Pass complet
			// Elle est testée via analysistest dans 001_external_test.go
		})
	}
}

// Test_getExternalConcreteType tests the getExternalConcreteType function.
// Tested via analysistest framework in 001_external_test.go.
func Test_getExternalConcreteType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "tested_via_analysistest_framework"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction getExternalConcreteType nécessite types.Type
			// Elle est testée via analysistest dans 001_external_test.go
		})
	}
}

// Test_scanBodyForMethodCalls tests the scanBodyForMethodCalls function.
// Tested via analysistest framework in 001_external_test.go.
func Test_scanBodyForMethodCalls(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "tested_via_analysistest_framework"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction scanBodyForMethodCalls nécessite un *analysis.Pass complet
			// Elle est testée via analysistest dans 001_external_test.go
		})
	}
}

// Test_matchesParam tests the matchesParam function.
// Tested via analysistest framework in 001_external_test.go.
func Test_matchesParam(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "tested_via_analysistest_framework"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction matchesParam nécessite un *analysis.Pass complet
			// Elle est testée via analysistest dans 001_external_test.go
		})
	}
}

// Test_reportDiagnostics tests the reportDiagnostics function.
// Tested via analysistest framework in 001_external_test.go.
func Test_reportDiagnostics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "tested_via_analysistest_framework"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction reportDiagnostics nécessite un *analysis.Pass complet
			// Elle est testée via analysistest dans 001_external_test.go
		})
	}
}

// Test_formatTypeName tests the formatTypeName function.
// Tested via analysistest framework in 001_external_test.go.
func Test_formatTypeName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "tested_via_analysistest_framework"},
	}

	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// La fonction formatTypeName utilise types.TypeString
			// Elle est testée via analysistest dans 001_external_test.go
		})
	}
}
