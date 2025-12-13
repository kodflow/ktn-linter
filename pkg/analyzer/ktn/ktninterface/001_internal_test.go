// Internal tests for analyzer 001 in ktninterface package.
package ktninterface

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runInterface001 tests the private runInterface001 function
func Test_runInterface001(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantErrs int
	}{
		{
			name: "unused interface",
			code: `package test
type unusedInterface interface {
	Method()
}`,
			wantErrs: 1,
		},
		{
			name: "used interface in function parameter",
			code: `package test
type UsedInterface interface {
	Method()
}
func UseIt(u UsedInterface) {}`,
			wantErrs: 0,
		},
		{
			name: "used interface in struct field",
			code: `package test
type UsedInterface interface {
	Method()
}
type MyStruct struct {
	field UsedInterface
}`,
			wantErrs: 0,
		},
		{
			name: "struct interface pattern - should not report",
			code: `package test
type User struct {
	Name string
}
type UserInterface interface {
	GetName() string
}`,
			wantErrs: 0,
		},
		{
			name: "interface with same name as struct - should not report",
			code: `package test
type Service struct {
	Name string
}
type Service interface {
	GetName() string
}`,
			wantErrs: 0,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					t.Logf("Report: %s at %s", d.Message, fset.Position(d.Pos))
				},
				ResultOf: make(map[*analysis.Analyzer]any),
			}

			// Run inspect analyzer first
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)
			pass.ResultOf[inspect.Analyzer] = inspectResult

			// Track reported errors
			errorCount := 0
			pass.Report = func(d analysis.Diagnostic) {
				errorCount++
			}

			// Run analyzer
			_, err = runInterface001(pass)
			// Check for execution errors
			if err != nil {
				t.Fatalf("runInterface001 failed: %v", err)
			}

			// Check error count matches expectation
			if errorCount != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, errorCount)
			}
		})
	}
}

// Test_collectDeclarations tests the private collectDeclarations function
func Test_collectDeclarations(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "collect interface and struct declarations",
			code: `package test
type MyInterface interface {
	Method()
}
type MyStruct struct {
	Field string
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			interfaces := make(map[string]*ast.TypeSpec)
			structNames := make(map[string]bool)

			cfg := config.Get()
			collectDeclarations(pass, cfg, interfaces, structNames)

			// Check interface was collected
			if _, exists := interfaces["MyInterface"]; !exists {
				t.Error("expected MyInterface to be collected")
			}

			// Check struct name was collected
			if !structNames["MyStruct"] {
				t.Error("expected MyStruct to be in structNames")
			}
		})
	}
}

// Test_isStructInterfacePattern tests the isStructInterfacePattern function
func Test_isStructInterfacePattern(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		structs       map[string]bool
		want          bool
	}{
		{
			name:          "matching pattern",
			interfaceName: "UserInterface",
			structs:       map[string]bool{"User": true},
			want:          true,
		},
		{
			name:          "no matching struct",
			interfaceName: "UserInterface",
			structs:       map[string]bool{"Post": true},
			want:          false,
		},
		{
			name:          "not ending with Interface",
			interfaceName: "Reader",
			structs:       map[string]bool{"Read": true},
			want:          false,
		},
		{
			name:          "too short name",
			interfaceName: "Interface",
			structs:       map[string]bool{},
			want:          false,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isStructInterfacePattern(tt.interfaceName, tt.structs)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("isStructInterfacePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_collectTypeSpecs tests the collectTypeSpecs private function.
func Test_collectTypeSpecs(t *testing.T) {
	tests := []struct {
		name              string
		code              string
		wantInterfaces    int
		wantStructs       int
	}{
		{
			name: "error case - empty specs",
			code: `package test`,
			wantInterfaces: 0,
			wantStructs: 0,
		},
		{
			name: "with imports - non TypeSpec should be skipped",
			code: `package test
import "fmt"
type MyInterface interface {
	Method()
}`,
			wantInterfaces: 1,
			wantStructs: 0,
		},
		{
			name: "with const declaration - non TypeSpec should be skipped",
			code: `package test
const MaxSize = 100
type Reader interface {
	Read() error
}
type Writer struct {
	Name string
}`,
			wantInterfaces: 1,
			wantStructs: 1,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			interfaces := make(map[string]*ast.TypeSpec)
			structNames := make(map[string]bool)

			ast.Inspect(file, func(n ast.Node) bool {
				// Check if general declaration
				if genDecl, ok := n.(*ast.GenDecl); ok {
					collectTypeSpecs(genDecl.Specs, interfaces, structNames)
				}
				// Continue traversal
				return true
			})

			// Verify interfaces count
			if len(interfaces) != tt.wantInterfaces {
				t.Errorf("expected %d interfaces, got %d", tt.wantInterfaces, len(interfaces))
			}

			// Verify structs count
			if len(structNames) != tt.wantStructs {
				t.Errorf("expected %d structs, got %d", tt.wantStructs, len(structNames))
			}
		})
	}
}

// Test_findInterfaceUsages tests the findInterfaceUsages private function.
func Test_findInterfaceUsages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - no interface usage",
			code: `package test
func NoInterfaces() {}`,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			usedInterfaces := make(map[string]bool)
			findInterfaceUsages(pass, usedInterfaces)

			// Verify no interfaces found
			if len(usedInterfaces) != 0 {
				t.Errorf("expected 0 used interfaces, got %d", len(usedInterfaces))
			}
		})
	}
}

// Test_checkNodeForInterfaceUsage tests the checkNodeForInterfaceUsage private function.
func Test_checkNodeForInterfaceUsage(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - non-interface node",
			code: `package test
const x = 1`,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usedInterfaces := make(map[string]bool)
			// Test with nil node
			checkNodeForInterfaceUsage(nil, usedInterfaces)
			// Verify no interfaces marked
			if len(usedInterfaces) != 0 {
				t.Errorf("expected 0 interfaces, got %d", len(usedInterfaces))
			}
		})
	}
}

// Test_checkFuncDeclForInterfaces tests the checkFuncDeclForInterfaces private function.
func Test_checkFuncDeclForInterfaces(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - function with no interfaces",
			code: `func Simple() {}`,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			usedInterfaces := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if function declaration
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					checkFuncDeclForInterfaces(funcDecl, usedInterfaces)
				}
				// Continue traversal
				return true
			})

			// Verify no interfaces found
			if len(usedInterfaces) != 0 {
				t.Errorf("expected 0 interfaces, got %d", len(usedInterfaces))
			}
		})
	}
}

// Test_checkEmbeddedInterfaces tests the checkEmbeddedInterfaces private function.
func Test_checkEmbeddedInterfaces(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		wantUsedCount int
	}{
		{
			name: "error case - interface with no embedded interfaces",
			code: `package test
type Simple interface {
	Method()
}`,
			wantUsedCount: 0,
		},
		{
			name: "interface with embedded interface",
			code: `package test
type Reader interface {
	Read() error
}
type Writer interface {
	Write() error
}
type ReadWriter interface {
	Reader
	Writer
}`,
			wantUsedCount: 2,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			usedInterfaces := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check if interface type
				if ifaceType, ok := n.(*ast.InterfaceType); ok {
					checkEmbeddedInterfaces(ifaceType, usedInterfaces)
				}
				// Continue traversal
				return true
			})

			// Verify expected count
			if len(usedInterfaces) != tt.wantUsedCount {
				t.Errorf("expected %d embedded interfaces, got %d", tt.wantUsedCount, len(usedInterfaces))
			}
		})
	}

	// Test special case - interface with nil Methods field
	t.Run("interface with nil Methods field", func(t *testing.T) {
		usedInterfaces := make(map[string]bool)
		interfaceType := &ast.InterfaceType{
			Methods: nil,
		}
		checkEmbeddedInterfaces(interfaceType, usedInterfaces)
		// Verify no interfaces marked
		if len(usedInterfaces) != 0 {
			t.Errorf("expected 0 interfaces with nil Methods, got %d", len(usedInterfaces))
		}
	})
}

// Test_reportUnusedInterfaces tests the reportUnusedInterfaces private function.
func Test_reportUnusedInterfaces(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - no unused interfaces",
			code: `package test`,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			interfaces := make(map[string]*ast.TypeSpec)
			usedInterfaces := make(map[string]bool)
			structNames := make(map[string]bool)
			compileTimeChecks := make(map[string]bool)

			reportUnusedInterfaces(pass, interfaces, usedInterfaces, structNames, compileTimeChecks)

			// Verify no reports for empty interfaces
			if reportCount != 0 {
				t.Errorf("expected 0 reports, got %d", reportCount)
			}
		})
	}
}

// Test_collectCompileTimeChecks tests the private collectCompileTimeChecks function.
func Test_collectCompileTimeChecks(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "compile time check found",
			code: `package test
type MyInterface interface { Method() }
var _ MyInterface = (*MyStruct)(nil)
type MyStruct struct{}`,
			expected: 1,
		},
		{
			name:     "no compile time check",
			code:     `package test; type MyInterface interface { Method() }`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parse result
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			checks := collectCompileTimeChecks(pass)
			// Verify check count
			if len(checks) != tt.expected {
				t.Errorf("expected %d checks, got %d", tt.expected, len(checks))
			}
		})
	}
}

// Test_extractInterfaceFromCheck tests the private extractInterfaceFromCheck function.
func Test_extractInterfaceFromCheck(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "extract interface from check",
			code:     `package test; var _ MyInterface = (*S)(nil)`,
			expected: "MyInterface",
		},
		{
			name:     "no interface in check",
			code:     `package test; var x int = 5`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			// Check parse result
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Find first var spec
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				// Check if GenDecl
				if !ok || genDecl.Tok != token.VAR {
					continue
				}
				// Check first spec
				if len(genDecl.Specs) > 0 {
					spec := genDecl.Specs[0].(*ast.ValueSpec)
					result := extractInterfaceFromCheck(spec)
					// Verify result
					if result != tt.expected {
						t.Errorf("expected %q, got %q", tt.expected, result)
					}
				}
			}
		})
	}
}

// Test_extractInterfaceNameFromExpr tests the private extractInterfaceNameFromExpr function.
func Test_extractInterfaceNameFromExpr(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{name: "identifier extraction"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with identifier
			ident := &ast.Ident{Name: "MyInterface"}
			result := extractInterfaceNameFromExpr(ident)
			// Verify result
			if result != "MyInterface" {
				t.Errorf("expected MyInterface, got %q", result)
			}

			// Test with nil
			result = extractInterfaceNameFromExpr(nil)
			// Verify nil result
			if result != "" {
				t.Errorf("expected empty, got %q", result)
			}
		})
	}
}

// Test_reportUnusedInterface tests the private reportUnusedInterface function.
func Test_reportUnusedInterface(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "report exported interface"},
		{name: "report private interface"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			reportCount := 0
			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			typeSpec := &ast.TypeSpec{Name: &ast.Ident{Name: "TestInterface"}}
			reportUnusedInterface(pass, typeSpec, "TestInterface")

			// Verify report was generated
			if reportCount != 1 {
				t.Errorf("expected 1 report, got %d", reportCount)
			}
		})
	}
}

// Test_checkFieldList tests the checkFieldList private function.
func Test_checkFieldList(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - empty field list",
			code: `package test
func Empty() {}`,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			used := make(map[string]bool)
			fieldList := &ast.FieldList{List: []*ast.Field{}}
			checkFieldList(fieldList, used)
			// Verify no types marked
			if len(used) != 0 {
				t.Errorf("expected 0 types, got %d", len(used))
			}
		})
	}
}

// Test_checkType tests the checkType private function.
func Test_checkType(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want int
	}{
		{
			name: "error case - nil expression",
			expr: nil,
			want: 0,
		},
		{
			name: "Ident type",
			expr: &ast.Ident{Name: "MyInterface"},
			want: 1,
		},
		{
			name: "StarExpr pointer type",
			expr: &ast.StarExpr{
				X: &ast.Ident{Name: "Reader"},
			},
			want: 1,
		},
		{
			name: "ArrayType slice type",
			expr: &ast.ArrayType{
				Elt: &ast.Ident{Name: "Writer"},
			},
			want: 1,
		},
		{
			name: "MapType with key and value",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "Handler"},
			},
			want: 2,
		},
		{
			name: "ChanType channel",
			expr: &ast.ChanType{
				Value: &ast.Ident{Name: "Message"},
			},
			want: 1,
		},
		{
			name: "SelectorExpr qualified type",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "io"},
				Sel: &ast.Ident{Name: "Reader"},
			},
			want: 1,
		},
		{
			name: "nested pointer to slice",
			expr: &ast.StarExpr{
				X: &ast.ArrayType{
					Elt: &ast.Ident{Name: "Service"},
				},
			},
			want: 1,
		},
		{
			name: "unknown expression type",
			expr: &ast.BasicLit{Value: "42"},
			want: 0,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			used := make(map[string]bool)
			checkType(tt.expr, used)
			// Verify expected count
			if len(used) != tt.want {
				t.Errorf("expected %d types, got %d", tt.want, len(used))
			}
		})
	}
}

// Test_hasCorrespondingStruct tests the hasCorrespondingStruct private function.
func Test_hasCorrespondingStruct(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		structs       map[string]bool
		want          bool
	}{
		{
			name:          "struct exists with same name",
			interfaceName: "UserService",
			structs:       map[string]bool{"UserService": true},
			want:          true,
		},
		{
			name:          "struct does not exist",
			interfaceName: "UserService",
			structs:       map[string]bool{"OrderService": true},
			want:          false,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasCorrespondingStruct(tt.interfaceName, tt.structs)
			// Verify result matches expectation
			if got != tt.want {
				t.Errorf("hasCorrespondingStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_checkValueSpec tests the checkValueSpec private function.
func Test_checkValueSpec(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int
	}{
		{
			name: "variable with interface type",
			code: `package test
var x MyInterface`,
			want: 1,
		},
		{
			name: "variable without explicit type",
			code: `package test
var x = 42`,
			want: 0,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for ValueSpec
				if vs, ok := n.(*ast.ValueSpec); ok {
					checkValueSpec(vs, used)
				}
				// Continue traversal
				return true
			})
			// Verify count
			if len(used) != tt.want {
				t.Errorf("expected %d types, got %d", tt.want, len(used))
			}
		})
	}
}

// Test_checkTypeAssert tests the checkTypeAssert private function.
func Test_checkTypeAssert(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int
	}{
		{
			name: "type assertion with interface",
			code: `package test
func f(x interface{}) {
	_ = x.(MyInterface)
}`,
			want: 1,
		},
		{
			name: "type assertion with concrete type",
			code: `package test
func f(x interface{}) {
	_ = x.(string)
}`,
			want: 1,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for TypeAssertExpr
				if ta, ok := n.(*ast.TypeAssertExpr); ok {
					checkTypeAssert(ta, used)
				}
				// Continue traversal
				return true
			})
			// Verify count
			if len(used) != tt.want {
				t.Errorf("expected %d types, got %d", tt.want, len(used))
			}
		})
	}
}

// Test_checkTypeSwitch tests the checkTypeSwitch private function.
func Test_checkTypeSwitch(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int
	}{
		{
			name: "type switch with multiple cases",
			code: `package test
func f(x interface{}) {
	switch x.(type) {
	case MyInterface:
	case OtherInterface:
	}
}`,
			want: 2,
		},
		{
			name: "type switch with nil body",
			code: `package test
func f(x interface{}) {
	switch x.(type) {
	}
}`,
			want: 0,
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Check parsing success
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			used := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Check for TypeSwitchStmt
				if ts, ok := n.(*ast.TypeSwitchStmt); ok {
					checkTypeSwitch(ts, used)
				}
				// Continue traversal
				return true
			})
			// Verify count
			if len(used) != tt.want {
				t.Errorf("expected %d types, got %d", tt.want, len(used))
			}
		})
	}
}
