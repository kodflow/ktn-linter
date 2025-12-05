// Internal tests for analyzer 001 in ktninterface package.
package ktninterface

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

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
type UnusedInterface interface {
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

			collectDeclarations(pass, interfaces, structNames)

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
		name string
		code string
	}{
		{
			name: "error case - empty specs",
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

			// Verify no interfaces or structs collected for empty code
			if len(interfaces) != 0 {
				t.Errorf("expected 0 interfaces, got %d", len(interfaces))
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
		name string
		code string
	}{
		{
			name: "error case - interface with no embedded interfaces",
			code: `package test
type Simple interface {
	Method()
}`,
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

			// Verify no embedded interfaces found
			if len(usedInterfaces) != 0 {
				t.Errorf("expected 0 embedded interfaces, got %d", len(usedInterfaces))
			}
		})
	}
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

			reportUnusedInterfaces(pass, interfaces, usedInterfaces, structNames)

			// Verify no reports for empty interfaces
			if reportCount != 0 {
				t.Errorf("expected 0 reports, got %d", reportCount)
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
