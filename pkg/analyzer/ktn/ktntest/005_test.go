package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestIsTestFunc tests isTestFunc helper
func TestIsTestFunc(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		want     bool
	}{
		{"test function", "TestSomething", true},
		{"benchmark function", "BenchmarkSomething", true},
		{"fuzz function", "FuzzSomething", true},
		{"regular function", "Something", false},
		{"helper function", "testHelper", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: tt.funcName},
			}

			got := isTestFunc(funcDecl)
			if got != tt.want {
				t.Errorf("isTestFunc(%q) = %v, want %v", tt.funcName, got, tt.want)
			}
		})
	}
}

// TestIsAssertionMethod tests isAssertionMethod helper
func TestIsAssertionMethod(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		want       bool
	}{
		{"Error", "Error", true},
		{"Errorf", "Errorf", true},
		{"Fatal", "Fatal", true},
		{"Fatalf", "Fatalf", true},
		{"Fail", "Fail", true},
		{"FailNow", "FailNow", true},
		{"Log", "Log", false},
		{"Run", "Run", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAssertionMethod(tt.methodName)
			if got != tt.want {
				t.Errorf("isAssertionMethod(%q) = %v, want %v", tt.methodName, got, tt.want)
			}
		})
	}
}

// TestIsTestsVariableName tests isTestsVariableName helper
func TestIsTestsVariableName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"tests", true},
		{"testcases", true},
		{"cases", true},
		{"Tests", true},
		{"TestCases", true},
		{"myTests", false},
		{"data", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTestsVariableName(tt.name)
			if got != tt.want {
				t.Errorf("isTestsVariableName(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// TestCheckAssignStmt tests checkAssignStmt helper
func TestCheckAssignStmt(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"tests variable", "tests := []struct{}{}", true},
		{"cases variable", "cases := []struct{}{}", true},
		{"other variable", "data := []struct{}{}", false},
	}

	fset := token.NewFileSet()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\nfunc TestFoo() {\n" + tt.code + "\n}"
			file, err := parser.ParseFile(fset, "", code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			var assignStmt *ast.AssignStmt
			ast.Inspect(file, func(n ast.Node) bool {
				if as, ok := n.(*ast.AssignStmt); ok {
					assignStmt = as
					// Early return from iteration.
					return false
				}
				// Continue traversal
				return true
			})

			if assignStmt == nil {
				t.Fatal("No assignment found")
			}

			got := checkAssignStmt(assignStmt)
			if got != tt.want {
				t.Errorf("checkAssignStmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCheckRangeStmt tests checkRangeStmt helper
func TestCheckRangeStmt(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"range over tests", "for _, tt := range tests {}", true},
		{"range over cases", "for _, c := range cases {}", true},
		{"range over data", "for _, d := range data {}", false},
	}

	fset := token.NewFileSet()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\nfunc TestFoo() {\n" + tt.code + "\n}"
			file, err := parser.ParseFile(fset, "", code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			var rangeStmt *ast.RangeStmt
			ast.Inspect(file, func(n ast.Node) bool {
				if rs, ok := n.(*ast.RangeStmt); ok {
					rangeStmt = rs
					// Early return from iteration.
					return false
				}
				// Continue traversal
				return true
			})

			if rangeStmt == nil {
				t.Fatal("No range statement found")
			}

			got := checkRangeStmt(rangeStmt)
			if got != tt.want {
				t.Errorf("checkRangeStmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasTableDrivenPattern tests hasTableDrivenPattern helper
func TestHasTableDrivenPattern(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "with table-driven pattern",
			code: `
func TestFoo(t *testing.T) {
	tests := []struct{
		name string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
	}
}`,
			want: true,
		},
		{
			name: "without table-driven pattern",
			code: `
func TestFoo(t *testing.T) {
	if true {
		t.Error("fail")
	}
}`,
			want: false,
		},
	}

	fset := token.NewFileSet()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					// Early return from iteration.
					return false
				}
				// Continue traversal
				return true
			})

			if funcDecl == nil {
				t.Fatal("No function found")
			}

			got := hasTableDrivenPattern(funcDecl)
			if got != tt.want {
				t.Errorf("hasTableDrivenPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
