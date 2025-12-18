// Internal tests for analyzer 001 in ktninterface package.
package ktninterface

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

// Test_runInterface001 tests the main analyzer function with various scenarios.
func Test_runInterface001(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErrs int
	}{
		// === Private interface unused => report ===
		{
			name: "private interface unused",
			code: `package test
type myInterface interface {
	Method()
}`,
			wantErrs: 1,
		},

		// === Private interface used as struct field => no report ===
		{
			name: "private interface used as struct field",
			code: `package test
type myInterface interface {
	Method()
}
type S struct {
	field myInterface
}`,
			wantErrs: 0,
		},

		// === Private interface used as func param => no report ===
		{
			name: "private interface used as func param",
			code: `package test
type myInterface interface {
	Method()
}
func f(x myInterface) {}`,
			wantErrs: 0,
		},

		// === Private interface used as func return => no report ===
		{
			name: "private interface used as func return",
			code: `package test
type myInterface interface {
	Method()
}
func f() myInterface { return nil }`,
			wantErrs: 0,
		},

		// === Private interface used in var type => no report ===
		{
			name: "private interface used in var type",
			code: `package test
type myInterface interface {
	Method()
}
var x myInterface`,
			wantErrs: 0,
		},

		// === Private interface used in compile-time check => no report ===
		{
			name: "private interface used in compile-time check",
			code: `package test
type myInterface interface {
	Method()
}
type S struct{}
func (S) Method() {}
var _ myInterface = (*S)(nil)`,
			wantErrs: 0,
		},

		// === Exported interface unused => no report (default) ===
		{
			name: "exported interface unused - no report",
			code: `package test
type MyInterface interface {
	Method()
}`,
			wantErrs: 0,
		},

		// === Nested types: slice ===
		{
			name: "private interface used in slice",
			code: `package test
type myInterface interface { Method() }
type S struct { xs []myInterface }`,
			wantErrs: 0,
		},

		// === Nested types: map value ===
		{
			name: "private interface used in map value",
			code: `package test
type myInterface interface { Method() }
type S struct { m map[string]myInterface }`,
			wantErrs: 0,
		},

		// === Nested types: map key ===
		{
			name: "private interface used in map key",
			code: `package test
type myInterface interface { Method() }
type S struct { m map[myInterface]string }`,
			wantErrs: 0,
		},

		// === Nested types: pointer ===
		{
			name: "private interface used as pointer",
			code: `package test
type myInterface interface { Method() }
func f(x *myInterface) {}`,
			wantErrs: 0,
		},

		// === Nested types: channel ===
		{
			name: "private interface used in channel",
			code: `package test
type myInterface interface { Method() }
func f(c chan myInterface) {}`,
			wantErrs: 0,
		},

		// === Nested types: slice of slice ===
		{
			name: "private interface used in nested slice",
			code: `package test
type myInterface interface { Method() }
type S struct { xs [][]myInterface }`,
			wantErrs: 0,
		},

		// === Type assertion ===
		{
			name: "private interface used in type assertion",
			code: `package test
type myInterface interface { Method() }
func f(x interface{}) {
	_ = x.(myInterface)
}`,
			wantErrs: 0,
		},

		// === Multiple interfaces, one unused ===
		{
			name: "multiple interfaces one unused",
			code: `package test
type usedInterface interface { Method() }
type unusedInterface interface { Other() }
func f(x usedInterface) {}`,
			wantErrs: 1,
		},

		// === Interface embedding ===
		{
			name: "private interface used via embedding",
			code: `package test
type myInterface interface { Method() }
type composite interface { myInterface }`,
			wantErrs: 1, // myInterface is used (embedded), composite is unused
		},

		// === Method receiver does not count as usage ===
		{
			name: "method on interface type - still unused",
			code: `package test
type myInterface interface { Method() }
// No actual usage of myInterface as a type`,
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config
			config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Count diagnostics
			var diagnostics []analysis.Diagnostic
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					diagnostics = append(diagnostics, d)
				},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run analyzer
			_, runErr := runInterface001(pass)
			if runErr != nil {
				t.Fatalf("analyzer error: %v", runErr)
			}

			// Verify diagnostic count
			if len(diagnostics) != tt.wantErrs {
				t.Errorf("got %d diagnostics, want %d", len(diagnostics), tt.wantErrs)
				for _, d := range diagnostics {
					t.Logf("  diagnostic: %s", d.Message)
				}
			}
		})
	}
}

// Test_extractTypeIdents tests the recursive type extraction function.
func Test_extractTypeIdents(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		field    string
		expected []string
	}{
		{
			name:     "simple ident",
			code:     `package test; type S struct { x MyType }`,
			field:    "x",
			expected: []string{"MyType"},
		},
		{
			name:     "pointer",
			code:     `package test; type S struct { x *MyType }`,
			field:    "x",
			expected: []string{"MyType"},
		},
		{
			name:     "slice",
			code:     `package test; type S struct { x []MyType }`,
			field:    "x",
			expected: []string{"MyType"},
		},
		{
			name:     "map",
			code:     `package test; type S struct { x map[K]V }`,
			field:    "x",
			expected: []string{"K", "V"},
		},
		{
			name:     "channel",
			code:     `package test; type S struct { x chan MyType }`,
			field:    "x",
			expected: []string{"MyType"},
		},
		{
			name:     "ellipsis",
			code:     `package test; func f(x ...MyType) {}`,
			field:    "x",
			expected: []string{"MyType"},
		},
		{
			name:     "selector",
			code:     `package test; type S struct { x pkg.MyType }`,
			field:    "x",
			expected: []string{"MyType"},
		},
		{
			name:     "paren",
			code:     `package test; type S struct { x (MyType) }`,
			field:    "x",
			expected: []string{"MyType"},
		},
		{
			name:     "nil expr",
			code:     `package test; type S struct { x int }`,
			field:    "x",
			expected: []string{"int"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			// Find the field
			var fieldType ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				if f, ok := n.(*ast.Field); ok {
					for _, name := range f.Names {
						if name.Name == tt.field {
							fieldType = f.Type
							return false
						}
					}
				}
				return true
			})

			if fieldType == nil {
				t.Fatalf("field %q not found", tt.field)
			}

			result := extractTypeIdents(fieldType)

			// Check expected types are present
			for _, exp := range tt.expected {
				found := false
				for _, r := range result {
					if r == exp {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected %q in result, got %v", exp, result)
				}
			}
		})
	}
}

// Test_extractTypeIdents_nil tests nil handling.
func Test_extractTypeIdents_nil(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected []string
	}{
		{
			name:     "nil expr returns empty",
			expr:     nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTypeIdents(tt.expr)
			if len(result) != len(tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_collectInterfaces tests interface collection.
func Test_collectInterfaces(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expected    int
		excludeFile bool
	}{
		{
			name: "single interface",
			code: `package test
type MyInterface interface { Method() }`,
			expected:    1,
			excludeFile: false,
		},
		{
			name: "multiple interfaces",
			code: `package test
type A interface { Method() }
type B interface { Other() }`,
			expected:    2,
			excludeFile: false,
		},
		{
			name: "no interfaces",
			code: `package test
type S struct { x int }`,
			expected:    0,
			excludeFile: false,
		},
		{
			name: "mixed types",
			code: `package test
type MyInterface interface { Method() }
type S struct { x int }`,
			expected:    1,
			excludeFile: false,
		},
		{
			name: "file excluded - no interfaces collected",
			code: `package test
type MyInterface interface { Method() }`,
			expected:    0,
			excludeFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Rules: make(map[string]*config.RuleConfig),
			}

			// Configure file exclusion
			if tt.excludeFile {
				cfg.Rules["KTN-INTERFACE-001"] = &config.RuleConfig{
					Exclude: []string{"test.go"},
				}
			}
			config.Set(cfg)
			defer config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			interfaces := collectInterfaces(pass, cfg)

			if len(interfaces) != tt.expected {
				t.Errorf("got %d interfaces, want %d", len(interfaces), tt.expected)
			}
		})
	}
}

// Test_findUsages tests usage detection.
func Test_findUsages(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		ifaceName string
		wantUsed  bool
	}{
		{
			name: "used in struct field",
			code: `package test
type myI interface { M() }
type S struct { x myI }`,
			ifaceName: "myI",
			wantUsed:  true,
		},
		{
			name: "used in func param",
			code: `package test
type myI interface { M() }
func f(x myI) {}`,
			ifaceName: "myI",
			wantUsed:  true,
		},
		{
			name: "not used",
			code: `package test
type myI interface { M() }`,
			ifaceName: "myI",
			wantUsed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			// First collect interfaces
			interfaces := collectInterfaces(pass, config.Get())

			// Then find usages
			used := findUsages(pass, interfaces)

			if used[tt.ifaceName] != tt.wantUsed {
				t.Errorf("interface %q used=%v, want %v", tt.ifaceName, used[tt.ifaceName], tt.wantUsed)
			}
		})
	}
}

// Test_reportUnused tests the reporting logic.
func Test_reportUnused(t *testing.T) {
	tests := []struct {
		name       string
		ifaceName  string
		used       bool
		wantReport bool
	}{
		{
			name:       "private unused - report",
			ifaceName:  "myInterface",
			used:       false,
			wantReport: true,
		},
		{
			name:       "private used - no report",
			ifaceName:  "myInterface",
			used:       true,
			wantReport: false,
		},
		{
			name:       "exported unused - no report",
			ifaceName:  "MyInterface",
			used:       false,
			wantReport: false,
		},
		{
			name:       "exported used - no report",
			ifaceName:  "MyInterface",
			used:       true,
			wantReport: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
			}

			interfaces := map[string]*ast.TypeSpec{
				tt.ifaceName: {Name: &ast.Ident{Name: tt.ifaceName}},
			}

			used := map[string]bool{}
			if tt.used {
				used[tt.ifaceName] = true
			}

			reportUnused(pass, interfaces, used)

			gotReport := reportCount > 0
			if gotReport != tt.wantReport {
				t.Errorf("report=%v, want %v", gotReport, tt.wantReport)
			}
		})
	}
}

// Test_runInterface001_disabled tests that analyzer respects disabled config.
func Test_runInterface001_disabled(t *testing.T) {
	tests := []struct {
		name       string
		enabled    bool
		wantReport int
	}{
		{
			name:       "rule disabled",
			enabled:    false,
			wantReport: 0,
		},
		{
			name:       "rule enabled",
			enabled:    true,
			wantReport: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-INTERFACE-001": {Enabled: config.Bool(tt.enabled)},
				},
			})
			defer config.Reset()

			fset := token.NewFileSet()
			file, _ := parser.ParseFile(fset, "test.go", `package test
type unusedInterface interface { Method() }`, 0)

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(_ analysis.Diagnostic) {
					reportCount++
				},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			runInterface001(pass)

			if reportCount != tt.wantReport {
				t.Errorf("expected %d reports, got %d", tt.wantReport, reportCount)
			}
		})
	}
}

// Test_extractTypesFromNode tests the node type extraction.
func Test_extractTypesFromNode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "field node",
			code:     `package test; type S struct { x MyType }`,
			expected: 1,
		},
		{
			name:     "valuespec node with type",
			code:     `package test; var x MyType`,
			expected: 1,
		},
		{
			name:     "valuespec node without type",
			code:     `package test; var x = 42`,
			expected: 0,
		},
		{
			name:     "type assertion with type",
			code:     `package test; func f(x interface{}) { _ = x.(MyType) }`,
			expected: 1,
		},
		{
			name:     "type assertion type switch",
			code:     `package test; func f(x interface{}) { switch x.(type) {} }`,
			expected: 0,
		},
		{
			name:     "interface embedding",
			code:     `package test; type I interface { MyEmbed }`,
			expected: 1,
		},
		{
			name:     "type switch case clause",
			code:     `package test; func f(x interface{}) { switch x.(type) { case MyType: } }`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				types := extractTypesFromNode(n)
				count += len(types)
				return true
			})

			if count < tt.expected {
				t.Errorf("got %d types, want at least %d", count, tt.expected)
			}
		})
	}
}

// Test_extractCaseClauseTypes tests case clause type extraction.
func Test_extractCaseClauseTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "type switch with cases",
			code: `package test
func f(x interface{}) {
	switch x.(type) {
	case int, string:
	}
}`,
			expected: 2,
		},
		{
			name: "empty case",
			code: `package test
func f(x interface{}) {
	switch x.(type) {
	default:
	}
}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if cc, ok := n.(*ast.CaseClause); ok {
					types := extractCaseClauseTypes(cc)
					count += len(types)
				}
				return true
			})

			if count != tt.expected {
				t.Errorf("got %d types, want %d", count, tt.expected)
			}
		})
	}
}

// Test_extractEmbeddedTypes tests embedded type extraction.
func Test_extractEmbeddedTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "interface with embedding",
			code:     `package test; type I interface { MyEmbed; OtherEmbed }`,
			expected: 2,
		},
		{
			name:     "interface without embedding",
			code:     `package test; type I interface { Method() }`,
			expected: 0,
		},
		{
			name:     "empty interface",
			code:     `package test; type I interface {}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if it, ok := n.(*ast.InterfaceType); ok {
					types := extractEmbeddedTypes(it)
					count += len(types)
				}
				return true
			})

			if count != tt.expected {
				t.Errorf("got %d types, want %d", count, tt.expected)
			}
		})
	}
}

// Test_extractEmbeddedTypes_nil tests nil handling in extractEmbeddedTypes.
func Test_extractEmbeddedTypes_nil(t *testing.T) {
	// Test with nil Methods
	nilMethodsIface := &ast.InterfaceType{
		Methods: nil,
	}
	result := extractEmbeddedTypes(nilMethodsIface)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil Methods, got %v", result)
	}
}

// Test_extractMapTypes tests map type extraction.
func Test_extractMapTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "simple map",
			code:     `package test; type S struct { m map[K]V }`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if mt, ok := n.(*ast.MapType); ok {
					types := extractMapTypes(mt)
					count += len(types)
				}
				return true
			})

			if count != tt.expected {
				t.Errorf("got %d types, want %d", count, tt.expected)
			}
		})
	}
}

// Test_extractFuncTypes tests function type extraction.
func Test_extractFuncTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "func with params and returns",
			code:     `package test; type S struct { f func(A, B) (C, D) }`,
			expected: 4,
		},
		{
			name:     "func with no params",
			code:     `package test; type S struct { f func() R }`,
			expected: 1,
		},
		{
			name:     "func with no returns",
			code:     `package test; type S struct { f func(P) }`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if ft, ok := n.(*ast.FuncType); ok {
					types := extractFuncTypes(ft)
					count += len(types)
				}
				return true
			})

			if count != tt.expected {
				t.Errorf("got %d types, want %d", count, tt.expected)
			}
		})
	}
}

// Test_extractIndexTypes tests generic index type extraction.
func Test_extractIndexTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "generic with single param",
			code:     `package test; type S struct { x Container[T] }`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if ie, ok := n.(*ast.IndexExpr); ok {
					types := extractIndexTypes(ie)
					count += len(types)
				}
				return true
			})

			if count != tt.expected {
				t.Errorf("got %d types, want %d", count, tt.expected)
			}
		})
	}
}

// Test_extractIndexListTypes tests generic multi-param type extraction.
func Test_extractIndexListTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "generic with multiple params",
			code:     `package test; type S struct { x Container[T, U, V] }`,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if ile, ok := n.(*ast.IndexListExpr); ok {
					types := extractIndexListTypes(ile)
					count += len(types)
				}
				return true
			})

			if count != tt.expected {
				t.Errorf("got %d types, want %d", count, tt.expected)
			}
		})
	}
}

// Test_collectInterfacesFromFile tests interface collection from a file.
func Test_collectInterfacesFromFile(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "single interface",
			code:     `package test; type I interface { M() }`,
			expected: 1,
		},
		{
			name:     "no interfaces",
			code:     `package test; type S struct {}`,
			expected: 0,
		},
		{
			name:     "multiple interfaces",
			code:     `package test; type A interface {}; type B interface {}`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			interfaces := make(map[string]*ast.TypeSpec)
			collectInterfacesFromFile(file, interfaces)

			if len(interfaces) != tt.expected {
				t.Errorf("got %d interfaces, want %d", len(interfaces), tt.expected)
			}
		})
	}
}

// Test_findUsagesInFile tests usage detection in a file.
func Test_findUsagesInFile(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		iface    string
		expected bool
	}{
		{
			name:     "interface used in field",
			code:     `package test; type I interface {}; type S struct { x I }`,
			iface:    "I",
			expected: true,
		},
		{
			name:     "interface not used",
			code:     `package test; type I interface {}`,
			iface:    "I",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			interfaces := map[string]*ast.TypeSpec{
				tt.iface: {},
			}
			used := make(map[string]bool)

			findUsagesInFile(file, interfaces, used)

			if used[tt.iface] != tt.expected {
				t.Errorf("got used=%v, want %v", used[tt.iface], tt.expected)
			}
		})
	}
}

// Test_extractSimpleType tests simple type extraction.
func Test_extractSimpleType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "ident type",
			code:     `package test; type S struct { x MyType }`,
			expected: 1,
		},
		{
			name:     "selector type",
			code:     `package test; type S struct { x pkg.MyType }`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if f, ok := n.(*ast.Field); ok && f.Type != nil {
					result := extractSimpleType(f.Type)
					count += len(result)
				}
				return true
			})

			if count < tt.expected {
				t.Errorf("got %d types, want at least %d", count, tt.expected)
			}
		})
	}
}

// Test_extractRecursiveType tests recursive type extraction.
func Test_extractRecursiveType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "pointer type",
			code:     `package test; type S struct { x *MyType }`,
			expected: 1,
		},
		{
			name:     "array type",
			code:     `package test; type S struct { x []MyType }`,
			expected: 1,
		},
		{
			name:     "channel type",
			code:     `package test; type S struct { x chan MyType }`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if f, ok := n.(*ast.Field); ok && f.Type != nil {
					result := extractRecursiveType(f.Type)
					count += len(result)
				}
				return true
			})

			if count < tt.expected {
				t.Errorf("got %d types, want at least %d", count, tt.expected)
			}
		})
	}
}

// Test_extractCompositeType tests composite type extraction.
func Test_extractCompositeType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "map type",
			code:     `package test; type S struct { x map[K]V }`,
			expected: 2,
		},
		{
			name:     "func type",
			code:     `package test; type S struct { x func(A) B }`,
			expected: 2,
		},
		{
			name:     "generic single param",
			code:     `package test; type S struct { x Container[T] }`,
			expected: 2,
		},
		{
			name:     "generic multiple params",
			code:     `package test; type S struct { x Container[T, U] }`,
			expected: 3,
		},
		{
			name:     "unknown type returns empty",
			code:     `package test; type S struct { x int }`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var count int
			ast.Inspect(file, func(n ast.Node) bool {
				if f, ok := n.(*ast.Field); ok && f.Type != nil {
					result := extractCompositeType(f.Type)
					count += len(result)
				}
				return true
			})

			if count < tt.expected {
				t.Errorf("got %d types, want at least %d", count, tt.expected)
			}
		})
	}
}
